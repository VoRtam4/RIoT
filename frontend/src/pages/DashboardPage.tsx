import { useNavigate } from "react-router-dom";

import { useUserConfig } from "../modules/auth/hooks/useUserConfig";
import { useKpiDefinition } from "../modules/kpi/hooks/useKpiDefinition";
import { useSdInstance } from "../modules/sdInstances/hooks/useSdInstance";

import KpiCard from "../modules/kpi/components/KpiCard";
import SdInstanceCard from "../modules/sdInstances/components/SdInstanceCard";
import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import VirtualizedList from "../components/virtualization/VirtualizedList";

export default function DashboardPage() {
  const navigate = useNavigate();

  const { config, loading } = useUserConfig();
  const { sdTypeMap } = useSdTypes();

  if (loading) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  return (
    <div className="container mt-3">
      <div className="row g-3">
        {/* SD INSTANCES */}
        <div className="col-md-6">
          <div className="card p-3 h-100">
            <h5 className="mb-3 form-label">Favourite devices</h5>

            <div
              style={{
                height: "min(60vh, 520px)",
              }}
            >
              <VirtualizedList
                items={config.favoriteSdInstances}
                rowHeight={136}
                threshold={20}
                style={{ height: "100%" }}
                emptyState={<div className="text-muted form-label">No results</div>}
                renderItem={(id) => (
                  <SdInstanceItem key={id} id={id} navigate={navigate} />
                )}
              />
            </div>
          </div>
        </div>

        {/* KPI */}
        <div className="col-md-6">
          <div className="card p-3 h-100">
            <h5 className="mb-3 form-label">Favourite KPI</h5>

            <div
              style={{
                height: "min(60vh, 520px)",
              }}
            >
              <VirtualizedList
                items={config.favoriteKpis}
                rowHeight={156}
                threshold={20}
                style={{ height: "100%" }}
                emptyState={<div className="text-muted form-label">No results</div>}
                renderItem={(id) => (
                  <KpiItem
                    key={id}
                    id={id}
                    navigate={navigate}
                    sdTypeMap={sdTypeMap}
                  />
                )}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function KpiItem({
  id,
  sdTypeMap,
}: {
  id: string;
  navigate: any;
  sdTypeMap: Map<string, string>;
}) {
  const { kpi, loading } = useKpiDefinition(id);

  if (loading || !kpi) return null;

  return (
    <KpiCard
      kpi={kpi}
      sdTypeMap={sdTypeMap}
    />
  );
}

function SdInstanceItem({
  id,
}: {
  id: string;
  navigate: any;
}) {
  const { sdInstance, loading } = useSdInstance(id);

  if (loading || !sdInstance) return null;

  return (
    <SdInstanceCard
      instance={sdInstance}
    />
  );
}
