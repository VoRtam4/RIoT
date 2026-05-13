/**
 * @file SdInstanceFilters.tsx
 * @brief Filtrační ovládání seznamu sledovaných instancí.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import Select, { type SingleValue } from "react-select";
import {
  buildOptionSearchText,
  filterSelectOption,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

type SdType = {
  id: string;
  label?: string | null;
  uid?: string | null;
};

type SortOption = "label_asc" | "label_desc";

type Props = {
  sdTypes: SdType[];
  selectedSdType: string | null;
  onSdTypeChange: (id: string | null) => void;

  search: string;
  onSearchChange: (v: string) => void;

  sort: SortOption;
  onSortChange: (v: SortOption) => void;
};

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

export default function SdInstanceFilters({
  sdTypes,
  selectedSdType,
  onSdTypeChange,
  search,
  onSearchChange,
  sort,
  onSortChange,
}: Props) {
  const sdTypeOptions: Option[] = sdTypes.map((t) => ({
    value: String(t.id),
    label: t.label ?? t.uid ?? "",
    searchText: buildOptionSearchText(t.label, t.uid, String(t.id)),
  }));

  const sortOptions: Option[] = [
    { value: "label_asc", label: "Name ↑" },
    { value: "label_desc", label: "Name ↓" },
  ];

  return (
    <div className="card p-3 mb-3">
      <div className="row g-3">
        {/* SEARCH */}
        <div className="col-md-4">
          <label className="form-label">Search</label>
          <input
            className="form-control"
            placeholder="Name or UID..."
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
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
            onChange={(v: SingleValue<Option>) =>
              onSdTypeChange(v?.value ?? null)
            }
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        {/* SORT */}
        <div className="col-md-4">
          <label className="form-label">Sort</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={sortOptions}
            value={sortOptions.find((o) => o.value === sort) ?? null}
            onChange={(v: SingleValue<Option>) =>
              onSortChange((v?.value as SortOption) ?? "label_asc")
            }
            isClearable={false}
            {...virtualizedSelectProps}
          />
        </div>
      </div>
    </div>
  );
}
