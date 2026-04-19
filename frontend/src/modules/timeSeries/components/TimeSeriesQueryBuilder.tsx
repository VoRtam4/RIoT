import { useState } from "react";
import {
  QueryBuilder,
  type RuleGroupType,
  type RuleGroupTypeAny,
  type Field,
} from "react-querybuilder";
import { mapToFilterNode } from "../utils/mapQueryToFilterNode";
import { QueryBuilderBootstrap } from "@react-querybuilder/bootstrap";
import { defaultValidator } from "react-querybuilder";
import { QueryBuilderDnD } from "@react-querybuilder/dnd";

import * as ReactDnD from "react-dnd";
import * as ReactDndHtml5Backend from "react-dnd-html5-backend";

import "bootstrap/dist/css/bootstrap.min.css";
import "react-querybuilder/dist/query-builder.css";

type TagParameter = {
  denotation: string;
  label: string;
};

type Props = {
  parameters: TagParameter[];
  onChange?: (filter: any) => void;
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

export default function TimeSeriesQueryBuilder({
  parameters,
  onChange,
}: Props) {
  const [query, setQuery] = useState<RuleGroupType>({
    combinator: "or",
    rules: [],
  });

  const fields = mapFields(parameters);

  const handleChange = (q: RuleGroupTypeAny) => {
    const typed = q as RuleGroupType;
    setQuery(typed);

    const mapped = mapToFilterNode(typed);
    onChange?.(mapped);
  };

  const validator = (query: RuleGroupTypeAny) => {
    const base = defaultValidator(query);

    const result = typeof base === "boolean" ? {} : { ...base };

    const validateRule = (rule: any) => {
      if (!rule.field || !rule.operator) {
        return { valid: false };
      }

      if (!rule.value || rule.value.trim() === "") {
        return { valid: false };
      }

      return { valid: true };
    };

    const walk = (node: any) => {
      if (node.rules) {
        node.rules.forEach(walk);
      } else {
        result[node.id] = validateRule(node);
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
        <QueryBuilderDnD dnd={{ ...ReactDnD, ...ReactDndHtml5Backend }}>
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
