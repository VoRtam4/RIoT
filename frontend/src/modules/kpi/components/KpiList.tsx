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
} from "../../../utils/reactSelectSearch";

type Raw = {
  id: string;
  label?: string | null;
  sdTypeID: string;
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
  sort: SortOption;
  mode: ModeFilter;
  onSearchChange: (value: string) => void;
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
  sort,
  mode,
  onSearchChange,
  onSortChange,
  onModeChange,
  onSdTypeChange,
  loading,
}: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  const debouncedSearch = useDebouncedValue(search, 150);

  const base =
    sort === "label_desc"
      ? entry?.rawSortedDesc ?? []
      : entry?.rawSortedAsc ?? [];

  const filtered = useMemo(() => {
    if (!base.length) return [];

    return base.filter((kpi) => {
      if (mode !== "ALL_MODES" && kpi.sdInstanceMode !== mode) {
        return false;
      }

      return matchesSearchText(
        debouncedSearch,
        buildOptionSearchText(kpi.label, kpi.id, kpi.sdTypeID),
      );
    });
  }, [base, debouncedSearch, mode]);

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
        sort={sort}
        mode={mode}
        onSearchChange={onSearchChange}
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
            <KpiCard key={kpi.id} kpi={kpi} />
          )}
        />
      </div>
    </>
  );
}
