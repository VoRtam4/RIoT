/**
 * @file APIKeySidebar.tsx
 * @brief Postranní panel detailu a akcí nad vybraným API klíčem.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useMemo } from "react";
import Select from "react-select";
import { matchesSearchText } from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";
import EmptyStateNotice from "../../../components/EmptyStateNotice";
import UrlSyncedSearchInput from "../../../components/UrlSyncedSearchInput";
import type {
  ApiKeysFilter,
  ApiKeysSort,
} from "../state/apiKeysPageState";

import APIKeyCard from "./APIKeyCard";

type Props = {
  apiKeys: any[];
  selectedId: string | null;
  loading?: any;
  search: string;
  filter: ApiKeysFilter;
  sort: ApiKeysSort;
  onSearchChange: (value: string) => void;
  onFilterChange: (value: ApiKeysFilter) => void;
  onSortChange: (value: ApiKeysSort) => void;
  onSelect: (id: string) => void;
};

type FilterOption = {
  value: ApiKeysFilter;
  label: string;
};

type SortOption = {
  value: ApiKeysSort;
  label: string;
};

export default function APIKeySidebar({
  apiKeys,
  selectedId,
  loading,
  search,
  filter,
  sort,
  onSearchChange,
  onFilterChange,
  onSortChange,
  onSelect,
}: Props) {
  const isActive = (k: any) =>
    !k.revoked && (!k.expiresAt || new Date(k.expiresAt) > new Date());

  const filterOptions: FilterOption[] = [
    { value: "all", label: "All" },
    { value: "active", label: "Active" },
    { value: "inactive", label: "Inactive" },
  ];

  const sortOptions: SortOption[] = [
    { value: "label_asc", label: "Name ↑" },
    { value: "label_desc", label: "Name ↓" },
  ];

  const filtered = useMemo(() => {
    let result = [...apiKeys];

    if (search.trim()) {
      result = result.filter((k) =>
        matchesSearchText(search, k.label),
      );
    }

    if (filter === "active") {
      result = result.filter((k) => isActive(k));
    }

    if (filter === "inactive") {
      result = result.filter((k) => !isActive(k));
    }

    result.sort((a, b) => {
      const aLabel = (a.label ?? "").toLowerCase();
      const bLabel = (b.label ?? "").toLowerCase();

      if (sort === "label_asc") return aLabel.localeCompare(bLabel);
      if (sort === "label_desc") return bLabel.localeCompare(aLabel);

      return 0;
    });

    return result;
  }, [apiKeys, search, filter, sort]);

  const selectedIndex = useMemo(
    () => filtered.findIndex((k) => String(k.id) === String(selectedId)),
    [filtered, selectedId],
  );

  if (loading) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" role="status" />
      </div>
    );
  }

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        gap: 8,
        height: "100%",
        minHeight: 0,
        overflow: "hidden",
      }}
    >
      {/* FILTER */}
      <div className="mb-2">
        <label className="form-label">Filter</label>
        <Select
          classNamePrefix="react-select"
          options={filterOptions}
          value={filterOptions.find((f) => f.value === filter)}
          onChange={(v) => v && onFilterChange(v.value)}
          isClearable={false}
          {...virtualizedSelectProps}
        />
      </div>

      {/* SORT */}
      <div className="mb-3">
        <label className="form-label">Sort</label>
        <Select
          classNamePrefix="react-select"
          options={sortOptions}
          value={sortOptions.find((s) => s.value === sort)}
          onChange={(v) => v && onSortChange(v.value)}
          isClearable={false}
          {...virtualizedSelectProps}
        />
      </div>

      {/* SEARCH */}
      <div className="mb-2">
        <label className="form-label">Search</label>
        <UrlSyncedSearchInput
          placeholder="Hledat..."
          value={search}
          onChange={onSearchChange}
        />
      </div>

      {/* LIST */}
      <VirtualizedList
        items={filtered}
        rowHeight={84}
        itemSpacing={5}
        scrollToIndex={selectedIndex}
        pinnedIndex={selectedIndex}
        style={{ flex: 1, minHeight: 0 }}
        emptyState={
          <EmptyStateNotice
            compact
            title="No matching keys"
            description="Adjust the search or filter settings."
          />
        }
        renderItem={(k) => (
          <APIKeyCard
            key={k.id}
            apiKey={k}
            selected={String(selectedId) === String(k.id)}
            onClick={() => onSelect(String(k.id))}
          />
        )}
      />
    </div>
  );
}
