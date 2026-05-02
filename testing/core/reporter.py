from __future__ import annotations

import csv
from collections import defaultdict
from pathlib import Path
from statistics import mean, median

from testing.core.models import ScenarioRunResult, ScenarioSummary


class Reporter:
    def summarize(self, results: list[ScenarioRunResult]) -> list[ScenarioSummary]:
        grouped: dict[tuple[str, str], list[ScenarioRunResult]] = defaultdict(list)
        for result in results:
            grouped[(result.experiment_id, result.scenario_id)].append(result)

        summaries = []
        for (experiment_id, scenario_id), entries in grouped.items():
            metric_names = sorted({name for entry in entries for name in entry.metrics})
            metrics_avg = {}
            metrics_min = {}
            metrics_max = {}
            metrics_median = {}
            for metric_name in metric_names:
                values = [entry.metrics[metric_name] for entry in entries if metric_name in entry.metrics]
                if not values:
                    continue
                metrics_avg[metric_name] = mean(values)
                metrics_min[metric_name] = min(values)
                metrics_max[metric_name] = max(values)
                metrics_median[metric_name] = median(values)
            summaries.append(
                ScenarioSummary(
                    experiment_id=experiment_id,
                    scenario_id=scenario_id,
                    repetitions=len(entries),
                    metrics_avg=metrics_avg,
                    metrics_min=metrics_min,
                    metrics_max=metrics_max,
                    metrics_median=metrics_median,
                )
            )
        return sorted(summaries, key=lambda item: (item.experiment_id, item.scenario_id))

    def print_console_summary(self, summaries: list[ScenarioSummary]) -> None:
        if not summaries:
            print("No results found.")
            return
        for summary in summaries:
            print(f"{summary.experiment_id}/{summary.scenario_id} ({summary.repetitions} repetitions)")
            for metric_name, avg_value in sorted(summary.metrics_avg.items()):
                print(
                    f"  - {metric_name}: avg={avg_value:.4f} "
                    f"min={summary.metrics_min[metric_name]:.4f} "
                    f"max={summary.metrics_max[metric_name]:.4f}"
                )

    def save_markdown(self, summaries: list[ScenarioSummary], path: Path) -> None:
        path.parent.mkdir(parents=True, exist_ok=True)
        with path.open("w", encoding="utf-8") as handle:
            handle.write("| Experiment | Scenario | Repetitions | Metrics |\n")
            handle.write("|---|---|---:|---|\n")
            for summary in summaries:
                metrics = ", ".join(
                    f"{name}={value:.4f}" for name, value in sorted(summary.metrics_avg.items())
                )
                handle.write(
                    f"| {summary.experiment_id} | {summary.scenario_id} | {summary.repetitions} | {metrics} |\n"
                )

    def save_csv(self, summaries: list[ScenarioSummary], path: Path) -> None:
        path.parent.mkdir(parents=True, exist_ok=True)
        with path.open("w", encoding="utf-8", newline="") as handle:
            writer = csv.writer(handle)
            writer.writerow(["experiment_id", "scenario_id", "repetitions", "metric_name", "avg", "min", "max", "median"])
            for summary in summaries:
                for metric_name, avg_value in sorted(summary.metrics_avg.items()):
                    writer.writerow(
                        [
                            summary.experiment_id,
                            summary.scenario_id,
                            summary.repetitions,
                            metric_name,
                            avg_value,
                            summary.metrics_min[metric_name],
                            summary.metrics_max[metric_name],
                            summary.metrics_median[metric_name],
                        ]
                    )
