import type { PageStateCodec } from "../../../app/navigation/usePageState";

export type KpiDetailSidebarSort = "label_asc" | "label_desc";

export type KpiDetailQueryState = {
  instQ: string;
  instSort: KpiDetailSidebarSort;
  inst: string | null;
  from: string | null;
  to: string | null;
};

const sortValues = new Set<KpiDetailSidebarSort>([
  "label_asc",
  "label_desc",
]);

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const kpiDetailPageStateCodec: PageStateCodec<
  KpiDetailQueryState,
  null
> = {
  decodeQuery(searchParams) {
    const instSort = searchParams.get("instSort");

    return {
      instQ: searchParams.get("instQ") ?? "",
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
