/**
 * @file TimeSeriesQueryBuilder.tsx
 * @brief Builder filtrů pro dotazování historických raw a KPI dat.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useMemo, useRef, useState } from "react";
import {
  QueryBuilder,
  type RuleType,
  type RuleGroupType,
  type RuleGroupTypeAny,
  type Field,
} from "react-querybuilder";
import { QueryBuilderBootstrap } from "@react-querybuilder/bootstrap";
import { defaultValidator } from "react-querybuilder";
import { QueryBuilderDnD } from "@react-querybuilder/dnd";

import * as ReactDnD from "react-dnd";
import * as ReactDndHtml5Backend from "react-dnd-html5-backend";
import type {
  FilterNodeInput,
  FilterOperator,
} from "../../../generated/graphql";

import "bootstrap/dist/css/bootstrap.min.css";
import "react-querybuilder/dist/query-builder.css";

type TagParameter = {
  denotation: string;
  label: string;
};

type Props = {
  parameters: TagParameter[];
  value?: FilterNodeInput;
  onChange?: (filter: FilterNodeInput | undefined) => void;
};

const operators = [
  { name: "eq", label: "=" },
  { name: "neq", label: "≠" },
  { name: "contains", label: "contains" },
  { name: "in", label: "in" },
  { name: "prefix", label: "starts with" },
  { name: "suffix", label: "ends with" },
  { name: "regex", label: "regex" },
];

const mapFields = (params: TagParameter[]): Field[] =>
  params.map((p) => ({
    name: p.denotation,
    label: p.label ?? p.denotation,
    operators,
    inputType: "text",
  }));

function filterKey(filter?: FilterNodeInput) {
  return filter ? JSON.stringify(filter) : "";
}

function emptyQuery(): RuleGroupType {
  return {
    combinator: "or",
    rules: [],
  };
}

function mapFilterNodeToQuery(filter?: FilterNodeInput): RuleGroupType {
  if (!filter || filter.type !== "logical" || !filter.nodes?.length) {
    return emptyQuery();
  }

  const mapNode = (
    node: FilterNodeInput,
  ): RuleGroupType | RuleType | null => {
    if (node.type === "logical") {
      return {
        combinator: node.operator === "and" ? "and" : "or",
        rules: (node.nodes ?? [])
          .map((child) => mapNode(child))
          .filter(Boolean) as Array<RuleGroupType | RuleType>,
      };
    }

    if (node.type === "rule" && node.rule) {
      return {
        field: node.rule.tag,
        operator: node.rule.operator,
        value: node.rule.value,
      };
    }

    return null;
  };

  const mapped = mapNode(filter);

  return mapped && "rules" in mapped ? mapped : emptyQuery();
}

function mapQueryToFilterNode(query: RuleGroupType): FilterNodeInput | undefined {
  if (!query.rules.length) return undefined;

  const nodes = query.rules
    .map((rule) => {
      if ("rules" in rule) {
        return mapQueryToFilterNode(rule);
      }

      if (!rule.field || !rule.operator || !String(rule.value ?? "").trim()) {
        return null;
      }

      return {
        type: "rule",
        rule: {
          tag: String(rule.field),
          operator: rule.operator as FilterOperator,
          value: String(rule.value),
        },
      } satisfies FilterNodeInput;
    })
    .filter(Boolean) as FilterNodeInput[];

  if (!nodes.length) {
    return undefined;
  }

  return {
    type: "logical",
    operator: query.combinator === "and" ? "and" : "or",
    nodes,
  };
}

export default function TimeSeriesQueryBuilder({
  parameters,
  value,
  onChange,
}: Props) {
  const [query, setQuery] = useState<RuleGroupType>(mapFilterNodeToQuery(value));
  const lastAppliedValueKeyRef = useRef(filterKey(value));
  const lastEmittedValueKeyRef = useRef<string | null>(null);

  const fields = useMemo(() => mapFields(parameters), [parameters]);
  const dnd = useMemo(
    () => ({ ...ReactDnD, ...ReactDndHtml5Backend }),
    [],
  );

  useEffect(() => {
    const nextValueKey = filterKey(value);

    if (nextValueKey === lastAppliedValueKeyRef.current) {
      return;
    }

    if (nextValueKey === lastEmittedValueKeyRef.current) {
      lastAppliedValueKeyRef.current = nextValueKey;
      return;
    }

    queueMicrotask(() => {
      setQuery(mapFilterNodeToQuery(value));
    });
    lastAppliedValueKeyRef.current = nextValueKey;
  }, [value]);

  const handleChange = (q: RuleGroupTypeAny) => {
    const typed = q as RuleGroupType;
    const filter = mapQueryToFilterNode(typed);

    setQuery(typed);
    lastEmittedValueKeyRef.current = filterKey(filter);
    lastAppliedValueKeyRef.current = lastEmittedValueKeyRef.current;
    onChange?.(filter);
  };

  const validator = (query: RuleGroupTypeAny) => {
    const base = defaultValidator(query);

    const result: Record<string, { valid: boolean }> =
      typeof base === "boolean"
        ? {}
        : { ...(base as Record<string, { valid: boolean }>) };

    const validateRule = (rule: RuleType) => {
      if (!rule.field || !rule.operator) {
        return { valid: false };
      }

      if (!String(rule.value ?? "").trim()) {
        return { valid: false };
      }

      return { valid: true };
    };

    const walk = (node: RuleGroupTypeAny | RuleType) => {
      if ("rules" in node && Array.isArray(node.rules)) {
        node.rules.forEach((child) => walk(child as RuleGroupTypeAny | RuleType));
      } else {
        const rule = node as RuleType;
        if (rule.id) {
          result[rule.id] = validateRule(rule);
        }
      }
    };

    walk(query);

    return result;
  };

  return (
    <div>
      <style>
        {`
          .ruleGroup-notToggle .form-check-label {
            margin-left: 2px;
            color: white;
          }
          .drag-handle {
            color: white;
          }
        `}
      </style>

      <QueryBuilderBootstrap>
        <QueryBuilderDnD dnd={dnd}>
          <QueryBuilder
            fields={fields}
            query={query}
            onQueryChange={handleChange}
            validator={validator}
            enableDragAndDrop
            showNotToggle
          />
        </QueryBuilderDnD>
      </QueryBuilderBootstrap>
    </div>
  );
}
