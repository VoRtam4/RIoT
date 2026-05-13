"""
@file models.py
@brief Datové modely konfigurace, experimentů, událostí a naměřených výsledků.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Literal, Optional


ExperimentType = Literal["ingest", "reprocess", "history_read", "storage"]
FilterMode = Literal["single_instance", "multi_instance", "no_filter"]
LoadProfile = Literal["normal", "high", "stress"]
DataType = Literal["raw", "kpi", "aggregate_kpi"]


def utcnow() -> datetime:
    return datetime.now(timezone.utc)


@dataclass
class EnvConfig:
    base_url: str
    graphql_url: str
    rest_url: str
    api_key: str
    timeout_seconds: int = 30
    poll_interval_seconds: float = 2.0
    default_repetitions: int = 3
    output_dir: str = "testing/outputs"
    rabbitmq_amqp_url: Optional[str] = None
    rabbitmq_management_url: Optional[str] = None
    rabbitmq_user: Optional[str] = None
    rabbitmq_password: Optional[str] = None
    prometheus_url: Optional[str] = None


@dataclass
class SourceSpec:
    name: str
    payload: dict[str, Any]


@dataclass
class KPIDefinitionSpec:
    id: str
    label: str
    source: str
    mode: Literal["ALL", "SELECTED"]
    complexity: Literal["simple", "medium"]
    selected_instance_labels: list[str] = field(default_factory=list)
    definition: dict[str, Any] = field(default_factory=dict)


@dataclass
class DatasetDefinition:
    id: str
    label: str
    history_duration_hours: int
    sources: list[str]
    build_mode: str = "incremental"
    load_profile: LoadProfile = "normal"
    purpose: str = ""


@dataclass
class ScenarioDefinition:
    id: str
    experiment_id: str
    experiment_type: ExperimentType
    repetitions: int = 1
    dataset_id: Optional[str] = None
    load_profile: Optional[LoadProfile] = None
    duration_minutes: Optional[int] = None
    kpi_id: Optional[str] = None
    source: Optional[str] = None
    data_type: Optional[DataType] = None
    interval_hours: Optional[int] = None
    filter_mode: Optional[FilterMode] = None
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class ExperimentDefinition:
    id: str
    type: ExperimentType
    enabled: bool
    repetitions: int
    description: str
    scenarios: list[ScenarioDefinition]


@dataclass
class TestingConfig:
    env: EnvConfig
    sources: dict[str, SourceSpec]
    kpis: dict[str, KPIDefinitionSpec]
    extended_kpis: dict[str, KPIDefinitionSpec]
    datasets: dict[str, DatasetDefinition]
    experiments: dict[str, ExperimentDefinition]


@dataclass
class RunDescriptor:
    run_id: str
    started_at: datetime
    output_dir: Path


@dataclass
class ScenarioRunResult:
    run_id: str
    experiment_id: str
    scenario_id: str
    repetition: int
    started_at: datetime
    finished_at: datetime
    duration_seconds: float
    metrics: dict[str, float]
    metadata: dict[str, Any] = field(default_factory=dict)
    success: bool = True
    error: Optional[str] = None

    def to_dict(self) -> dict[str, Any]:
        out = asdict(self)
        out["started_at"] = self.started_at.isoformat()
        out["finished_at"] = self.finished_at.isoformat()
        return out


@dataclass
class ScenarioSummary:
    experiment_id: str
    scenario_id: str
    repetitions: int
    metrics_avg: dict[str, float]
    metrics_min: dict[str, float]
    metrics_max: dict[str, float]
    metrics_median: dict[str, float]

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass
class DatasetRecord:
    dataset_id: str
    built_at: datetime
    raw_points: int
    kpi_points: int
    storage_size_bytes: int
    history_duration_hours: int = 0
    source_stats: dict[str, dict[str, float]] = field(default_factory=dict)

    def to_dict(self) -> dict[str, Any]:
        out = asdict(self)
        out["built_at"] = self.built_at.isoformat()
        return out
