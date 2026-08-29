/**
 * @file mapQueryToKPINodes.ts
 * @brief Převod uložené KPI podmínky na uzly vizuálního editoru.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
type Group = any;

type SDParam = {
  id?: string | number;
  denotation: string;
  type: "number" | "string" | "boolean";
};

type ParameterReferenceValue = {
  mode: "parameter";
  comparedSDParameterSpecification?: string;
  comparedRecordOffset?: number;
};

export function mapQueryToKPINodes(query: Group, parameters: SDParam[]) {
  const nodes: any[] = [];
  let idCounter = 1;

  const nextId = () => String(idCounter++);

  const normalizeField = (field: string) =>
    field?.startsWith("tag.") ? field.slice(4) : field;

  const findParam = (field: string) =>
    parameters.find((p) => p.denotation === normalizeField(field));

  const isParameterReferenceValue = (
    value: unknown,
  ): value is ParameterReferenceValue =>
    !!value &&
    typeof value === "object" &&
    (value as ParameterReferenceValue).mode === "parameter";

  const mapOperatorToNodeType = (operator: string, type: SDParam["type"]) => {
    switch (type) {
      case "string":
        switch (operator) {
          case "=":
            return "StringEQAtom";
          case "!=":
          case "≠":
            return "StringNEQAtom";
          case "exists":
          case "∃":
            return "StringExistsAtom";
          case "notExists":
          case "∄":
            return "StringNotExistsAtom";
          default:
            return "StringEQAtom";
        }

      case "boolean":
        switch (operator) {
          case "=":
            return "BooleanEQAtom";
          case "!=":
          case "≠":
            return "BooleanNEQAtom";
          case "exists":
          case "∃":
            return "BooleanExistsAtom";
          case "notExists":
          case "∄":
            return "BooleanNotExistsAtom";
          default:
            return "BooleanEQAtom";
        }

      case "number":
        switch (operator) {
          case "=":
            return "NumericEQAtom";
          case "!=":
          case "≠":
            return "NumericNEQAtom";
          case ">":
            return "NumericGTAtom";
          case ">=":
          case "≥":
            return "NumericGEQAtom";
          case "<":
            return "NumericLTAtom";
          case "<=":
          case "≤":
            return "NumericLEQAtom";
          case "exists":
          case "∃":
            return "NumericExistsAtom";
          case "notExists":
          case "∄":
            return "NumericNotExistsAtom";
          default:
            return "NumericEQAtom";
        }
    }
  };

  const pushRuleNode = (rule: any, parentNodeID: string | null) => {
    const param = findParam(rule.field);
    if (!param) return;

    const nodeID = nextId();
    const nodeType = mapOperatorToNodeType(rule.operator, param.type);

    const base = {
      id: nodeID,
      parentNodeID,
      type: nodeType,
      sdParameterSpecification: param.denotation,
    };

    if (
      rule.operator === "exists" ||
      rule.operator === "notExists" ||
      rule.operator === "∃" ||
      rule.operator === "∄"
    ) {
      nodes.push(base);
      return;
    }

    const parameterReference = isParameterReferenceValue(rule.value)
      ? rule.value
      : null;

    const referenceFields = parameterReference
      ? {
          referenceMode: "parameter",
          comparedSDParameterSpecification:
            parameterReference.comparedSDParameterSpecification ?? "",
          comparedRecordOffset: Number(
            parameterReference.comparedRecordOffset ?? 0,
          ),
        }
      : {
          referenceMode: "literal",
        };

    if (parameterReference) {
      nodes.push({
        ...base,
        ...referenceFields,
      });
      return;
    }

    if (param.type === "number") {
      nodes.push({
        ...base,
        ...referenceFields,
        numericReferenceValue:
          rule.value === "" || rule.value === null || rule.value === undefined
            ? 0
            : Number(rule.value),
      });
      return;
    }

    if (param.type === "boolean") {
      nodes.push({
        ...base,
        ...referenceFields,
        booleanReferenceValue: rule.value === true || rule.value === "true",
      });
      return;
    }

    nodes.push({
      ...base,
      ...referenceFields,
      stringReferenceValue: String(rule.value ?? ""),
    });
  };

  const walk = (group: Group, parentID: string | null) => {
    if (!group || !Array.isArray(group.rules)) return;

    let currentParentID = parentID;

    if (group.not) {
      const notID = nextId();

      nodes.push({
        id: notID,
        parentNodeID: parentID,
        type: "LogicalOperation",
        logicalOperationType: "not",
      });

      currentParentID = notID;
    }

    const groupID = nextId();

    nodes.push({
      id: groupID,
      parentNodeID: currentParentID,
      type: "LogicalOperation",
      logicalOperationType: group.combinator || "and",
    });

    for (const rule of group.rules) {
      if (rule?.rules) {
        walk(rule, groupID);
        continue;
      }

      pushRuleNode(rule, groupID);
    }
  };

  walk(query, null);

  return nodes;
}
