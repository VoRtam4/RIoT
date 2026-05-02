from __future__ import annotations

from typing import Any

from testing.generators.base_generator import BaseGenerator
from testing.generators.instance_pool import InstanceState
from testing.generators.profiles import WazeDailyProfile


class WazeGenerator(BaseGenerator):
    def __init__(self, source_name: str, payload: dict[str, Any], timeline, load_profile: str, seed: int = 42):
        super().__init__(source_name, payload, timeline, load_profile, seed=seed)
        self.profile = WazeDailyProfile()

    def target_active_count(self) -> int:
        base = 175.0 * self.profile.multiplier_for_minute_of_day(self.timeline.minute_of_day()) * self.load_multiplier()
        return min(self.pool.max_instances, max(25, int(round(base + self.random.randint(-10, 10)))))

    def desired_new_count(self) -> int:
        base = 12.0 * self.profile.multiplier_for_minute_of_day(self.timeline.minute_of_day()) * self.load_multiplier()
        return max(0, min(30, int(round(base + self.random.randint(-3, 4)))))

    def lifetime_ticks(self) -> int:
        return self.random.randint(3, 15)

    def initial_tags(self, index: int) -> dict[str, str]:
        city = ["Praha", "Brno", "Ostrava", "Plzen", "Olomouc"][index % 5]
        street = f"Street {index % 750:03d}"
        is_forward = "true" if index % 2 else "false"
        return {
            "country": "CZ",
            "city": city,
            "street": street,
            "roadType": ["highway", "arterial", "local"][index % 3],
            "segmentId": f"{index:05d}",
            "fromNode": f"{100000 + index * 2}",
            "toNode": f"{100001 + index * 2}",
            "isForward": is_forward,
        }

    def initial_fields(self, index: int) -> dict[str, Any]:
        _ = index
        return {
            "delay": 0.0,
            "speed": 0.0,
            "length": 0.0,
            "level": 0.0,
            "speedKPH": 0.0,
            "jamCount": 0.0,
            "pubMillisLatest": 0.0,
            "rawJams": "[]",
        }

    def inactive_fields(self, instance: InstanceState) -> dict[str, Any]:
        _ = instance
        return {
            "delay": 0.0,
            "length": 0.0,
            "level": 0.0,
            "speed": 0.0,
            "speedKPH": 0.0,
            "jamCount": 0.0,
            "pubMillisLatest": 0.0,
            "rawJams": "[]",
        }

    def instance_uid(self, index: int) -> str:
        suffix = "FWD" if index % 2 else "REV"
        return f"WAZE_JAM_{index}_{suffix}"

    def instance_label(self, index: int) -> str:
        tags = self.initial_tags(index)
        direction = "Forward" if tags["isForward"] == "true" else "Reverse"
        return f'{tags["city"]}, {tags["street"]}: {tags["segmentId"]} {direction}'

    def update_active_fields(self, instance: InstanceState) -> None:
        current_level = int(instance.fields.get("level", 0))
        level_shift = self.random.choice([-1, 0, 0, 1])
        level = max(1, min(5, current_level + level_shift if current_level else self.random.randint(1, 4)))
        speed = max(4.0, 65.0 - level * 9.5 + self.random.uniform(-4.0, 4.0))
        speed_kph = speed * 1.60934
        delay = max(5.0, level * 45.0 + self.random.uniform(-10.0, 15.0))
        length = max(80.0, instance.fields.get("length", 150.0) + self.random.uniform(-20.0, 20.0))
        jam_count = max(1.0, round(level + self.random.uniform(0.0, 2.0)))
        pub_millis = float(int(self.timeline.now().timestamp() * 1000))
        instance.fields.update(
            {
                "delay": round(delay, 2),
                "speed": round(speed, 2),
                "length": round(length, 2),
                "level": float(level),
                "speedKPH": round(speed_kph, 2),
                "jamCount": jam_count,
                "pubMillisLatest": pub_millis,
                "rawJams": '[{"uuid":"synthetic","pubMillis":%d}]' % int(pub_millis),
            }
        )
