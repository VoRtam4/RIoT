import { useState, useEffect } from "react";
import {
  QueryBuilder,
  type ValueSelectorProps,
  type RuleGroupType,
  type RuleGroupTypeAny,
  type Field,
} from "react-querybuilder";
import { QueryBuilderBootstrap } from "@react-querybuilder/bootstrap";
import { defaultValidator } from "react-querybuilder";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-icons/font/bootstrap-icons.css";
import "react-querybuilder/dist/query-builder.css";
import { QueryBuilderDnD } from "@react-querybuilder/dnd";
import * as ReactDnD from "react-dnd";
import * as ReactDndHtml5Backend from "react-dnd-html5-backend";

export type SDParameter = {
  id: string;
  label?: string;
  type: "string" | "number" | "boolean";
  role?: "field" | "tag";
  denotation?: string;
};

type Props = {
  parameters: SDParameter[];
  onChange?: (query: RuleGroupType) => void;
  initialQuery?: RuleGroupType | null;
};

const MyCombinatorSelector = (props: ValueSelectorProps) => {
  return (
    <select
      className="form-select form-select-sm"
      value={props.value}
      onChange={(e) => props.handleOnChange(e.target.value)}
    >
      <option value="and">AND</option>
      <option value="or">OR</option>
      <option value="nor">NOR</option>
    </select>
  );
};

const getOperators = (type: SDParameter["type"]) => {
  const common = [
    { name: "=", label: "=" },
    { name: "!=", label: "≠" },
    { name: "exists", label: "∃", arity: "unary" },
    { name: "notExists", label: "∄", arity: "unary" },
  ];

  switch (type) {
    case "number":
      return [
        ...common,
        { name: "<", label: "<" },
        { name: ">", label: ">" },
        { name: "<=", label: "≤" },
        { name: ">=", label: "≥" },
      ];
    default:
      return common;
  }
};

const mapFields = (params: SDParameter[]): Field[] =>
  params.map((p) => ({
    name:
      p.role === "tag" ? `tag.${p.denotation ?? p.id}` : (p.denotation ?? p.id),

    label: p.label ?? p.denotation ?? p.id,

    operators: getOperators(p.type),
    inputType: p.type === "number" ? "number" : "text",
    valueEditorType: p.type === "boolean" ? "select" : "text",
    values:
      p.type === "boolean"
        ? [
            { name: "true", label: "True" },
            { name: "false", label: "False" },
          ]
        : undefined,
  }));

export default function KpiQueryBuilder({
  parameters,
  onChange,
  initialQuery,
}: Props) {
  const [query, setQuery] = useState<RuleGroupType>({
    combinator: "and",
    rules: [],
  });

  useEffect(() => {
    if (initialQuery && Array.isArray(initialQuery.rules)) {
      setQuery(initialQuery);
    }
  }, [initialQuery]);

  const fields = mapFields(parameters);

  const handleChange = (q: RuleGroupTypeAny) => {
    const typed = q as RuleGroupType;
    setQuery(typed);
    onChange?.(typed);
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
            validator={defaultValidator}
            enableDragAndDrop
            showNotToggle
            controlElements={{
              combinatorSelector: MyCombinatorSelector,
            }}
          />
        </QueryBuilderDnD>
      </QueryBuilderBootstrap>
    </div>
  );
}
