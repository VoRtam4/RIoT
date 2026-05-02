import { create } from "zustand";
import { apolloClient } from "../../../app/apollo";
import {
  SdTypesDocument,
  type SdTypesQuery,
} from "../../../generated/graphql";

const TTL = 5 * 60_000;

type Raw = SdTypesQuery["sdTypes"][number];

type Entry = {
  rawSortedAsc: Raw[];
  rawSortedDesc: Raw[];
  lastFetchedAt: number | null;
  isLoading: boolean;
  error: Error | null;
};

type Store = {
  entry: Entry | null;
  ensure: () => Promise<void>;
  refresh: () => Promise<void>;
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

export const useSdTypesStore = create<Store>((set, get) => ({
  entry: null,

  ensure: async () => {
    const entry = get().entry;

    if (entry?.lastFetchedAt && Date.now() - entry.lastFetchedAt < TTL) {
      return;
    }

    await get().refresh();
  },

  refresh: async () => {
    const existing = get().entry;

    if (existing?.isLoading) return;

    set({
      entry: {
        ...existing,
        rawSortedAsc: existing?.rawSortedAsc ?? [],
        rawSortedDesc: existing?.rawSortedDesc ?? [],
        isLoading: true,
        error: null,
        lastFetchedAt: existing?.lastFetchedAt ?? null,
      },
    });

    try {
      const { data } = await apolloClient.query<SdTypesQuery>({
        query: SdTypesDocument,
        fetchPolicy: "network-only",
      });

      const { rawSortedAsc, rawSortedDesc } = sort(data?.sdTypes ?? []);

      set({
        entry: {
          rawSortedAsc,
          rawSortedDesc,
          lastFetchedAt: Date.now(),
          isLoading: false,
          error: null,
        },
      });
    } catch (e) {
      set({
        entry: {
          rawSortedAsc: [],
          rawSortedDesc: [],
          lastFetchedAt: null,
          isLoading: false,
          error: e as Error,
        },
      });
    }
  },
}));