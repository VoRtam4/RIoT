from __future__ import annotations

import json
import uuid
from dataclasses import dataclass
from pathlib import Path

from testing.clients.graphql_client import GraphQLClient
from testing.clients.http_client import HTTPClient
from testing.clients.metrics_client import MetricsClient
from testing.clients.rabbitmq_amqp_client import RabbitMQAMQPClient
from testing.clients.rabbitmq_management_client import RabbitMQManagementClient
from testing.clients.rest_client import RESTClient
from testing.core.models import EnvConfig, RunDescriptor, utcnow
from testing.core.state_manager import StateManager


@dataclass
class RuntimeContext:
    repo_root: Path
    config: object
    env: EnvConfig
    run: RunDescriptor
    state_manager: StateManager
    http_client: HTTPClient
    graphql_client: GraphQLClient
    rest_client: RESTClient
    metrics_client: MetricsClient
    rabbitmq_amqp_client: RabbitMQAMQPClient
    rabbitmq_management_client: RabbitMQManagementClient


def build_runtime_context(repo_root: Path, config) -> RuntimeContext:
    env = config.env
    output_dir = repo_root / env.output_dir
    (output_dir / "raw").mkdir(parents=True, exist_ok=True)
    (output_dir / "reports").mkdir(parents=True, exist_ok=True)
    (repo_root / "testing" / ".state").mkdir(parents=True, exist_ok=True)

    run = RunDescriptor(
        run_id=str(uuid.uuid4()),
        started_at=utcnow(),
        output_dir=output_dir,
    )
    state_manager = StateManager(repo_root / "testing" / ".state")
    http_client = HTTPClient(env.base_url, env.api_key, env.timeout_seconds)
    graphql_client = GraphQLClient(http_client, env.graphql_url)
    rest_client = RESTClient(http_client, env.rest_url)
    metrics_client = MetricsClient(
        rabbitmq_url=env.rabbitmq_management_url,
        rabbitmq_user=env.rabbitmq_user,
        rabbitmq_password=env.rabbitmq_password,
        prometheus_url=env.prometheus_url,
    )
    rabbitmq_amqp_client = RabbitMQAMQPClient(
        amqp_url=env.rabbitmq_amqp_url,
        timeout_seconds=env.timeout_seconds,
    )
    rabbitmq_management_client = RabbitMQManagementClient(
        base_url=env.rabbitmq_management_url,
        username=env.rabbitmq_user,
        password=env.rabbitmq_password,
        timeout_seconds=env.timeout_seconds,
    )
    return RuntimeContext(
        repo_root=repo_root,
        config=config,
        env=env,
        run=run,
        state_manager=state_manager,
        http_client=http_client,
        graphql_client=graphql_client,
        rest_client=rest_client,
        metrics_client=metrics_client,
        rabbitmq_amqp_client=rabbitmq_amqp_client,
        rabbitmq_management_client=rabbitmq_management_client,
    )
