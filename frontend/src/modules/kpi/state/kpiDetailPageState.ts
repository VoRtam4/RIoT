/**
 * @file kpiDetailPageState.ts
 * @brief Perzistentní stav detailu KPI, vybrané instance a historického panelu.
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

export type KpiDetailSidebarSort = "label_asc" | "label_desc";

export type KpiDetailQueryState = {
  instQ: string;
  instQMode: SearchMode;
  instSort: KpiDetailSidebarSort;
  inst: string | null;
  from: string | null;
  to: string | null;
};

const sortValues = new Set<KpiDetailSidebarSort>(["label_asc", "label_desc"]);
function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const kpiDetailPageStateCodec: PageStateCodec<
  KpiDetailQueryState,
  null
> = {
  decodeQuery(searchParams) {
    const instSort = searchParams.get("instSort");
    const instQMode = searchParams.get("instQMode");

    return {
      instQ: searchParams.get("instQ") ?? "",
      instQMode: instQMode && isSearchMode(instQMode) ? instQMode : "or",
      instSort:
        instSort && sortValues.has(instSort as KpiDetailSidebarSort)
          ? (instSort as KpiDetailSidebarSort)
          : "label_asc",
      inst: nonEmpty(searchParams.get("inst")),
      from: nonEmpty(searchParams.get("from")),
      to: nonEmpty(searchParams.get("to")),
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.instQ.trim()) searchParams.set("instQ", query.instQ);
    if (query.instQMode !== "or") {
      searchParams.set("instQMode", query.instQMode);
    }
    if (query.instSort !== "label_asc") {
      searchParams.set("instSort", query.instSort);
    }
    if (query.inst) searchParams.set("inst", query.inst);
    if (query.from) searchParams.set("from", query.from);
    if (query.to) searchParams.set("to", query.to);

    return searchParams;
  },
  decodeEntry() {
    return null;
  },
  encodeEntry() {
    return null;
  },
};
