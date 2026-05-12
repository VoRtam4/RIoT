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
import {
  minusDaysLocal,
  nowLocal,
} from "../modules/kpi/utils/dateTimeUtils";

export default function SdInstanceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const defaultRangeRef = useRef({
    from: minusDaysLocal(1),
    to: nowLocal(),
  });

  const { config, toggleSdInstance } = useUserConfig();

  const isFavorite =
    id && config.favoriteSdInstances.includes(String(id));

  const { sdInstance: instance, loading: instanceLoading } =
    useSdInstance(id ?? null);

  const { kpiDefinitions, loading: kpiLoading } = useKpiDefinitionsBySdInstance(id ?? null);
  const { query, setPageState } = usePageState(sdInstanceDetailPageStateCodec);

  const { sdType } = useSdType(instance?.type?.id);
  const selectedKpiId = query.kpi;
  const from = query.from ?? defaultRangeRef.current.from;
  const to = query.to ?? defaultRangeRef.current.to;

  const loading = instanceLoading || kpiLoading;

  useEffect(() => {
    if (
      !loading &&
      query.kpi &&
      !(kpiDefinitions ?? []).some((kpi) => String(kpi.id) === String(query.kpi))
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
                if (!id) return;
                toggleSdInstance(String(id));
              }}
            >
              <i
                className={
                  isFavorite
                    ? "fa-solid fa-heart"
                    : "fa-regular fa-heart"
                }
              />
            </button>

            {/* HISTORY */}
            <button
              className="btn btn-outline-light"
              onClick={() =>
                {
                  const searchParams = new URLSearchParams();
                  searchParams.set("type", "raw");
                  searchParams.set("auto", "1");
                  if (instance.type?.id) {
                    searchParams.set("sdType", String(instance.type.id));
                  }

                  navigate(
                    {
                      pathname: "/history",
                      search: `?${searchParams.toString()}`,
                    },
                    {
                      state: {
                        v: 1,
                        sdInstanceIDs: id ? [String(id)] : [],
                        kpiDefinitionIDs: [],
                      },
                    },
                  );
                }
              }
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
              selectedKpiId={selectedKpiId}
              loading={loading}
              search={query.kpiQ}
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
              onSelect={(nextKpiId) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      kpi:
                        String(selectedKpiId) === String(nextKpiId)
                          ? null
                          : String(nextKpiId),
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onOpenDetail={() => {
                if (!selectedKpiId) return;
                navigate(`/kpi/${selectedKpiId}`);
              }}
            />
          </div>
        </div>

        {/* RIGHT */}
        <div
          className="me-3"
          style={{
            display: "grid",
            gridTemplateRows: selectedKpiId ? "auto auto auto" : "auto auto",
            gap: 16,
            alignContent: "start",
          }}
        >

          {/* KPI HISTORY */}
          {selectedKpiId ? (
            <div className="card p-3">
              <KpiResultHistoryPanel
                kpiDefinitionID={selectedKpiId ?? undefined}
                sdInstanceID={id}
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
              sdInstanceID={id ?? null}
              sdType={sdType}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
