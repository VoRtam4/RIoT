import { useNavigate, useParams } from "react-router-dom";
import { useState } from "react";
import { useUserConfig } from "../modules/auth/hooks/useUserConfig";
import { useSdInstance } from "../modules/sdInstances/hooks/useSdInstance";
import { useKpiDefinitionsBySdInstance } from "../modules/kpi/hooks/useKpiDefinitionsBySdInstance";
import { useSdType } from "../modules/sdTypes/hooks/useSdType";

import SdInstanceKpiSidebar from "../modules/sdInstances/components/SdInstanceKpiSidebar";
import KpiSdTypeDataPanel from "../modules/kpi/components/KpiSdTypeDataPanel";
import KpiResultHistoryPanel from "../modules/kpi/components/KpiResultHistoryPanel";

export default function SdInstanceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { config, toggleSdInstance } = useUserConfig();

  const isFavorite =
    id && config.favoriteSdInstances.includes(String(id));

  const { sdInstance: instance, loading: instanceLoading } =
    useSdInstance(id ?? null);

  const { kpiDefinitions, loading: kpiLoading } = useKpiDefinitionsBySdInstance(id ?? null);

  const [selectedKpiId, setSelectedKpiId] = useState<string | null>(null);

  const { sdType } = useSdType(instance?.type?.id);

  const loading = instanceLoading || kpiLoading;

  if (!instance) {
    return <div className="container mt-3">Not found</div>;
  }

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
              {instance.label || instance.uid || instance.id}
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
                navigate(
                  `/history?type=kpi&sdInstanceID=${id}${
                    selectedKpiId
                      ? `&kpiDefinitionID=${selectedKpiId}`
                      : ""
                  }`
                )
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
            height: "calc(100vh - 180px)",
            overflow: "hidden",
          }}
        >
          <div style={{ height: "100%", overflow: "auto" }}>
            <SdInstanceKpiSidebar
              kpis={kpiDefinitions}
              selectedKpiId={selectedKpiId}
              loading={loading}
              onSelect={setSelectedKpiId}
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
            gridTemplateRows: "auto auto auto",
            gap: 16,
            alignContent: "start",
          }}
        >

          {/* KPI HISTORY */}
          <div className="card p-3">
            <KpiResultHistoryPanel
              kpiDefinitionID={selectedKpiId ?? undefined}
              sdInstanceID={id}
            />
          </div>

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