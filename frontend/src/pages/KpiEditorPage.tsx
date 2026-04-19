import { useParams } from "react-router-dom";
import KpiEditorForm from "../modules/kpi/components/KpiEditorForm";

export default function KpiEditorPage() {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="container mt-3">
      <KpiEditorForm kpiId={id} />
    </div>
  );
}
