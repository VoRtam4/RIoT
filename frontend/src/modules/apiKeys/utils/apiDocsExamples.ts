/**
 * @file apiDocsExamples.ts
 * @brief Generování ukázkových GraphQL, REST a WebSocket požadavků pro API dokumentaci.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { ApiActionDoc, ApiTechnology } from "../data/apiDocs";
import {
  formatExampleUrl,
  getApiInterfaceDisplayUrl,
  getGraphQLSubscriptionDisplayUrl,
} from "../../../app/apiEndpoints";

export function getResponseExample(
  action: ApiActionDoc,
  variant: ApiActionDoc["variants"][number],
) {
  const title = action.title.toLowerCase();

  if (variant.technology === "graphql") {
    if (title.includes("list")) {
      return `{
  "data": {
    "<queryName>": [
      {
        "id": "1",
        "label": "Example item"
      }
    ]
  }
}`;
    }

    if (title.includes("subscribe")) {
      return `{
  "data": {
    "<subscriptionName>": {
      "id": "1",
      "eventTime": "2026-04-14T12:30:00Z",
      "payload": {}
    }
  }
}`;
    }

    if (
      title.includes("create") ||
      title.includes("update") ||
      title.includes("detail") ||
      title.includes("get ")
    ) {
      return `{
  "data": {
    "<operationName>": {
      "id": "1",
      "label": "Example item"
    }
  }
}`;
    }

    if (title.includes("delete")) {
      return `{
  "data": {
    "<operationName>": true
  }
}`;
    }
  }

  if (variant.technology === "websocket") {
    if (title.includes("subscribe")) {
      return `{
  "type": "event",
  "id": "stream-1",
  "topic": "<event_topic>",
  "payload": {
    "id": 1,
    "eventTime": "2026-04-14T12:30:00Z"
  }
}`;
    }

    return `{
  "type": "response",
  "id": "request-1",
  "success": true,
  "payload": ${
    title.includes("list")
      ? `[
    {
      "id": 1,
      "label": "Example item"
    }
  ]`
      : `{
    "id": 1,
    "label": "Example item"
  }`
  }
}`;
  }

  if (title.includes("subscribe")) {
    return `event: <event_topic>
data: {"id":1,"eventTime":"2026-04-14T12:30:00Z","payload":{}}`;
  }

  if (title.includes("delete")) {
    return `HTTP/1.1 200 OK

true`;
  }

  if (title.includes("list")) {
    return `[
  {
    "id": 1,
    "label": "Example item"
  }
]`;
  }

  return `{
  "id": 1,
  "label": "Example item"
}`;
}

export function getAuthExample(technology: ApiTechnology) {
  if (technology === "graphql") {
    return `HTTP queries and mutations

POST /graphql
Content-Type: application/json
X-API-Key: riot_live_abc123_secret_value

{
  "query": "query { apiKeys { id label } }"
}

JavaScript example

fetch("${getApiInterfaceDisplayUrl("graphql")}", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "X-API-Key": "riot_live_abc123_secret_value"
  },
  body: JSON.stringify({
    query: "query { apiKeys { id label } }"
  })
})

GraphQL subscriptions use ${getGraphQLSubscriptionDisplayUrl()} with the GraphQL WebSocket protocol.
For API key auth, send the key in the connection_init payload:

{
  "X-API-Key": "riot_live_abc123_secret_value"
}

Cookie-based login also works when the browser sends the authenticated session automatically.`;
  }

  if (technology === "websocket") {
    return `WebSocket API key authentication happens during the HTTP upgrade request.
Send the API key in the X-API-Key header of the handshake.

wscat example

wscat -c ws://<frontend-host>/ws -H "X-API-Key: riot_live_abc123_secret_value"

After the connection is established, send normal request or subscribe messages.

Important:
Browser WebSocket APIs do not let you set arbitrary handshake headers directly.
For browser clients, use session-cookie authentication or a proxy that injects X-API-Key.`;
  }

  return `REST and SSE requests accept API key authentication through the X-API-Key header.

curl example

curl -H "X-API-Key: riot_live_abc123_secret_value" \\
  ${formatExampleUrl(getApiInterfaceDisplayUrl("rest"), "/api-keys")}`;
}
