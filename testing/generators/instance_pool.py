"""
@file instance_pool.py
@brief Správa poolu syntetických instancí zdrojů dat pro generování událostí.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime
import random
from typing import Any


@dataclass
class InstanceState:
    uid: str
    source: str
    label: str
    tags: dict[str, str]
    fields: dict[str, Any]
    active: bool = False
    created_at: datetime | None = None
    activated_at: datetime | None = None
    deactivated_at: datetime | None = None
    last_emitted_at: datetime | None = None
    active_ticks: int = 0
    lifetime_target_ticks: int = 0


class InstancePool:
    def __init__(self, source_name: str, max_instances: int, seed: int = 42):
        self.source_name = source_name
        self.max_instances = max_instances
        self.random = random.Random(seed)
        self._instances: dict[str, InstanceState] = {}
        self._active: set[str] = set()
        self._inactive_seen: set[str] = set()
        self._never_activated: set[str] = set()

    def register(self, instance: InstanceState) -> None:
        self._instances[instance.uid] = instance
        self._never_activated.add(instance.uid)

    def iter_all(self) -> list[InstanceState]:
        return list(self._instances.values())

    def active_instances(self) -> list[InstanceState]:
        return [self._instances[uid] for uid in self._active]

    @property
    def active_count(self) -> int:
        return len(self._active)

    def activate(self, uid: str, current_time: datetime, lifetime_ticks: int) -> InstanceState:
        instance = self._instances[uid]
        instance.active = True
        instance.activated_at = current_time
        instance.active_ticks = 0
        instance.lifetime_target_ticks = lifetime_ticks
        if instance.created_at is None:
            instance.created_at = current_time
        self._active.add(uid)
        self._inactive_seen.discard(uid)
        self._never_activated.discard(uid)
        return instance

    def deactivate(self, uid: str, current_time: datetime) -> InstanceState:
        instance = self._instances[uid]
        instance.active = False
        instance.deactivated_at = current_time
        self._active.discard(uid)
        self._inactive_seen.add(uid)
        return instance

    def take_never_activated(self, count: int) -> list[str]:
        return self._take_random(self._never_activated, count)

    def take_returning(self, count: int) -> list[str]:
        return self._take_random(self._inactive_seen, count)

    def _take_random(self, source: set[str], count: int) -> list[str]:
        if count <= 0 or not source:
            return []
        if count >= len(source):
            return list(source)
        return self.random.sample(list(source), count)
