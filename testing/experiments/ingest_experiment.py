from __future__ import annotations

import csv
import threading
import time
from datetime import timedelta
from pathlib import Path

from testing.core.models import ScenarioRunResult, utcnow
from testing.experiments.base_experiment import BaseExperiment
from testing.generators.generator_manager import GeneratorManager
from testing.generators.transport.direct_isc_adapter import DirectISCAdapter
from testing.setup.payload_builders import build_history_export_payload


INGEST_QUEUE_NAMES = [
    "sd-instance-registration-requests",
    "kpi-fulfillment-check-requests",
    "time-series-raw-data",
    "time-series-kpi-results",
]
INGEST_QUEUE_SAMPLING_INTERVAL_SECONDS = 0.1


class IngestExperiment(BaseExperiment):
    def prepare(self, ctx, scenario) -> None:
        if not ctx.rabbitmq_management_client.is_configured():
            raise RuntimeError("RabbitMQ management client is required for ingest benchmarking")

    def run_once(self, ctx, scenario, repetition):
        started = utcnow()
        duration_seconds = float((scenario.duration_minutes or 1) * 60)
        profile = scenario.load_profile or "normal"
        manager = GeneratorManager.from_config(ctx.config.sources, seed=42 + repetition)
        simulated_window = manager.simulate_window(
            duration_minutes=scenario.duration_minutes or 1,
            load_profile=profile,
            sources=["waze", "ndic", "mhd"],
            rate_based=True,
        )
        runtime_state = ctx.state_manager.load_runtime_state()
        adapter = DirectISCAdapter(ctx.rabbitmq_amqp_client, runtime_state)

        queue_sampler = _QueueDepthSampler(
            ctx.rabbitmq_management_client,
            INGEST_QUEUE_NAMES,
            poll_interval_seconds=INGEST_QUEUE_SAMPLING_INTERVAL_SECONDS,
        )
        publish_started = time.perf_counter()
        queue_sampler.start()
        source_stats: dict[str, dict[str, float]] = {}
        event_times = []
        input_total = 0.0
        registrations_total = 0.0
        try:
            for source_name in ("waze", "ndic", "mhd"):
                events = manager.generate_window_events(
                    duration_minutes=scenario.duration_minutes or 1,
                    load_profile=profile,
                    source_name=source_name,
                    include_bootstrap=False,
                    rate_based=True,
                )
                event_times.extend(event.event_time for event in events)
                result = adapter.push_events(
                    source_name=source_name,
                    sd_type_uid=ctx.config.sources[source_name].payload["sd_type_uid"],
                    events=events,
                )
                source_stats[source_name] = {
                    "generated_events": len(events),
                    "registrations": result["registrations"],
                    "published_events": result["published_events"],
                    "configured_events_per_second": float(
                        ctx.config.sources[source_name].payload.get("profiles", {}).get(profile, {}).get("events_per_second", 0.0)
                    ),
                }
                input_total += float(result["published_events"])
                registrations_total += float(result["registrations"])
            publish_finished = time.perf_counter()
            idle_snapshot = ctx.rabbitmq_management_client.wait_for_queues_idle(
                INGEST_QUEUE_NAMES,
                poll_interval_seconds=INGEST_QUEUE_SAMPLING_INTERVAL_SECONDS,
                consecutive_idle_polls=3,
                timeout_seconds=1800.0,
                require_rates_idle=False,
            )
        finally:
            queue_sampler.stop()
        end_finished = time.perf_counter()
        ctx.state_manager.save_runtime_state(runtime_state)

        raw_total = float(simulated_window["raw_points"])
        kpi_total = float(simulated_window["kpi_points"])
        pipeline_time = end_finished - publish_started
        publish_time = publish_finished - publish_started
        queue_sampler_summary = queue_sampler.summary()
        idle_max_messages = idle_snapshot.get("_meta", {}).get("max_messages_seen", {})
        max_queue_depth = max(
            [
                *queue_sampler_summary.get("max_messages_seen", {}).values(),
                *idle_max_messages.values(),
            ],
            default=0,
        )
        metrics = {
            "input_msgs_total": input_total,
            "input_msgs_per_s": input_total / duration_seconds if duration_seconds else 0.0,
            "configured_input_rate_per_s": input_total / duration_seconds if duration_seconds else 0.0,
            "actual_publish_throughput_per_s": input_total / publish_time if publish_time else 0.0,
            "actual_pipeline_throughput_per_s": input_total / pipeline_time if pipeline_time else 0.0,
            "raw_points_per_s": raw_total / duration_seconds if duration_seconds else 0.0,
            "actual_raw_points_per_s": raw_total / pipeline_time if pipeline_time else 0.0,
            "kpi_points_per_s": kpi_total / duration_seconds if duration_seconds else 0.0,
            "avg_latency_ms": (pipeline_time / input_total * 1000.0) if input_total else 0.0,
            "max_queue_depth": float(max_queue_depth),
            "publish_time_s": publish_time,
            "pipeline_completion_time_s": pipeline_time,
            "registrations_total": registrations_total,
            "estimated_raw_points_total": raw_total,
            "estimated_kpi_points_total": kpi_total,
        }
        kpi_counts_path = ctx.run.output_dir / "raw" / f"kpi_result_counts_{scenario.experiment_id}_{scenario.id}_{repetition}.csv"
        kpi_counts = _count_kpi_results_by_alias(ctx, runtime_state, event_times, kpi_counts_path)
        finished = utcnow()
        return ScenarioRunResult(
            run_id=ctx.run.run_id,
            experiment_id=scenario.experiment_id,
            scenario_id=scenario.id,
            repetition=repetition,
            started_at=started,
            finished_at=finished,
            duration_seconds=duration_seconds,
            metrics=metrics,
            metadata={
                "profile": profile,
                "simulated": False,
                "source_stats": source_stats,
                "queue_depth_sampling": queue_sampler_summary,
                "queue_idle_snapshot": idle_snapshot,
                "actual_kpi_points_by_alias": kpi_counts,
                "missing_kpi_result_aliases": [alias for alias, count in kpi_counts.items() if count == 0],
                "kpi_counts_csv": str(kpi_counts_path),
            },
        )


