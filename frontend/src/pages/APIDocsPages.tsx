/**
 * @file APIDocsPages.tsx
 * @brief Stránka interaktivní dokumentace GraphQL, REST a WebSocket API pro API klíče.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useMemo } from "react";

import { apiFeatures } from "../modules/apiKeys/data/apiDocs";
import type {
  ApiFeatureDoc,
  ApiTechnology,
} from "../modules/apiKeys/data/apiDocs";
import { matchesSearchText } from "../utils/reactSelectSearch";

import ApiDocsSidebar from "../modules/apiKeys/components/APIDocsSidebar";
import FeatureDetail from "../modules/apiKeys/components/APIFeatureDetail";
import { usePageState } from "../app/navigation/usePageState";
import { apiDocsPageStateCodec } from "../modules/apiKeys/state/apiDocsPageState";
import EmptyStateNotice from "../components/EmptyStateNotice";

export default function ApiDocsPage() {
  const { query, setPageState } = usePageState(apiDocsPageStateCodec);

  const filteredFeatures = useMemo(() => {
    if (!query.q.trim()) {
      return apiFeatures;
    }

    return apiFeatures.filter((feature: ApiFeatureDoc) => {
      return matchesSearchText(
        query.q,
        query.qMode,
        feature.title,
        feature.summary,
        ...feature.actions.map((action) => action.title),
        ...feature.actions.flatMap((action) =>
          action.variants.map((variant) => variant.label),
        ),
      );
    });
  }, [query.q, query.qMode]);

  const selectedFeature =
    filteredFeatures.find((feature) => feature.id === query.feature) ?? null;

  const featureTechnologies = useMemo(() => {
    if (!selectedFeature) return [];

    return (["graphql", "rest", "websocket"] as ApiTechnology[]).filter(
      (technology) =>
        selectedFeature.actions.some((action) =>
          action.variants.some((variant) => variant.technology === technology),
        ),
    );
  }, [selectedFeature]);

  const selectedTechnology =
    query.tech && featureTechnologies.includes(query.tech)
      ? query.tech
      : featureTechnologies[0];

  const actionsForTechnology = useMemo(() => {
    if (!selectedFeature || !selectedTechnology) return [];

    return selectedFeature.actions.filter((action) =>
      action.variants.some(
        (variant) => variant.technology === selectedTechnology,
      ),
    );
  }, [selectedFeature, selectedTechnology]);

  const selectedAction =
    actionsForTechnology.find((action) => action.id === query.op) ??
    actionsForTechnology[0] ??
    null;

  useEffect(() => {
    const nextFeature = selectedFeature?.id ?? null;
    const nextTech = selectedFeature ? (selectedTechnology ?? null) : null;
    const nextOp = selectedFeature ? (selectedAction?.id ?? null) : null;

    if (
      nextFeature === query.feature &&
      nextTech === query.tech &&
      nextOp === query.op
    ) {
      return;
    }

    setPageState(
      {
        query: {
          ...query,
          feature: nextFeature,
          tech: nextTech,
          op: nextOp,
        },
        entry: null,
      },
      { replace: true },
    );
  }, [
    query,
    selectedAction,
    selectedFeature,
    selectedTechnology,
    setPageState,
  ]);

  return (
    <div className="container-fluid mt-3">
      <div className="api-docs-layout">
        <div
          style={{
            position: "sticky",
            top: 16,
            height: "calc(100vh - 110px)",
          }}
        >
          <ApiDocsSidebar
            features={filteredFeatures}
            selectedFeatureId={selectedFeature?.id ?? null}
            search={query.q}
            searchMode={query.qMode}
            onSearchChange={(q) =>
              setPageState(
                {
                  query: {
                    ...query,
                    q,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
            onSearchModeChange={(qMode) =>
              setPageState(
                {
                  query: {
                    ...query,
                    qMode,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
            onSelectFeature={(featureId) =>
              setPageState(
                {
                  query: {
                    ...query,
                    feature: featureId,
                    tech: null,
                    op: null,
                  },
                  entry: null,
                },
                { replace: true },
              )
            }
          />
        </div>

        <section className="api-docs-content">
          {selectedFeature && selectedTechnology && selectedAction ? (
            <FeatureDetail
              key={selectedFeature.id}
              feature={selectedFeature}
              selectedTechnology={selectedTechnology}
              selectedActionId={selectedAction.id}
              onTechnologyChange={(tech) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      tech,
                      op:
                        selectedFeature.actions.find((action) =>
                          action.variants.some(
                            (variant) => variant.technology === tech,
                          ),
                        )?.id ?? null,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onActionChange={(op) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      op,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
            />
          ) : (
            <div className="panel p-4">
              <EmptyStateNotice
                title="Select a feature"
                description="Choose a feature from the sidebar to view its operations, interface variants and examples."
              />
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
