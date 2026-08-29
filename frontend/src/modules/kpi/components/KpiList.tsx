/**
 * @file KpiList.tsx
 * @brief Seznam KPI definic s podporou výběru a virtualizovaného vykreslení.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useMemo, useRef } from "react";
import type { SdTypesQuery } from "../../../generated/graphql";
import KpiCard from "./KpiCard";
import KpiFilters, { type ModeFilter, type SortOption } from "./KpiFilters";
import VirtualizedCardGrid from "../../../components/virtualization/VirtualizedCardGrid";
import { useDebouncedValue } from "../../../utils/useDebouncedValue";
import EmptyStateNotice from "../../../components/EmptyStateNotice";
import {
  buildOptionSearchText,
  matchesSearchText,
  type SearchMode,
} from "../../../utils/reactSelectSearch";

type Raw = {
  uid?: string | null;
  label?: string | null;
  sdTypeUID?: string | null;
  sdInstanceMode?: string | null;
};

type Props = {
  entry?: {
    rawSortedAsc: Raw[];
    rawSortedDesc: Raw[];
  };
  sdTypes: SdTypesQuery["sdTypes"];
  selectedSdType: string | null;
  search: string;
  searchMode: SearchMode;
  sort: SortOption;
  mode: ModeFilter;
  onSearchChange: (value: string) => void;
  onSearchModeChange: (value: SearchMode) => void;
  onSortChange: (value: SortOption) => void;
  onModeChange: (value: ModeFilter) => void;
  onSdTypeChange: (value: string | null) => void;
  loading?: boolean;
};

export default function KpiList({
  entry,
  sdTypes,
  selectedSdType,
  search,
  searchMode,
  sort,
  mode,
  onSearchChange,
  onSearchModeChange,
  onSortChange,
  onModeChange,
  onSdTypeChange,
  loading,
}: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  const debouncedSearch = useDebouncedValue(search, 150);

  const filtered = useMemo(() => {
    const base =
      sort === "label_desc"
        ? (entry?.rawSortedDesc ?? [])
        : (entry?.rawSortedAsc ?? []);

    if (!base.length) return [];

    return base.filter((kpi) => {
      if (mode !== "ALL_MODES" && kpi.sdInstanceMode !== mode) {
        return false;
      }

      return matchesSearchText(
        debouncedSearch,
        searchMode,
        buildOptionSearchText(kpi.label, kpi.uid, kpi.sdTypeUID),
      );
    });
  }, [entry, sort, debouncedSearch, searchMode, mode]);

  if (loading) {
    return (
      <div className="d-flex h-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  return (
    <>
      <KpiFilters
        sdTypes={sdTypes}
        selectedSdType={selectedSdType}
        search={search}
        searchMode={searchMode}
        sort={sort}
        mode={mode}
        onSearchChange={onSearchChange}
        onSearchModeChange={onSearchModeChange}
        onSortChange={onSortChange}
        onModeChange={onModeChange}
        onSdTypeChange={onSdTypeChange}
      />

      <div ref={containerRef}>
        <VirtualizedCardGrid
          items={filtered}
          rowHeight={156}
          minColumnWidth={300}
          style={{ height: "calc(100vh - 220px)" }}
          emptyState={
            <EmptyStateNotice
              title="No KPI match the current filter"
              description="Try a different search phrase, model or mode."
            />
          }
          renderItem={(kpi) => (
            <KpiCard key={kpi.uid ?? kpi.label ?? ""} kpi={kpi} />
          )}
        />
      </div>
    </>
  );
}
