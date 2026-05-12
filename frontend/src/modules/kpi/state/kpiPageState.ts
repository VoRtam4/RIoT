import type { PageStateCodec } from "../../../app/navigation/usePageState";
import type {
  ModeFilter,
  SortOption,
} from "../components/KpiFilters";

export type KpiPageQueryState = {
  q: string;
  sdType: string | null;
  sort: SortOption;
  mode: ModeFilter;
};

const sortValues = new Set<SortOption>(["label_asc", "label_desc"]);
const modeValues = new Set<ModeFilter>(["ALL_MODES", "all", "selected"]);

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const kpiPageStateCodec: PageStateCodec<KpiPageQueryState, null> = {
  decodeQuery(searchParams) {
    const sort = searchParams.get("sort");
    const mode = searchParams.get("mode");

    return {
      q: searchParams.get("q") ?? "",
      sdType: nonEmpty(searchParams.get("sdType")),
      sort:
        sort && sortValues.has(sort as SortOption)
          ? (sort as SortOption)
          : "label_asc",
      mode:
        mode && modeValues.has(mode as ModeFilter)
          ? (mode as ModeFilter)
          : "ALL_MODES",
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.q.trim()) searchParams.set("q", query.q);
    if (query.sdType) searchParams.set("sdType", query.sdType);
    if (query.sort !== "label_asc") searchParams.set("sort", query.sort);
    if (query.mode !== "ALL_MODES") searchParams.set("mode", query.mode);

    return searchParams;
  },
  decodeEntry() {
    return null;
  },
  encodeEntry() {
    return null;
  },
};
