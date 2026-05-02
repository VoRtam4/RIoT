import { useEffect } from "react";
import { useKpiDefinitionsBySdTypeStore } from "../stores/kpiDefinitionsBySdTypeStore";

export const useKpiDefinitionsBySdType = (id: string | null) => {
  const ensure = useKpiDefinitionsBySdTypeStore((s) => s.ensure);
  const refresh = useKpiDefinitionsBySdTypeStore((s) => s.refresh);

  const entry = useKpiDefinitionsBySdTypeStore((s) =>
    id ? s.byTypeId[id] : undefined,
  );

  useEffect(() => {
    if (!id) return;

    if (!entry) {
      void refresh(id);
    } else {
      void ensure(id);
    }
  }, [id, entry, ensure, refresh]);

  return {
    entry,
    kpiDefinitions: entry?.rawSortedAsc ?? [],
    loading: entry?.isLoading ?? true,
    error: entry?.error ?? null,
  };
};