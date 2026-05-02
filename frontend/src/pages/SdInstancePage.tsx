import { useNavigate } from "react-router-dom";
import { useMemo, useState, useEffect } from "react";

import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import { useSdInstancesByType } from "../modules/sdInstances/hooks/useSdInstancesByType";

import SdInstanceList from "../modules/sdInstances/components/SdInstanceList";
import SdInstanceFilters from "../modules/sdInstances/components/SdInstanceFilters";
import {
  buildOptionSearchText,
  matchesSearchText,
} from "../utils/reactSelectSearch";
import { useDebouncedValue } from "../utils/useDebouncedValue";

type SortOption = "label_asc" | "label_desc";

export default function SdInstancesPage() {
  const navigate = useNavigate();

  const { sdTypes, loading: typesLoading } = useSdTypes();

  const [selectedSdType, setSelectedSdType] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortOption>("label_asc");

  const debouncedSearch = useDebouncedValue(search, 150);

  useEffect(() => {
    if (!sdTypes.length) return;
    if (selectedSdType) return;

    setSelectedSdType(String(sdTypes[0].id));
  }, [sdTypes, selectedSdType]);

  const { sdInstances, loading: instancesLoading } =
    useSdInstancesByType(selectedSdType);

  const filteredInstances = useMemo(() => {
    let result = (sdInstances ?? []).map((instance: any) => ({
      instance,
      searchText: buildOptionSearchText(
        instance.label,
        instance.uid,
        String(instance.id),
      ).toLowerCase(),
      sortLabel: (instance.label ?? instance.uid ?? "").toLowerCase(),
    }));

    if (debouncedSearch.trim()) {
      result = result.filter(({ searchText }) =>
        matchesSearchText(debouncedSearch, searchText),
      );
    }

    if (sort === "label_asc") {
      result.sort((a, b) =>
        a.sortLabel < b.sortLabel ? -1 : 1,
      );
    } else {
      result.sort((a, b) =>
        a.sortLabel > b.sortLabel ? -1 : 1,
      );
    }

    return result.map((r) => r.instance);
  }, [sdInstances, debouncedSearch, sort]);

  const loading = typesLoading || (selectedSdType !== null && instancesLoading);

  return (
    <div className="container mt-3">
      <SdInstanceFilters
        sdTypes={sdTypes}
        selectedSdType={selectedSdType}
        onSdTypeChange={setSelectedSdType}
        search={search}
        onSearchChange={setSearch}
        sort={sort}
        onSortChange={setSort}
      />

      <SdInstanceList
        instances={filteredInstances}
        loading={loading}
        onOpen={(id) => navigate(`/sd-instance/${id}`)}
      />
    </div>
  );
}
