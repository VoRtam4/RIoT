import { useNavigate } from "react-router-dom";
import { useKpis } from "../modules/kpi/hooks/useKpis";
import { useSdTypes } from "../modules/sdTypes/hooks/useSdTypes";
import KpiList from "../modules/kpi/components/KpiList";

export default function KpiPage() {
  const navigate = useNavigate();
  const { kpis, loading: kpisLoading, error } = useKpis();
  const { sdTypes, sdTypeMap, loading: typesLoading } = useSdTypes();
  const loading = kpisLoading || typesLoading;

  if (error) {
    console.error(error);
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="container mt-3">
      <KpiList kpis={kpis} sdTypes={sdTypes} sdTypeMap={sdTypeMap} loading={loading}/>

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
