from __future__ import annotations

import csv
from datetime import datetime, timedelta
import math
from pathlib import Path
import time

from testing.core.models import ScenarioRunResult, utcnow
from testing.experiments.base_experiment import BaseExperiment
from testing.setup.payload_builders import build_history_export_payload
from testing.setup.setup_entities import EntitySetupService


REPROCESS_QUEUE_NAMES = [
    "kpi-reprocess-requests",
    "time-series-reprocess-read-request",
    "time-series-reprocess-read-response",
    "time-series-delete-requests-tsdb",
    "time-series-delete-requests-mpu",
    "kpi-fulfillment-check-results",
    "time-series-kpi-results",
]
KPI_REPROCESS_QUEUE = "kpi-reprocess-requests"
MPU_CONNECTION_NOTIFICATION_QUEUE = "message-processing-unit-connection-notifications"


class ReprocessExperiment(BaseExperiment):
    def prepare(self, ctx, scenario) -> None:
        registry = ctx.state_manager.load_dataset_registry()
        if scenario.dataset_id not in registry:
            raise RuntimeError(f"Dataset {scenario.dataset_id} is not prepared")
        if not ctx.rabbitmq_management_client.is_configured():
            raise RuntimeError("RabbitMQ management client is required for reprocess benchmarking")
        if not ctx.rabbitmq_amqp_client.is_configured():
            raise RuntimeError("RabbitMQ AMQP client is required for direct reprocess benchmarking")
        if scenario.kpi_id not in ctx.config.kpis:
            raise RuntimeError(f"Unknown KPI alias {scenario.kpi_id}")

    def run_once(self, ctx, scenario, repetition):
        started = utcnow()
        full_wall_start = time.perf_counter()
        stale_snapshot = _get_non_idle_reprocess_snapshot(ctx)
        if stale_snapshot:
            raise RuntimeError(
                "Reprocess queues are not clean before benchmark start. "
                f"Snapshot={stale_snapshot!r}. Restart/purge the reprocess pipeline first."
            )
        registry = ctx.state_manager.load_dataset_registry()
        dataset = registry[scenario.dataset_id]
        spec = ctx.config.kpis[scenario.kpi_id]
        source_state = ctx.state_manager.get_source_state(spec.source)
        if spec.mode == "SELECTED" and not source_state.get("instance_ids"):
            refreshed_state = ctx.state_manager.load_runtime_state()
            entity_service = EntitySetupService()
            if spec.source not in refreshed_state.get("sources", {}):
                refreshed_state = entity_service.ensure_sd_types(ctx, ctx.config)
            else:
                entity_service.refresh_instances_for_source(ctx, refreshed_state, spec.source)
            ctx.state_manager.save_runtime_state(refreshed_state)
            source_state = refreshed_state["sources"][spec.source]
        kpi_definition_id = ctx.state_manager.get_kpi_id(scenario.kpi_id)
        selected_instance_ids = _resolve_selected_instance_ids(ctx, source_state, spec, kpi_definition_id)
        if spec.mode == "SELECTED" and not selected_instance_ids:
            raise RuntimeError(f"KPI {scenario.kpi_id} requires selected instances but none were resolved")

        interval = _resolve_dataset_interval(dataset)
        raw_points = _estimate_or_count_raw_points(
            ctx=ctx,
            dataset=dataset,
            source_name=spec.source,
            sd_type_id=int(source_state["sd_type_id"]),
            selected_instance_ids=selected_instance_ids,
            interval=interval,
            scenario=scenario,
            repetition=repetition,
        )
        if raw_points <= 0:
            raise RuntimeError(
                f"Scenario {scenario.id} selected no raw points for dataset {scenario.dataset_id} "
                f"and KPI {scenario.kpi_id}. Check selected instances and dataset contents."
            )
        before_count, before_count_mode = _safe_count_kpi_points(
            ctx=ctx,
            dataset=dataset,
            source_name=spec.source,
            sd_type_id=int(source_state["sd_type_id"]),
            kpi_definition_id=kpi_definition_id,
            interval=interval,
            suffix=f"{scenario.id}_{repetition}_before",
        )
        before_export_time = time.perf_counter() - full_wall_start

        _request_kpi_config_refresh(ctx)
        reprocess_wait_start = time.perf_counter()
        _enqueue_reprocess_request(
            ctx=ctx,
            kpi_definition_id=kpi_definition_id,
            sd_type_uid=source_state["sd_type_uid"],
            selected_instance_uids=_resolve_selected_instance_uids(source_state, selected_instance_ids),
        )
        timeout_seconds = _resolve_reprocess_timeout(raw_points)
        telemetry_path = ctx.run.output_dir / "raw" / f"reprocess_trace_{scenario.id}_{repetition}.csv"
        telemetry_rows: list[dict[str, float | int | str]] = []
        idle_snapshot = {}

        def _on_poll(snapshot, elapsed, idle_count, all_idle):
            telemetry_rows.append(_flatten_queue_snapshot(snapshot, elapsed, idle_count, all_idle))

        try:
            idle_snapshot = _wait_for_reprocess_activity_and_idle(
                ctx=ctx,
                queue_names=REPROCESS_QUEUE_NAMES,
                poll_interval_seconds=min(ctx.env.poll_interval_seconds, 0.25),
                consecutive_idle_polls=3,
                timeout_seconds=timeout_seconds,
                on_poll=_on_poll,
            )
        finally:
            _write_telemetry_csv(telemetry_path, telemetry_rows)
        reprocess_wait_time = time.perf_counter() - reprocess_wait_start
        after_export_start = time.perf_counter()
        after_count, after_count_mode = _safe_count_kpi_points(
            ctx=ctx,
            dataset=dataset,
            source_name=spec.source,
            sd_type_id=int(source_state["sd_type_id"]),
            kpi_definition_id=kpi_definition_id,
            interval=interval,
            suffix=f"{scenario.id}_{repetition}_after",
        )
        after_export_time = time.perf_counter() - after_export_start
        total_time = time.perf_counter() - full_wall_start
        _validate_reprocess_result(
            scenario=scenario,
            raw_points=raw_points,
            before_count=before_count,
            after_count=after_count,
            idle_snapshot=idle_snapshot,
        )
        kpi_points = after_count if not math.isnan(after_count) else float(dataset["source_stats"][spec.source]["kpi_points"])
        kpi_points_written = max(after_count - before_count, 0.0)

        metrics = {
            "raw_points_processed": raw_points,
            "kpi_points_written": kpi_points_written,
            "kpi_points_before": before_count,
            "kpi_points_after": after_count,
            "total_time_s": total_time,
            "before_export_time_s": before_export_time,
            "reprocess_wait_time_s": reprocess_wait_time,
            "after_export_time_s": after_export_time,
            "raw_points_per_s": (raw_points / reprocess_wait_time) if reprocess_wait_time else 0.0,
            "kpi_points_per_s": (kpi_points / reprocess_wait_time) if reprocess_wait_time else 0.0,
        }
        finished = utcnow()
        return ScenarioRunResult(
            run_id=ctx.run.run_id,
            experiment_id=scenario.experiment_id,
            scenario_id=scenario.id,
            repetition=repetition,
            started_at=started,
            finished_at=finished,
            duration_seconds=total_time,
            metrics=metrics,
            metadata={
                "dataset_id": scenario.dataset_id,
                "kpi_id": scenario.kpi_id,
                "source": spec.source,
                "selected_instance_count": len(selected_instance_ids),
                "trigger": "direct_queue_reprocess_request",
                "queue_idle_snapshot": idle_snapshot,
                "queue_idle_timeout_s": timeout_seconds,
                "queue_trace_csv": str(telemetry_path),
                "queue_wait_elapsed_s": reprocess_wait_time,
                "before_count_mode": before_count_mode,
                "after_count_mode": after_count_mode,
                "simulated": False,
            },
        )


