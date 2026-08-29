/**
 * @file APIDocsSidebar.tsx
 * @brief Postranní navigace kategorií a akcí ve frontendové API dokumentaci.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { ApiFeatureDoc } from "../data/apiDocs";
import ApiDocsFeatureCard from "./APIDocsFeatureCard";
import VirtualizedList from "../../../components/virtualization/VirtualizedList";
import EmptyStateNotice from "../../../components/EmptyStateNotice";
import UrlSyncedSearchInput from "../../../components/UrlSyncedSearchInput";
import type { SearchMode } from "../../../utils/reactSelectSearch";

type Props = {
  features: ApiFeatureDoc[];
  selectedFeatureId: string | null;
  search: string;
  searchMode: SearchMode;
  onSearchChange: (value: string) => void;
  onSearchModeChange: (value: SearchMode) => void;
  onSelectFeature: (id: string | null) => void;
};

export default function ApiDocsSidebar({
  features,
  selectedFeatureId,
  search,
  searchMode,
  onSearchChange,
  onSearchModeChange,
  onSelectFeature,
}: Props) {
  const selectedIndex = features.findIndex(
    (feature) => feature.id === selectedFeatureId,
  );

  return (
    <div
      className="card p-3 d-flex flex-column"
      style={{ height: "100%", minHeight: 0 }}
    >
      <div className="mb-3">
        <label className="form-label">Search Features</label>
        <UrlSyncedSearchInput
          placeholder="For example KPI, API keys, history..."
          value={search}
          onChange={onSearchChange}
          mode={searchMode}
          onModeChange={onSearchModeChange}
        />
      </div>

      <VirtualizedList
        items={features}
        rowHeight={100}
        itemSpacing={5}
        scrollToIndex={selectedIndex}
        pinnedIndex={selectedIndex}
        style={{ flex: 1 }}
        emptyState={
          <EmptyStateNotice
            compact
            title="No matching features"
            description="Try a different search phrase."
          />
        }
        renderItem={(feature) => (
          <ApiDocsFeatureCard
            key={feature.id}
            feature={feature}
            selected={feature.id === selectedFeatureId}
            onClick={() =>
              onSelectFeature(
                feature.id === selectedFeatureId ? null : feature.id,
              )
            }
          />
        )}
      />
    </div>
  );
}
