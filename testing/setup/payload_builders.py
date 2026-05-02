from __future__ import annotations

from typing import Any


NODE_TYPE_MAP = {
    "AND": ("LogicalOperation", "and"),
    "OR": ("LogicalOperation", "or"),
    "NOR": ("LogicalOperation", "nor"),
    "NOT": ("LogicalOperation", "not"),
    "STRING_EQ": ("StringEQAtom", None),
    "STRING_NEQ": ("StringNEQAtom", None),
    "STRING_EXISTS": ("StringExistsAtom", None),
    "STRING_NOT_EXISTS": ("StringNotExistsAtom", None),
    "BOOLEAN_EQ": ("BooleanEQAtom", None),
    "BOOLEAN_NEQ": ("BooleanNEQAtom", None),
    "BOOLEAN_EXISTS": ("BooleanExistsAtom", None),
    "BOOLEAN_NOT_EXISTS": ("BooleanNotExistsAtom", None),
    "NUMERIC_EQ": ("NumericEQAtom", None),
    "NUMERIC_NEQ": ("NumericNEQAtom", None),
    "NUMERIC_GT": ("NumericGTAtom", None),
    "NUMERIC_GEQ": ("NumericGEQAtom", None),
    "NUMERIC_LT": ("NumericLTAtom", None),
    "NUMERIC_LEQ": ("NumericLEQAtom", None),
    "NUMERIC_EXISTS": ("NumericExistsAtom", None),
    "NUMERIC_NOT_EXISTS": ("NumericNotExistsAtom", None),
}


def build_sd_type_payload(source_name: str, source_payload: dict[str, Any]) -> dict[str, Any]:
    parameters = []
    for role in ("tags", "fields"):
        for item in source_payload["parameters"][role]:
            parameters.append(
                {
                    "label": item["label"],
                    "denotation": item["denotation"],
                    "type": item["type"],
                    "role": "tag" if role == "tags" else "field",
                }
            )
    return {
        "uid": source_payload["sd_type_uid"],
        "label": source_payload["sd_type_label"],
        "parameters": parameters,
    }


def build_kpi_payload(
    spec,
    sd_type_id: str,
    sd_type_uid: str,
    selected_instance_ids: list[str],
    parameter_ids_by_denotation: dict[str, int],
) -> dict[str, Any]:
    # This payload shape is intentionally centralized here.
    # If the final GraphQL schema differs, update this function only.
    return {
        "label": spec.label,
        "sdTypeID": sd_type_id,
        "sdTypeUID": sd_type_uid,
        "userIdentifier": "testing",
        "sdInstanceMode": "selected" if spec.mode == "SELECTED" else "all",
        "selectedSDInstanceIDs": selected_instance_ids,
        "nodes": _map_kpi_nodes(spec.definition.get("nodes", []), parameter_ids_by_denotation),
    }


def build_kpi_payload_variant(
    spec,
    sd_type_id: str,
    sd_type_uid: str,
    selected_instance_ids: list[str],
    parameter_ids_by_denotation: dict[str, int],
    variant: int,
) -> dict[str, Any]:
    definition = {"nodes": _variant_nodes(spec.definition.get("nodes", []), variant)}
    spec_like = type(
        "SpecLike",
        (),
        {
            "label": spec.label,
            "mode": spec.mode,
            "definition": definition,
        },
    )
    return build_kpi_payload(
        spec_like,
        sd_type_id=sd_type_id,
        sd_type_uid=sd_type_uid,
        selected_instance_ids=selected_instance_ids,
        parameter_ids_by_denotation=parameter_ids_by_denotation,
    )


def build_history_export_payload(
    *,
    data_type: str,
    sd_type_id: int,
    sd_instance_ids: list[int] | None,
    kpi_definition_ids: list[int] | None,
    from_iso: str,
    to_iso: str,
    aggregate_seconds: int | None = None,
    sort_desc: bool | None = None,
) -> dict[str, Any]:
    payload = {
        "type": data_type,
        "sdTypeID": sd_type_id,
        "sdInstanceIDs": sd_instance_ids or [],
        "kpiDefinitionIDs": kpi_definition_ids or [],
        "from": from_iso,
        "to": to_iso,
    }
    if aggregate_seconds is not None:
        payload["aggregateSeconds"] = aggregate_seconds
    if sort_desc is not None:
        payload["sortDesc"] = sort_desc
    return payload


def _map_kpi_nodes(nodes: list[dict[str, Any]], parameter_ids_by_denotation: dict[str, int]) -> list[dict[str, Any]]:
    id_map = {
        node["node_id"]: str(index)
        for index, node in enumerate(nodes, start=1)
    }
    mapped = []
    for node in nodes:
        mapped_type, logical_operation = NODE_TYPE_MAP[node["node_type"]]
        out = {
            "id": id_map[node["node_id"]],
            "type": mapped_type,
        }
        if "parent_node_id" in node:
            out["parentNodeID"] = id_map[node["parent_node_id"]]
        if logical_operation is not None:
            out["logicalOperationType"] = logical_operation
        if "parameter_denotation" in node:
            denotation = node["parameter_denotation"]
            out["sdParameterSpecification"] = denotation
            if denotation not in parameter_ids_by_denotation:
                raise KeyError(f"Missing parameter id for denotation '{denotation}'")
            out["sdParameterID"] = parameter_ids_by_denotation[denotation]
        if "numeric_reference_value" in node:
            out["numericReferenceValue"] = node["numeric_reference_value"]
        if "boolean_reference_value" in node:
            out["booleanReferenceValue"] = node["boolean_reference_value"]
        if "string_reference_value" in node:
            out["stringReferenceValue"] = node["string_reference_value"]
        mapped.append(out)
    return mapped


def _variant_nodes(nodes: list[dict[str, Any]], variant: int) -> list[dict[str, Any]]:
    if variant % 2 == 0 or not nodes:
        return [dict(node) for node in nodes]

    root = next((node for node in nodes if node["node_id"] == "root"), None)
    if root is None:
        return [dict(node) for node in nodes]

    child_nodes = [node for node in nodes if node.get("parent_node_id") == "root"]
    if root["node_type"] in {"AND", "OR"} and len(child_nodes) >= 2:
        reordered_children = list(reversed(child_nodes))
        others = [node for node in nodes if node not in child_nodes]
        return [dict(node) for node in others] + [dict(node) for node in reordered_children]
    # For single-atom or otherwise non-reorderable trees, use a semantically
    # equivalent double negation instead of wrapping the root in AND with a
    # single child. The latter proved fragile in the backend update path.
    wrapped_root = {
        "node_id": "root",
        "node_type": "NOT",
    }
    inner_not = {
        "node_id": "n1",
        "parent_node_id": "root",
        "node_type": "NOT",
    }
    original_root = dict(root)
    original_root["node_id"] = "n2"
    original_root["parent_node_id"] = "n1"

    remapped_nodes = [wrapped_root, inner_not, original_root]
    for node in nodes:
        if node["node_id"] == "root":
            continue
        copied = dict(node)
        if copied.get("parent_node_id") == "root":
            copied["parent_node_id"] = "n2"
        remapped_nodes.append(copied)
    return remapped_nodes
