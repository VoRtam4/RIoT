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

export const useSdInstancesByType = (typeId: string | null) => {
  const ensure = useSdInstancesStore((s) => s.ensure);
  const refresh = useSdInstancesStore((s) => s.refresh);

  const entry = useSdInstancesStore((s) =>
    typeId ? s.byType[typeId] : undefined,
  );

  useEffect(() => {
    if (!typeId) return;

    if (!entry) {
      void refresh(typeId);
    } else {
      void ensure(typeId);
    }
  }, [typeId, entry, ensure, refresh]);

  return {
    entry,
    sdInstances: entry?.rawSortedAsc ?? [],
    loading: typeId ? entry?.isLoading ?? true : false,
    error: entry?.error ?? null,
  };
};