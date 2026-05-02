from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml

from testing.core.models import (
    DatasetDefinition,
    EnvConfig,
    ExperimentDefinition,
    KPIDefinitionSpec,
    ScenarioDefinition,
    SourceSpec,
    TestingConfig,
)


def _load_yaml(path: Path) -> dict[str, Any]:
    with path.open("r", encoding="utf-8") as handle:
        return yaml.safe_load(handle) or {}


def load_testing_config(config_dir: Path) -> TestingConfig:
    env_raw = _load_yaml(config_dir / "env.yaml")
    sources_raw = _load_yaml(config_dir / "sources.yaml")
    kpis_raw = _load_yaml(config_dir / "kpis.yaml")
    datasets_raw = _load_yaml(config_dir / "datasets.yaml")
    experiments_raw = _load_yaml(config_dir / "experiments.yaml")

    env = EnvConfig(**env_raw)
    sources = {
        name: SourceSpec(name=name, payload=payload)
        for name, payload in sources_raw.get("sources", {}).items()
    }
    kpis = {
        item["id"]: KPIDefinitionSpec(**item)
        for item in kpis_raw.get("kpis", [])
    }
    extended_kpis = {
        item["id"]: KPIDefinitionSpec(**item)
        for item in kpis_raw.get("extended_kpis", [])
    }
    datasets = {
        item["id"]: DatasetDefinition(**item)
        for item in datasets_raw.get("datasets", [])
    }

    experiments: dict[str, ExperimentDefinition] = {}
    for raw_experiment in experiments_raw.get("experiments", []):
        scenario_objects = []
        for raw_scenario in raw_experiment.get("scenarios", []):
            scenario_objects.append(
                ScenarioDefinition(
                    experiment_id=raw_experiment["id"],
                    experiment_type=raw_experiment["type"],
                    repetitions=raw_experiment["repetitions"],
                    **raw_scenario,
                )
            )
        experiments[raw_experiment["id"]] = ExperimentDefinition(
            id=raw_experiment["id"],
            type=raw_experiment["type"],
            enabled=raw_experiment["enabled"],
            repetitions=raw_experiment["repetitions"],
            description=raw_experiment.get("description", ""),
            scenarios=scenario_objects,
        )

    return TestingConfig(
        env=env,
        sources=sources,
        kpis=kpis,
        extended_kpis=extended_kpis,
        datasets=datasets,
        experiments=experiments,
    )