def _resolve_dataset_interval(dataset_payload: dict) -> tuple[datetime, datetime]:
    to_time = datetime.fromisoformat(dataset_payload["built_at"])
    from_time = to_time - timedelta(hours=int(dataset_payload.get("history_duration_hours", 0)))
    return from_time, to_time


def _resolve_selected_instance_ids(ctx, source_state: dict, spec, kpi_definition_id: int) -> list[str]:
    if spec.mode != "SELECTED":
        return []
    mapping = {
        uid: str(instance_id)
        for uid, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
    }
    resolved = [mapping[label] for label in spec.selected_instance_labels if label in mapping]
    if resolved:
        return resolved
    try:
        existing = ctx.graphql_client.get_kpi(kpi_definition_id)
    except Exception:
        existing = None
    if isinstance(existing, dict):
        selected_ids = existing.get("selectedSDInstanceIDs") or []
        if selected_ids:
            return [str(value) for value in selected_ids]
    fallback_count = min(len(spec.selected_instance_labels), len(source_state.get("instance_ids", [])))
    return [str(value) for value in source_state.get("instance_ids", [])[:fallback_count]]


def _resolve_selected_instance_uids(source_state: dict, selected_instance_ids: list[str]) -> list[str]:
    if not selected_instance_ids:
        return []
    mapping = {
        str(instance_id): uid
        for uid, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
    }
    return [mapping[value] for value in selected_instance_ids if value in mapping]


