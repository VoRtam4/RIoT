from __future__ import annotations

from typing import Any

from testing.generators.base_generator import BaseGenerator
from testing.generators.instance_pool import InstanceState
from testing.generators.profiles import FlatProfile


class NDICGenerator(BaseGenerator):
    def __init__(self, source_name: str, payload: dict[str, Any], timeline, load_profile: str, seed: int = 42):
        super().__init__(source_name, payload, timeline, load_profile, seed=seed)
        self.profile = FlatProfile()

    def target_active_count(self) -> int:
        base = 78.0 * self.load_multiplier()
        return min(self.pool.max_instances, max(4, int(round(base + self.random.randint(-6, 6)))))

    def desired_new_count(self) -> int:
        scale = self.load_multiplier()
        weights = [0.60, 0.28, 0.10, 0.02]
        if scale > 1.0:
            weights = [0.48, 0.30, 0.16, 0.06]
        return self.random.choices([0, 1, 2, 3], weights=weights, k=1)[0]

    def lifetime_ticks(self) -> int:
        return self.random.randint(60, 120)

    def initial_tags(self, index: int) -> dict[str, str]:
        region = ["PHA", "JHC", "JHM", "MSK", "OLK", "VYS"][index % 6]
        road_number = f"D{index % 80:02d}"
        return {
            "sourceIdentification": f"SRC-{index:05d}",
            "primaryLocationCode": f"PLC-{index:05d}",
            "secondaryLocationCode": f"SLC-{index:05d}",
            "tmcLocationCode": f"TMC-{index:05d}",
            "tmcPointName": f"Point {index:05d}",
            "tmcAreaRef": region,
            "tmcAreaName": region,
            "tmcRoadLCD": f"LCD-{index % 80:03d}",
            "tmcSegmentLCD": f"SEG-{index:05d}",
            "tmcRoadNumber": road_number,
            "tmcRoadName": road_number,
        }

    def initial_fields(self, index: int) -> dict[str, Any]:
        latitude = 49.0 + (index % 100) * 0.01
        longitude = 16.0 + (index % 100) * 0.01
        return {
            "trafficSpeedAnyVehicle": 0.0,
            "travelTimeAnyVehicle": 0.0,
            "trafficLevelAnyVehicle": 0.0,
            "isInactive": True,
            "tmcLatitude": round(latitude, 6),
            "tmcLongitude": round(longitude, 6),
        }

    def inactive_fields(self, instance: InstanceState) -> dict[str, Any]:
        return {
            "trafficSpeedAnyVehicle": instance.fields.get("trafficSpeedAnyVehicle", 0.0),
            "travelTimeAnyVehicle": instance.fields.get("travelTimeAnyVehicle", 0.0),
            "trafficLevelAnyVehicle": instance.fields.get("trafficLevelAnyVehicle", 0.0),
            "isInactive": True,
            "tmcLatitude": instance.fields.get("tmcLatitude", 0.0),
            "tmcLongitude": instance.fields.get("tmcLongitude", 0.0),
        }

    def instance_uid(self, index: int) -> str:
        return f"NDIC_TRAFFIC_SRC-{index:05d}"

    def update_active_fields(self, instance: InstanceState) -> None:
        previous_speed = float(instance.fields.get("trafficSpeedAnyVehicle", 50.0)) or 50.0
        speed = max(5.0, min(130.0, previous_speed + self.random.uniform(-6.0, 6.0)))
        level = max(1.0, min(5.0, round((140.0 - speed) / 30.0 + self.random.uniform(-0.4, 0.4))))
        travel_time = max(30.0, min(1200.0, 90.0 + (130.0 - speed) * 3.5 + self.random.uniform(-20.0, 20.0)))
        instance.fields.update(
            {
                "trafficSpeedAnyVehicle": round(speed, 2),
                "travelTimeAnyVehicle": round(travel_time, 2),
                "trafficLevelAnyVehicle": float(level),
                "isInactive": False,
            }
        )
