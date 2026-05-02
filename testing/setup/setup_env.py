from __future__ import annotations


class EnvironmentSetupService:
    def validate(self, ctx) -> dict[str, bool]:
        checks = {"graphql": False, "rest_sd_types": False, "rabbitmq_amqp": False, "rabbitmq_management": False}
        try:
            checks["graphql"] = ctx.graphql_client.healthcheck()
        except Exception:
            checks["graphql"] = False
        try:
            response = ctx.http_client.get("/rest/sd-types")
            checks["rest_sd_types"] = response.status_code < 400
        except Exception:
            checks["rest_sd_types"] = False
        try:
            checks["rabbitmq_amqp"] = (
                ctx.rabbitmq_amqp_client.healthcheck()
                if ctx.rabbitmq_amqp_client.is_configured()
                else False
            )
        except Exception:
            checks["rabbitmq_amqp"] = False
        try:
            checks["rabbitmq_management"] = (
                ctx.rabbitmq_management_client.healthcheck()
                if ctx.rabbitmq_management_client.is_configured()
                else False
            )
        except Exception:
            checks["rabbitmq_management"] = False
        return checks
