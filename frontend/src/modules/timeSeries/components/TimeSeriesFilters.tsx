import { useEffect, useMemo, useState } from "react";
import Select, { type MultiValue, type SingleValue } from "react-select";

import type { TimeSeriesReadInput } from "../../../generated/graphql";
import TimeSeriesQueryBuilder from "../../timeSeries/components/TimeSeriesQueryBuilder";
import { useSdTypes } from "../../sdTypes/hooks/useSdTypes";
import { useKpiDefinitionsBySdType } from "../../kpi/hooks/useKpiDefinitionsBySdType";
import { useSdInstancesByType } from "../../sdInstances/hooks/useSdInstancesByType";
import {
  buildOptionSearchText,
  filterSelectOption,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

type SortOption = {
  value: "DESC" | "ASC";
  label: "DESC" | "ASC";
};

type Props = {
  onSubmit: (input: TimeSeriesReadInput) => void;
  onExport: (input: TimeSeriesReadInput) => void;
  initialInput?: TimeSeriesReadInput | null;
};

export default function TimeSeriesFilters({
  onSubmit,
  onExport,
  initialInput,
}: Props) {
  const { sdTypes } = useSdTypes();

  const [input, setInput] = useState<TimeSeriesReadInput>(
    initialInput ?? {
      type: "raw",
      sortDesc: true,
      limit: 50,
      sdInstanceIDs: [],
    },
  );

  const [selectedSdType, setSelectedSdType] = useState<string | null>(
    initialInput?.sdTypeID ?? null,
  );

  const [selectedKpis, setSelectedKpis] = useState<string[]>(
    initialInput?.kpiDefinitionIDs ?? [],
  );

  const { kpiDefinitions } = useKpiDefinitionsBySdType(selectedSdType);
  const { sdInstances } = useSdInstancesByType(selectedSdType);

  useEffect(() => {
    if (!initialInput) return;

    setInput(initialInput);
    setSelectedSdType(initialInput.sdTypeID ?? null);
    setSelectedKpis(initialInput.kpiDefinitionIDs ?? []);
  }, [initialInput]);

  useEffect(() => {
    if (!sdTypes.length) return;
    if (selectedSdType) return;

    const first = sdTypes[0];
    const id = String(first.id);

    setSelectedSdType(id);

    setInput((s) => ({
      ...s,
      sdTypeID: id,
      sdInstanceIDs: [],
      kpiDefinitionIDs: s.type === "kpi" ? [] : undefined,
    }));
  }, [sdTypes, selectedSdType]);

  const sdTypeOptions: Option[] = useMemo(
    () =>
      sdTypes.map((t: any) => ({
        value: String(t.id),
        label: t.label ?? t.uid ?? "",
        searchText: buildOptionSearchText(t.label, t.uid, String(t.id)),
      })),
    [sdTypes],
  );

  const kpiOptions: Option[] = useMemo(
    () =>
      (kpiDefinitions ?? []).map((k: any) => ({
        value: String(k.id),
        label: k.label || k.id,
        searchText: buildOptionSearchText(k.label, String(k.id)),
      })),
    [kpiDefinitions],
  );

  const instancesOptions: Option[] = useMemo(
    () =>
      (sdInstances ?? []).map((i: any) => ({
        value: String(i.id),
        label: i.label || i.uid,
        searchText: buildOptionSearchText(i.label, i.uid, String(i.id)),
      })),
    [sdInstances],
  );

  const tagParameters = useMemo(() => {
    if (!selectedSdType) return [];

    const type = sdTypes.find((t: any) => String(t.id) === selectedSdType);
    if (!type?.parameters) return [];

    return type.parameters
      .filter((p: any) => p.role === "tag")
      .map((p: any) => ({
        denotation: p.denotation,
        label: p.label ?? p.denotation,
      }));
  }, [sdTypes, selectedSdType]);

  return (
    <div className="card p-3 mb-3">
      <div className="row g-3">
        {/* TYPE */}
        <div className="col-md-4">
          <label className="form-label">Type</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={[
              { value: "raw", label: "RAW" },
              { value: "kpi", label: "KPI" },
            ]}
            value={{
              value: input.type,
              label: input.type.toUpperCase(),
            }}
            onChange={(v: SingleValue<Option>) => {
              const type = (v?.value ?? "raw") as "raw" | "kpi";

              setInput((s) => ({
                type,
                sortDesc: true,
                limit: 50,
                sdInstanceIDs: [],
                from: s.from,
                to: s.to,
                sdTypeID: selectedSdType ?? undefined,
                kpiDefinitionIDs: type === "kpi" ? [] : undefined,
              }));

              setSelectedKpis([]);
            }}
            {...virtualizedSelectProps}
          />
        </div>

        {/* SD TYPE */}
        <div className="col-md-4">
          <label className="form-label">Model</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={sdTypeOptions}
            value={
              sdTypeOptions.find((o) => o.value === selectedSdType) ?? null
            }
            onChange={(v: SingleValue<Option>) => {
              const id = v?.value ?? null;

              setSelectedSdType(id);
              setSelectedKpis([]);

              setInput((s) => ({
                ...s,
                sdTypeID: id ?? undefined,
                sdInstanceIDs: [],
                kpiDefinitionIDs: s.type === "kpi" ? [] : undefined,
              }));
            }}
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        {/* KPI */}
        {input.type === "kpi" && (
          <div className="col-md-4">
            <label className="form-label">KPI</label>
            <Select<Option, true>
              isMulti
              classNamePrefix="react-select"
              options={kpiOptions}
              value={kpiOptions.filter((o) => selectedKpis.includes(o.value))}
              onChange={(v: MultiValue<Option>) => {
                const ids = v.map((i) => i.value);

                setSelectedKpis(ids);

                setInput((s) => ({
                  ...s,
                  kpiDefinitionIDs: ids,
                }));
              }}
              closeMenuOnSelect={false}
              filterOption={filterSelectOption}
              {...virtualizedSelectProps}
            />
          </div>
        )}

        {/* INSTANCES */}
        <div className="col-md-4">
          <label className="form-label">Devices</label>
          <Select<Option, true>
            isMulti
            classNamePrefix="react-select"
            closeMenuOnSelect={false}
            options={instancesOptions}
            value={instancesOptions.filter((o) =>
              (input.sdInstanceIDs ?? []).includes(o.value),
            )}
            onChange={(v: MultiValue<Option>) => {
              setInput((s) => ({
                ...s,
                sdInstanceIDs: v.map((i) => i.value),
              }));
            }}
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        {/* FROM */}
        <div className="col-md-4">
          <label className="form-label">From</label>
          <input
            type="datetime-local"
            step="1"
            className="form-control"
            value={
              input.from
                ? new Date(input.from).toISOString().slice(0, 19)
                : ""
            }
            onChange={(e) =>
              setInput((s) => ({
                ...s,
                from: e.target.value
                  ? new Date(e.target.value).toISOString()
                  : undefined,
              }))
            }
          />
        </div>

        {/* TO */}
        <div className="col-md-4">
          <label className="form-label">To</label>
          <input
            type="datetime-local"
            step="1"
            className="form-control"
            value={
              input.to
                ? new Date(input.to).toISOString().slice(0, 19)
                : ""
            }
            onChange={(e) =>
              setInput((s) => ({
                ...s,
                to: e.target.value
                  ? new Date(e.target.value).toISOString()
                  : undefined,
              }))
            }
          />
        </div>

        {/* SORT */}
        <div className="col-md-4">
          <label className="form-label">Sort</label>
          <Select<SortOption, false>
            classNamePrefix="react-select"
            options={[
              { value: "DESC", label: "DESC" },
              { value: "ASC", label: "ASC" },
            ]}
            value={{
              value: input.sortDesc ? "DESC" : "ASC",
              label: input.sortDesc ? "DESC" : "ASC",
            }}
            onChange={(v: SingleValue<SortOption>) =>
              setInput((s) => ({
                ...s,
                sortDesc: v?.value === "DESC",
              }))
            }
            {...virtualizedSelectProps}
          />
        </div>
      </div>

      {/* TAG FILTER */}
      {(input.type === "raw" || input.type === "kpi") && selectedSdType && tagParameters.length > 0 && (
        <div className="col-12 mt-3">
          <label className="form-label">Tag filters</label>

          <TimeSeriesQueryBuilder
            parameters={tagParameters}
            onChange={(filter) =>
              setInput((s) => ({
                ...s,
                filters: filter,
              }))
            }
          />
        </div>
      )}

      {/* ACTIONS */}
      <div className="d-flex justify-content-end gap-2 mt-3">
        <button
          className="btn btn-outline-light"
          onClick={() => onExport(input)}
        >
          Download
        </button>

        <button className="btn btn-primary" onClick={() => onSubmit(input)}>
          Show
        </button>
      </div>
    </div>
  );
}
