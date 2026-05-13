"""
@file timeline.py
@brief Tvorba časové osy událostí pro online ingest i historické seedování.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timedelta
import random


@dataclass
class SimulationTimeline:
    start_time: datetime
    tick_seconds: int = 60
    seed: int = 42

    def __post_init__(self) -> None:
        self.random = random.Random(self.seed)
        self.current_time = self.start_time
        self.last_event_time_by_instance: dict[str, datetime] = {}

    def now(self) -> datetime:
        return self.current_time

    def minute_of_day(self) -> int:
        return self.current_time.hour * 60 + self.current_time.minute

    def advance(self) -> datetime:
        self.current_time = self.current_time + timedelta(seconds=self.tick_seconds)
        return self.current_time

    def bootstrap_time(self, instance_uid: str) -> datetime:
        jitter_us = self.random.randint(-500_000, 500_000)
        candidate = self.start_time + timedelta(microseconds=jitter_us)
        return self._ensure_monotonic(instance_uid, candidate)

    def event_time_for_tick(self, instance_uid: str, max_back_jitter_seconds: int = 60) -> datetime:
        back_seconds = self.random.uniform(0, max_back_jitter_seconds)
        extra_us = self.random.randint(0, 999_999)
        candidate = self.current_time - timedelta(seconds=back_seconds, microseconds=extra_us)
        return self._ensure_monotonic(instance_uid, candidate)

    def _ensure_monotonic(self, instance_uid: str, candidate: datetime) -> datetime:
        previous = self.last_event_time_by_instance.get(instance_uid)
        if previous is not None and candidate <= previous:
            candidate = previous + timedelta(microseconds=1)
        self.last_event_time_by_instance[instance_uid] = candidate
        return candidate
