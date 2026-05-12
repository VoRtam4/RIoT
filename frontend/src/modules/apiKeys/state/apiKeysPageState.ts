import type { PageStateCodec } from "../../../app/navigation/usePageState";

export type ApiKeysFilter = "all" | "active" | "inactive";
export type ApiKeysSort = "label_asc" | "label_desc";

export type ApiKeysQueryState = {
  q: string;
  filter: ApiKeysFilter;
  sort: ApiKeysSort;
  selected: string | null;
};

const filterValues = new Set<ApiKeysFilter>(["all", "active", "inactive"]);
const sortValues = new Set<ApiKeysSort>(["label_asc", "label_desc"]);

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const apiKeysPageStateCodec: PageStateCodec<ApiKeysQueryState, null> = {
  decodeQuery(searchParams) {
    const filter = searchParams.get("filter");
    const sort = searchParams.get("sort");

    return {
      q: searchParams.get("q") ?? "",
      filter:
        filter && filterValues.has(filter as ApiKeysFilter)
          ? (filter as ApiKeysFilter)
          : "all",
      sort:
        sort && sortValues.has(sort as ApiKeysSort)
          ? (sort as ApiKeysSort)
          : "label_asc",
      selected: nonEmpty(searchParams.get("selected")),
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.q.trim()) searchParams.set("q", query.q);
    if (query.filter !== "all") searchParams.set("filter", query.filter);
    if (query.sort !== "label_asc") searchParams.set("sort", query.sort);
    if (query.selected) searchParams.set("selected", query.selected);

    return searchParams;
  },
  decodeEntry() {
    return null;
  },
  encodeEntry() {
    return null;
  },
};
