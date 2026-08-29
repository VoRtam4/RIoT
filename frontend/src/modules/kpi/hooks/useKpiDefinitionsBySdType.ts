/**
 * @file useKpiDefinitionsBySdType.ts
 * @brief Hook pro načítání KPI definic navázaných na konkrétní sledovaný typ.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { useKpiDefinitionsBySdTypeStore } from "../stores/kpiDefinitionsBySdTypeStore";

export const useKpiDefinitionsBySdType = (uid: string | null) => {
  const ensure = useKpiDefinitionsBySdTypeStore((s) => s.ensure);
  const refresh = useKpiDefinitionsBySdTypeStore((s) => s.refresh);

  const entry = useKpiDefinitionsBySdTypeStore((s) =>
    uid ? s.byTypeUID[uid] : undefined,
  );

  useEffect(() => {
    if (!uid) return;

    if (!entry) {
      void refresh(uid);
    } else {
      void ensure(uid);
    }
  }, [uid, entry, ensure, refresh]);

  return {
    entry,
    kpiDefinitions: entry?.rawSortedAsc ?? [],
    loading: entry?.isLoading ?? true,
    error: entry?.error ?? null,
  };
};
