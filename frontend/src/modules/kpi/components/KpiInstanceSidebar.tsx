/**
 * @file KpiInstanceSidebar.tsx
 * @brief Postranní panel instance otevřené v kontextu KPI definice.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import Select, { type SingleValue } from "react-select";
import {
  matchesSearchText,
  type SearchMode,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

import Box from "@mui/material/Box";

import KpiInstanceCard from "./KpiInstanceCard";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";
import type { KpiDetailSidebarSort } from "../state/kpiDetailPageState";
import EmptyStateNotice from "../../../components/EmptyStateNotice";
import UrlSyncedSearchInput from "../../../components/UrlSyncedSearchInput";

type Instance = {
  label?: string | null;
  uid?: string | null;
};

type Props = {
  instances?: Instance[];
  selectedInstanceUID?: string | null;
  loading?: boolean;
  search: string;
  searchMode: SearchMode;
  sort: KpiDetailSidebarSort;
  onSearchChange: (value: string) => void;
  onSearchModeChange: (value: SearchMode) => void;
  onSortChange: (value: KpiDetailSidebarSort) => void;
  onSelect: (uid: string) => void;
};

type SortOption = {
  value: KpiDetailSidebarSort;
  label: string;
};

export default function KpiInstanceSidebar({
  instances = [],
  selectedInstanceUID,
  loading,
  search,
  searchMode,
  sort,
  onSearchChange,
  onSearchModeChange,
  onSortChange,
  onSelect,
}: Props) {
  const navigate = useNavigate();

  const sortOptions: SortOption[] = [
    { value: "label_asc", label: "Name ↑" },
    { value: "label_desc", label: "Name ↓" },
  ];

  const selectedSortOption =
    sortOptions.find((o) => o.value === sort) ?? sortOptions[0];

  const filtered = useMemo(() => {
    if (!instances.length) return [];

    let result = instances;

    if (search.trim()) {
      result = result.filter((i) =>
        matchesSearchText(search, searchMode, i.label, i.uid),
      );
    }

    result = [...result].sort((a, b) => {
      const aLabel = (a.label ?? a.uid ?? "").toLowerCase();
      const bLabel = (b.label ?? b.uid ?? "").toLowerCase();

      if (sort === "label_asc") return aLabel.localeCompare(bLabel);
      if (sort === "label_desc") return bLabel.localeCompare(aLabel);

      return 0;
    });

    return result;
  }, [instances, search, searchMode, sort]);

  const selectedIndex = useMemo(
    () =>
      filtered.findIndex(
        (inst) => String(inst.uid) === String(selectedInstanceUID),
      ),
    [filtered, selectedInstanceUID],
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
      style={{
        height: "100%",
        minHeight: 0,
        overflow: "hidden",
      }}
    >
      <h5 className="mb-2">Devices</h5>

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
      <Box
        sx={{
          flex: 1,
          minHeight: 0,
          overflow: "hidden",
        }}
      >
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
              title="No matching devices"
              description="Try a different search or sort setting."
            />
          }
          renderItem={(inst) => (
            <KpiInstanceCard
              key={inst.uid ?? inst.label ?? ""}
              instance={inst}
              selected={String(inst.uid) === String(selectedInstanceUID)}
              onClick={() => inst.uid && onSelect(String(inst.uid))}
            />
          )}
        />
      </Box>

      {/* BUTTON */}
      <button
        className="btn btn-outline-light mt-3 w-100"
        disabled={!selectedInstanceUID}
        onClick={() => {
          if (!selectedInstanceUID) return;
          navigate(`/sd-instance/${selectedInstanceUID}`);
        }}
      >
        Device detail
      </button>
    </div>
  );
}
