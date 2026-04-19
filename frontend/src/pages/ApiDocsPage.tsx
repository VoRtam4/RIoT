import { useMemo, useState } from "react";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import {
  apiFeatures,
  technologyOrder,
  type ApiActionDoc,
  type ApiFeatureDoc,
  type ApiTechnology,
} from "../modules/apiKeys/data/apiDocs";

const technologyLabels: Record<ApiTechnology, string> = {
  graphql: "GraphQL",
  rest: "REST",
  websocket: "WebSocket",
};

export default function ApiDocsPage() {
  const [selectedFeatureId, setSelectedFeatureId] = useState(apiFeatures[0]?.id);
  const [search, setSearch] = useState("");

  const filteredFeatures = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) {
      return apiFeatures;
    }

    return apiFeatures.filter((feature) => {
      const haystack = [
        feature.title,
        feature.summary,
        ...feature.actions.map((action) => action.title),
        ...feature.actions.flatMap((action) =>
          action.variants.map((variant) => variant.label),
        ),
      ]
        .join(" ")
        .toLowerCase();

      return haystack.includes(query);
    });
  }, [search]);

  const selectedFeature =
    filteredFeatures.find((feature) => feature.id === selectedFeatureId) ??
    filteredFeatures[0] ??
    null;

  return (
    <div className="container-fluid mt-3">
      <div className="api-docs-layout">
        <div
          style={{
            position: "sticky",
            top: 16,
            height: "calc(100vh - 110px)",
            overflow: "hidden",
            minHeight: 0,
          }}
        >
          <ApiDocsSidebar
            features={filteredFeatures}
            selectedFeatureId={selectedFeature?.id ?? null}
            search={search}
            onSearchChange={setSearch}
            onSelectFeature={setSelectedFeatureId}
          />
        </div>

        <section className="api-docs-content">
          {selectedFeature ? (
            <FeatureDetail key={selectedFeature.id} feature={selectedFeature} />
          ) : (
            <div className="panel p-4">
              <div className="form-label">Nothing to display.</div>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

function ApiDocsSidebar({
  features,
  selectedFeatureId,
  search,
  onSearchChange,
  onSelectFeature,
}: {
  features: ApiFeatureDoc[];
  selectedFeatureId: string | null;
  search: string;
  onSearchChange: (value: string) => void;
  onSelectFeature: (id: string) => void;
}) {
  return (
    <div className="card p-3 d-flex flex-column" style={{ height: "100%", minHeight: 0 }}>
      <div className="mb-3">
        <label className="form-label">Search Features</label>
        <input
          className="form-control"
          placeholder="For example KPI, API keys, history..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
        />
      </div>

      <Box
        sx={{
          overflowY: "auto",
          flex: 1,
          minHeight: 0,
          display: "flex",
          flexDirection: "column",
          gap: 1,
        }}
      >
        {features.map((feature) => (
          <ApiDocsFeatureCard
            key={feature.id}
            feature={feature}
            selected={feature.id === selectedFeatureId}
            onClick={() => onSelectFeature(feature.id)}
          />
        ))}

        {features.length === 0 && (
          <div className="form-label small">No feature matches the filter.</div>
        )}
      </Box>
    </div>
  );
}

function ApiDocsFeatureCard({
  feature,
  selected,
  onClick,
}: {
  feature: ApiFeatureDoc;
  selected: boolean;
  onClick: () => void;
}) {
  return (
    <Card
      sx={{
        flexShrink: 0,
        backgroundColor: selected ? "var(--primary)" : "var(--bg-card)",
        transition: "all 0.15s",
        border: selected
          ? "1px solid var(--bs-primary)"
          : "1px solid transparent",
        "&:hover": {
          transform: "translateY(-1px)",
          boxShadow: 3,
        },
      }}
    >
      <CardActionArea
        onClick={onClick}
        sx={{
          backgroundColor: selected ? "rgba(0,123,255,0.1)" : "transparent",
          "&:hover": {
            backgroundColor: "rgba(0,123,255,0.2)",
          },
        }}
      >
        <CardContent sx={{ py: 1.5 }}>
          <Typography variant="body1" className="form-label" sx={{ mb: 0.5 }}>
            {feature.title}
          </Typography>
          <Typography
            variant="caption"
            className="form-label"
            sx={{ opacity: 0.8, display: "block" }}
          >
            {feature.summary}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}

function FeatureDetail({ feature }: { feature: ApiFeatureDoc }) {
  const featureTechnologies = technologyOrder.filter((technology) =>
    feature.actions.some((action) =>
      action.variants.some((variant) => variant.technology === technology),
    ),
  );

  const [selectedTechnology, setSelectedTechnology] = useState<ApiTechnology>(
    featureTechnologies[0],
  );

  const actionsForTechnology = feature.actions.filter((action) =>
    action.variants.some((variant) => variant.technology === selectedTechnology),
  );

  const [selectedActionId, setSelectedActionId] = useState<string>(
    actionsForTechnology[0]?.id ?? feature.actions[0].id,
  );

  const activeAction =
    actionsForTechnology.find((action) => action.id === selectedActionId) ??
    actionsForTechnology[0];

  const activeVariant = activeAction.variants.find(
    (variant) => variant.technology === selectedTechnology,
  );

  const selectTechnology = (technology: ApiTechnology) => {
    setSelectedTechnology(technology);
    const firstActionForTechnology = feature.actions.find((action) =>
      action.variants.some((variant) => variant.technology === technology),
    );
    if (firstActionForTechnology) {
      setSelectedActionId(firstActionForTechnology.id);
    }
  };

  return (
    <div className="d-flex flex-column gap-3">
      <div className="card p-4">
        <div className="d-flex justify-content-between align-items-start flex-wrap gap-3">
          <div>
            <h3 className="mb-2">{feature.title}</h3>
            <p className="form-label mb-3">{feature.summary}</p>
          </div>

          <div className="api-docs-badges">
            {featureTechnologies.map((technology) => (
              <button
                key={technology}
                type="button"
                className={`api-docs-chip-button ${
                  selectedTechnology === technology ? "is-active" : ""
                }`}
                onClick={() => selectTechnology(technology)}
              >
                {technologyLabels[technology]}
              </button>
            ))}
          </div>
        </div>

        <div>
          <div className="form-label mb-2">Actions</div>
          <div className="api-docs-actions">
            {actionsForTechnology.map((action) => (
              <button
                key={action.id}
                type="button"
                className={`api-docs-chip-button ${
                  activeAction?.id === action.id ? "is-active" : ""
                }`}
                onClick={() => setSelectedActionId(action.id)}
              >
                {action.title}
              </button>
            ))}
          </div>
        </div>
      </div>

      {activeAction && activeVariant ? (
        <ActionDetail action={activeAction} variant={activeVariant} />
      ) : null}
    </div>
  );
}

function ActionDetail({
  action,
  variant,
}: {
  action: ApiActionDoc;
  variant: ApiActionDoc["variants"][number];
}) {
  return (
    <div className="card p-4">
      <div className="d-flex justify-content-between align-items-start flex-wrap gap-3 mb-3">
        <div>
          <div className="form-label mb-1">{variant.label}</div>
          <h4 className="mb-2">{action.title}</h4>
          <p className="form-label mb-2">{action.summary}</p>
          <p className="form-label mb-0">{variant.summary}</p>
        </div>
      </div>

      <CodeSection
        title="Authentication"
        code={variant.authExample ?? getAuthExample(variant.technology)}
      />

      <CodeSection title="Request Example" code={variant.requestExample} />
      <CodeSection
        title="Expected Response"
        code={variant.responseExample ?? getResponseExample(action, variant)}
      />

      {variant.notes && variant.notes.length > 0 && (
        <div className="mt-3">
          <div className="form-label mb-2">Notes</div>
          <div className="api-docs-notes">
            {variant.notes.map((note) => (
              <div key={note} className="api-docs-note form-label">
                {note}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function getResponseExample(
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
    title.includes("list") ? `[
    {
      "id": 1,
      "label": "Example item"
    }
  ]` : `{
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

function getAuthExample(technology: ApiTechnology) {
  if (technology === "graphql") {
    return `HTTP queries and mutations

POST /graphql
Content-Type: application/json
X-API-Key: riot_live_abc123_secret_value

{
  "query": "query { apiKeys { id label } }"
}

JavaScript example

fetch("http://localhost:9090/graphql", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "X-API-Key": "riot_live_abc123_secret_value"
  },
  body: JSON.stringify({
    query: "query { apiKeys { id label } }"
  })
})

For GraphQL subscriptions, the backend authenticates during the WebSocket handshake.
Use a client that can send X-API-Key as a handshake header, or use cookie-based login.`;
  }

  if (technology === "websocket") {
    return `WebSocket API key authentication happens during the HTTP upgrade request.
Send the API key in the X-API-Key header of the handshake.

wscat example

wscat -c ws://localhost:9090/ws -H "X-API-Key: riot_live_abc123_secret_value"

After the connection is established, send normal request or subscribe messages.

Important:
Browser WebSocket APIs do not let you set arbitrary handshake headers directly.
For browser clients, use session-cookie authentication or a proxy that injects X-API-Key.`;
  }

  return `REST and SSE requests accept API key authentication through the X-API-Key header.

curl example

curl -H "X-API-Key: riot_live_abc123_secret_value" \\
  http://localhost:9090/rest/api-keys`;
}

function CodeSection({ title, code }: { title: string; code: string }) {
  return (
    <div className="mb-3">
      <div className="form-label mb-2">{title}</div>
      <pre className="api-docs-code-block">
        <code>{code}</code>
      </pre>
    </div>
  );
}
