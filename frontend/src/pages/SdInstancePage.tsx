import { useNavigate } from "react-router-dom";
import { useMemo, useEffect } from "react";

import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import { useSdInstancesByType } from "../modules/sdInstances/hooks/useSdInstancesByType";

import SdInstanceList from "../modules/sdInstances/components/SdInstanceList";
import SdInstanceFilters from "../modules/sdInstances/components/SdInstanceFilters";
import {
  buildOptionSearchText,
  matchesSearchText,
} from "../utils/reactSelectSearch";
import { useDebouncedValue } from "../utils/useDebouncedValue";
import { usePageState } from "../app/navigation/usePageState";
import { sdInstancePageStateCodec } from "../modules/sdInstances/state/sdInstancePageState";

export default function SdInstancesPage() {
  const navigate = useNavigate();

  const { sdTypes, loading: typesLoading } = useSdTypes();
  const { query, setPageState } = usePageState(sdInstancePageStateCodec);
  const debouncedSearch = useDebouncedValue(query.q, 150);

  useEffect(() => {
    if (!sdTypes.length) return;
    if (query.sdType) return;

    setPageState(
      {
        query: {
          ...query,
          sdType: String(sdTypes[0].id),
        },
        entry: null,
      },
      { replace: true },
    );
  }, [query, sdTypes, setPageState]);

  const { sdInstances, loading: instancesLoading } =
    useSdInstancesByType(query.sdType);

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

    if (query.sort === "label_asc") {
      result.sort((a, b) =>
        a.sortLabel < b.sortLabel ? -1 : 1,
      );
    } else {
      result.sort((a, b) =>
        a.sortLabel > b.sortLabel ? -1 : 1,
      );
    }

    return result.map((r) => r.instance);
  }, [sdInstances, debouncedSearch, query.sort]);

  const loading = typesLoading || (query.sdType !== null && instancesLoading);

  return (
    <div className="container mt-3">
      <SdInstanceFilters
        sdTypes={sdTypes}
        selectedSdType={query.sdType}
        onSdTypeChange={(sdType) =>
          setPageState(
            {
              query: {
                ...query,
                sdType,
              },
              entry: null,
            },
            { replace: true },
          )
        }
        search={query.q}
        onSearchChange={(q) =>
          setPageState(
            {
              query: {
                ...query,
                q,
              },
              entry: null,
            },
            { replace: true },
          )
        }
        sort={query.sort}
        onSortChange={(sort) =>
          setPageState(
            {
              query: {
                ...query,
                sort,
              },
              entry: null,
            },
            { replace: true },
          )
        }
      />

      <SdInstanceList
        instances={filteredInstances}
        loading={loading}
        onOpen={(id) => navigate(`/sd-instance/${id}`)}
      />
    </div>
  );
}
