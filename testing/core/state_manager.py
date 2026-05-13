"""
@file state_manager.py
@brief Správa lokálního stavu připravených entit, instancí a KPI mezi běhy experimentů.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any


class StateManager:
    def __init__(self, state_dir: Path):
        self.state_dir = state_dir
        self.runtime_state_path = state_dir / "runtime_state.json"
        self.dataset_registry_path = state_dir / "datasets.json"

    def load_runtime_state(self) -> dict[str, Any]:
        return self._read_json(self.runtime_state_path, default={"sources": {}, "kpis": {}})

    def save_runtime_state(self, payload: dict[str, Any]) -> None:
        self._write_json(self.runtime_state_path, payload)

    def load_dataset_registry(self) -> dict[str, Any]:
        return self._read_json(self.dataset_registry_path, default={})

    def save_dataset_registry(self, payload: dict[str, Any]) -> None:
        self._write_json(self.dataset_registry_path, payload)

    def get_source_state(self, source_name: str) -> dict[str, Any]:
        state = self.load_runtime_state()
        sources = state.get("sources", {})
        if source_name not in sources:
            raise KeyError(f"Unknown source state: {source_name}")
        return sources[source_name]

    def get_kpi_id(self, kpi_alias: str) -> int:
        state = self.load_runtime_state()
        kpis = state.get("kpis", {})
        if kpi_alias not in kpis:
            raise KeyError(f"Unknown KPI alias: {kpi_alias}")
        return int(kpis[kpi_alias])

    def select_instance_ids(self, source_name: str, filter_mode: str, multi_count: int = 10) -> list[int]:
        source_state = self.get_source_state(source_name)
        instance_ids = [int(value) for value in source_state.get("instance_ids", [])]
        if filter_mode == "single_instance":
            return instance_ids[:1]
        if filter_mode == "multi_instance":
            return instance_ids[:multi_count]
        if filter_mode == "no_filter":
            return []
        raise ValueError(f"Unsupported filter mode: {filter_mode}")

    def _read_json(self, path: Path, default: dict[str, Any]) -> dict[str, Any]:
        if not path.exists():
            return default
        with path.open("r", encoding="utf-8") as handle:
            return json.load(handle)

    def _write_json(self, path: Path, payload: dict[str, Any]) -> None:
        with path.open("w", encoding="utf-8") as handle:
            json.dump(payload, handle, ensure_ascii=True, indent=2, sort_keys=True)
