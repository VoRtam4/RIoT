/**
 * @file KpiQueryBuilder.tsx
 * @brief Vizuální builder logických a porovnávacích podmínek KPI.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useState, useEffect } from "react";
import {
  QueryBuilder,
  type ValueSelectorProps,
  type ValueEditorProps,
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

type ParameterReferenceValue = {
  mode: "parameter";
  comparedSDParameterSpecification?: string;
  comparedRecordOffset?: number;
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

const normalizeField = (field: string) =>
  field?.startsWith("tag.") ? field.slice(4) : field;

const isParameterReferenceValue = (
  value: unknown,
): value is ParameterReferenceValue =>
  !!value &&
  typeof value === "object" &&
  (value as ParameterReferenceValue).mode === "parameter";

const KpiValueEditor = ({
  parameters,
  ...props
}: ValueEditorProps & { parameters: SDParameter[] }) => {
  const selectedParameter = parameters.find(
    (p) => String(p.denotation ?? p.id) === normalizeField(String(props.field)),
  );
  const compatibleParameters = parameters.filter(
    (p) => p.type === selectedParameter?.type,
  );
  const parameterValue = isParameterReferenceValue(props.value)
    ? props.value
    : null;
  const referenceMode = parameterValue ? "parameter" : "literal";

  const setParameterReference = (
    comparedSpecification: string | undefined,
    comparedRecordOffset = parameterValue?.comparedRecordOffset ?? 0,
  ) => {
    const comparedParameter =
      compatibleParameters.find(
        (p) => String(p.denotation ?? p.id) === String(comparedSpecification),
      ) ??
      compatibleParameters[0] ??
      selectedParameter;

    props.handleOnChange({
      mode: "parameter",
      comparedSDParameterSpecification:
        comparedParameter?.denotation ?? String(comparedParameter?.id ?? ""),
      comparedRecordOffset,
    });
  };

  const setReferenceMode = (mode: "literal" | "parameter") => {
    if (mode === "parameter") {
      setParameterReference(parameterValue?.comparedSDParameterSpecification);
      return;
    }
    props.handleOnChange(selectedParameter?.type === "boolean" ? "true" : "");
  };

  const literalValue = parameterValue ? "" : (props.value ?? "");

  return (
    <div className="d-flex flex-grow-1 gap-2 align-items-center kpi-value-editor">
      <select
        className="form-select form-select-sm kpi-reference-mode-select"
        value={referenceMode}
        onChange={(e) =>
          setReferenceMode(e.target.value as "literal" | "parameter")
        }
      >
        <option value="literal">Value</option>
        <option value="parameter">Parameter</option>
      </select>

      {parameterValue ? (
        <>
          <select
            className="form-select form-select-sm"
            value={String(
              parameterValue.comparedSDParameterSpecification ?? "",
            )}
            onChange={(e) => setParameterReference(e.target.value)}
          >
            {compatibleParameters.map((p) => (
              <option
                key={String(p.denotation ?? p.id)}
                value={String(p.denotation ?? p.id)}
              >
                {p.label ?? p.denotation ?? p.id}
              </option>
            ))}
          </select>

          <select
            className="form-select form-select-sm kpi-record-offset-select"
            value={String(parameterValue.comparedRecordOffset ?? 0)}
            onChange={(e) =>
              setParameterReference(
                parameterValue.comparedSDParameterSpecification,
                Number(e.target.value),
              )
            }
          >
            <option value="0">Current</option>
            <option value="-1">Previous</option>
          </select>
        </>
      ) : selectedParameter?.type === "boolean" ? (
        <select
          className="form-select form-select-sm"
          value={String(literalValue)}
          onChange={(e) => props.handleOnChange(e.target.value)}
        >
          <option value="true">True</option>
          <option value="false">False</option>
        </select>
      ) : (
        <input
          className="form-control form-control-sm"
          type={selectedParameter?.type === "number" ? "number" : "text"}
          value={String(literalValue)}
          onChange={(e) => props.handleOnChange(e.target.value)}
        />
      )}
    </div>
  );
};

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
  const ValueEditor = (props: ValueEditorProps) => (
    <KpiValueEditor {...props} parameters={parameters} />
  );

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

          .kpi-value-editor {
            min-width: 260px;
          }

          .kpi-reference-mode-select {
            width: auto;
            flex: 0 0 auto;
          }

          .kpi-record-offset-select {
            width: auto;
            flex: 0 0 auto;
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
              valueEditor: ValueEditor,
            }}
          />
        </QueryBuilderDnD>
      </QueryBuilderBootstrap>
    </div>
  );
}
