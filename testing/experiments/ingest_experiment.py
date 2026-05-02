from __future__ import annotations

import time

from testing.core.models import ScenarioRunResult, utcnow
from testing.experiments.base_experiment import BaseExperiment
from testing.generators.generator_manager import GeneratorManager
from testing.generators.transport.direct_isc_adapter import DirectISCAdapter


INGEST_QUEUE_NAMES = [
    "sd-instance-registration-requests",
    "kpi-fulfillment-check-requests",
    "time-series-raw-data",
    "time-series-kpi-results",
]


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
        )
        runtime_state = ctx.state_manager.load_runtime_state()
        adapter = DirectISCAdapter(ctx.rabbitmq_amqp_client, runtime_state)

        publish_started = time.perf_counter()
        source_stats: dict[str, dict[str, float]] = {}
        input_total = 0.0
        registrations_total = 0.0
        for source_name in ("waze", "ndic", "mhd"):
            events = manager.generate_window_events(
                duration_minutes=scenario.duration_minutes or 1,
                load_profile=profile,
                source_name=source_name,
                include_bootstrap=False,
            )
            result = adapter.push_events(
                source_name=source_name,
                sd_type_uid=ctx.config.sources[source_name].payload["sd_type_uid"],
                events=events,
            )
            source_stats[source_name] = {
                "generated_events": len(events),
                "registrations": result["registrations"],
                "published_events": result["published_events"],
            }
            input_total += float(result["published_events"])
            registrations_total += float(result["registrations"])
        publish_finished = time.perf_counter()
        idle_snapshot = ctx.rabbitmq_management_client.wait_for_queues_idle(
            INGEST_QUEUE_NAMES,
            poll_interval_seconds=ctx.env.poll_interval_seconds,
            consecutive_idle_polls=3,
            timeout_seconds=1800.0,
        )
        end_finished = time.perf_counter()
        ctx.state_manager.save_runtime_state(runtime_state)

        raw_total = float(simulated_window["raw_points"])
        kpi_total = float(simulated_window["kpi_points"])
        pipeline_time = end_finished - publish_started
        publish_time = publish_finished - publish_started
        max_queue_depth = max(idle_snapshot.get("_meta", {}).get("max_messages_seen", {}).values(), default=0)
        metrics = {
            "input_msgs_total": input_total,
            "input_msgs_per_s": input_total / duration_seconds if duration_seconds else 0.0,
            "raw_points_per_s": raw_total / duration_seconds if duration_seconds else 0.0,
            "kpi_points_per_s": kpi_total / duration_seconds if duration_seconds else 0.0,
            "avg_latency_ms": (pipeline_time / input_total * 1000.0) if input_total else 0.0,
            "max_queue_depth": float(max_queue_depth),
            "publish_time_s": publish_time,
            "pipeline_completion_time_s": pipeline_time,
            "registrations_total": registrations_total,
            "estimated_raw_points_total": raw_total,
            "estimated_kpi_points_total": kpi_total,
        }
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
            metadata={"profile": profile, "simulated": False, "source_stats": source_stats, "queue_idle_snapshot": idle_snapshot},
        )
