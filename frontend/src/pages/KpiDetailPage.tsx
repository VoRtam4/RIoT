import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useMemo, useRef } from "react";
import { useUserConfig } from "../modules/auth/hooks/useUserConfig";
import Chip from "@mui/material/Chip";

import { useKpiDefinition } from "../modules/kpi/hooks/useKpiDefinition";
import { useSdType } from "../modules/sdTypes/hooks/useSdType";
import { useSdInstancesByType } from "../modules/sdInstances/hooks/useSdInstancesByType";
import { useSdInstancesByKpiDefinition } from "../modules/sdInstances/hooks/useSdInstancesByKpiDefinition";

import KpiInstanceSidebar from "../modules/kpi/components/KpiInstanceSidebar";
import KpiSdTypeDataPanel from "../modules/kpi/components/KpiSdTypeDataPanel";
import KpiResultHistoryPanel from "../modules/kpi/components/KpiResultHistoryPanel";
import { usePageState } from "../app/navigation/usePageState";
import { kpiDetailPageStateCodec } from "../modules/kpi/state/kpiDetailPageState";
import {
  minusDaysLocal,
  nowLocal,
} from "../modules/kpi/utils/dateTimeUtils";

export default function KpiDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const defaultRangeRef = useRef({
    from: minusDaysLocal(1),
    to: nowLocal(),
  });

  const { config, toggleKpi } = useUserConfig();
  const isFavorite = id && config.favoriteKpis.includes(String(id));

  const { kpi, loading: kpiLoading } = useKpiDefinition(id);
  const { sdType, loading: sdTypeLoading } = useSdType(kpi?.sdTypeID);

  const mode = kpi?.sdInstanceMode?.toLowerCase();

  const { entry: instancesEntryByType, loading: instancesByTypeLoading } =
    useSdInstancesByType(
      mode === "all" ? kpi?.sdTypeID ?? null : null,
    );

  const { sdInstances: instancesByKpi, loading: instancesByKpiLoading } =
    useSdInstancesByKpiDefinition(
      mode === "selected" ? id ?? null : null,
      mode === "selected",
    );
  const { query, setPageState } = usePageState(kpiDetailPageStateCodec);

  const loading =
    kpiLoading ||
    sdTypeLoading ||
    (mode === "all"
      ? instancesByTypeLoading
      : mode === "selected"
        ? instancesByKpiLoading
        : false);

  const instances = useMemo(() => {
    if (!kpi) return [];

    if (mode === "all") {
      return instancesEntryByType?.rawSortedAsc ?? [];
    }

    if (mode === "selected") {
      return [...(instancesByKpi ?? [])].sort((a, b) =>
        (a.label ?? a.uid ?? "")
          .toLowerCase()
          .localeCompare((b.label ?? b.uid ?? "").toLowerCase()),
      );
    }

    return [];
  }, [kpi, mode, instancesEntryByType, instancesByKpi]);

  const selectedInstanceId = query.inst;

  const from = query.from ?? defaultRangeRef.current.from;
  const to = query.to ?? defaultRangeRef.current.to;

  useEffect(() => {
    if (
      !loading &&
      query.inst &&
      !instances.some((instance) => String(instance.id) === String(query.inst))
    ) {
      setPageState(
        {
          query: {
            ...query,
            inst: null,
          },
          entry: null,
        },
        { replace: true },
      );
    }
  }, [instances, loading, query, setPageState]);

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
                {kpi?.label ?? "\u00A0"}
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

              {kpi?.sdInstanceMode && (
                <Chip
                  label={String(kpi.sdInstanceMode).toUpperCase()}
                  size="small"
                  color="primary"
                />
              )}
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
              {kpi?.id ?? "\u00A0"} ({sdType?.uid ?? "\u00A0"})
            </div>
          </div>

          <div className="d-flex gap-2">
            <button
              className={`btn ${isFavorite ? "btn-danger" : "btn-outline-light"}`}
              onClick={() => id && toggleKpi(String(id))}
            >
              <i className={isFavorite ? "fa-solid fa-heart" : "fa-regular fa-heart"} />
            </button>

            <button
              className="btn btn-outline-light"
              onClick={() =>
                {
                  const searchParams = new URLSearchParams();
                  searchParams.set("type", "kpi");
                  searchParams.set("auto", "1");
                  if (kpi?.sdTypeID) {
                    searchParams.set("sdType", String(kpi.sdTypeID));
                  }

                  navigate(
                    {
                      pathname: "/history",
                      search: `?${searchParams.toString()}`,
                    },
                    {
                      state: {
                        v: 1,
                        sdInstanceIDs: selectedInstanceId ? [selectedInstanceId] : [],
                        kpiDefinitionIDs: id ? [String(id)] : [],
                      },
                    },
                  );
                }
              }
            >
              History
            </button>

            <button
              className="btn btn-primary"
              onClick={() => navigate(`/kpi/edit/${id}`)}
            >
              Edit
            </button>
          </div>
        </div>
      </div>

      {/* LAYOUT */}
      <div style={{ display: "grid", gridTemplateColumns: "25% 75%", gap: 16 }}>
        <div
          style={{
            position: "sticky",
            top: 16,
            height: "calc(100vh - 200px)",
          }}
        >
          <KpiInstanceSidebar
            instances={instances}
            selectedInstanceId={selectedInstanceId}
            loading={loading}
            search={query.instQ}
            sort={query.instSort}
            onSearchChange={(instQ) =>
              setPageState(
                {
                  query: {
                    ...query,
                    instQ,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
            onSortChange={(instSort) =>
              setPageState(
                {
                  query: {
                    ...query,
                    instSort,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
            onSelect={(inst) =>
              setPageState(
                {
                  query: {
                    ...query,
                    inst: selectedInstanceId === inst ? null : inst,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
          />
        </div>

        <div className="me-3" style={{ display: "grid", gap: 16 }}>
          <div className="card p-3">
            <KpiResultHistoryPanel
              kpiDefinitionID={id}
              sdInstanceID={selectedInstanceId}
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

          <div className="card p-0">
            <KpiSdTypeDataPanel
              sdInstanceID={selectedInstanceId}
              sdType={sdType}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
