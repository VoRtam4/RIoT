/**
 * @file KpiEditorPage.tsx
 * @brief Editor KPI definic s vizuálním skládáním podmínek nad parametry zdrojů dat.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
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
