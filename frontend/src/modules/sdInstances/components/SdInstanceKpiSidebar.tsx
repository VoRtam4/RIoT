import { useMemo } from "react";
import Select, { type SingleValue } from "react-select";
import { matchesSearchText } from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

import Box from "@mui/material/Box";

import SdInstanceKpiCard from "./SdInstanceKpiCard";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";
import type { SdInstanceDetailSidebarSort } from "../state/sdInstanceDetailPageState";
import EmptyStateNotice from "../../../components/EmptyStateNotice";

type Kpi = {
  id: string;
  label?: string | null;
};

type Props = {
  kpis?: Kpi[];
  selectedKpiId?: string | null;
  loading?: any;
  search: string;
  sort: SdInstanceDetailSidebarSort;
  onSearchChange: (value: string) => void;
  onSortChange: (value: SdInstanceDetailSidebarSort) => void;
  onSelect: (id: string) => void;
  onOpenDetail: () => void;
};

type SortOption = {
  value: SdInstanceDetailSidebarSort;
  label: string;
};

export default function SdInstanceKpiSidebar({
  kpis = [],
  selectedKpiId,
  loading,
  search,
  sort,
  onSearchChange,
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
        return matchesSearchText(search, k.label, String(k.id));
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
  }, [kpis, search, sort]);

  const selectedIndex = useMemo(
    () =>
      filtered.findIndex(
        (kpi) => String(kpi.id) === String(selectedKpiId),
      ),
    [filtered, selectedKpiId],
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
      <h5 className="mb-3">KPIs</h5>

      {/* SORT */}
      <div className="mb-3">
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
      <div className="mb-3">
        <label className="form-label">Search</label>
        <input
          className="form-control"
          placeholder="Search..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
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
              key={kpi.id}
              kpi={kpi}
              selected={String(kpi.id) === String(selectedKpiId)}
              onClick={() => onSelect(String(kpi.id))}
            />
          )}
        />
      </Box>

      {/* BUTTON */}
      <button
        className="btn btn-outline-light mt-3"
        disabled={!selectedKpiId}
        onClick={onOpenDetail}
      >
        KPI details
      </button>
    </div>
  );
}
