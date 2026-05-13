/**
 * @file sdInstanceDetailPageState.ts
 * @brief Perzistentní stav detailu instance, vybraných KPI a raw panelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { PageStateCodec } from "../../../app/navigation/usePageState";

export type SdInstanceDetailSidebarSort = "label_asc" | "label_desc";

export type SdInstanceDetailQueryState = {
  kpiQ: string;
  kpiSort: SdInstanceDetailSidebarSort;
  kpi: string | null;
  from: string | null;
  to: string | null;
};

const sortValues = new Set<SdInstanceDetailSidebarSort>([
  "label_asc",
  "label_desc",
]);

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const sdInstanceDetailPageStateCodec: PageStateCodec<
  SdInstanceDetailQueryState,
  null
> = {
  decodeQuery(searchParams) {
    const kpiSort = searchParams.get("kpiSort");

    return {
      kpiQ: searchParams.get("kpiQ") ?? "",
      kpiSort:
        kpiSort && sortValues.has(kpiSort as SdInstanceDetailSidebarSort)
          ? (kpiSort as SdInstanceDetailSidebarSort)
          : "label_asc",
      kpi: nonEmpty(searchParams.get("kpi")),
      from: nonEmpty(searchParams.get("from")),
      to: nonEmpty(searchParams.get("to")),
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.kpiQ.trim()) searchParams.set("kpiQ", query.kpiQ);
    if (query.kpiSort !== "label_asc") {
      searchParams.set("kpiSort", query.kpiSort);
    }
    if (query.kpi) searchParams.set("kpi", query.kpi);
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
