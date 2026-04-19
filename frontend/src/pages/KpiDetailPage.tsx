import { useNavigate, useParams } from "react-router-dom";
import { useState } from "react";
import { useUserConfig } from "../modules/auth/hooks/useUserConfig";

import { useKpiDefinition } from "../modules/kpi/hooks/useKpiDefinition";
import { useSdType } from "../modules/sdTypes/hooks/useSdType";
import { useSdInstancesByKpiDefinition } from "../modules/sdInstances/hooks/useSdInstancesByKpiDefinition";

import KpiInstanceSidebar from "../modules/kpi/components/KpiInstanceSidebar";
import KpiSdTypeDataPanel from "../modules/kpi/components/KpiSdTypeDataPanel";
import KpiResultHistoryPanel from "../modules/kpi/components/KpiResultHistoryPanel";

export default function KpiDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { config, toggleKpi } = useUserConfig();
  const isFavorite =
    id && config.favoriteKpis.includes(String(id));

  const { kpi, loading: kpiLoading } = useKpiDefinition(id);
  const { sdType, loading: sdTypeLoading } = useSdType(kpi?.sdTypeID);

  const { sdInstances, loading: instancesLoading } = useSdInstancesByKpiDefinition(id ?? null);

  const [selectedInstanceId, setSelectedInstanceId] = useState<string | null>(null);

  const loading = kpiLoading || sdTypeLoading || instancesLoading;

  return (
    <div className="container-fluid mt-3">
      {/* HEADER */}
      <div className="card p-3 mb-3">
        <div className="d-flex justify-content-between align-items-center flex-wrap gap-3">
          <div className="d-flex align-items-baseline gap-2">
            <span
              className="form-label"
              style={{ fontSize: "1.25rem", fontWeight: 600 }}
            >
              {kpi?.label || "\u00A0"}
            </span>

            <span
              className="form-label"
              style={{ fontSize: "0.9rem", opacity: 0.7 }}
            >
              ({sdType?.label || sdType?.uid || "\u00A0"})
            </span>
          </div>

          <div className="d-flex gap-2">
            {/* FAVORITE */}
            <button
              className={`btn ${
                isFavorite ? "btn-danger" : "btn-outline-light"
              }`}
              onClick={() => {
                if (!id) return;
                toggleKpi(String(id));
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
                navigate(
                  `/history?type=kpi&kpiDefinitionID=${id}${
                    selectedInstanceId
                      ? `&sdInstanceID=${selectedInstanceId}`
                      : ""
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
            height: "calc(100vh - 180px)",
            overflow: "hidden",
          }}
        >
          <div style={{ height: "100%", overflow: "auto" }}>
            <KpiInstanceSidebar
              instances={sdInstances}
              selectedInstanceId={selectedInstanceId}
              loading={loading}
              onSelect={setSelectedInstanceId}
            />
          </div>
        </div>

        {/* RIGHT */}
        <div
          className="me-3"
          style={{
            display: "grid",
            gridTemplateRows: "auto auto",
            gap: 16,
          }}
        >
          {/* RESULT */}
          <div className="card p-3">
            <KpiResultHistoryPanel
              kpiDefinitionID={id}
              sdInstanceID={selectedInstanceId}
            />
          </div>

          {/* DATA */}
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