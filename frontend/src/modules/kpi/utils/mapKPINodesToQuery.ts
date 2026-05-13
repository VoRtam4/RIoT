/**
 * @file mapKPINodesToQuery.ts
 * @brief Převod uzlů vizuálního editoru zpět na dotazovou strukturu KPI.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { RuleGroupType } from "react-querybuilder";

type Node = any;

type Param = {
  id: string | number;
  denotation: string;
  role?: string;
};

function mapNodeTypeToOperator(type: string): string {
  switch (type) {
    case "StringEQAtom":
    case "BooleanEQAtom":
    case "NumericEQAtom":
      return "=";

    case "StringNEQAtom":
    case "BooleanNEQAtom":
    case "NumericNEQAtom":
      return "!=";

    case "NumericGTAtom":
      return ">";
    case "NumericGEQAtom":
      return ">=";
    case "NumericLTAtom":
      return "<";
    case "NumericLEQAtom":
      return "<=";

    case "StringExistsAtom":
    case "BooleanExistsAtom":
    case "NumericExistsAtom":
      return "exists";

    case "StringNotExistsAtom":
    case "BooleanNotExistsAtom":
    case "NumericNotExistsAtom":
      return "notExists";

    default:
      return "=";
  }
}

function extractValue(node: Node): any {
  if (node.numericReferenceValue !== undefined) {
    return String(node.numericReferenceValue);
  }

  if (node.booleanReferenceValue !== undefined) {
    return node.booleanReferenceValue ? "true" : "false";
  }

  if (node.stringReferenceValue !== undefined) {
    return node.stringReferenceValue;
  }

  return "";
}

function buildTree(nodes: Node[]) {
  const map = new Map<string, any>();

  nodes.forEach((n) => {
    map.set(String(n.id), {
      ...n,
      id: String(n.id),
      parentNodeID:
        n.parentNodeID === null || n.parentNodeID === undefined
          ? null
          : String(n.parentNodeID),
      children: [],
    });
  });

  nodes.forEach((n) => {
    const currentId = String(n.id);
    const parentId =
      n.parentNodeID === null || n.parentNodeID === undefined
        ? null
        : String(n.parentNodeID);

    if (!parentId) return;

    const parent = map.get(parentId);
    const child = map.get(currentId);

    if (parent && child) {
      parent.children.push(child);
    }
  });

  return [...map.values()].filter((n) => !n.parentNodeID);
}

function resolveFieldName(
  sdParameterSpecification: string,
  parameters: Param[],
) {
  const param = parameters.find(
    (p) => String(p.denotation) === String(sdParameterSpecification),
  );

  if (!param) return sdParameterSpecification;

  if (String(param.role).toLowerCase() === "tag") {
    return `tag.${param.denotation}`;
  }

  return param.denotation;
}

function nodeToQuery(node: any, parameters: Param[]): any {
  if (node.nodeType === "LogicalOperation") {
    const op = String(node.type).toLowerCase();

    if (op === "not") {
      const child = node.children?.[0];
      if (!child) return null;

      const mappedChild = nodeToQuery(child, parameters);
      if (!mappedChild) return null;

      if (mappedChild.rules) {
        return {
          ...mappedChild,
          not: true,
        };
      }

      return {
        combinator: "and",
        not: true,
        rules: [mappedChild],
      };
    }

    return {
      combinator: op === "or" || op === "nor" || op === "and" ? op : "and",
      rules: (node.children ?? [])
        .map((child: any) => nodeToQuery(child, parameters))
        .filter(Boolean),
    };
  }

  return {
    field: resolveFieldName(node.sdParameterSpecification, parameters),
    operator: mapNodeTypeToOperator(node.nodeType),
    value: extractValue(node),
  };
}

export function mapKPINodesToQuery(
  nodes: Node[],
  parameters: Param[],
): RuleGroupType {
  if (!nodes || nodes.length === 0) {
    return {
      combinator: "and",
      rules: [],
    };
  }

  const roots = buildTree(nodes);

  if (roots.length === 0) {
    return {
      combinator: "and",
      rules: [],
    };
  }

  if (roots.length === 1) {
    const root = roots[0];

    if (root.nodeType === "LogicalOperation") {
      const mapped = nodeToQuery(root, parameters);

      if (mapped?.rules) {
        return mapped;
      }

      return {
        combinator: "and",
        rules: mapped ? [mapped] : [],
      };
    }

    const mappedRule = nodeToQuery(root, parameters);

    return {
      combinator: "and",
      rules: mappedRule ? [mappedRule] : [],
    };
  }

  return {
    combinator: "and",
    rules: roots.map((root) => nodeToQuery(root, parameters)).filter(Boolean),
  };
}
