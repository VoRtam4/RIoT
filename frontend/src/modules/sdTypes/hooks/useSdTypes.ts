/**
 * @file useSdTypes.ts
 * @brief Hook pro načítání seznamu sledovaných typů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { useSdTypesStore } from "../stores/sdTypesStore";

export const useSdTypes = () => {
  const ensure = useSdTypesStore((s) => s.ensure);
  const refresh = useSdTypesStore((s) => s.refresh);
  const entry = useSdTypesStore((s) => s.entry);

  useEffect(() => {
    if (!entry) {
      void refresh();
    } else {
      void ensure();
    }
  }, [entry, ensure, refresh]);

  return {
    entry,
    sdTypes: entry?.rawSortedAsc ?? [],
    loading: entry?.isLoading ?? true,
    error: entry?.error ?? null,
  };
};
