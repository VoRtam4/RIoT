import { useMemo, useRef, useState } from "react";
import type { KpiDefinitionsQuery } from "../../../generated/graphql";
import KpiCard from "./KpiCard";
import KpiFilters, { type ModeFilter, type SortOption } from "./KpiFilters";
import VirtualizedCardGrid from "../../../components/virtualization/VirtualizedCardGrid";
import { buildOptionSearchText } from "../../../utils/reactSelectSearch";
import { useDebouncedValue } from "../../../utils/useDebouncedValue";

type Props = {
  kpis: KpiDefinitionsQuery["kpiDefinitions"];
  sdTypes: Array<{
    id: string | number;
    label?: string | null;
    uid?: string | null;
  }>;
  sdTypeMap: Map<string, string>;
  loading?: any;
};

export default function KpiList({ kpis, sdTypes, sdTypeMap, loading }: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortOption>("label_asc");
  const [mode, setMode] = useState<ModeFilter>("ALL_MODES");
  const [sdTypeFilter, setSdTypeFilter] = useState<string | null>(null);
  const debouncedSearch = useDebouncedValue(search, 150);
  const activeSdTypeFilter =
    sdTypeFilter ?? (sdTypes[0] ? String(sdTypes[0].id) : null);

  const filtered = useMemo(() => {
    const searchTerms = debouncedSearch
      .trim()
      .toLowerCase()
      .split(/\s+/)
      .filter(Boolean);

    let result = kpis.map((kpi) => ({
      kpi,
      searchText: buildOptionSearchText(
        kpi.label,
        String(kpi.id),
        kpi.sdTypeUID,
        sdTypeMap.get(kpi.sdTypeID),
      ).toLowerCase(),
      labelSort: kpi.label.toLowerCase(),
      typeSort: (kpi.sdTypeUID ?? "").toLowerCase(),
    }));

    if (searchTerms.length > 0) {
      result = result.filter(({ searchText }) =>
        searchTerms.every((term) => searchText.includes(term)),
      );
    }

    if (activeSdTypeFilter) {
      result = result.filter(
        ({ kpi }) => String(kpi.sdTypeID) === activeSdTypeFilter,
      );
    }

    if (mode !== "ALL_MODES") {
      result = result.filter(({ kpi }) => kpi.sdInstanceMode === mode);
    }

    result.sort((a, b) => {
      if (sort === "label_asc") return a.labelSort.localeCompare(b.labelSort);
      if (sort === "label_desc") return b.labelSort.localeCompare(a.labelSort);

      if (sort === "type_asc") return a.typeSort.localeCompare(b.typeSort);

      if (sort === "type_desc") return b.typeSort.localeCompare(a.typeSort);

      return 0;
    });

    return result.map(({ kpi }) => kpi);
  }, [activeSdTypeFilter, debouncedSearch, kpis, mode, sdTypeMap, sort]);

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
        selectedSdType={activeSdTypeFilter}
        onSearchChange={setSearch}
        onSortChange={setSort}
        onModeChange={setMode}
        onSdTypeChange={setSdTypeFilter}
      />

      <div ref={containerRef}>
        <VirtualizedCardGrid
          items={filtered}
          rowHeight={156}
          minColumnWidth={300}
          style={{ height: "calc(100vh - 220px)" }}
          emptyState={<div className="text-muted form-label">No results</div>}
          renderItem={(kpi) => (
            <KpiCard key={kpi.id} kpi={kpi} sdTypeMap={sdTypeMap} />
          )}
        />
      </div>
    </>
  );
}
