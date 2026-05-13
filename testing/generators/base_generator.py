"""
@file base_generator.py
@brief Společný základ generátorů syntetických časových řad a dopravních událostí.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass
from datetime import datetime
import random
from typing import Any

from testing.generators.event_model import SourceEvent
from testing.generators.instance_pool import InstancePool, InstanceState
from testing.generators.timeline import SimulationTimeline


@dataclass
class GeneratorStats:
    source: str
    sent_messages: int
    raw_points: int
    kpi_points: int
    bootstrap_events: int
    inactive_events: int
    active_events: int
    unique_instances_seen: int
    activations: int
    max_active_instances: int

    def as_source_stats(self) -> dict[str, float]:
        return {
            "sent_messages": self.sent_messages,
            "raw_points": self.raw_points,
            "kpi_points": self.kpi_points,
            "bootstrap_events": self.bootstrap_events,
            "inactive_events": self.inactive_events,
            "active_events": self.active_events,
            "unique_instances_seen": self.unique_instances_seen,
            "activations": self.activations,
            "max_active_instances": self.max_active_instances,
        }


class BaseGenerator(ABC):
    source_name: str

    def __init__(self, source_name: str, payload: dict[str, Any], timeline: SimulationTimeline, load_profile: str, seed: int = 42):
        self.source_name = source_name
        self.payload = payload
        self.timeline = timeline
        self.load_profile = load_profile
        self.random = random.Random(seed)
        self.pool = InstancePool(source_name, int(payload["instance_count"]), seed=seed)
        self.bootstrap_events = 0
        self.inactive_events = 0
        self.active_events = 0
        self.activations = 0
        self.unique_instances_seen: set[str] = set()
        self.max_active_instances = 0
        self._emitted_events: list[SourceEvent] = []
        self._bootstrap_pool()

    def initialize(self, include_bootstrap: bool = True) -> None:
        if include_bootstrap:
            self._emit_bootstrap()

    def advance_tick(self) -> None:
        self._run_tick()
        self.timeline.advance()

    def finalize(self) -> None:
        self._finalize()

    def snapshot_stats(self) -> GeneratorStats:
        raw_points = self.bootstrap_events + self.active_events + self.inactive_events
        kpi_ratio = self.kpi_ratio()
        kpi_points = int(raw_points * kpi_ratio)
        return GeneratorStats(
            source=self.source_name,
            sent_messages=raw_points,
            raw_points=raw_points,
            kpi_points=kpi_points,
            bootstrap_events=self.bootstrap_events,
            inactive_events=self.inactive_events,
            active_events=self.active_events,
            unique_instances_seen=len(self.unique_instances_seen),
            activations=self.activations,
            max_active_instances=self.max_active_instances,
        )

    def simulate(self, duration_hours: float, include_bootstrap: bool = True) -> GeneratorStats:
        self.initialize(include_bootstrap=include_bootstrap)
        total_ticks = max(1, int(round(duration_hours * 60)))
        for _ in range(total_ticks):
            self.advance_tick()
        self.finalize()

        return self.snapshot_stats()

    def generate_rate_limited(self, duration_seconds: int, events_per_second: float, include_bootstrap: bool = False) -> GeneratorStats:
        self.initialize(include_bootstrap=include_bootstrap)
        remaining_fraction = 0.0
        total_ticks = max(1, int(duration_seconds))
        for _ in range(total_ticks):
            remaining_fraction += max(0.0, events_per_second)
            events_this_tick = int(remaining_fraction)
            remaining_fraction -= events_this_tick
            target_active = max(0, self.target_active_count(), events_this_tick)
            self._fill_active_set(target_active)
            self._ensure_minimum_active_count(events_this_tick)
            self._emit_sampled_active_snapshots(events_this_tick)
            self.timeline.advance()

        return self.snapshot_stats()

    def pop_emitted_events(self) -> list[SourceEvent]:
        events = self._emitted_events
        self._emitted_events = []
        return events

    def _run_tick(self) -> None:
        self._deactivate_completed()
        target_active = max(0, self.target_active_count())
        self._fill_active_set(target_active)
        self._emit_active_snapshots()

    def _emit_bootstrap(self) -> None:
        for instance in self.pool.iter_all():
            event = SourceEvent(
                source=self.source_name,
                instance_uid=instance.uid,
                label=instance.label,
                simulated_tick=self.timeline.start_time,
                event_time=self.timeline.bootstrap_time(instance.uid),
                active=False,
                tags=instance.tags,
                fields=self.inactive_fields(instance),
                metadata={"phase": "bootstrap"},
            )
            self._consume_event(event)
            self.bootstrap_events += 1

    def _deactivate_completed(self) -> None:
        for instance in list(self.pool.active_instances()):
            if instance.active_ticks >= instance.lifetime_target_ticks:
                self.pool.deactivate(instance.uid, self.timeline.now())
                self._emit_inactive(instance, phase="lifecycle")

    def _fill_active_set(self, target_active: int) -> None:
        deficit = max(0, target_active - self.pool.active_count)
        if deficit <= 0:
            return

        new_quota = min(deficit, self.desired_new_count())
        for uid in self.pool.take_never_activated(new_quota):
            instance = self.pool.activate(uid, self.timeline.now(), self.lifetime_ticks())
            self.activations += 1
            self.unique_instances_seen.add(instance.uid)

        remaining = max(0, target_active - self.pool.active_count)
        if remaining <= 0:
            return

        for uid in self.pool.take_returning(remaining):
            instance = self.pool.activate(uid, self.timeline.now(), self.lifetime_ticks())
            self.activations += 1
            self.unique_instances_seen.add(instance.uid)

    def _emit_active_snapshots(self) -> None:
        self.max_active_instances = max(self.max_active_instances, self.pool.active_count)
        for instance in self.pool.active_instances():
            self.update_active_fields(instance)
            event = SourceEvent(
                source=self.source_name,
                instance_uid=instance.uid,
                label=instance.label,
                simulated_tick=self.timeline.now(),
                event_time=self.timeline.event_time_for_tick(instance.uid, max_back_jitter_seconds=60),
                active=True,
                tags=instance.tags,
                fields=dict(instance.fields),
                metadata={"phase": "active"},
            )
            self._consume_event(event)
            instance.active_ticks += 1
            instance.last_emitted_at = event.event_time
            self.active_events += 1

    def _emit_sampled_active_snapshots(self, count: int) -> None:
        self.max_active_instances = max(self.max_active_instances, self.pool.active_count)
        if count <= 0:
            return
        active_instances = self.pool.active_instances()
        if not active_instances:
            return
        if count <= len(active_instances):
            selected_instances = self.random.sample(active_instances, count)
        else:
            selected_instances = self.random.choices(active_instances, k=count)
        for instance in selected_instances:
            self.update_active_fields(instance)
            event = SourceEvent(
                source=self.source_name,
                instance_uid=instance.uid,
                label=instance.label,
                simulated_tick=self.timeline.now(),
                event_time=self.timeline.event_time_for_tick(instance.uid, max_back_jitter_seconds=max(1, self.timeline.tick_seconds)),
                active=True,
                tags=instance.tags,
                fields=dict(instance.fields),
                metadata={"phase": "active", "generation_mode": "rate_limited"},
            )
            self._consume_event(event)
            instance.active_ticks += 1
            instance.last_emitted_at = event.event_time
            self.active_events += 1

    def _ensure_minimum_active_count(self, minimum_active: int) -> None:
        deficit = max(0, minimum_active - self.pool.active_count)
        if deficit <= 0:
            return
        for uid in self.pool.take_never_activated(deficit):
            instance = self.pool.activate(uid, self.timeline.now(), self.lifetime_ticks())
            self.activations += 1
            self.unique_instances_seen.add(instance.uid)
        remaining = max(0, minimum_active - self.pool.active_count)
        if remaining <= 0:
            return
        for uid in self.pool.take_returning(remaining):
            instance = self.pool.activate(uid, self.timeline.now(), self.lifetime_ticks())
            self.activations += 1
            self.unique_instances_seen.add(instance.uid)

    def _emit_inactive(self, instance: InstanceState, phase: str) -> None:
        event = SourceEvent(
            source=self.source_name,
            instance_uid=instance.uid,
            label=instance.label,
            simulated_tick=self.timeline.now(),
            event_time=self.timeline.event_time_for_tick(instance.uid, max_back_jitter_seconds=5),
            active=False,
            tags=instance.tags,
            fields=self.inactive_fields(instance),
            metadata={"phase": phase},
        )
        self._consume_event(event)
        self.inactive_events += 1

    def _finalize(self) -> None:
        for instance in list(self.pool.active_instances()):
            self.pool.deactivate(instance.uid, self.timeline.now())
            self._emit_inactive(instance, phase="finalize")

    def _consume_event(self, event: SourceEvent) -> None:
        self._emitted_events.append(event)

    def _bootstrap_pool(self) -> None:
        for index in range(1, self.pool.max_instances + 1):
            uid = self.instance_uid(index)
            instance = InstanceState(
                uid=uid,
                source=self.source_name,
                label=self.instance_label(index),
                tags=self.initial_tags(index),
                fields=self.initial_fields(index),
            )
            self.pool.register(instance)

    def load_multiplier(self) -> float:
        return {
            "normal": 1.0,
            "high": 1.35,
            "stress": 1.75,
        }.get(self.load_profile, 1.0)

    def kpi_ratio(self) -> float:
        return {
            "waze": 0.14,
            "ndic": 0.10,
            "mhd": 0.12,
        }.get(self.source_name, 0.12)

    @abstractmethod
    def target_active_count(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def desired_new_count(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def lifetime_ticks(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def initial_tags(self, index: int) -> dict[str, str]:
        raise NotImplementedError

    def instance_uid(self, index: int) -> str:
        return f'{self.payload["identifier_prefix"]}_{index:05d}'

    def instance_label(self, index: int) -> str:
        return self.instance_uid(index)

    @abstractmethod
    def initial_fields(self, index: int) -> dict[str, Any]:
        raise NotImplementedError

    @abstractmethod
    def inactive_fields(self, instance: InstanceState) -> dict[str, Any]:
        raise NotImplementedError

    @abstractmethod
    def update_active_fields(self, instance: InstanceState) -> None:
        raise NotImplementedError
