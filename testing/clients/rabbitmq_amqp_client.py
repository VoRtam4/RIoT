from __future__ import annotations

import json
from typing import Any

import pika
from pika.exceptions import AMQPError, StreamLostError


class RabbitMQAMQPClient:
    def __init__(self, amqp_url: str | None, timeout_seconds: int = 30):
        self.amqp_url = (amqp_url or "").strip()
        self.timeout_seconds = timeout_seconds
        self._connection: pika.BlockingConnection | None = None
        self._channel = None

    def is_configured(self) -> bool:
        return bool(self.amqp_url)

    def healthcheck(self) -> bool:
        if not self.is_configured():
            return False
        try:
            channel = self._get_channel()
            return bool(channel.is_open)
        except Exception:
            return False

    def publish_to_queue(self, queue_name: str, payload: Any) -> None:
        if not self.is_configured():
            raise RuntimeError("RabbitMQ AMQP client is not configured")
        body = json.dumps(payload, ensure_ascii=True).encode("utf-8")
        try:
            self._basic_publish(queue_name, body)
        except (StreamLostError, AMQPError, OSError):
            self.close()
            self._basic_publish(queue_name, body)

    def close(self) -> None:
        if self._channel is not None:
            try:
                self._channel.close()
            except Exception:
                pass
            self._channel = None
        if self._connection is not None:
            try:
                self._connection.close()
            except Exception:
                pass
            self._connection = None

    def _get_channel(self):
        if self._channel is not None and self._channel.is_open:
            return self._channel
        parameters = pika.URLParameters(self.amqp_url)
        parameters.socket_timeout = self.timeout_seconds
        parameters.blocked_connection_timeout = self.timeout_seconds
        parameters.heartbeat = self.timeout_seconds
        self._connection = pika.BlockingConnection(parameters)
        self._channel = self._connection.channel()
        return self._channel

    def _basic_publish(self, queue_name: str, body: bytes) -> None:
        channel = self._get_channel()
        channel.basic_publish(
            exchange="",
            routing_key=queue_name,
            body=body,
            properties=pika.BasicProperties(
                content_type="application/json",
                delivery_mode=2,
            ),
            mandatory=False,
        )
