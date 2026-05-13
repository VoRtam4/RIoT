"""
@file setup_entities.py
@brief Příprava sledovaných typů, instancí a skupin před spuštěním experimentů.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from testing.setup.payload_builders import build_sd_type_payload


class EntitySetupService:
    def ensure_sd_types(self, ctx, config) -> dict:
        state = ctx.state_manager.load_runtime_state()
        state.setdefault("sources", {})

        if ctx.rabbitmq_amqp_client.is_configured() and ctx.rabbitmq_management_client.is_configured() and ctx.rabbitmq_management_client.healthcheck():
            self._register_sd_types_via_rabbitmq(ctx, config)

        existing = {item["uid"]: item for item in ctx.graphql_client.list_sd_types()}
        for source_name, source_spec in config.sources.items():
            payload = source_spec.payload
            uid = payload["sd_type_uid"]
            item = existing.get(uid)
            if item is None:
                item = ctx.graphql_client.create_sd_type(
                    build_sd_type_payload(source_name, payload)
                )
                item["parameters"] = build_sd_type_payload(source_name, payload)["parameters"]
            state["sources"][source_name] = {
                "sd_type_id": item["id"],
                "sd_type_uid": item["uid"],
                "sd_type_label": item["label"],
                "parameter_ids": {
                    param["denotation"]: int(param["id"])
                    for param in item.get("parameters", [])
                    if "id" in param
                },
                "instance_ids": [],
                "instance_uids": [],
            }
            self.refresh_instances_for_source(ctx, state, source_name)
        return state

    def refresh_instances_for_source(self, ctx, state: dict, source_name: str) -> None:
        source_state = state["sources"][source_name]
        instances = ctx.graphql_client.list_sd_instances_by_type(source_state["sd_type_id"])
        if not instances:
            instances = [
                instance
                for instance in ctx.graphql_client.list_sd_instances()
                if str(instance.get("type", {}).get("uid", "")) == str(source_state["sd_type_uid"])
            ]
        source_state["instance_ids"] = [instance["id"] for instance in instances]
        source_state["instance_uids"] = [instance["uid"] for instance in instances]

    def _register_sd_types_via_rabbitmq(self, ctx, config) -> None:
        queue_name = "sd-type-registration-requests"
        ctx.rabbitmq_management_client.wait_for_queue_consumers(
            {
                "sd-type-registration-requests": 1,
                "set-of-sd-types-updates": 1,
            },
            poll_interval_seconds=1.0,
            timeout_seconds=120.0,
        )
        messages: list[dict] = []
        for source_name, source_spec in config.sources.items():
            payload = build_sd_type_payload(source_name, source_spec.payload)
            messages.append(
                {
                    "sdTypeUID": payload["uid"],
                    "label": payload["label"],
                    "parameters": payload["parameters"],
                }
            )
        if messages:
            ctx.rabbitmq_amqp_client.publish_to_queue(queue_name, messages)
        ctx.rabbitmq_management_client.wait_for_queues_idle(
            [queue_name, "set-of-sd-types-updates"],
            poll_interval_seconds=1.0,
            consecutive_idle_polls=2,
            timeout_seconds=60.0,
        )
