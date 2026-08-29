/**
 * @file apiKeysStore.ts
 * @brief Lokální store stavu seznamu API klíčů a vybrané položky.
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
import { ApiKeysDocument, type ApiKeysQuery } from "../../../generated/graphql";

const TTL = 5 * 60_000;
const ERROR_RETRY_COOLDOWN = 5_000;

type Raw = ApiKeysQuery["apiKeys"][number];

type Entry = {
  raw: Raw[];
  lastFetchedAt: number | null;
  lastAttemptAt: number | null;
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
        lastAttemptAt: Date.now(),
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
          lastAttemptAt: Date.now(),
          isLoading: false,
          error: null,
        },
      });
    } catch (e) {
      set({
        entry: {
          raw: [],
          lastFetchedAt: null,
          lastAttemptAt: Date.now(),
          isLoading: false,
          error: e as Error,
        },
      });
    }
  },
}));
