import type { PageStateCodec } from "../../../app/navigation/usePageState";

export type ApiDocsQueryState = {
  q: string;
  feature: string | null;
  op: string | null;
  tech: "graphql" | "rest" | "websocket" | null;
};

function nonEmpty(value: string | null) {
  return value && value.trim() ? value : null;
}

export const apiDocsPageStateCodec: PageStateCodec<ApiDocsQueryState, null> = {
  decodeQuery(searchParams) {
    const tech = searchParams.get("tech");

    return {
      q: searchParams.get("q") ?? "",
      feature: nonEmpty(searchParams.get("feature")),
      op: nonEmpty(searchParams.get("op")),
      tech:
        tech === "graphql" || tech === "rest" || tech === "websocket"
          ? tech
          : null,
    };
  },
  encodeQuery(query) {
    const searchParams = new URLSearchParams();

    if (query.q.trim()) searchParams.set("q", query.q);
    if (query.feature) searchParams.set("feature", query.feature);
    if (query.op) searchParams.set("op", query.op);
    if (query.tech) searchParams.set("tech", query.tech);

    return searchParams;
  },
  decodeEntry() {
    return null;
  },
  encodeEntry() {
    return null;
  },
};