def _estimate_or_count_raw_points(ctx, dataset: dict, source_name: str, sd_type_id: int, selected_instance_ids: list[str], interval, scenario, repetition: int) -> float:
    if not selected_instance_ids:
        return float(dataset["source_stats"][source_name]["raw_points"])
    from_time, to_time = interval
    payload = build_history_export_payload(
        data_type="raw",
        sd_type_id=sd_type_id,
        sd_instance_ids=[int(value) for value in selected_instance_ids],
        kpi_definition_ids=[],
        from_iso=from_time.isoformat(),
        to_iso=to_time.isoformat(),
    )
    target_path = ctx.run.output_dir / "raw" / f"{scenario.experiment_id}_{scenario.id}_{repetition}_raw_count.csv"
    ctx.rest_client.timed_export_download(payload, target_path, aggregate=False)
    return float(_count_csv_rows(target_path))


def _count_kpi_points(ctx, sd_type_id: int, kpi_definition_id: int, interval: tuple[datetime, datetime], suffix: str) -> float:
    from_time, to_time = interval
    payload = build_history_export_payload(
        data_type="kpi",
        sd_type_id=sd_type_id,
        sd_instance_ids=[],
        kpi_definition_ids=[kpi_definition_id],
        from_iso=from_time.isoformat(),
        to_iso=to_time.isoformat(),
    )
    target_path = ctx.run.output_dir / "raw" / f"reprocess_kpi_count_{sd_type_id}_{kpi_definition_id}_{suffix}.csv"
    ctx.rest_client.timed_export_download(payload, target_path, aggregate=False)
    return float(_count_csv_rows(target_path))


def _safe_count_kpi_points(ctx, dataset: dict, source_name: str, sd_type_id: int, kpi_definition_id: int, interval: tuple[datetime, datetime], suffix: str) -> tuple[float, str]:
    try:
        return _count_kpi_points(ctx, sd_type_id, kpi_definition_id, interval, suffix), "export"
    except Exception:
        fallback = float(dataset["source_stats"][source_name]["kpi_points"])
        return fallback, "dataset_estimate"


def _count_csv_rows(path: Path) -> int:
    with path.open("r", encoding="utf-8", newline="") as handle:
        reader = csv.reader(handle)
        row_count = sum(1 for _ in reader)
    return max(row_count - 1, 0)


def _resolve_reprocess_timeout(raw_points: float) -> float:
    # Reprocess keeps the source queues unacked for the entire streaming job, so
    # small fixed timeouts produce false negatives on larger datasets.
    # Bound the timeout to keep failures finite while still scaling with volume.
    estimated = raw_points / 750.0
    return max(1800.0, min(3600.0, estimated))


def _request_kpi_config_refresh(ctx) -> None:
    ctx.rabbitmq_amqp_client.publish_to_queue(MPU_CONNECTION_NOTIFICATION_QUEUE, {})
    ctx.rabbitmq_management_client.wait_for_queues_idle(
        [MPU_CONNECTION_NOTIFICATION_QUEUE],
        poll_interval_seconds=min(ctx.env.poll_interval_seconds, 0.25),
        consecutive_idle_polls=2,
        timeout_seconds=60.0,
        require_rates_idle=False,
    )


def _enqueue_reprocess_request(ctx, kpi_definition_id: int, sd_type_uid: str, selected_instance_uids: list[str]) -> None:
    payload = {
        "jobId": "",
        "wait": False,
        "kpiDefinitionID": int(kpi_definition_id),
        "sdTypeUID": sd_type_uid,
        "to": utcnow().isoformat(),
    }
    if selected_instance_uids:
        payload["sdInstanceUIDs"] = selected_instance_uids
    ctx.rabbitmq_amqp_client.publish_to_queue(KPI_REPROCESS_QUEUE, payload)


def _validate_reprocess_result(
    scenario,
    raw_points: float,
    before_count: float,
    after_count: float,
    idle_snapshot: dict[str, dict],
) -> None:
    if raw_points <= 0:
        raise RuntimeError(
            f"Scenario {scenario.id} did not process any raw points. "
            "The scenario is not a valid reprocess benchmark."
        )
    saw_activity = bool(idle_snapshot.get("_meta", {}).get("saw_activity", False))
    if not saw_activity:
        raise RuntimeError(
            f"Scenario {scenario.id} did not observe any queue activity during reprocess waiting. "
            "The benchmark likely measured an idle pipeline instead of a real reprocess run."
        )
    if raw_points > 0 and after_count <= 0:
        raise RuntimeError(
            f"Scenario {scenario.id} lost KPI points during reprocess "
            f"(before={before_count}, after={after_count}). "
            "This indicates that the KPI update removed existing points but no valid recomputation followed."
        )


