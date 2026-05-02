from __future__ import annotations

from datetime import datetime, timezone

from testing.generators.event_model import SourceEvent


SD_INSTANCE_REGISTRATION_QUEUE = "sd-instance-registration-requests"
TIME_SERIES_RAW_QUEUE = "time-series-raw-data"


class TSDBSeedAdapter:
    def __init__(self, publisher_client, runtime_state: dict[str, object]):
        self.publisher_client = publisher_client
        self.runtime_state = runtime_state
        self.runtime_state.setdefault("synthetic", {})

    def register_instances(self, source_name: str, sd_type_uid: str, instances: list[dict]) -> int:
        source_state = self.runtime_state["synthetic"].setdefault(source_name, {})
        known_instances = set(source_state.get("registered_instance_uids", []))
        registrations = 0
        batch: list[dict] = []
        for instance in instances:
            uid = instance["uid"]
            if uid in known_instances:
                continue
            batch.append(
                {
                    "eventTime": self._iso(instance["event_time"]),
                    "label": instance["label"],
                    "sdInstanceUID": uid,
                    "sdTypeUID": sd_type_uid,
                }
            )
            known_instances.add(uid)
            registrations += 1
            if len(batch) >= 500:
                self.publisher_client.publish_to_queue(SD_INSTANCE_REGISTRATION_QUEUE, batch)
                batch = []
        if batch:
            self.publisher_client.publish_to_queue(SD_INSTANCE_REGISTRATION_QUEUE, batch)
        source_state["registered_instance_uids"] = sorted(known_instances)
        return registrations

    def push_events(self, source_name: str, sd_type_uid: str, events: list[SourceEvent], raw_batch_size: int = 2000) -> dict[str, int]:
        published_events = 0
        published_batches = 0
        batch: list[dict] = []

        for event in events:
            batch.append(
                {
                    "eventTime": self._iso(event.event_time),
                    "sdInstanceUID": event.instance_uid,
                    "sdTypeID": sd_type_uid,
                    "fields": dict(event.fields),
                    "tags": dict(event.tags),
                }
            )
            if len(batch) >= raw_batch_size:
                self.publisher_client.publish_to_queue(TIME_SERIES_RAW_QUEUE, batch)
                published_events += len(batch)
                published_batches += 1
                batch = []

        if batch:
            self.publisher_client.publish_to_queue(TIME_SERIES_RAW_QUEUE, batch)
            published_events += len(batch)
            published_batches += 1

        return {
            "registrations": 0,
            "published_events": published_events,
            "published_batches": published_batches,
        }

    def _iso(self, value: datetime) -> str:
        if value.tzinfo is None:
            value = value.replace(tzinfo=timezone.utc)
        return value.isoformat()
