/**
 * @file useSdInstancesByType.ts
 * @brief Hook pro načítání instancí podle sledovaného typu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { useSdInstancesStore } from "../stores/sdInstancesStore";

export const useSdInstancesByType = (typeUID: string | null) => {
  const ensure = useSdInstancesStore((s) => s.ensure);
  const refresh = useSdInstancesStore((s) => s.refresh);

  const entry = useSdInstancesStore((s) =>
    typeUID ? s.byType[typeUID] : undefined,
  );

  useEffect(() => {
    if (!typeUID) return;

    if (!entry) {
      void refresh(typeUID);
    } else {
      void ensure(typeUID);
    }
  }, [typeUID, entry, ensure, refresh]);

  return {
    entry,
    sdInstances: entry?.rawSortedAsc ?? [],
    loading: typeUID ? (entry?.isLoading ?? true) : false,
    error: entry?.error ?? null,
  };
};
