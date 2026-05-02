import { create } from "zustand";
import { apolloClient } from "../../../app/apollo";
import {
  ApiKeysDocument,
  type ApiKeysQuery,
} from "../../../generated/graphql";

const TTL = 5 * 60_000;

type Raw = ApiKeysQuery["apiKeys"][number];

type Entry = {
  raw: Raw[];
  lastFetchedAt: number | null;
  isLoading: boolean;
  error: Error | null;
};

type Store = {
  entry: Entry | null;
  ensure: () => Promise<void>;
  refresh: () => Promise<void>;
};

export const useApiKeysStore = create<Store>((set, get) => ({
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
        raw: existing?.raw ?? [],
        isLoading: true,
        error: null,
        lastFetchedAt: existing?.lastFetchedAt ?? null,
      },
    });

    try {
      const { data } = await apolloClient.query<ApiKeysQuery>({
        query: ApiKeysDocument,
        fetchPolicy: "network-only",
      });

      set({
        entry: {
          raw: data?.apiKeys ?? [],
          lastFetchedAt: Date.now(),
          isLoading: false,
          error: null,
        },
      });
    } catch (e) {
      set({
        entry: {
          raw: [],
          lastFetchedAt: null,
          isLoading: false,
          error: e as Error,
        },
      });
    }
  },
}));