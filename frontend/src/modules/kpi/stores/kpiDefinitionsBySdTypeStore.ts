/**
 * @file kpiDefinitionsBySdTypeStore.ts
 * @brief Lokální cache KPI definic seskupených podle sledovaného typu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { create } from "zustand";
import { apolloClient } from "../../../app/apollo";
import {
  KpiDefinitionsBySdTypeDocument,
  type KpiDefinitionsBySdTypeQuery,
  type KpiDefinitionsBySdTypeQueryVariables,
} from "../../../generated/graphql";

const TTL = 60_000;
const ERROR_RETRY_COOLDOWN = 5_000;

type Raw = KpiDefinitionsBySdTypeQuery["kpiDefinitionsBySdType"][number];

type Entry = {
  rawSortedAsc: Raw[];
  rawSortedDesc: Raw[];
  lastFetchedAt: number | null;
  lastAttemptAt: number | null;
  isLoading: boolean;
  error: Error | null;
};

type Store = {
  byTypeUID: Record<string, Entry>;
  ensure: (uid: string) => Promise<void>;
  refresh: (uid: string) => Promise<void>;
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
  byTypeUID: {},

  ensure: async (uid) => {
    const entry = get().byTypeUID[uid];
    const now = Date.now();

    if (entry?.lastFetchedAt && now - entry.lastFetchedAt < TTL) {
      return;
    }

    if (
      entry?.error &&
      entry.lastAttemptAt &&
      now - entry.lastAttemptAt < ERROR_RETRY_COOLDOWN
    ) {
      return;
    }

    await get().refresh(uid);
  },

  refresh: async (uid) => {
    const existing = get().byTypeUID[uid];

    if (existing?.isLoading) return;

    set((s) => ({
      byTypeUID: {
        ...s.byTypeUID,
        [uid]: {
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
        variables: { uid },
        fetchPolicy: "network-only",
      });

      const { rawSortedAsc, rawSortedDesc } = sort(
        data?.kpiDefinitionsBySdType ?? [],
      );

      set((s) => ({
        byTypeUID: {
          ...s.byTypeUID,
          [uid]: {
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
        byTypeUID: {
          ...s.byTypeUID,
          [uid]: {
            ...(s.byTypeUID[uid] ?? {}),
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
