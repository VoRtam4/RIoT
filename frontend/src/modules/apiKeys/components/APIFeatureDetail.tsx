import {
  technologyOrder,
  type ApiActionDoc,
  type ApiFeatureDoc,
  type ApiTechnology,
} from "../data/apiDocs";

import ActionDetail from "./APIActionDetail";

const technologyLabels: Record<ApiTechnology, string> = {
  graphql: "GraphQL",
  rest: "REST",
  websocket: "WebSocket",
};

type Props = {
  feature: ApiFeatureDoc;
  selectedTechnology: ApiTechnology;
  selectedActionId: string;
  onTechnologyChange: (technology: ApiTechnology) => void;
  onActionChange: (actionId: string) => void;
};

export default function FeatureDetail({
  feature,
  selectedTechnology,
  selectedActionId,
  onTechnologyChange,
  onActionChange,
}: Props) {
  const featureTechnologies = technologyOrder.filter((technology) =>
    feature.actions.some((action) =>
      action.variants.some((variant) => variant.technology === technology),
    ),
  );

  const actionsForTechnology = feature.actions.filter((action) =>
    action.variants.some((variant) => variant.technology === selectedTechnology),
  );

  const activeAction =
    actionsForTechnology.find((action) => action.id === selectedActionId) ??
    actionsForTechnology[0];

  const activeVariant = activeAction?.variants.find(
    (variant) => variant.technology === selectedTechnology,
  );

  const selectTechnology = (technology: ApiTechnology) => {
    onTechnologyChange(technology);
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
            {actionsForTechnology.map((action: ApiActionDoc) => (
              <button
                key={action.id}
                type="button"
                className={`api-docs-chip-button ${
                  activeAction?.id === action.id ? "is-active" : ""
                }`}
                onClick={() => onActionChange(action.id)}
              >
                {action.title}
              </button>
            ))}
          </div>
        </div>
      </div>

      {activeAction && activeVariant ? (
        <ActionDetail
          feature={feature}
          action={activeAction}
          variant={activeVariant}
        />
      ) : null}
    </div>
  );
}
