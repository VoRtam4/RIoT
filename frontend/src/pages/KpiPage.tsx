import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";

import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import { useKpiDefinitionsBySdType } from "../modules/kpi/hooks/useKpiDefinitionsBySdType";
import KpiList from "../modules/kpi/components/KpiList";

export default function KpiPage() {
  const navigate = useNavigate();

  const {
    sdTypes,
    loading: sdTypesLoading,
  } = useSdTypes();

  const [selectedSdType, setSelectedSdType] = useState<string | null>(null);

  const activeSdType = useMemo(() => {
    return selectedSdType ?? (sdTypes[0] ? String(sdTypes[0].id) : null);
  }, [sdTypes, selectedSdType]);

  const { entry, loading: kpisLoading, } = useKpiDefinitionsBySdType(activeSdType);

  const loading = sdTypesLoading || (selectedSdType !== null && kpisLoading);

  return (
    <div className="container mt-3">
      <KpiList
        entry={entry}
        sdTypes={sdTypes}
        selectedSdType={activeSdType}
        onSdTypeChange={setSelectedSdType}
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