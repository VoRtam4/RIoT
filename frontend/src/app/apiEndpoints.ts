import type { ApiTechnology } from "../modules/apiKeys/data/apiDocs";

const defaultApiOrigin = import.meta.env.VITE_API_ORIGIN?.trim() || "";

function trimTrailingSlash(value: string) {
  return value.replace(/\/+$/, "");
}

function joinUrl(base: string, path: string) {
  return `${trimTrailingSlash(base)}${path.startsWith("/") ? path : `/${path}`}`;
}

function toHttpUrl(explicitUrl: string | undefined, fallbackPath: string) {
  const url = explicitUrl?.trim();
  if (url) {
    return url;
  }

  if (defaultApiOrigin) {
    return joinUrl(defaultApiOrigin, fallbackPath);
  }

  return fallbackPath;
}

function toWebSocketUrl(explicitUrl: string | undefined, fallbackPath: string) {
  const url = explicitUrl?.trim();
  if (url) {
    return url;
  }

  if (defaultApiOrigin) {
    return joinUrl(
      defaultApiOrigin.replace(/^http:/i, "ws:").replace(/^https:/i, "wss:"),
      fallbackPath,
    );
  }

  if (typeof window !== "undefined") {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    return `${protocol}//${window.location.host}${fallbackPath}`;
  }

  return fallbackPath;
}

export const apiEndpoints = {
  graphqlHttp: toHttpUrl(import.meta.env.VITE_GRAPHQL_URL, "/graphql"),
  graphqlWebSocket: toWebSocketUrl(
    import.meta.env.VITE_GRAPHQL_WS_URL,
    "/graphql",
  ),
  rest: toHttpUrl(import.meta.env.VITE_REST_URL, "/rest"),
  websocket: toWebSocketUrl(import.meta.env.VITE_WS_URL, "/ws"),
};

function toDisplayUrl(explicitUrl: string | undefined, fallbackPath: string) {
  const url = explicitUrl?.trim();
  if (url) {
    return url;
  }

  if (defaultApiOrigin) {
    return joinUrl(defaultApiOrigin, fallbackPath);
  }

  if (typeof window !== "undefined") {
    return joinUrl(window.location.origin, fallbackPath);
  }

  return fallbackPath;
}

export function getApiInterfaceUrl(technology: ApiTechnology) {
  switch (technology) {
    case "graphql":
      return apiEndpoints.graphqlHttp;
    case "rest":
      return apiEndpoints.rest;
    case "websocket":
      return apiEndpoints.websocket;
  }
}

export function getApiInterfaceDisplayUrl(technology: ApiTechnology) {
  switch (technology) {
    case "graphql":
      return toDisplayUrl(import.meta.env.VITE_GRAPHQL_URL, "/graphql");
    case "rest":
      return toDisplayUrl(import.meta.env.VITE_REST_URL, "/rest");
    case "websocket":
      return toDisplayUrl(import.meta.env.VITE_WS_URL, "/ws");
  }
}

export function getGraphQLSubscriptionUrl() {
  return apiEndpoints.graphqlWebSocket;
}

export function getGraphQLSubscriptionDisplayUrl() {
  const explicitUrl = import.meta.env.VITE_GRAPHQL_WS_URL?.trim();
  if (explicitUrl) {
    return explicitUrl;
  }

  if (defaultApiOrigin) {
    return joinUrl(defaultApiOrigin, "/graphql");
  }

  if (typeof window !== "undefined") {
    return joinUrl(window.location.origin, "/graphql");
  }

  return "/graphql";
}

export function formatExampleUrl(url: string, suffix = "") {
  return `${trimTrailingSlash(url)}${suffix}`;
}
