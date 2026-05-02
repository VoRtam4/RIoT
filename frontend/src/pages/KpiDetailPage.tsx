import { useNavigate, useParams } from "react-router-dom";
import { useState, useMemo } from "react";
import { useUserConfig } from "../modules/auth/hooks/useUserConfig";

import { useKpiDefinition } from "../modules/kpi/hooks/useKpiDefinition";
import { useSdType } from "../modules/sdTypes/hooks/useSdType";
import { useSdInstancesByType } from "../modules/sdInstances/hooks/useSdInstancesByType";
import { useSdInstancesByKpiDefinition } from "../modules/sdInstances/hooks/useSdInstancesByKpiDefinition";

import KpiInstanceSidebar from "../modules/kpi/components/KpiInstanceSidebar";
import KpiSdTypeDataPanel from "../modules/kpi/components/KpiSdTypeDataPanel";
import KpiResultHistoryPanel from "../modules/kpi/components/KpiResultHistoryPanel";

export default function KpiDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

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

  const [selectedInstanceId, setSelectedInstanceId] = useState<string | null>(null);

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

  const loading =
    kpiLoading ||
    sdTypeLoading ||
    (mode === "all"
      ? instancesByTypeLoading
      : mode === "selected"
        ? instancesByKpiLoading
        : false);

  return (
    <div className="container-fluid mt-3">
      {/* HEADER */}
      <div className="card p-3 mb-3">
        <div className="d-flex justify-content-between align-items-center flex-wrap gap-3">
          <div className="d-flex align-items-baseline gap-2">
            <span className="form-label" style={{ fontSize: "1.25rem", fontWeight: 600 }}>
              {kpi?.label || "\u00A0"}
            </span>

            <span className="form-label" style={{ fontSize: "0.9rem", opacity: 0.7 }}>
              ({sdType?.label || sdType?.uid || "\u00A0"})
            </span>
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
                navigate(
                  `/history?type=kpi&kpiDefinitionID=${id}${
                    selectedInstanceId ? `&sdInstanceID=${selectedInstanceId}` : ""
                  }`
                )
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
            height: "calc(100vh - 180px)",
          }}
        >
          <KpiInstanceSidebar
            instances={instances}
            selectedInstanceId={selectedInstanceId}
            loading={loading}
            onSelect={setSelectedInstanceId}
          />
        </div>

        <div className="me-3" style={{ display: "grid", gap: 16 }}>
          <div className="card p-3">
            <KpiResultHistoryPanel
              kpiDefinitionID={id}
              sdInstanceID={selectedInstanceId}
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