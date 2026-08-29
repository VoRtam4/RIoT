/**
 * @file KpiPage.tsx
 * @brief Stránka seznamu KPI definic s filtrováním podle typu a instance.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useMemo } from "react";
import { useNavigate } from "react-router-dom";

import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import { useKpiDefinitionsBySdType } from "../modules/kpi/hooks/useKpiDefinitionsBySdType";
import KpiList from "../modules/kpi/components/KpiList";
import { usePageState } from "../app/navigation/usePageState";
import { kpiPageStateCodec } from "../modules/kpi/state/kpiPageState";

export default function KpiPage() {
  const navigate = useNavigate();

  const { sdTypes, loading: sdTypesLoading } = useSdTypes();
  const { query, setPageState } = usePageState(kpiPageStateCodec);

  const activeSdType = useMemo(() => {
    return query.sdType ?? (sdTypes[0] ? String(sdTypes[0].uid) : null);
  }, [query.sdType, sdTypes]);

  useEffect(() => {
    if (!sdTypes.length || query.sdType) return;

    setPageState(
      {
        query: {
          ...query,
          sdType: String(sdTypes[0].uid),
        },
        entry: null,
      },
      { replace: true },
    );
  }, [query, sdTypes, setPageState]);

  const { entry, loading: kpisLoading } =
    useKpiDefinitionsBySdType(activeSdType);

  const loading = sdTypesLoading || (query.sdType !== null && kpisLoading);

  return (
    <div className="container mt-3">
      <KpiList
        entry={entry}
        sdTypes={sdTypes}
        selectedSdType={activeSdType}
        search={query.q}
        searchMode={query.qMode}
        sort={query.sort}
        mode={query.mode}
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
        onSearchModeChange={(qMode) =>
          setPageState(
            {
              query: {
                ...query,
                qMode,
              },
              entry: null,
            },
            { replace: true },
          )
        }
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
        onModeChange={(mode) =>
          setPageState(
            {
              query: {
                ...query,
                mode,
              },
              entry: null,
            },
            { replace: true },
          )
        }
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
        loading={loading}
      />

      <button
        className="btn btn-primary"
        style={{
          position: "fixed",
          bottom: 20,
          right: "calc(max(50px, 50vw - 600px))",
        }}
        onClick={() => navigate("/kpi/new")}
      >
        New
      </button>
    </div>
  );
}
