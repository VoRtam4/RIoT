from __future__ import annotations

from datetime import datetime, timedelta, timezone

from testing.core.models import SourceSpec
from testing.generators.base_generator import BaseGenerator, GeneratorStats
from testing.generators.mhd_generator import MHDGenerator
from testing.generators.ndic_generator import NDICGenerator
from testing.generators.timeline import SimulationTimeline
from testing.generators.waze_generator import WazeGenerator


class GeneratorManager:
    GENERATOR_MAP = {
        "waze": WazeGenerator,
        "mhd": MHDGenerator,
        "ndic": NDICGenerator,
    }

    def __init__(self, sources: dict[str, SourceSpec], seed: int = 42):
        self.sources = sources
        self.seed = seed

    @classmethod
    def from_config(cls, sources: dict[str, SourceSpec], seed: int = 42):
        return cls(sources, seed=seed)

    def build_generators(
        self,
        *,
        sources: list[str],
        start_time: datetime,
        load_profile: str,
    ) -> dict[str, BaseGenerator]:
        generators: dict[str, BaseGenerator] = {}
        for index, source_name in enumerate(sources):
            spec = self.sources[source_name]
            if not spec.payload.get("enabled", True):
                continue
            generator_class = self.GENERATOR_MAP[source_name]
            timeline = SimulationTimeline(start_time=start_time, tick_seconds=60, seed=self.seed + index)
            generators[source_name] = generator_class(
                source_name=source_name,
                payload=spec.payload,
                timeline=timeline,
                load_profile=load_profile,
                seed=self.seed + index,
            )
        return generators

    def simulate_dataset(self, dataset) -> dict:
        dataset_end = datetime.now(timezone.utc)
        dataset_start = dataset_end - timedelta(hours=dataset.history_duration_hours)

        source_stats: dict[str, dict[str, float]] = {}
        raw_points = 0
        kpi_points = 0
        storage_size = 0

        for index, source_name in enumerate(dataset.sources):
            spec = self.sources[source_name]
            if not spec.payload.get("enabled", True):
                continue
            generator_class = self.GENERATOR_MAP[source_name]
            timeline = SimulationTimeline(start_time=dataset_start, tick_seconds=60, seed=self.seed + index)
            generator: BaseGenerator = generator_class(
                source_name=source_name,
                payload=spec.payload,
                timeline=timeline,
                load_profile=dataset.load_profile,
                seed=self.seed + index,
            )
            stats: GeneratorStats = generator.simulate(dataset.history_duration_hours)
            source_stats[source_name] = stats.as_source_stats()
            raw_points += stats.raw_points
            kpi_points += stats.kpi_points
            storage_size += stats.raw_points * 110 + stats.kpi_points * 76

        return {
            "source_stats": source_stats,
            "raw_points": raw_points,
            "kpi_points": kpi_points,
            "storage_size_bytes": storage_size,
            "built_at": dataset_end,
        }

    def simulate_checkpoints(self, datasets: list) -> dict[str, dict]:
        if not datasets:
            return {}

        ordered = sorted(datasets, key=lambda item: item.history_duration_hours)
        max_hours = ordered[-1].history_duration_hours
        dataset_end = datetime.now(timezone.utc)
        dataset_start = dataset_end - timedelta(hours=max_hours)

        generators = self.build_generators(
            sources=ordered[-1].sources,
            start_time=dataset_start,
            load_profile=ordered[-1].load_profile,
        )
        for generator in generators.values():
            generator.initialize(include_bootstrap=True)

        checkpoints_by_tick: dict[int, list] = {}
        for dataset in ordered:
            checkpoints_by_tick.setdefault(dataset.history_duration_hours * 60, []).append(dataset)

        records: dict[str, dict] = {}
        total_ticks = max_hours * 60
        for tick in range(1, total_ticks + 1):
            for generator in generators.values():
                generator.advance_tick()

            if tick not in checkpoints_by_tick:
                continue

            built_at = dataset_start + timedelta(minutes=tick)
            for dataset in checkpoints_by_tick[tick]:
                records[dataset.id] = self._collect_checkpoint_record(
                    dataset=dataset,
                    built_at=built_at,
                    generators=generators,
                )

        for generator in generators.values():
            generator.finalize()

        final_tick = total_ticks
        for dataset in ordered:
            if dataset.id not in records and dataset.history_duration_hours * 60 == final_tick:
                records[dataset.id] = self._collect_checkpoint_record(
                    dataset=dataset,
                    built_at=dataset_end,
                    generators=generators,
                )

        return records

    def simulate_window(self, duration_minutes: int, load_profile: str, sources: list[str] | None = None) -> dict:
        dataset_end = datetime.now(timezone.utc)
        dataset_start = dataset_end - timedelta(minutes=duration_minutes)

        source_names = sources or list(self.sources.keys())
        source_stats: dict[str, dict[str, float]] = {}
        raw_points = 0
        kpi_points = 0

        for index, source_name in enumerate(source_names):
            spec = self.sources[source_name]
            if not spec.payload.get("enabled", True):
                continue
            generator_class = self.GENERATOR_MAP[source_name]
            timeline = SimulationTimeline(start_time=dataset_start, tick_seconds=60, seed=self.seed + index)
            generator: BaseGenerator = generator_class(
                source_name=source_name,
                payload=spec.payload,
                timeline=timeline,
                load_profile=load_profile,
                seed=self.seed + index,
            )
            stats: GeneratorStats = generator.simulate(duration_minutes / 60.0, include_bootstrap=False)
            source_stats[source_name] = stats.as_source_stats()
            raw_points += stats.raw_points
            kpi_points += stats.kpi_points

        return {
            "source_stats": source_stats,
            "raw_points": raw_points,
            "kpi_points": kpi_points,
            "built_at": dataset_end,
        }

    def generate_window_events(
        self,
        duration_minutes: int,
        load_profile: str,
        source_name: str,
        include_bootstrap: bool = False,
    ) -> list:
        dataset_end = datetime.now(timezone.utc)
        dataset_start = dataset_end - timedelta(minutes=duration_minutes)
        spec = self.sources[source_name]
        generator_class = self.GENERATOR_MAP[source_name]
        timeline = SimulationTimeline(start_time=dataset_start, tick_seconds=60, seed=self.seed)
        generator: BaseGenerator = generator_class(
            source_name=source_name,
            payload=spec.payload,
            timeline=timeline,
            load_profile=load_profile,
            seed=self.seed,
        )
        generator.initialize(include_bootstrap=include_bootstrap)
        events = generator.pop_emitted_events()
        total_ticks = max(1, duration_minutes)
        for _ in range(total_ticks):
            generator.advance_tick()
            events.extend(generator.pop_emitted_events())
        generator.finalize()
        events.extend(generator.pop_emitted_events())
        return events

    def _collect_checkpoint_record(self, dataset, built_at: datetime, generators: dict[str, BaseGenerator]) -> dict:
        source_stats: dict[str, dict[str, float]] = {}
        raw_points = 0
        kpi_points = 0
        storage_size = 0
        for source_name in dataset.sources:
            generator = generators[source_name]
            stats = generator.snapshot_stats()
            source_stats[source_name] = stats.as_source_stats()
            raw_points += stats.raw_points
            kpi_points += stats.kpi_points
            storage_size += stats.raw_points * 110 + stats.kpi_points * 76
        return {
            "source_stats": source_stats,
            "raw_points": raw_points,
            "kpi_points": kpi_points,
            "storage_size_bytes": storage_size,
            "built_at": built_at,
        }
