/**
 * @file useApiKeys.ts
 * @brief Hook pro načítání seznamu API klíčů aktuálního uživatele.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { useApiKeysStore } from "../stores/apiKeysStore";

export const useApiKeys = () => {
  const ensure = useApiKeysStore((s) => s.ensure);
  const refresh = useApiKeysStore((s) => s.refresh);
  const entry = useApiKeysStore((s) => s.entry);

  useEffect(() => {
    if (!entry) {
      void refresh();
    } else {
      void ensure();
    }
  }, [entry, ensure, refresh]);

  return {
    entry,
    apiKeys: entry?.raw ?? [],
    loading: entry?.isLoading ?? true,
    error: entry?.error ?? null,
  };
};