def _wait_for_reprocess_activity_and_idle(
    ctx,
    queue_names: list[str],
    poll_interval_seconds: float,
    consecutive_idle_polls: int,
    timeout_seconds: float,
    on_poll=None,
) -> dict[str, dict]:
    import time

    started = time.perf_counter()
    idle_count = 0
    saw_activity = False
    last_snapshot: dict[str, dict] = {}
    max_messages_seen = {queue_name: 0 for queue_name in queue_names}

    while True:
        snapshot = ctx.rabbitmq_management_client.snapshot_queues(queue_names)
        last_snapshot = snapshot
        all_idle = True

        for queue_name in queue_names:
            queue_snapshot = snapshot[queue_name]
            messages = int(queue_snapshot.get("messages", 0))
            messages_ready = int(queue_snapshot.get("messages_ready", 0))
            messages_unack = int(queue_snapshot.get("messages_unacknowledged", 0))
            publish_rate = float(queue_snapshot.get("publish_rate", 0.0))
            deliver_rate = float(queue_snapshot.get("deliver_rate", 0.0))

            max_messages_seen[queue_name] = max(max_messages_seen[queue_name], messages)

            if (
                messages > 0
                or messages_ready > 0
                or messages_unack > 0
                or publish_rate > 0.01
                or deliver_rate > 0.01
            ):
                saw_activity = True

            if messages > 0 or messages_ready > 0 or messages_unack > 0:
                all_idle = False

        elapsed_seconds = time.perf_counter() - started
        if on_poll is not None:
            on_poll(snapshot, elapsed_seconds, idle_count, all_idle)

        if saw_activity and all_idle:
            idle_count += 1
            if idle_count >= consecutive_idle_polls:
                last_snapshot["_meta"] = {
                    "max_messages_seen": max_messages_seen,
                    "saw_activity": True,
                }
                return last_snapshot
        elif not all_idle:
            idle_count = 0

        if elapsed_seconds > timeout_seconds:
            raise TimeoutError(
                f"Reprocess did not reach active-and-idle completion within {timeout_seconds} seconds. "
                f"saw_activity={saw_activity}, last_snapshot={last_snapshot!r}, max_messages_seen={max_messages_seen!r}"
            )

        time.sleep(poll_interval_seconds)


def _flatten_queue_snapshot(snapshot: dict[str, dict], elapsed_seconds: float, idle_count: int, all_idle: bool) -> dict[str, float | int | str]:
    row: dict[str, float | int | str] = {
        "elapsed_seconds": round(elapsed_seconds, 3),
        "idle_count": idle_count,
        "all_idle": int(all_idle),
        "messages_total": 0,
        "messages_ready_total": 0,
        "messages_unack_total": 0,
        "publish_rate_total": 0.0,
        "deliver_rate_total": 0.0,
    }
    for queue_name, values in snapshot.items():
        prefix = queue_name.replace("-", "_")
        messages = int(values.get("messages", 0))
        messages_ready = int(values.get("messages_ready", 0))
        messages_unack = int(values.get("messages_unacknowledged", 0))
        publish_rate = float(values.get("publish_rate", 0.0))
        deliver_rate = float(values.get("deliver_rate", 0.0))
        row["messages_total"] += messages
        row["messages_ready_total"] += messages_ready
        row["messages_unack_total"] += messages_unack
        row["publish_rate_total"] += publish_rate
        row["deliver_rate_total"] += deliver_rate
        row[f"{prefix}_messages"] = messages
        row[f"{prefix}_messages_ready"] = messages_ready
        row[f"{prefix}_messages_unack"] = messages_unack
        row[f"{prefix}_publish_rate"] = round(publish_rate, 3)
        row[f"{prefix}_deliver_rate"] = round(deliver_rate, 3)
    row["publish_rate_total"] = round(float(row["publish_rate_total"]), 3)
    row["deliver_rate_total"] = round(float(row["deliver_rate_total"]), 3)
    return row


def _write_telemetry_csv(path: Path, rows: list[dict[str, float | int | str]]) -> None:
    if not rows:
        return
    fieldnames: list[str] = []
    seen = set()
    for row in rows:
        for key in row.keys():
            if key not in seen:
                seen.add(key)
                fieldnames.append(key)
    with path.open("w", encoding="utf-8", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)


def _get_non_idle_reprocess_snapshot(ctx) -> dict[str, dict]:
    snapshot: dict[str, dict] = {}
    for queue_name in REPROCESS_QUEUE_NAMES:
        details = ctx.rabbitmq_management_client.get_queue_details(queue_name)
        messages = int(details.get("messages", 0))
        messages_ready = int(details.get("messages_ready", 0))
        messages_unacknowledged = int(details.get("messages_unacknowledged", 0))
        if messages > 0 or messages_ready > 0 or messages_unacknowledged > 0:
            snapshot[queue_name] = {
                "messages": messages,
                "messages_ready": messages_ready,
                "messages_unacknowledged": messages_unacknowledged,
            }
    return snapshot
