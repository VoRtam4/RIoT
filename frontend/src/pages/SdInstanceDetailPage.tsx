/**
 * @file SdInstanceDetailPage.tsx
 * @brief Detail sledované instance s raw daty, KPI a historickými výsledky.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useRef } from "react";
import { useUserConfig } from "../modules/auth/hooks/useUserConfig";
import { useSdInstance } from "../modules/sdInstances/hooks/useSdInstance";
import { useKpiDefinitionsBySdInstance } from "../modules/kpi/hooks/useKpiDefinitionsBySdInstance";
import { useSdType } from "../modules/sdTypes/hooks/useSdType";

import SdInstanceKpiSidebar from "../modules/sdInstances/components/SdInstanceKpiSidebar";
import KpiSdTypeDataPanel from "../modules/kpi/components/KpiSdTypeDataPanel";
import KpiResultHistoryPanel from "../modules/kpi/components/KpiResultHistoryPanel";
import { usePageState } from "../app/navigation/usePageState";
import { sdInstanceDetailPageStateCodec } from "../modules/sdInstances/state/sdInstanceDetailPageState";
import { minusDaysLocal, nowLocal } from "../modules/kpi/utils/dateTimeUtils";

export default function SdInstanceDetailPage() {
  const { uid } = useParams<{ uid: string }>();
  const navigate = useNavigate();
  const defaultRangeRef = useRef({
    from: minusDaysLocal(1),
    to: nowLocal(),
  });

  const { config, toggleSdInstance } = useUserConfig();

  const isFavorite = uid && config.favoriteSdInstances.includes(String(uid));

  const { sdInstance: instance, loading: instanceLoading } = useSdInstance(
    uid ?? null,
  );

  const { kpiDefinitions, loading: kpiLoading } = useKpiDefinitionsBySdInstance(
    uid ?? null,
  );
  const { query, setPageState } = usePageState(sdInstanceDetailPageStateCodec);

  const { sdType } = useSdType(instance?.type?.uid);
  const selectedKpiUID = query.kpi;
  const selectedKpi =
    kpiDefinitions.find((kpi) => String(kpi.uid) === String(selectedKpiUID)) ??
    null;
  const from = query.from ?? defaultRangeRef.current.from;
  const to = query.to ?? defaultRangeRef.current.to;

  const loading = instanceLoading || kpiLoading;

  useEffect(() => {
    if (
      !loading &&
      query.kpi &&
      !(kpiDefinitions ?? []).some(
        (kpi) => String(kpi.uid) === String(query.kpi),
      )
    ) {
      setPageState(
        {
          query: {
            ...query,
            kpi: null,
          },
          entry: null,
        },
        { replace: true },
      );
    }
  }, [kpiDefinitions, loading, query, setPageState]);

  if (!instance) {
    return <div className="container mt-3">Not found</div>;
  }

  return (
    <div className="container-fluid mt-3">
      {/* HEADER */}
      <div className="card p-3 mb-3">
        <div className="d-flex justify-content-between align-items-start flex-wrap gap-3">
          <div style={{ minWidth: 0, flex: "1 1 320px" }}>
            <div
              className="d-flex align-items-center flex-wrap gap-2"
              style={{ minWidth: 0 }}
            >
              <span
                className="form-label mb-0"
                style={{
                  fontSize: "1.25rem",
                  fontWeight: 600,
                  minWidth: 0,
                  overflowWrap: "anywhere",
                  wordBreak: "break-word",
                }}
              >
                {instance.label ?? "\u00A0"}
              </span>

              <span
                className="form-label mb-0"
                style={{
                  fontSize: "0.9rem",
                  opacity: 0.7,
                  minWidth: 0,
                  overflowWrap: "anywhere",
                  wordBreak: "break-word",
                }}
              >
                ({sdType?.label ?? "\u00A0"})
              </span>
            </div>

            <div
              className="form-label mt-1 mb-0"
              style={{
                fontSize: "0.9rem",
                opacity: 0.7,
                minWidth: 0,
                overflowWrap: "anywhere",
                wordBreak: "break-word",
              }}
            >
              {instance.uid ?? "\u00A0"} ({sdType?.uid ?? "\u00A0"})
            </div>
          </div>

          <div className="d-flex gap-2">
            {/* FAVORITE */}
            <button
              className={`btn ${
                isFavorite ? "btn-danger" : "btn-outline-light"
              }`}
              onClick={() => {
                if (!uid) return;
                toggleSdInstance(String(uid));
              }}
            >
              <i
                className={
                  isFavorite ? "fa-solid fa-heart" : "fa-regular fa-heart"
                }
              />
            </button>

            {/* HISTORY */}
            <button
              className="btn btn-outline-light"
              onClick={() => {
                const searchParams = new URLSearchParams();
                searchParams.set("type", "raw");
                searchParams.set("auto", "1");
                if (instance.type?.uid) {
                  searchParams.set("sdType", String(instance.type.uid));
                }

                navigate(
                  {
                    pathname: "/history",
                    search: `?${searchParams.toString()}`,
                  },
                  {
                    state: {
                      v: 1,
                      sdInstanceUIDs: instance.uid
                        ? [String(instance.uid)]
                        : [],
                      kpiDefinitionUIDs: [],
                    },
                  },
                );
              }}
            >
              History
            </button>
          </div>
        </div>
      </div>

      {/* LAYOUT */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "25% 75%",
          gap: 16,
        }}
      >
        {/* LEFT */}
        <div
          style={{
            position: "sticky",
            top: 16,
            height: "calc(100vh - 200px)",
            overflow: "hidden",
          }}
        >
          <div style={{ height: "100%", overflow: "auto" }}>
            <SdInstanceKpiSidebar
              kpis={kpiDefinitions}
              selectedKpiUID={selectedKpiUID}
              loading={loading}
              search={query.kpiQ}
              searchMode={query.kpiQMode}
              sort={query.kpiSort}
              onSearchChange={(kpiQ) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      kpiQ,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onSearchModeChange={(kpiQMode) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      kpiQMode,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onSortChange={(kpiSort) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      kpiSort,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onSelect={(nextKpiUID) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      kpi:
                        String(selectedKpiUID) === String(nextKpiUID)
                          ? null
                          : String(nextKpiUID),
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onOpenDetail={() => {
                if (!selectedKpiUID) return;
                navigate(`/kpi/${selectedKpiUID}`);
              }}
            />
          </div>
        </div>

        {/* RIGHT */}
        <div
          className="me-3"
          style={{
            display: "grid",
            gridTemplateRows: selectedKpiUID ? "auto auto auto" : "auto auto",
            gap: 16,
            alignContent: "start",
          }}
        >
          {/* KPI HISTORY */}
          {selectedKpiUID ? (
            <div className="card p-3">
              <KpiResultHistoryPanel
                kpiDefinitionUID={selectedKpi?.uid}
                sdInstanceUID={instance.uid}
                sdTypeUID={instance.type?.uid}
                from={from}
                to={to}
                onFromChange={(nextFrom) =>
                  setPageState(
                    {
                      query: {
                        ...query,
                        from: nextFrom,
                      },
                      entry: null,
                    },
                    { replace: true },
                  )
                }
                onToChange={(nextTo) =>
                  setPageState(
                    {
                      query: {
                        ...query,
                        to: nextTo,
                      },
                      entry: null,
                    },
                    { replace: true },
                  )
                }
              />
            </div>
          ) : null}

          {/* RAW DATA */}
          <div className="card p-0">
            <KpiSdTypeDataPanel
              sdInstanceUID={instance.uid ?? null}
              sdType={sdType}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
