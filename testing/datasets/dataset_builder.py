from __future__ import annotations

from datetime import datetime, timedelta, timezone

from testing.core.models import DatasetRecord
from testing.generators.generator_manager import GeneratorManager
from testing.generators.transport.tsdb_seed_adapter import TSDBSeedAdapter
from testing.setup.setup_entities import EntitySetupService


DATASET_BUILD_QUEUE_NAMES = [
    "sd-instance-registration-requests",
    "time-series-raw-data",
]


class DatasetBuilder:
    def build(self, ctx, config, dataset, rebuild: bool = False, resume: bool = False) -> DatasetRecord:
        registry = ctx.state_manager.load_dataset_registry()
        if dataset.id in registry and not rebuild and not resume:
            payload = registry[dataset.id]
            return DatasetRecord(
                dataset_id=dataset.id,
                built_at=__import__("datetime").datetime.fromisoformat(payload["built_at"]),
                raw_points=payload["raw_points"],
                kpi_points=payload["kpi_points"],
                storage_size_bytes=payload["storage_size_bytes"],
                history_duration_hours=payload.get("history_duration_hours", 0),
                source_stats=payload.get("source_stats", {}),
            )

        checkpoint_datasets = self._resolve_checkpoint_datasets(config, dataset)
        if self._can_build_live(ctx):
            checkpoint_stats = self._build_live_checkpoints(ctx, config, checkpoint_datasets)
        else:
            manager = GeneratorManager.from_config(config.sources)
            checkpoint_stats = manager.simulate_checkpoints(checkpoint_datasets)

        for checkpoint_dataset in checkpoint_datasets:
            stats = checkpoint_stats[checkpoint_dataset.id]
            record = DatasetRecord(
                dataset_id=checkpoint_dataset.id,
                built_at=stats["built_at"],
                raw_points=int(stats.get("raw_points", 0)),
                kpi_points=int(stats.get("kpi_points", 0)),
                storage_size_bytes=int(stats.get("storage_size_bytes", 0)),
                history_duration_hours=int(checkpoint_dataset.history_duration_hours),
                source_stats=stats.get("source_stats", {}),
            )
            registry[checkpoint_dataset.id] = record.to_dict()

        ctx.state_manager.save_dataset_registry(registry)
        payload = registry[dataset.id]
        return DatasetRecord(
            dataset_id=dataset.id,
            built_at=__import__("datetime").datetime.fromisoformat(payload["built_at"]),
            raw_points=payload["raw_points"],
            kpi_points=payload["kpi_points"],
            storage_size_bytes=payload["storage_size_bytes"],
            history_duration_hours=payload.get("history_duration_hours", 0),
            source_stats=payload.get("source_stats", {}),
        )

    def resolve_checkpoint_ids(self, config, dataset) -> list[str]:
        return [item.id for item in self._resolve_checkpoint_datasets(config, dataset)]

    def _resolve_checkpoint_datasets(self, config, dataset) -> list:
        candidates = []
        for candidate in config.datasets.values():
            if candidate.sources != dataset.sources:
                continue
            if candidate.load_profile != dataset.load_profile:
                continue
            if candidate.history_duration_hours <= dataset.history_duration_hours:
                candidates.append(candidate)
        return sorted(candidates, key=lambda item: item.history_duration_hours)

    def _can_build_live(self, ctx) -> bool:
        return (
            ctx.rabbitmq_amqp_client.is_configured()
            and ctx.rabbitmq_amqp_client.healthcheck()
            and ctx.rabbitmq_management_client.is_configured()
            and ctx.rabbitmq_management_client.healthcheck()
        )

    def _build_live_checkpoints(self, ctx, config, checkpoint_datasets: list) -> dict[str, dict]:
        if not checkpoint_datasets:
            return {}

        # Refresh SDType IDs from the live backend before we start publishing
        # registrations. After stack restarts the persisted testing state may
        # still contain stale numeric IDs from a previous database.
        entity_service = EntitySetupService()
        runtime_state = entity_service.ensure_sd_types(ctx, config)
        ctx.state_manager.save_runtime_state(runtime_state)

        ordered = sorted(checkpoint_datasets, key=lambda item: item.history_duration_hours)
        max_dataset = ordered[-1]
        max_hours = max_dataset.history_duration_hours
        dataset_end = datetime.now(timezone.utc)
        dataset_start = dataset_end - timedelta(hours=max_hours)
        checkpoints_by_tick: dict[int, list] = {}
        for item in ordered:
            checkpoints_by_tick.setdefault(item.history_duration_hours * 60, []).append(item)

        manager = GeneratorManager.from_config(config.sources)
        generators = manager.build_generators(
            sources=max_dataset.sources,
            start_time=dataset_start,
            load_profile=max_dataset.load_profile,
        )

        adapter = TSDBSeedAdapter(ctx.rabbitmq_amqp_client, runtime_state)
        buffered_events: dict[str, list] = {source_name: [] for source_name in generators.keys()}
        records: dict[str, dict] = {}
        batch_ticks = 720

        for source_name, generator in generators.items():
            adapter.register_instances(
                source_name=source_name,
                sd_type_uid=config.sources[source_name].payload["sd_type_uid"],
                instances=[
                    {
                        "uid": instance.uid,
                        "label": instance.label,
                        "event_time": dataset_start,
                    }
                    for instance in generator.pool.iter_all()
                ],
            )

        for source_name, generator in generators.items():
            generator.initialize(include_bootstrap=True)
            buffered_events[source_name].extend(generator.pop_emitted_events())

        total_ticks = max_hours * 60
        for tick in range(1, total_ticks + 1):
            for source_name, generator in generators.items():
                generator.advance_tick()
                buffered_events[source_name].extend(generator.pop_emitted_events())

            if tick in checkpoints_by_tick:
                built_at = dataset_start + timedelta(minutes=tick)
                for item in checkpoints_by_tick[tick]:
                    records[item.id] = self._collect_checkpoint_record(
                        dataset=item,
                        built_at=built_at,
                        generators=generators,
                    )

            if tick % batch_ticks == 0:
                self._flush_live_events(ctx, config, adapter, runtime_state, buffered_events, wait_for_idle=False)

        for source_name, generator in generators.items():
            generator.finalize()
            buffered_events[source_name].extend(generator.pop_emitted_events())

        self._flush_live_events(ctx, config, adapter, runtime_state, buffered_events, wait_for_idle=True)
        ctx.state_manager.save_runtime_state(runtime_state)

        if max_dataset.id not in records:
            records[max_dataset.id] = self._collect_checkpoint_record(
                dataset=max_dataset,
                built_at=dataset_end,
                generators=generators,
            )

        refreshed_state = ctx.state_manager.load_runtime_state()
        for source_name in max_dataset.sources:
            if source_name in refreshed_state.get("sources", {}):
                entity_service.refresh_instances_for_source(ctx, refreshed_state, source_name)
        ctx.state_manager.save_runtime_state(refreshed_state)
        return records

    def _flush_live_events(self, ctx, config, adapter, runtime_state: dict, buffered_events: dict[str, list], wait_for_idle: bool) -> None:
        for source_name, events in buffered_events.items():
            if not events:
                continue
            adapter.push_events(
                source_name=source_name,
                sd_type_uid=config.sources[source_name].payload["sd_type_uid"],
                events=events,
            )
            buffered_events[source_name] = []

        if wait_for_idle:
            ctx.rabbitmq_management_client.wait_for_queues_idle(
                DATASET_BUILD_QUEUE_NAMES,
                poll_interval_seconds=ctx.env.poll_interval_seconds,
                consecutive_idle_polls=2,
                timeout_seconds=1800.0,
            )
        ctx.state_manager.save_runtime_state(runtime_state)

    def _collect_checkpoint_record(self, dataset, built_at: datetime, generators: dict) -> dict:
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
