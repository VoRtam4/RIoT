from __future__ import annotations

import json
from urllib.parse import quote

import requests


class RabbitMQManagementClient:
    def __init__(self, base_url: str | None, username: str | None, password: str | None, timeout_seconds: int = 30):
        self.base_url = (base_url or "").rstrip("/")
        self.timeout_seconds = timeout_seconds
        self.session = requests.Session()
        if username is not None and password is not None:
            self.session.auth = (username, password)
        self.session.headers.update({"Content-Type": "application/json"})

    def is_configured(self) -> bool:
        return bool(self.base_url)

    def healthcheck(self) -> bool:
        if not self.is_configured():
            return False
        response = self.session.get(f"{self.base_url}/api/overview", timeout=self.timeout_seconds)
        return response.ok

    def publish_to_queue(self, queue_name: str, payload: dict) -> None:
        if not self.is_configured():
            raise RuntimeError("RabbitMQ management client is not configured")
        vhost = quote("/", safe="")
        exchange = quote("amq.default", safe="")
        endpoint = f"{self.base_url}/api/exchanges/{vhost}/{exchange}/publish"
        response = self.session.post(
            endpoint,
            json={
                "properties": {
                    "content_type": "application/json",
                    "delivery_mode": 2,
                },
                "routing_key": queue_name,
                "payload": json.dumps(payload, ensure_ascii=True),
                "payload_encoding": "string",
            },
            timeout=self.timeout_seconds,
        )
        if not response.ok:
            raise RuntimeError(f"Failed to publish to {queue_name}: {response.status_code} {response.text}")
        body = response.json()
        if not body.get("routed", False):
            raise RuntimeError(f"RabbitMQ did not route message to queue {queue_name}")

    def get_queue_details(self, queue_name: str) -> dict:
        if not self.is_configured():
            raise RuntimeError("RabbitMQ management client is not configured")
        vhost = quote("/", safe="")
        queue = quote(queue_name, safe="")
        endpoint = f"{self.base_url}/api/queues/{vhost}/{queue}"
        response = self.session.get(endpoint, timeout=self.timeout_seconds)
        if response.status_code == 404:
            return {
                "name": queue_name,
                "messages": 0,
                "messages_ready": 0,
                "messages_unacknowledged": 0,
                "message_stats": {},
            }
        if not response.ok:
            raise RuntimeError(f"Failed to fetch queue {queue_name}: {response.status_code} {response.text}")
        return response.json()

    def get_queue_depths(self, queue_names: list[str]) -> dict[str, int]:
        depths: dict[str, int] = {}
        for queue_name in queue_names:
            details = self.get_queue_details(queue_name)
            depths[queue_name] = int(details.get("messages", 0))
        return depths

    def snapshot_queues(self, queue_names: list[str]) -> dict[str, dict]:
        snapshot: dict[str, dict] = {}
        for queue_name in queue_names:
            details = self.get_queue_details(queue_name)
            snapshot[queue_name] = {
                "messages": int(details.get("messages", 0)),
                "messages_ready": int(details.get("messages_ready", 0)),
                "messages_unacknowledged": int(details.get("messages_unacknowledged", 0)),
                "consumers": int(details.get("consumers", 0)),
                "publish_rate": float(details.get("message_stats", {}).get("publish_details", {}).get("rate", 0.0) or 0.0),
                "deliver_rate": float(details.get("message_stats", {}).get("deliver_get_details", {}).get("rate", 0.0) or 0.0),
            }
        return snapshot

    def wait_for_queue_consumers(
        self,
        required_consumers: dict[str, int],
        poll_interval_seconds: float = 2.0,
        timeout_seconds: float = 180.0,
        on_poll=None,
    ) -> dict[str, dict]:
        import time

        started = time.perf_counter()
        last_snapshot: dict[str, dict] = {}
        queue_names = list(required_consumers.keys())
        while True:
            snapshot = self.snapshot_queues(queue_names)
            last_snapshot = snapshot
            all_ready = True
            for queue_name, minimum_consumers in required_consumers.items():
                consumers = int(snapshot[queue_name]["consumers"])
                if consumers < minimum_consumers:
                    all_ready = False
                    break

            elapsed_seconds = time.perf_counter() - started
            if on_poll is not None:
                on_poll(snapshot, elapsed_seconds, all_ready)

            if all_ready:
                return snapshot

            if elapsed_seconds > timeout_seconds:
                raise TimeoutError(
                    f"Queues did not reach required consumer counts within {timeout_seconds} seconds. "
                    f"Required={required_consumers!r}, last_snapshot={last_snapshot!r}"
                )

            time.sleep(poll_interval_seconds)

    def wait_for_queues_idle(
        self,
        queue_names: list[str],
        poll_interval_seconds: float = 2.0,
        consecutive_idle_polls: int = 3,
        timeout_seconds: float = 1800.0,
        require_rates_idle: bool = True,
        on_poll=None,
    ) -> dict[str, dict]:
        import time

        started = time.perf_counter()
        idle_count = 0
        last_snapshot: dict[str, dict] = {}

        max_messages_seen = {queue_name: 0 for queue_name in queue_names}
        while True:
            all_idle = True
            snapshot = self.snapshot_queues(queue_names)
            for queue_name in queue_names:
                queue_snapshot = snapshot[queue_name]
                messages = int(queue_snapshot["messages"])
                publish_rate = float(queue_snapshot["publish_rate"])
                deliver_rate = float(queue_snapshot["deliver_rate"])
                max_messages_seen[queue_name] = max(max_messages_seen[queue_name], messages)
                rates_active = require_rates_idle and (publish_rate > 0.01 or deliver_rate > 0.01)
                if messages > 0 or rates_active:
                    all_idle = False
            last_snapshot = snapshot
            elapsed_seconds = time.perf_counter() - started
            if on_poll is not None:
                on_poll(snapshot, elapsed_seconds, idle_count, all_idle)

            if all_idle:
                idle_count += 1
                if idle_count >= consecutive_idle_polls:
                    last_snapshot["_meta"] = {"max_messages_seen": max_messages_seen}
                    return last_snapshot
            else:
                idle_count = 0

            if time.perf_counter() - started > timeout_seconds:
                raise TimeoutError(
                    f"Queues did not become idle within {timeout_seconds} seconds. "
                    f"Last snapshot={last_snapshot!r}, max_messages_seen={max_messages_seen!r}"
                )

            time.sleep(poll_interval_seconds)