class _QueueDepthSampler:
    def __init__(self, rabbitmq_management_client, queue_names: list[str], poll_interval_seconds: float) -> None:
        self.rabbitmq_management_client = rabbitmq_management_client
        self.queue_names = queue_names
        self.poll_interval_seconds = poll_interval_seconds
        self._stop_event = threading.Event()
        self._lock = threading.Lock()
        self._thread: threading.Thread | None = None
        self._samples = 0
        self._max_messages_seen = {queue_name: 0 for queue_name in queue_names}
        self._max_ready_seen = {queue_name: 0 for queue_name in queue_names}
        self._max_unacknowledged_seen = {queue_name: 0 for queue_name in queue_names}

    def start(self) -> None:
        self._thread = threading.Thread(target=self._run, name="ingest-queue-depth-sampler", daemon=True)
        self._thread.start()

    def stop(self) -> None:
        self._stop_event.set()
        if self._thread is not None:
            self._thread.join(timeout=max(1.0, self.poll_interval_seconds * 2))

    def summary(self) -> dict[str, object]:
        with self._lock:
            return {
                "poll_interval_seconds": self.poll_interval_seconds,
                "samples": self._samples,
                "max_messages_seen": dict(self._max_messages_seen),
                "max_ready_seen": dict(self._max_ready_seen),
                "max_unacknowledged_seen": dict(self._max_unacknowledged_seen),
            }

    def _run(self) -> None:
        while not self._stop_event.is_set():
            try:
                snapshot = self.rabbitmq_management_client.snapshot_queues(self.queue_names)
                with self._lock:
                    self._samples += 1
                    for queue_name, values in snapshot.items():
                        self._max_messages_seen[queue_name] = max(
                            self._max_messages_seen[queue_name],
                            int(values.get("messages", 0)),
                        )
                        self._max_ready_seen[queue_name] = max(
                            self._max_ready_seen[queue_name],
                            int(values.get("messages_ready", 0)),
                        )
                        self._max_unacknowledged_seen[queue_name] = max(
                            self._max_unacknowledged_seen[queue_name],
                            int(values.get("messages_unacknowledged", 0)),
                        )
            except Exception:
                # The benchmark should fail on the primary idle wait if RabbitMQ is unavailable.
                pass
            self._stop_event.wait(self.poll_interval_seconds)


def _count_kpi_results_by_alias(ctx, runtime_state: dict, event_times: list, target_path: Path) -> dict[str, int]:
    if not event_times:
        return {}
    from_time = min(event_times) - timedelta(seconds=1)
    to_time = max(event_times) + timedelta(seconds=1)
    target_path.parent.mkdir(parents=True, exist_ok=True)
    rows: list[dict[str, str | int]] = []
    counts: dict[str, int] = {}
    for alias, kpi_id in sorted(runtime_state.get("kpis", {}).items()):
        spec = ctx.config.kpis.get(alias)
        if spec is None:
            continue
        source_state = runtime_state.get("sources", {}).get(spec.source, {})
        sd_type_id = source_state.get("sd_type_id")
        if sd_type_id is None:
            continue
        export_path = target_path.with_name(f"{target_path.stem}_{alias}.csv")
        payload = build_history_export_payload(
            data_type="kpi",
            sd_type_id=int(sd_type_id),
            sd_instance_ids=[],
            kpi_definition_ids=[int(kpi_id)],
            from_iso=from_time.isoformat(),
            to_iso=to_time.isoformat(),
        )
        try:
            ctx.rest_client.timed_export_download(payload, export_path, aggregate=False)
            count = _count_csv_rows(export_path)
            status = "ok"
        except Exception as exc:
            count = 0
            status = f"error: {exc}"
        counts[alias] = count
        rows.append(
            {
                "alias": alias,
                "kpi_id": int(kpi_id),
                "label": spec.label,
                "source": spec.source,
                "points": count,
                "status": status,
            }
        )
    with target_path.open("w", encoding="utf-8", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=["alias", "kpi_id", "label", "source", "points", "status"])
        writer.writeheader()
        writer.writerows(rows)
    return counts


def _count_csv_rows(path: Path) -> int:
    with path.open("r", encoding="utf-8", newline="") as handle:
        reader = csv.reader(handle)
        row_count = sum(1 for _ in reader)
    return max(row_count - 1, 0)
