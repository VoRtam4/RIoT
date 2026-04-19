import { useNavigate } from "react-router-dom";
import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import { useSdInstancesByType } from "../modules/sdInstances/hooks/useSdInstancesByType";

import SdInstanceList from "../modules/sdInstances/components/SdInstanceList";
import SdInstanceFilters from "../modules/sdInstances/components/SdInstanceFilters";
import { buildOptionSearchText } from "../utils/reactSelectSearch";
import { useDebouncedValue } from "../utils/useDebouncedValue";
import { useMemo, useState } from "react";

type SortOption = "label_asc" | "label_desc";
const SEARCH_DEBOUNCE_MS = 150;

export default function SdInstancesPage() {
  const navigate = useNavigate();

  const { sdTypes, loading: typesLoading } = useSdTypes();
  const firstSdTypeId = sdTypes[0] ? String(sdTypes[0].id) : null;

  const [selectedSdType, setSelectedSdType] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortOption>("label_asc");
  const debouncedSearch = useDebouncedValue(search, SEARCH_DEBOUNCE_MS);
  const activeSdType = selectedSdType ?? firstSdTypeId;

  const { sdInstances, loading: instancesLoading } =
    useSdInstancesByType(activeSdType);

  const filteredInstances = useMemo(() => {
    const searchTerms = debouncedSearch
      .trim()
      .toLowerCase()
      .split(/\s+/)
      .filter(Boolean);

    let result = (sdInstances ?? []).map((instance: any) => ({
      instance,
      searchText: buildOptionSearchText(
        instance.label,
        instance.uid,
        String(instance.id),
      ).toLowerCase(),
      sortLabel: (instance.label ?? instance.uid ?? "").toLowerCase(),
    }));

    if (searchTerms.length > 0) {
      result = result.filter(({ searchText }) =>
        searchTerms.every((term) => searchText.includes(term))
      );
    }

    if (sort === "label_asc") {
      result = [...result].sort((a, b) =>
        a.sortLabel < b.sortLabel ? -1 : 1
      );
    } else if (sort === "label_desc") {
      result = [...result].sort((a, b) =>
        a.sortLabel > b.sortLabel ? -1 : 1
      );
    }

    return result.map(({ instance }) => instance);
  }, [debouncedSearch, sdInstances, sort]);

  const loading = typesLoading || (activeSdType !== null && instancesLoading);

  return (
    <div className="container mt-3">
      <SdInstanceFilters
        sdTypes={sdTypes}
        selectedSdType={activeSdType}
        onSdTypeChange={setSelectedSdType}
        search={search}
        onSearchChange={setSearch}
        sort={sort}
        onSortChange={setSort}
      />

      <SdInstanceList
        instances={filteredInstances}
        loading={loading}
        onOpen={(id: string) => navigate(`/sd-instance/${id}`)}
      />
    </div>
  );
}
