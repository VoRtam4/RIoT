from __future__ import annotations

from typing import Any

from testing.generators.base_generator import BaseGenerator
from testing.generators.instance_pool import InstanceState
from testing.generators.profiles import MHDDailyProfile


class MHDGenerator(BaseGenerator):
    def __init__(self, source_name: str, payload: dict[str, Any], timeline, load_profile: str, seed: int = 42):
        super().__init__(source_name, payload, timeline, load_profile, seed=seed)
        self.profile = MHDDailyProfile()

    def target_active_count(self) -> int:
        base = 440.0 * self.profile.multiplier_for_minute_of_day(self.timeline.minute_of_day()) * self.load_multiplier()
        return min(self.pool.max_instances, max(0, int(round(base + self.random.randint(-20, 20)))))

    def desired_new_count(self) -> int:
        if self.target_active_count() == 0:
            return 0
        base = 16.0 * self.profile.multiplier_for_minute_of_day(self.timeline.minute_of_day()) * self.load_multiplier()
        return max(0, min(28, int(round(base + self.random.randint(-3, 4)))))

    def lifetime_ticks(self) -> int:
        return self.random.randint(30, 50)

    def initial_tags(self, index: int) -> dict[str, str]:
        return {
            "lineid": f"L{index % 120:03d}",
            "routeid": f"LR{index % 700:03d}",
            "service_date": "2026-04-23",
            "trip_id": f"TRIP-{index:05d}",
            "route_id": f"GTFS-R{index % 700:03d}",
            "direction_id": str(index % 2),
            "departure_time": f"{(index % 24):02d}:{(index % 60):02d}:00",
            "from_stop_id": f"STOP-{index % 500:04d}",
            "to_stop_id": f"STOP-{(index + 3) % 500:04d}",
            "segment_index": str(index % 25),
            "segment_from_stop_id": f"STOP-{index % 500:04d}",
            "segment_to_stop_id": f"STOP-{(index + 1) % 500:04d}",
            "finalstopid": f"STOP-{(index + 5) % 500:04d}",
            "vtype": ["tram", "bus", "trolleybus"][index % 3],
            "serviceDays": "[1,2,3,4,5]",
        }

    def initial_fields(self, index: int) -> dict[str, Any]:
        return {
            "globalid": f"G-{index:05d}",
            "id": f"VR-{index:05d}",
            "linename": f"Line {index % 120:03d}",
            "course": f"C-{index % 30:02d}",
            "ltype": ["tram", "bus", "night"][index % 3],
            "lf": "true",
            "laststopid": f"STOP-{index % 500:04d}",
            "lastupdate": 0.0,
            "lat": 49.2 + (index % 100) * 0.001,
            "lng": 16.6 + (index % 100) * 0.001,
            "bearing": 0.0,
            "delay": 0.0,
            "segment_progress": 0.0,
            "isinactive": True,
        }

    def inactive_fields(self, instance: InstanceState) -> dict[str, Any]:
        return {
            "globalid": instance.fields.get("globalid", ""),
            "id": instance.fields.get("id", ""),
            "linename": instance.fields.get("linename", ""),
            "course": instance.fields.get("course", ""),
            "ltype": instance.fields.get("ltype", ""),
            "lf": instance.fields.get("lf", "true"),
            "laststopid": instance.fields.get("laststopid", ""),
            "lastupdate": instance.fields.get("lastupdate", 0.0),
            "lat": instance.fields.get("lat", 0.0),
            "lng": instance.fields.get("lng", 0.0),
            "bearing": instance.fields.get("bearing", 0.0),
            "delay": instance.fields.get("delay", 0.0),
            "segment_progress": instance.fields.get("segment_progress", 0.0),
            "isinactive": True,
        }

    def update_active_fields(self, instance: InstanceState) -> None:
        previous_delay = float(instance.fields.get("delay", 0.0))
        delay_delta = self.random.uniform(-45.0, 45.0)
        delay = max(-120.0, min(900.0, previous_delay + delay_delta))
        segment_progress = min(1.0, max(0.0, float(instance.fields.get("segment_progress", 0.0)) + self.random.uniform(0.02, 0.15)))
        bearing = float((float(instance.fields.get("bearing", 0.0)) + self.random.uniform(-20.0, 20.0)) % 360.0)
        lat = float(instance.fields.get("lat", 49.2)) + self.random.uniform(-0.0007, 0.0007)
        lng = float(instance.fields.get("lng", 16.6)) + self.random.uniform(-0.0007, 0.0007)
        instance.fields.update(
            {
                "delay": round(delay, 2),
                "lastupdate": float(int(self.timeline.now().timestamp() * 1000)),
                "lat": round(lat, 6),
                "lng": round(lng, 6),
                "bearing": round(bearing, 2),
                "segment_progress": round(segment_progress, 3),
                "isinactive": False,
                "laststopid": instance.tags.get("segment_to_stop_id", instance.tags.get("to_stop_id", "")),
            }
        )
