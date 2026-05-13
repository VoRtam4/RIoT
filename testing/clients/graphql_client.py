"""
@file graphql_client.py
@brief GraphQL klient používaný pro správu entit, KPI a čtení výsledků přes Backend Core.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass
class GraphQLResult:
    data: dict[str, Any] | None
    errors: list[dict[str, Any]] | None
    raw_body: Any = None
    status_code: int | None = None


class GraphQLClient:
    def __init__(self, http_client, graphql_url: str):
        self.http = http_client
        self.graphql_url = graphql_url

    def execute(self, query: str, variables: dict | None = None) -> GraphQLResult:
        path = self._to_path(self.graphql_url)
        payload = {"query": query, "variables": variables or {}}
        response = self.http.post(path, json=payload)
        body = response.body if isinstance(response.body, dict) else None
        return GraphQLResult(
            data=(body or {}).get("data"),
            errors=(body or {}).get("errors"),
            raw_body=response.body,
            status_code=response.status_code,
        )

    def healthcheck(self) -> bool:
        query = """
        query Healthcheck {
          sdTypes {
            id
          }
        }
        """
        result = self.execute(query)
        return result.errors is None and result.data is not None

    def list_sd_types(self) -> list[dict]:
        query = """
        query ListSDTypes {
          sdTypes {
            id
            uid
            label
            parameters {
              id
              denotation
              label
              type
              role
            }
          }
        }
        """
        result = self.execute(query)
        self._raise_if_invalid(result, "sdTypes")
        return result.data["sdTypes"]

    def list_sd_instances(self) -> list[dict]:
        query = """
        query ListSDInstances {
          sdInstances {
            id
            uid
            label
            type {
              id
              uid
            }
          }
        }
        """
        result = self.execute(query)
        self._raise_if_invalid(result, "sdInstances")
        return result.data["sdInstances"]

    def list_sd_instances_by_type(self, sd_type_id: int | str) -> list[dict]:
        query = """
        query ListSDInstancesByType($id: ID!) {
          sdInstancesByType(id: $id) {
            id
            uid
            label
            confirmedByUser
            type {
              id
              uid
            }
          }
        }
        """
        result = self.execute(query, {"id": str(sd_type_id)})
        self._raise_if_invalid(result, "sdInstancesByType")
        return result.data["sdInstancesByType"]

    def list_kpis(self) -> list[dict]:
        query = """
        query ListKPI {
          kpiDefinitions {
            id
            label
            sdTypeID
            sdTypeUID
            sdInstanceMode
            selectedSDInstanceIDs
          }
        }
        """
        result = self.execute(query)
        self._raise_if_invalid(result, "kpiDefinitions")
        return result.data["kpiDefinitions"]

    def create_sd_type(self, sd_type_input: dict) -> dict:
        query = """
        mutation CreateSDType($input: SDTypeInput!) {
          createSDType(input: $input) {
            id
            uid
            label
          }
        }
        """
        result = self.execute(query, {"input": sd_type_input})
        self._raise_if_invalid(result, "createSDType")
        return result.data["createSDType"]

    def create_kpi(self, kpi_input: dict) -> dict:
        query = """
        mutation CreateKPI($input: KPIDefinitionInput!) {
          createKPIDefinition(input: $input) {
            id
            label
            sdTypeID
            sdTypeUID
          }
        }
        """
        result = self.execute(query, {"input": kpi_input})
        self._raise_if_invalid(result, "createKPIDefinition")
        return result.data["createKPIDefinition"]

    def update_kpi(self, kpi_id: int | str, kpi_input: dict) -> dict:
        query = """
        mutation UpdateKPI($id: ID!, $input: KPIDefinitionInput!) {
          updateKPIDefinition(id: $id, input: $input) {
            id
            label
            sdTypeID
            sdTypeUID
          }
        }
        """
        result = self.execute(query, {"id": str(kpi_id), "input": kpi_input})
        self._raise_if_invalid(result, "updateKPIDefinition")
        return result.data["updateKPIDefinition"]

    def get_kpi(self, kpi_id: int | str) -> dict:
        query = """
        query GetKPI($id: ID!) {
          kpiDefinition(id: $id) {
            id
            label
            sdTypeID
            sdTypeUID
            userIdentifier
            sdInstanceMode
            selectedSDInstanceIDs
            nodes {
              id
              parentNodeID
              nodeType
            }
          }
        }
        """
        result = self.execute(query, {"id": str(kpi_id)})
        self._raise_if_invalid(result, "kpiDefinition")
        return result.data["kpiDefinition"]

    def trigger_reprocess_via_update(self, kpi_id: int | str, payload: dict) -> dict:
        return self.update_kpi(kpi_id, payload)

    def _raise_if_invalid(self, result: GraphQLResult, root_field: str) -> None:
        if result.errors:
            raise RuntimeError(f"GraphQL errors for {root_field}: {result.errors}")
        if not isinstance(result.data, dict):
            raise RuntimeError(
                f"GraphQL response for {root_field} has no data. status={result.status_code}, body={result.raw_body!r}"
            )
        if root_field not in result.data:
            raise RuntimeError(
                f"GraphQL response missing '{root_field}'. status={result.status_code}, body={result.raw_body!r}"
            )

    def _to_path(self, url_or_path: str) -> str:
        if url_or_path.startswith("http://") or url_or_path.startswith("https://"):
            marker = self.http.base_url
            if url_or_path.startswith(marker):
                return url_or_path[len(marker) :]
        return url_or_path
