/**
 * @file SdInstanceKpiSidebar.tsx
 * @brief Postranní panel KPI detailu otevřeného ze sledované instance.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useMemo } from "react";
import Select, { type SingleValue } from "react-select";
import {
  matchesSearchText,
  type SearchMode,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

import Box from "@mui/material/Box";

import SdInstanceKpiCard from "./SdInstanceKpiCard";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";
import type { SdInstanceDetailSidebarSort } from "../state/sdInstanceDetailPageState";
import EmptyStateNotice from "../../../components/EmptyStateNotice";
import UrlSyncedSearchInput from "../../../components/UrlSyncedSearchInput";

type Kpi = {
  uid?: string | null;
  label?: string | null;
};

type Props = {
  kpis?: Kpi[];
  selectedKpiUID?: string | null;
  loading?: any;
  search: string;
  searchMode: SearchMode;
  sort: SdInstanceDetailSidebarSort;
  onSearchChange: (value: string) => void;
  onSearchModeChange: (value: SearchMode) => void;
  onSortChange: (value: SdInstanceDetailSidebarSort) => void;
  onSelect: (uid: string) => void;
  onOpenDetail: () => void;
};

type SortOption = {
  value: SdInstanceDetailSidebarSort;
  label: string;
};

export default function SdInstanceKpiSidebar({
  kpis = [],
  selectedKpiUID,
  loading,
  search,
  searchMode,
  sort,
  onSearchChange,
  onSearchModeChange,
  onSortChange,
  onSelect,
  onOpenDetail,
}: Props) {
  const sortOptions: SortOption[] = [
    { value: "label_asc", label: "Name ↑" },
    { value: "label_desc", label: "Name ↓" },
  ];

  const selectedSortOption =
    sortOptions.find((o) => o.value === sort) ?? sortOptions[0];

  const filtered = useMemo(() => {
    let result = [...kpis];

    if (search.trim()) {
      result = result.filter((k) => {
        return matchesSearchText(search, searchMode, k.label, k.uid);
      });
    }

    result.sort((a, b) => {
      const aLabel = (a.label ?? "").toLowerCase();
      const bLabel = (b.label ?? "").toLowerCase();

      if (sort === "label_asc") return aLabel.localeCompare(bLabel);
      if (sort === "label_desc") return bLabel.localeCompare(aLabel);

      return 0;
    });

    return result;
  }, [kpis, search, searchMode, sort]);

  const selectedIndex = useMemo(
    () =>
      filtered.findIndex((kpi) => String(kpi.uid) === String(selectedKpiUID)),
    [filtered, selectedKpiUID],
  );

  if (loading) {
    return (
      <div className="d-flex h-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  return (
    <div
      className="card p-3 d-flex flex-column"
      style={{ height: "100%", minHeight: 0 }}
    >
      <h5 className="mb-2">KPIs</h5>

      {/* SORT */}
      <div className="mb-2">
        <label className="form-label">Sort</label>
        <Select<SortOption, false>
          classNamePrefix="react-select"
          options={sortOptions}
          value={selectedSortOption}
          onChange={(v: SingleValue<SortOption>) => {
            if (!v) return;
            onSortChange(v.value);
          }}
          isClearable={false}
          {...virtualizedSelectProps}
        />
      </div>

      {/* SEARCH */}
      <div className="mb-2">
        <label className="form-label">Search</label>
        <UrlSyncedSearchInput
          placeholder="Search..."
          value={search}
          onChange={onSearchChange}
          mode={searchMode}
          onModeChange={onSearchModeChange}
        />
      </div>

      {/* LIST */}
      <Box sx={{ flex: 1, minHeight: 0 }}>
        <VirtualizedList
          items={filtered}
          rowHeight={84}
          itemSpacing={5}
          scrollToIndex={selectedIndex}
          pinnedIndex={selectedIndex}
          style={{ height: "100%" }}
          emptyState={
            <EmptyStateNotice
              compact
              title="No matching KPI"
              description="Try a different search or sort setting."
            />
          }
          renderItem={(kpi) => (
            <SdInstanceKpiCard
              key={kpi.uid ?? kpi.label ?? ""}
              kpi={kpi}
              selected={String(kpi.uid) === String(selectedKpiUID)}
              onClick={() => kpi.uid && onSelect(String(kpi.uid))}
            />
          )}
        />
      </Box>

      {/* BUTTON */}
      <button
        className="btn btn-outline-light mt-3"
        disabled={!selectedKpiUID}
        onClick={onOpenDetail}
      >
        KPI details
      </button>
    </div>
  );
}
