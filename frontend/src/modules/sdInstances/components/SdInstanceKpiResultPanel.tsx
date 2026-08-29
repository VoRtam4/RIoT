/**
 * @file SdInstanceKpiResultPanel.tsx
 * @brief Panel aktuálních výsledků KPI pro vybranou sledovanou instanci.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useKpiResult } from "../../kpi/hooks/useKpiResult";
import { useKpiSubscription } from "../../kpi/hooks/useKpiSubscription";

import { colors } from "../../../theme/colors";

type Props = {
  kpiDefinitionUID?: string | null;
  sdInstanceUID?: string | null;
};

export default function SdInstanceKpiResultPanel({
  kpiDefinitionUID,
  sdInstanceUID,
}: Props) {
  const safeInstanceUID = sdInstanceUID ?? undefined;

  const { result } = useKpiResult(
    kpiDefinitionUID ?? undefined,
    safeInstanceUID,
  );

  const { latest } = useKpiSubscription(
    kpiDefinitionUID ?? undefined,
    safeInstanceUID,
  );

  const liveResult = latest ?? result;

  if (!kpiDefinitionUID) {
    return <div className="text-muted form-label">Vyber KPI definici</div>;
  }

  return (
    <div className="d-flex align-items-center gap-3">
      <div
        style={{
          width: 14,
          height: 14,
          borderRadius: "50%",
          backgroundColor: liveResult?.fulfilled
            ? colors.success
            : colors.error,
        }}
      />

      <span className="form-label">
        {liveResult?.fulfilled ? "True" : "False"}
      </span>

      {liveResult?.eventTime && (
        <small style={{ opacity: 0.6 }}>
          ({new Date(liveResult.eventTime).toLocaleString("cs-CZ")})
        </small>
      )}
    </div>
  );
}
