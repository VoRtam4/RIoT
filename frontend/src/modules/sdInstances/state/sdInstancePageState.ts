/**
 * @file sdInstancePageState.ts
 * @brief Perzistentní stav seznamu sledovaných instancí a filtrů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { PageStateCodec } from "../../../app/navigation/usePageState";
import {
  isSearchMode,
  type SearchMode,
} from "../../../utils/reactSelectSearch";

export type SdInstanceSort = "label_asc" | "label_desc";

export type SdInstancePageQueryState = {
  q: string;
  qMode: SearchMode;
  sdType: string | null;
  sort: SdInstanceSort;
};

const sortValues = new Set<SdInstanceSort>(["label_asc", "label_desc"]);
function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const sdInstancePageStateCodec: PageStateCodec<
  SdInstancePageQueryState,
  null
> = {
  decodeQuery(searchParams) {
    const sort = searchParams.get("sort");
    const qMode = searchParams.get("qMode");

    return {
      q: searchParams.get("q") ?? "",
      qMode: qMode && isSearchMode(qMode) ? qMode : "or",
      sdType: nonEmpty(searchParams.get("sdType")),
      sort:
        sort && sortValues.has(sort as SdInstanceSort)
          ? (sort as SdInstanceSort)
          : "label_asc",
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.q.trim()) searchParams.set("q", query.q);
    if (query.qMode !== "or") searchParams.set("qMode", query.qMode);
    if (query.sdType) searchParams.set("sdType", query.sdType);
    if (query.sort !== "label_asc") searchParams.set("sort", query.sort);

    return searchParams;
  },
  decodeEntry() {
    return null;
  },
  encodeEntry() {
    return null;
  },
};
