"""
@file storage_experiment.py
@brief Experiment ověřující ukládání raw a KPI dat do historické časové databáze.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from testing.core.models import ScenarioRunResult, utcnow
from testing.experiments.base_experiment import BaseExperiment


class StorageExperiment(BaseExperiment):
    def prepare(self, ctx, scenario) -> None:
        registry = ctx.state_manager.load_dataset_registry()
        if scenario.dataset_id not in registry:
            raise RuntimeError(f"Dataset {scenario.dataset_id} is not prepared")

    def run_once(self, ctx, scenario, repetition):
        registry = ctx.state_manager.load_dataset_registry()
        dataset = registry[scenario.dataset_id]
        started = utcnow()
        raw_points = float(dataset["raw_points"])
        kpi_points = float(dataset["kpi_points"])
        storage = float(dataset["storage_size_bytes"])
        estimated = False

        if scenario.metadata.get("apply_extended_kpi_set"):
            base_kpi_count = max(len(ctx.config.kpis), 1)
            extended_kpi_count = len(ctx.config.extended_kpis)
            kpi_multiplier = (base_kpi_count + extended_kpi_count) / base_kpi_count
            storage_multiplier = 1.0 + ((kpi_multiplier - 1.0) * 0.6)
            kpi_points *= kpi_multiplier
            storage *= storage_multiplier
            estimated = True

        metrics = {
            "raw_points": raw_points,
            "kpi_points": kpi_points,
            "storage_size_bytes": storage,
            "bytes_per_raw_point": storage / raw_points if raw_points else 0.0,
            "bytes_per_kpi_point": storage / kpi_points if kpi_points else 0.0,
        }
        finished = utcnow()
        return ScenarioRunResult(
            run_id=ctx.run.run_id,
            experiment_id=scenario.experiment_id,
            scenario_id=scenario.id,
            repetition=repetition,
            started_at=started,
            finished_at=finished,
            duration_seconds=0.0,
            metrics=metrics,
            metadata={"dataset_id": scenario.dataset_id, "simulated": estimated, "estimated": estimated},
        )
