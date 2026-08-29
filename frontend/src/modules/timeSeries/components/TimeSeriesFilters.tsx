/**
 * @file TimeSeriesFilters.tsx
 * @brief Filtrační panel historických dat pro výběr typu, instance, času a parametrů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useMemo } from "react";
import Select, { type MultiValue, type SingleValue } from "react-select";

import type { FilterNodeInput } from "../../../generated/graphql";
import TimeSeriesQueryBuilder from "../../timeSeries/components/TimeSeriesQueryBuilder";
import { useSdTypes } from "../../sdTypes/hooks/useSdTypes";
import { useKpiDefinitionsBySdType } from "../../kpi/hooks/useKpiDefinitionsBySdType";
import { useSdInstancesByType } from "../../sdInstances/hooks/useSdInstancesByType";
import {
  buildOptionSearchText,
  filterSelectOption,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";
import {
  minusDaysLocal,
  nowLocal,
  toLocalInputValue,
  toUTCString,
} from "../../kpi/utils/dateTimeUtils";
import type { TimeSeriesDraftState } from "../state/timeSeriesPageState";

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

type SortOption = {
  value: "DESC" | "ASC" | "NONE";
  label: "DESC" | "ASC" | "NONE";
};

type Props = {
  value: TimeSeriesDraftState;
  onChange: (next: TimeSeriesDraftState) => void;
  onSubmit: () => void;
  onExport: () => void;
};

export default function TimeSeriesFilters({
  value,
  onChange,
  onSubmit,
  onExport,
}: Props) {
  const { sdTypes } = useSdTypes();
  const selectedSdType = useMemo(
    () => sdTypes.find((t: any) => String(t.uid) === value.sdTypeUID) ?? null,
    [sdTypes, value.sdTypeUID],
  );
  const { entry: kpiEntry } = useKpiDefinitionsBySdType(value.sdTypeUID);
  const { entry: instancesEntry } = useSdInstancesByType(value.sdTypeUID);

  useEffect(() => {
    if (!sdTypes.length) return;
    if (value.sdTypeUID) return;

    onChange({
      ...value,
      sdTypeUID: String(sdTypes[0].uid),
      from: value.from ?? toUTCString(minusDaysLocal(1)),
      to: value.to ?? toUTCString(nowLocal()),
    });
  }, [onChange, sdTypes, value]);

  const sdTypeOptions: Option[] = useMemo(
    () =>
      sdTypes.map((t: any) => ({
        value: String(t.uid),
        label: t.label ?? t.uid ?? "",
        searchText: buildOptionSearchText(t.label, t.uid),
      })),
    [sdTypes],
  );

  const kpiOptions: Option[] = useMemo(
    () =>
      (kpiEntry?.rawSortedAsc ?? []).map((k) => ({
        value: String(k.uid),
        label: k.label || k.uid || "",
        searchText: buildOptionSearchText(k.label, k.uid),
      })),
    [kpiEntry],
  );

  const instancesOptions: Option[] = useMemo(
    () =>
      (instancesEntry?.rawSortedAsc ?? []).map((i) => ({
        value: String(i.uid),
        label: i.label || i.uid,
        searchText: buildOptionSearchText(i.label, i.uid),
      })),
    [instancesEntry],
  );

  const tagParameters = useMemo(() => {
    if (!value.sdTypeUID) return [];

    if (!selectedSdType?.parameters) return [];

    return selectedSdType.parameters
      .filter((p: any) => p.role === "tag")
      .map((p: any) => ({
        denotation: p.denotation,
        label: p.label ?? p.denotation,
      }));
  }, [selectedSdType, value.sdTypeUID]);

  const updateValue = (patch: Partial<TimeSeriesDraftState>) => {
    onChange({
      ...value,
      ...patch,
    });
  };

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
              value: value.type,
              label: value.type.toUpperCase(),
            }}
            onChange={(v: SingleValue<Option>) => {
              const type = (v?.value ?? "raw") as "raw" | "kpi";

              onChange({
                ...value,
                type,
                sort: "DESC",
                limit: 100,
                sdInstanceUIDs: [],
                kpiDefinitionUIDs: type === "kpi" ? [] : [],
              });
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
              sdTypeOptions.find((o) => o.value === value.sdTypeUID) ?? null
            }
            onChange={(v: SingleValue<Option>) => {
              const uid = v?.value ?? null;

              onChange({
                ...value,
                sdTypeUID: uid,
                sdInstanceUIDs: [],
                kpiDefinitionUIDs: value.type === "kpi" ? [] : [],
              });
            }}
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        {/* KPI */}
        {value.type === "kpi" && (
          <div className="col-md-4">
            <label className="form-label">KPI</label>
            <Select<Option, true>
              isMulti
              classNamePrefix="react-select"
              options={kpiOptions}
              value={kpiOptions.filter((o) =>
                value.kpiDefinitionUIDs.includes(o.value),
              )}
              onChange={(v: MultiValue<Option>) => {
                updateValue({
                  kpiDefinitionUIDs: v.map((i) => i.value),
                });
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
              value.sdInstanceUIDs.includes(o.value),
            )}
            onChange={(v: MultiValue<Option>) => {
              updateValue({
                sdInstanceUIDs: v.map((i) => i.value),
              });
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
            value={value.from ? toLocalInputValue(new Date(value.from)) : ""}
            onChange={(e) =>
              updateValue({
                from: e.target.value ? toUTCString(e.target.value) : undefined,
              })
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
            value={value.to ? toLocalInputValue(new Date(value.to)) : ""}
            onChange={(e) =>
              updateValue({
                to: e.target.value ? toUTCString(e.target.value) : undefined,
              })
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
              { value: "NONE", label: "NONE" },
            ]}
            value={
              value.sort === "NONE"
                ? { value: "NONE", label: "NONE" }
                : value.sort === "DESC"
                  ? { value: "DESC", label: "DESC" }
                  : { value: "ASC", label: "ASC" }
            }
            onChange={(v: SingleValue<SortOption>) => {
              if (v?.value === "NONE") {
                updateValue({
                  sort: "NONE",
                  limit: undefined,
                });
                return;
              }

              updateValue({
                sort: v?.value ?? "DESC",
                limit: value.limit ?? 100,
              });
            }}
            {...virtualizedSelectProps}
          />
        </div>
      </div>

      {(value.type === "raw" || value.type === "kpi") &&
        value.sdTypeUID &&
        tagParameters.length > 0 && (
          <div className="col-12 mt-3">
            <label className="form-label">Tag filters</label>

            <TimeSeriesQueryBuilder
              parameters={tagParameters}
              value={value.filters as FilterNodeInput | undefined}
              onChange={(filter) => updateValue({ filters: filter })}
            />
          </div>
        )}

      {/* ACTIONS */}
      <div className="d-flex justify-content-end gap-2 mt-3">
        <button className="btn btn-outline-light" onClick={onExport}>
          Download
        </button>

        <button className="btn btn-primary" onClick={onSubmit}>
          Show
        </button>
      </div>
    </div>
  );
}
