from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any


@dataclass
class SourceEvent:
    source: str
    instance_uid: str
    label: str
    simulated_tick: datetime
    event_time: datetime
    active: bool
    tags: dict[str, str]
    fields: dict[str, Any]
    metadata: dict[str, Any] = field(default_factory=dict)
