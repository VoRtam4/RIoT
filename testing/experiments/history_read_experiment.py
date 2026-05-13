"""
@file history_read_experiment.py
@brief Experiment měřící čtení historických raw a KPI dat přes API.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

import csv
from datetime import datetime, timedelta
from pathlib import Path

from testing.core.models import ScenarioRunResult, utcnow
from testing.experiments.base_experiment import BaseExperiment
from testing.setup.payload_builders import build_history_export_payload


class HistoryReadExperiment(BaseExperiment):
    def prepare(self, ctx, scenario) -> None:
        registry = ctx.state_manager.load_dataset_registry()
        if scenario.dataset_id not in registry:
            raise RuntimeError(f"Dataset {scenario.dataset_id} is not prepared")
        source_state = ctx.state_manager.get_source_state(scenario.source)
        if not source_state.get("sd_type_id"):
            raise RuntimeError(f"Missing SD type mapping for source {scenario.source}")
        if scenario.filter_mode != "no_filter" and not source_state.get("instance_ids"):
            raise RuntimeError(
                f"No instances found for source {scenario.source}. "
                "Run ingest first, then refresh setup and retry."
            )

    def run_once(self, ctx, scenario, repetition):
        registry = ctx.state_manager.load_dataset_registry()
        dataset = registry[scenario.dataset_id]
        started = utcnow()
        source_state = ctx.state_manager.get_source_state(scenario.source)
        interval_hours = int(scenario.interval_hours or 1)
        to_time = _resolve_dataset_end_time(dataset)
        from_time = to_time - timedelta(hours=interval_hours)
        kpi_definition_ids: list[int] = []
        kpi_alias: str | None = None
        if scenario.data_type in {"kpi", "aggregate_kpi"}:
            kpi_alias = scenario.metadata.get("kpi_id")
            if not kpi_alias:
                raise RuntimeError(f"Scenario {scenario.id} requires metadata.kpi_id")
            kpi_definition_ids = [ctx.state_manager.get_kpi_id(kpi_alias)]
        selected_instance_ids = _resolve_history_read_instance_ids(
            ctx=ctx,
            scenario=scenario,
            source_state=source_state,
            sd_type_id=int(source_state["sd_type_id"]),
            kpi_definition_ids=kpi_definition_ids,
            kpi_alias=kpi_alias,
            from_time=from_time,
            to_time=to_time,
        )

        payload = build_history_export_payload(
            data_type="kpi" if scenario.data_type in {"kpi", "aggregate_kpi"} else "raw",
            sd_type_id=int(source_state["sd_type_id"]),
            sd_instance_ids=selected_instance_ids,
            kpi_definition_ids=kpi_definition_ids,
            from_iso=from_time.isoformat(),
            to_iso=to_time.isoformat(),
            aggregate_seconds=scenario.metadata.get("aggregate_seconds"),
            sort_desc=scenario.metadata.get("sort_desc"),
        )
        target_path = ctx.run.output_dir / "raw" / f"{scenario.experiment_id}_{scenario.id}_{repetition}.csv"
        timing = ctx.rest_client.timed_export_download(
            payload,
            target_path,
            aggregate=(scenario.data_type == "aggregate_kpi"),
        )
        returned_points = float(_count_csv_rows(target_path))
        total_time = float(timing["total_time_s"])
        file_size = float(target_path.stat().st_size)
        metrics = {
            "returned_points": returned_points,
            "file_size_bytes": file_size,
            "total_time_s": total_time,
            "points_per_s": returned_points / total_time if total_time else 0.0,
            "bytes_per_s": file_size / total_time if total_time else 0.0,
            "export_ready_time_s": float(timing["export_ready_time_s"]),
            "download_time_s": float(timing["download_time_s"]),
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
                "source": scenario.source,
                "data_type": scenario.data_type,
                "filter_mode": scenario.filter_mode,
                "instance_count": len(selected_instance_ids),
                "download_url": timing["download_url"],
                "simulated": False,
            },
        )


def _resolve_dataset_end_time(dataset_payload: dict) -> datetime:
    built_at = dataset_payload.get("built_at")
    if built_at:
        return datetime.fromisoformat(built_at)
    return utcnow()


def _count_csv_rows(path: Path) -> int:
    with path.open("r", encoding="utf-8", newline="") as handle:
        reader = csv.reader(handle)
        row_count = sum(1 for _ in reader)
    return max(row_count - 1, 0)


def _resolve_history_read_instance_ids(ctx, scenario, source_state: dict, sd_type_id: int, kpi_definition_ids: list[int], kpi_alias: str | None, from_time: datetime, to_time: datetime) -> list[int]:
    filter_mode = scenario.filter_mode or "single_instance"
    if filter_mode == "no_filter":
        return []

    explicit_instance_uid = scenario.metadata.get("instance_uid")
    if explicit_instance_uid:
        uid_to_id = {
            str(uid): int(instance_id)
            for uid, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
        }
        if explicit_instance_uid not in uid_to_id:
            raise RuntimeError(
                f"Scenario {scenario.id} targets unknown instance UID {explicit_instance_uid!r}"
            )
        return [uid_to_id[explicit_instance_uid]]

    if scenario.data_type in {"kpi", "aggregate_kpi"} and kpi_definition_ids:
        if kpi_alias:
            selected_ids = _resolve_selected_kpi_instance_ids(ctx, source_state, kpi_alias, filter_mode)
            if selected_ids:
                return selected_ids
        instance_ids = _resolve_kpi_backed_instance_ids(
            ctx=ctx,
            source_state=source_state,
            sd_type_id=sd_type_id,
            kpi_definition_ids=kpi_definition_ids,
            from_time=from_time,
            to_time=to_time,
            filter_mode=filter_mode,
        )
        if instance_ids:
            return instance_ids

    return ctx.state_manager.select_instance_ids(scenario.source, filter_mode)


def _resolve_selected_kpi_instance_ids(ctx, source_state: dict, kpi_alias: str, filter_mode: str) -> list[int]:
    spec = ctx.config.kpis.get(kpi_alias)
    if spec is None or spec.mode != "SELECTED":
        return []
    uid_to_id = {
        str(uid): int(instance_id)
        for uid, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
    }
    matched_ids = [uid_to_id[label] for label in spec.selected_instance_labels if label in uid_to_id]
    if filter_mode == "single_instance":
        return matched_ids[:1]
    if filter_mode == "multi_instance":
        return matched_ids[:10]
    return matched_ids


def _resolve_kpi_backed_instance_ids(ctx, source_state: dict, sd_type_id: int, kpi_definition_ids: list[int], from_time: datetime, to_time: datetime, filter_mode: str) -> list[int]:
    payload = {
        "type": "kpi",
        "sdTypeID": sd_type_id,
        "kpiDefinitionIDs": kpi_definition_ids,
        "from": from_time.isoformat(),
        "to": to_time.isoformat(),
        "tag": "sdInstanceUID",
    }
    candidate_uids = ctx.rest_client.distinct_tag_values(payload)
    if not candidate_uids:
        return []

    uid_to_id = {
        str(uid): int(instance_id)
        for uid, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
    }
    if isinstance(candidate_uids, dict):
        candidate_uids = list(candidate_uids.keys())
    matched_ids = [uid_to_id[uid] for uid in candidate_uids if uid in uid_to_id]
    if filter_mode == "single_instance":
        return matched_ids[:1]
    if filter_mode == "multi_instance":
        return matched_ids[:10]
    return matched_ids
