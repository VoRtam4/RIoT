from __future__ import annotations

from datetime import datetime, timezone
import time

from testing.setup.payload_builders import build_kpi_payload
from testing.setup.setup_entities import EntitySetupService


SD_INSTANCE_REGISTRATION_QUEUE = "sd-instance-registration-requests"


class KPISetupService:
    def ensure_kpis(self, ctx, config, state: dict) -> dict:
        state.setdefault("kpis", {})
        existing = {item["label"]: item for item in ctx.graphql_client.list_kpis()}

        for spec in config.kpis.values():
            source_state = state["sources"][spec.source]
            if not source_state.get("instance_ids"):
                # Refresh once more in case ingest created instances after initial
                # setup. Use the shared refresh helper so we also benefit from its
                # fallback through the full SDInstance list when the by-type query
                # is empty.
                EntitySetupService().refresh_instances_for_source(ctx, state, spec.source)
            if spec.mode == "SELECTED" and not source_state.get("instance_ids"):
                self._ensure_selected_instances_registered(ctx, spec, source_state)
                EntitySetupService().refresh_instances_for_source(ctx, state, spec.source)
            selected_ids = self._resolve_selected_instance_ids(source_state, spec.selected_instance_labels)
            found = existing.get(spec.label)
            if spec.mode == "SELECTED" and not selected_ids:
                # After stack/DB restarts the local testing state may have no
                # instance mapping yet, but the KPI can already exist in the
                # live backend with valid selected instance IDs. In that case,
                # sync the live KPI ID into the local state instead of silently
                # leaving a stale mapping behind.
                if found is not None and found.get("selectedSDInstanceIDs"):
                    state["kpis"][spec.id] = found["id"]
                continue
            if found is None:
                payload = build_kpi_payload(
                    spec,
                    sd_type_id=source_state["sd_type_id"],
                    sd_type_uid=source_state["sd_type_uid"],
                    selected_instance_ids=selected_ids,
                    parameter_ids_by_denotation=source_state.get("parameter_ids", {}),
                )
                found = ctx.graphql_client.create_kpi(payload)
            elif spec.mode == "SELECTED" and selected_ids:
                live_selected = [str(value) for value in found.get("selectedSDInstanceIDs", [])]
                desired_selected = [str(value) for value in selected_ids]
                if live_selected != desired_selected:
                    payload = build_kpi_payload(
                        spec,
                        sd_type_id=source_state["sd_type_id"],
                        sd_type_uid=source_state["sd_type_uid"],
                        selected_instance_ids=selected_ids,
                        parameter_ids_by_denotation=source_state.get("parameter_ids", {}),
                    )
                    found = ctx.graphql_client.update_kpi(found["id"], payload)
            state["kpis"][spec.id] = found["id"]
        return state

    def _ensure_selected_instances_registered(self, ctx, spec, source_state: dict) -> None:
        if not spec.selected_instance_labels:
            return
        if not ctx.rabbitmq_amqp_client.is_configured():
            return
        now = datetime.now(timezone.utc).isoformat()
        messages = [
            {
                "eventTime": now,
                "label": label,
                "sdInstanceUID": label,
                "sdTypeUID": source_state["sd_type_uid"],
            }
            for label in spec.selected_instance_labels
        ]
        ctx.rabbitmq_amqp_client.publish_to_queue(SD_INSTANCE_REGISTRATION_QUEUE, messages)
        if ctx.rabbitmq_management_client.is_configured():
            ctx.rabbitmq_management_client.wait_for_queues_idle(
                [SD_INSTANCE_REGISTRATION_QUEUE],
                poll_interval_seconds=1.0,
                consecutive_idle_polls=2,
                timeout_seconds=60.0,
            )
        if not self._wait_until_selected_instances_visible(ctx, spec, source_state["sd_type_uid"]):
            raise TimeoutError(
                "Selected SD instances were published for registration, but they did not become "
                f"visible in backend-core within the timeout. source={spec.source}, "
                f"labels={spec.selected_instance_labels!r}"
            )

    def _wait_until_selected_instances_visible(self, ctx, spec, sd_type_uid: str, timeout_seconds: float = 60.0) -> bool:
        started = time.perf_counter()
        wanted = set(spec.selected_instance_labels)
        entity_service = EntitySetupService()
        while time.perf_counter() - started <= timeout_seconds:
            state = ctx.state_manager.load_runtime_state()
            if spec.source not in state.get("sources", {}):
                state = entity_service.ensure_sd_types(ctx, ctx.config)
            else:
                entity_service.refresh_instances_for_source(ctx, state, spec.source)
            ctx.state_manager.save_runtime_state(state)
            source_state = state["sources"][spec.source]
            if str(source_state.get("sd_type_uid", "")) == str(sd_type_uid):
                visible = set(source_state.get("instance_uids", []))
                if wanted.issubset(visible):
                    return True
            time.sleep(1.0)
        return False

    def _resolve_selected_instance_ids(self, source_state: dict, selected_labels: list[str]) -> list[str]:
        if not selected_labels:
            return []
        mapping = {
            label: str(instance_id)
            for label, instance_id in zip(source_state.get("instance_uids", []), source_state.get("instance_ids", []))
        }
        resolved = [mapping[label] for label in selected_labels if label in mapping]
        if resolved:
            return resolved
        fallback_count = min(len(selected_labels), len(source_state.get("instance_ids", [])))
        return [str(instance_id) for instance_id in source_state.get("instance_ids", [])[:fallback_count]]
