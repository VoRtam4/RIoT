import { useMemo, useRef, useState } from "react";
import type { SdTypesQuery } from "../../../generated/graphql";
import KpiCard from "./KpiCard";
import KpiFilters, { type ModeFilter, type SortOption } from "./KpiFilters";
import VirtualizedCardGrid from "../../../components/virtualization/VirtualizedCardGrid";
import { useDebouncedValue } from "../../../utils/useDebouncedValue";
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
  onSdTypeChange: (value: string | null) => void;
  loading?: boolean;
};

export default function KpiList({
  entry,
  sdTypes,
  selectedSdType,
  onSdTypeChange,
  loading,
}: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortOption>("label_asc");
  const [mode, setMode] = useState<ModeFilter>("ALL_MODES");

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
        onSearchChange={setSearch}
        onSortChange={setSort}
        onModeChange={setMode}
        onSdTypeChange={onSdTypeChange}
      />

      <div ref={containerRef}>
        <VirtualizedCardGrid
          items={filtered}
          rowHeight={156}
          minColumnWidth={300}
          style={{ height: "calc(100vh - 220px)" }}
          emptyState={<div className="text-muted form-label">No results</div>}
          renderItem={(kpi) => (
            <KpiCard key={kpi.id} kpi={kpi} />
          )}
        />
      </div>
    </>
  );
}
