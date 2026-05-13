"""
@file metrics_client.py
@brief Klient pro čtení metrik z monitoringu a jejich připojení k výsledkům experimentů.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations


class MetricsClient:
    def __init__(
        self,
        rabbitmq_url: str | None = None,
        rabbitmq_user: str | None = None,
        rabbitmq_password: str | None = None,
        prometheus_url: str | None = None,
    ):
        self.rabbitmq_url = rabbitmq_url
        self.rabbitmq_user = rabbitmq_user
        self.rabbitmq_password = rabbitmq_password
        self.prometheus_url = prometheus_url

    def get_queue_depths(self, queue_names: list[str]) -> dict[str, int]:
        # Placeholder: integrate RabbitMQ HTTP API here when needed.
        return {queue_name: 0 for queue_name in queue_names}

    def sample_ingest_metrics(self) -> dict[str, float]:
        return {
            "queue_depth_max": 0.0,
            "rabbitmq_publish_rate": 0.0,
            "rabbitmq_consume_rate": 0.0,
            "backend_cpu": 0.0,
            "mpu_cpu": 0.0,
            "tsdb_cpu": 0.0,
        }
