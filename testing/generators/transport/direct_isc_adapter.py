"""
@file direct_isc_adapter.py
@brief Transportní adaptér pro přímé odesílání generovaných událostí do ISC RabbitMQ front.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from datetime import datetime, timezone

from testing.generators.event_model import SourceEvent


SD_INSTANCE_REGISTRATION_QUEUE = "sd-instance-registration-requests"
KPI_FULFILLMENT_QUEUE = "kpi-fulfillment-check-requests"


class DirectISCAdapter:
    def __init__(self, publisher_client, runtime_state: dict[str, object]):
        self.publisher_client = publisher_client
        self.runtime_state = runtime_state
        self.runtime_state.setdefault("synthetic", {})

    def push_events(self, source_name: str, sd_type_uid: str, events: list[SourceEvent]) -> dict[str, int]:
        source_state = self.runtime_state["synthetic"].setdefault(source_name, {})
        known_instances = set(source_state.get("registered_instance_uids", []))
        registrations = 0
        published = 0
        registration_batch: list[dict] = []
        request_batch: list[dict] = []

        for event in events:
            if event.instance_uid not in known_instances:
                registration_batch.append(
                    {
                        "eventTime": self._iso(event.event_time),
                        "label": event.label,
                        "sdInstanceUID": event.instance_uid,
                        "sdTypeUID": sd_type_uid,
                    }
                )
                known_instances.add(event.instance_uid)
                registrations += 1
                if len(registration_batch) >= 500:
                    self.publisher_client.publish_to_queue(SD_INSTANCE_REGISTRATION_QUEUE, registration_batch)
                    registration_batch = []

            params = dict(event.tags)
            params.update(event.fields)
            request_batch.append(
                {
                    "eventTime": self._iso(event.event_time),
                    "sdInstanceUID": event.instance_uid,
                    "sdTypeUID": sd_type_uid,
                    "parameters": params,
                }
            )
            published += 1
            if len(request_batch) >= 500:
                self.publisher_client.publish_to_queue(KPI_FULFILLMENT_QUEUE, request_batch)
                request_batch = []

        if registration_batch:
            self.publisher_client.publish_to_queue(SD_INSTANCE_REGISTRATION_QUEUE, registration_batch)
        if request_batch:
            self.publisher_client.publish_to_queue(KPI_FULFILLMENT_QUEUE, request_batch)

        source_state["registered_instance_uids"] = sorted(known_instances)
        return {"registrations": registrations, "published_events": published}

    def _iso(self, value: datetime) -> str:
        if value.tzinfo is None:
            value = value.replace(tzinfo=timezone.utc)
        return value.isoformat()
