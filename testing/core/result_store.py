"""
@file result_store.py
@brief Ukládání surových i agregovaných výsledků experimentů do výstupních souborů.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Iterable, Optional

from testing.core.models import ScenarioRunResult, ScenarioSummary


class ResultStore:
    def __init__(self, ctx):
        self.raw_dir = ctx.run.output_dir / "raw"

    def save_results(self, results: Iterable[ScenarioRunResult]) -> None:
        for result in results:
            target = self.raw_dir / f"{result.experiment_id}_{result.scenario_id}_{result.repetition}.json"
            with target.open("w", encoding="utf-8") as handle:
                json.dump(result.to_dict(), handle, ensure_ascii=True, indent=2, sort_keys=True)

    def load_results(self, experiment_id: Optional[str] = None) -> list[ScenarioRunResult]:
        results: list[ScenarioRunResult] = []
        for path in sorted(self.raw_dir.glob("*.json")):
            if experiment_id is not None and not path.name.startswith(f"{experiment_id}_"):
                continue
            with path.open("r", encoding="utf-8") as handle:
                payload = json.load(handle)
            results.append(
                ScenarioRunResult(
                    run_id=payload["run_id"],
                    experiment_id=payload["experiment_id"],
                    scenario_id=payload["scenario_id"],
                    repetition=payload["repetition"],
                    started_at=__import__("datetime").datetime.fromisoformat(payload["started_at"]),
                    finished_at=__import__("datetime").datetime.fromisoformat(payload["finished_at"]),
                    duration_seconds=payload["duration_seconds"],
                    metrics=payload["metrics"],
                    metadata=payload.get("metadata", {}),
                    success=payload.get("success", True),
                    error=payload.get("error"),
                )
            )
        return results

    def save_summaries_json(self, summaries: list[ScenarioSummary], path: Path) -> None:
        with path.open("w", encoding="utf-8") as handle:
            json.dump([summary.to_dict() for summary in summaries], handle, ensure_ascii=True, indent=2, sort_keys=True)
