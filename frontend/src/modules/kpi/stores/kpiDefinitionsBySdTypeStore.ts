import { create } from "zustand";
import { apolloClient } from "../../../app/apollo";
import {
  KpiDefinitionsBySdTypeDocument,
  type KpiDefinitionsBySdTypeQuery,
  type KpiDefinitionsBySdTypeQueryVariables,
} from "../../../generated/graphql";

const TTL = 60_000;
const ERROR_RETRY_COOLDOWN = 5_000;

type Raw =
  KpiDefinitionsBySdTypeQuery["kpiDefinitionsBySdType"][number];

type Entry = {
  rawSortedAsc: Raw[];
  rawSortedDesc: Raw[];
  lastFetchedAt: number | null;
  lastAttemptAt: number | null;
  isLoading: boolean;
  error: Error | null;
};

type Store = {
  byTypeId: Record<string, Entry>;
  ensure: (id: string) => Promise<void>;
  refresh: (id: string) => Promise<void>;
};

function sort(items: Raw[]) {
  const sortedAsc = [...items].sort((a, b) =>
    (a.label ?? "").toLowerCase().localeCompare((b.label ?? "").toLowerCase()),
  );

  return {
    rawSortedAsc: sortedAsc,
    rawSortedDesc: [...sortedAsc].reverse(),
  };
}

export const useKpiDefinitionsBySdTypeStore = create<Store>((set, get) => ({
  byTypeId: {},

  ensure: async (id) => {
    const entry = get().byTypeId[id];
    const now = Date.now();

    if (entry?.lastFetchedAt && now - entry.lastFetchedAt < TTL) {
      return;
    }

    if (entry?.error && entry.lastAttemptAt && now - entry.lastAttemptAt < ERROR_RETRY_COOLDOWN) {
      return;
    }

    await get().refresh(id);
  },

  refresh: async (id) => {
    const existing = get().byTypeId[id];

    if (existing?.isLoading) return;

    set((s) => ({
      byTypeId: {
        ...s.byTypeId,
        [id]: {
          ...existing,
          rawSortedAsc: existing?.rawSortedAsc ?? [],
          rawSortedDesc: existing?.rawSortedDesc ?? [],
          isLoading: true,
          error: null,
          lastFetchedAt: existing?.lastFetchedAt ?? null,
          lastAttemptAt: Date.now(),
        },
      },
    }));

    try {
      const { data } = await apolloClient.query<
        KpiDefinitionsBySdTypeQuery,
        KpiDefinitionsBySdTypeQueryVariables
      >({
        query: KpiDefinitionsBySdTypeDocument,
        variables: { id },
        fetchPolicy: "network-only",
      });

      const { rawSortedAsc, rawSortedDesc } = sort(
        data?.kpiDefinitionsBySdType ?? [],
      );

      set((s) => ({
        byTypeId: {
          ...s.byTypeId,
          [id]: {
            rawSortedAsc,
            rawSortedDesc,
            lastFetchedAt: Date.now(),
            lastAttemptAt: Date.now(),
            isLoading: false,
            error: null,
          },
        },
      }));
    } catch (error) {
      set((s) => ({
        byTypeId: {
          ...s.byTypeId,
          [id]: {
            ...(s.byTypeId[id] ?? {}),
            lastFetchedAt: null,
            lastAttemptAt: Date.now(),
            isLoading: false,
            error: error as Error,
          },
        },
      }));
    }
  },
}));
