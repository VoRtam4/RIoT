import { useMemo, useState } from "react";

import { apiFeatures } from "../modules/apiKeys/data/apiDocs";
import type { ApiFeatureDoc } from "../modules/apiKeys/data/apiDocs";
import { matchesSearchText } from "../utils/reactSelectSearch";

import ApiDocsSidebar from "../modules/apiKeys/components/APIDocsSidebar";
import FeatureDetail from "../modules/apiKeys/components/APIFeatureDetail";

export default function ApiDocsPage() {
  const [selectedFeatureId, setSelectedFeatureId] = useState(apiFeatures[0]?.id);
  const [search, setSearch] = useState("");

  const filteredFeatures = useMemo(() => {
    if (!search.trim()) {
      return apiFeatures;
    }

    return apiFeatures.filter((feature: ApiFeatureDoc) => {
      return matchesSearchText(
        search,
        feature.title,
        feature.summary,
        ...feature.actions.map((action) => action.title),
        ...feature.actions.flatMap((action) =>
          action.variants.map((variant) => variant.label),
        ),
      );
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
