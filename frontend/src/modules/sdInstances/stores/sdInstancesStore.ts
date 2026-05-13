/**
 * @file sdInstancesStore.ts
 * @brief Lokální store seznamu sledovaných instancí a vybrané instance.
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
  SdInstancesByTypeDocument,
  type SdInstancesByTypeQuery,
  type SdInstancesByTypeQueryVariables,
} from "../../../generated/graphql";

const TTL = 5 * 60_000;
const ERROR_RETRY_COOLDOWN = 5_000;

type Raw = SdInstancesByTypeQuery["sdInstancesByType"][number];

type Entry = {
  rawSortedAsc: Raw[];
  rawSortedDesc: Raw[];
  lastFetchedAt: number | null;
  lastAttemptAt: number | null;
  isLoading: boolean;
  error: string | null;
};

type Store = {
  byType: Record<string, Entry>;
  ensure: (typeId: string) => Promise<void>;
  refresh: (typeId: string) => Promise<void>;
};

function sort(items: Raw[]) {
  const sortedAsc = [...items].sort((a, b) =>
    (a.label ?? a.uid ?? "")
      .toLowerCase()
      .localeCompare((b.label ?? b.uid ?? "").toLowerCase()),
  );

  return {
    rawSortedAsc: sortedAsc,
    rawSortedDesc: [...sortedAsc].reverse(),
  };
}

export const useSdInstancesStore = create<Store>((set, get) => ({
  byType: {},

  ensure: async (typeId) => {
    const entry = get().byType[typeId];
    const now = Date.now();

    if (entry?.lastFetchedAt && now - entry.lastFetchedAt < TTL) {
      return;
    }

    if (entry?.error && entry.lastAttemptAt && now - entry.lastAttemptAt < ERROR_RETRY_COOLDOWN) {
      return;
    }

    await get().refresh(typeId);
  },

  refresh: async (typeId) => {
    const existing = get().byType[typeId];

    if (existing?.isLoading) return;

    set((s) => ({
      byType: {
        ...s.byType,
        [typeId]: {
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
        SdInstancesByTypeQuery,
        SdInstancesByTypeQueryVariables
      >({
        query: SdInstancesByTypeDocument,
        variables: { id: typeId },
        fetchPolicy: "network-only",
      });

      const { rawSortedAsc, rawSortedDesc } = sort(
        data?.sdInstancesByType ?? [],
      );

      set((s) => ({
        byType: {
          ...s.byType,
          [typeId]: {
            rawSortedAsc,
            rawSortedDesc,
            lastFetchedAt: Date.now(),
            lastAttemptAt: Date.now(),
            isLoading: false,
            error: null,
          },
        },
      }));
    } catch {
      set((s) => ({
        byType: {
          ...s.byType,
          [typeId]: {
            ...(s.byType[typeId] ?? {}),
            lastFetchedAt: null,
            lastAttemptAt: Date.now(),
            isLoading: false,
            error: "Failed",
          },
        },
      }));
    }
  },
}));
