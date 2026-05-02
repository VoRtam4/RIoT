import { useMemo, useState } from "react";
import Select from "react-select";
import { matchesSearchText } from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";

import APIKeyCard from "./APIKeyCard";

type Props = {
  apiKeys: any[];
  selectedId: string | null;
  loading?: any;
  onSelect: (id: string) => void;
};

type FilterOption = {
  value: "all" | "active" | "inactive";
  label: string;
};

type SortOption = {
  value: "label_asc" | "label_desc";
  label: string;
};

export default function APIKeyList({
  apiKeys,
  selectedId,
  loading,
  onSelect,
}: Props) {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<FilterOption["value"]>("all");
  const [sort, setSort] = useState<SortOption["value"]>("label_asc");

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
    () => filtered.findIndex((k) => k.id === selectedId),
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
          onChange={(v) => v && setFilter(v.value)}
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
          onChange={(v) => v && setSort(v.value)}
          isClearable={false}
          {...virtualizedSelectProps}
        />
      </div>

      {/* SEARCH */}
      <div className="mb-2">
        <label className="form-label">Search</label>
        <input
          className="form-control"
          placeholder="Hledat..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {/* LIST */}
      <VirtualizedList
        items={filtered}
        rowHeight={92}
        scrollToIndex={selectedIndex}
        pinnedIndex={selectedIndex}
        style={{ flex: 1, minHeight: 0 }}
        emptyState={<div className="text-muted small form-label">No results</div>}
        renderItem={(k) => (
          <APIKeyCard
            key={k.id}
            apiKey={k}
            selected={selectedId === k.id}
            onClick={() => onSelect(k.id)}
          />
        )}
      />
    </div>
  );
}
