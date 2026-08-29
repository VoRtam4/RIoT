/**
 * @file timeSeriesPageState.ts
 * @brief Perzistentní stav stránky historických dat, filtrů, dotazu a exportů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { PageStateCodec } from "../../../app/navigation/usePageState";
import type {
  FilterNodeInput,
  TimeSeriesReadInput,
} from "../../../generated/graphql";
import {
  minusDaysLocal,
  nowLocal,
  toUTCString,
} from "../../kpi/utils/dateTimeUtils";

export type TimeSeriesSort = "DESC" | "ASC" | "NONE";

export type TimeSeriesDraftState = {
  type: "raw" | "kpi";
  sdTypeUID: string | null;
  sdInstanceUIDs: string[];
  kpiDefinitionUIDs: string[];
  from?: string;
  to?: string;
  sort: TimeSeriesSort;
  limit?: number;
  filters?: FilterNodeInput;
  autoRun: boolean;
};

export type TimeSeriesEntryState = {
  v: 1;
  sdInstanceUIDs: string[];
  kpiDefinitionUIDs: string[];
};

const sortValues = new Set<TimeSeriesSort>(["DESC", "ASC", "NONE"]);

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

function toBase64Url(value: string) {
  const bytes = new TextEncoder().encode(value);
  let binary = "";

  bytes.forEach((byte) => {
    binary += String.fromCharCode(byte);
  });

  return btoa(binary)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/g, "");
}

function fromBase64Url(value: string) {
  const normalized = value.replace(/-/g, "+").replace(/_/g, "/");
  const paddingLength = (4 - (normalized.length % 4)) % 4;
  const padded = normalized + "=".repeat(paddingLength);
  const binary = atob(padded);
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));

  return new TextDecoder().decode(bytes);
}

export function encodeTimeSeriesFilter(
  filter?: FilterNodeInput,
): string | undefined {
  if (!filter) return undefined;

  try {
    return toBase64Url(JSON.stringify(filter));
  } catch {
    return undefined;
  }
}

export function decodeTimeSeriesFilter(
  encoded?: string | null,
): FilterNodeInput | undefined {
  if (!encoded) return undefined;

  try {
    return JSON.parse(fromBase64Url(encoded)) as FilterNodeInput;
  } catch {
    return undefined;
  }
}

export function createDefaultTimeSeriesDraft(): TimeSeriesDraftState {
  return {
    type: "raw",
    sdTypeUID: null,
    sdInstanceUIDs: [],
    kpiDefinitionUIDs: [],
    from: toUTCString(minusDaysLocal(1)),
    to: toUTCString(nowLocal()),
    sort: "DESC",
    limit: 100,
    filters: undefined,
    autoRun: false,
  };
}

function normalizeEntry(state: unknown): TimeSeriesEntryState {
  if (!state || typeof state !== "object") {
    return {
      v: 1,
      sdInstanceUIDs: [],
      kpiDefinitionUIDs: [],
    };
  }

  const candidate = state as Partial<TimeSeriesEntryState>;

  return {
    v: 1,
    sdInstanceUIDs: Array.isArray(candidate.sdInstanceUIDs)
      ? candidate.sdInstanceUIDs.map(String)
      : [],
    kpiDefinitionUIDs: Array.isArray(candidate.kpiDefinitionUIDs)
      ? candidate.kpiDefinitionUIDs.map(String)
      : [],
  };
}

export function mergeTimeSeriesDraft(
  query: Omit<TimeSeriesDraftState, "sdInstanceUIDs" | "kpiDefinitionUIDs">,
  entry: TimeSeriesEntryState,
): TimeSeriesDraftState {
  return {
    ...query,
    sdInstanceUIDs: entry.sdInstanceUIDs,
    kpiDefinitionUIDs: entry.kpiDefinitionUIDs,
  };
}

export function draftToTimeSeriesInput(
  draft: TimeSeriesDraftState,
): TimeSeriesReadInput {
  return {
    type: draft.type,
    sdTypeUID: draft.sdTypeUID ?? undefined,
    sdInstanceUIDs: draft.sdInstanceUIDs,
    kpiDefinitionUIDs:
      draft.type === "kpi" ? draft.kpiDefinitionUIDs : undefined,
    from: draft.from,
    to: draft.to,
    filters: draft.filters,
    sortDesc: draft.sort === "NONE" ? undefined : draft.sort === "DESC",
    limit: draft.sort === "NONE" ? undefined : (draft.limit ?? 100),
  };
}

export const timeSeriesPageStateCodec: PageStateCodec<
  Omit<TimeSeriesDraftState, "sdInstanceUIDs" | "kpiDefinitionUIDs">,
  TimeSeriesEntryState
> = {
  decodeQuery(searchParams) {
    const defaults = createDefaultTimeSeriesDraft();
    const type = searchParams.get("type");
    const sort = searchParams.get("sort");
    const limitValue = searchParams.get("limit");
    const limit = limitValue ? Number(limitValue) : defaults.limit;

    return {
      ...defaults,
      type: type === "kpi" ? "kpi" : "raw",
      sdTypeUID: nonEmpty(searchParams.get("sdType")),
      from: nonEmpty(searchParams.get("from")) ?? defaults.from,
      to: nonEmpty(searchParams.get("to")) ?? defaults.to,
      sort:
        sort && sortValues.has(sort as TimeSeriesSort)
          ? (sort as TimeSeriesSort)
          : "DESC",
      limit:
        typeof limit === "number" && Number.isFinite(limit)
          ? limit
          : defaults.limit,
      filters: decodeTimeSeriesFilter(searchParams.get("tags")),
      autoRun: searchParams.get("auto") === "1",
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    searchParams.set("type", query.type);
    if (query.sdTypeUID) searchParams.set("sdType", query.sdTypeUID);
    if (query.from) searchParams.set("from", query.from);
    if (query.to) searchParams.set("to", query.to);
    if (query.sort !== "DESC") searchParams.set("sort", query.sort);
    if (
      typeof query.limit === "number" &&
      Number.isFinite(query.limit) &&
      query.limit !== 100
    ) {
      searchParams.set("limit", String(query.limit));
    }

    const encodedFilter = encodeTimeSeriesFilter(query.filters);
    if (encodedFilter) searchParams.set("tags", encodedFilter);
    if (query.autoRun) searchParams.set("auto", "1");

    return searchParams;
  },
  decodeEntry: normalizeEntry,
  encodeEntry(entry) {
    return {
      v: 1,
      sdInstanceUIDs: entry.sdInstanceUIDs,
      kpiDefinitionUIDs: entry.kpiDefinitionUIDs,
    };
  },
};
