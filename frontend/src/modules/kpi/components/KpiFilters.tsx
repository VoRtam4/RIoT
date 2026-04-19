import { useState, useMemo } from "react";
import Select from "react-select";
import {
  buildOptionSearchText,
  filterSelectOption,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

export type SortOption =
  | "label_asc"
  | "label_desc"
  | "type_asc"
  | "type_desc";

export type ModeFilter = "ALL_MODES" | "all" | "selected";

type Props = {
  sdTypes: Array<{
    id: string | number;
    label?: string | null;
    uid?: string | null;
  }>;
  selectedSdType: string | null;
  onSearchChange: (v: string) => void;
  onSortChange: (v: SortOption) => void;
  onModeChange: (v: ModeFilter) => void;
  onSdTypeChange: (v: string | null) => void;
};

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

export default function KpiFilters({
  sdTypes,
  selectedSdType,
  onSearchChange,
  onSortChange,
  onModeChange,
  onSdTypeChange,
}: Props) {
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortOption>("label_asc");
  const [mode, setMode] = useState<ModeFilter>("ALL_MODES");

  const sdTypeOptions: Option[] = useMemo(() => {
    return sdTypes.map((sdType) => ({
      value: String(sdType.id),
      label: sdType.label ?? sdType.uid ?? String(sdType.id),
      searchText: buildOptionSearchText(
        sdType.label,
        sdType.uid,
        String(sdType.id),
      ),
    }));
  }, [sdTypes]);

  const sortOptions: Option[] = [
    { value: "label_asc", label: "Name ↑" },
    { value: "label_desc", label: "Name ↓" },
    { value: "type_asc", label: "Type ↑" },
    { value: "type_desc", label: "Type ↓" },
  ];

  const modeOptions: Option[] = [
    { value: "ALL_MODES", label: "\u00A0" },
    { value: "ALL", label: "ALL" },
    { value: "SELECTED", label: "SELECTED" },
  ];

  return (
    <div className="card p-3 mb-3">
      <div className="row g-3">
        {/* SEARCH */}
        <div className="col-md-3">
          <label className="form-label">Search</label>
          <input
            className="form-control"
            placeholder="Name, KPI ID or type..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              onSearchChange(e.target.value);
            }}
          />
        </div>

        {/* SD TYPE */}
        <div className="col-md-3">
          <label className="form-label">Model</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={sdTypeOptions}
            value={sdTypeOptions.find((o) => o.value === selectedSdType) || null}
            onChange={(v) => {
              if (!v) return;
              onSdTypeChange(v.value);
            }}
            isClearable={false}
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        {/* SORT */}
        <div className="col-md-3">
          <label className="form-label">Sort</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={sortOptions}
            value={sortOptions.find((o) => o.value === sort) || null}
            onChange={(v) => {
              if (!v) return;
              const value = v.value as SortOption;
              setSort(value);
              onSortChange(value);
            }}
            isClearable={false}
            {...virtualizedSelectProps}
          />
        </div>

        {/* MODE */}
        <div className="col-md-3">
          <label className="form-label">Mode</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={modeOptions}
            value={modeOptions.find((o) => o.value === mode) || null}
            onChange={(v) => {
              if (!v) return;
              const value = v.value as ModeFilter;
              setMode(value);
              onModeChange(value);
            }}
            isClearable={false}
            {...virtualizedSelectProps}
          />
        </div>
      </div>
    </div>
  );
}
