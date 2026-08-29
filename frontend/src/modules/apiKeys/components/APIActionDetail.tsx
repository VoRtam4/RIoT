/**
 * @file APIActionDetail.tsx
 * @brief Detail konkrétní API akce s popisem požadavku, odpovědi a oprávnění.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import type { ApiActionDoc, ApiFeatureDoc } from "../data/apiDocs";
import { getApiInterfaceDisplayUrl } from "../../../app/apiEndpoints";

import CodeSection from "./APICodeSection";
import MetaSection from "./APIMetaSection";

import { getAuthExample, getResponseExample } from "../utils/apiDocsExamples";
import { getRequiredPermissionLabel } from "../utils/apiDocsPermissions";

type Props = {
  feature: ApiFeatureDoc;
  action: ApiActionDoc;
  variant: ApiActionDoc["variants"][number];
};

export default function ActionDetail({ feature, action, variant }: Props) {
  const interfaceUrl = getApiInterfaceDisplayUrl(variant.technology);
  const requiredPermissionLabel = getRequiredPermissionLabel(
    feature.id,
    action.id,
  );

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

      <div className="api-docs-meta-grid mb-3">
        <MetaSection title="Interface URL" value={interfaceUrl} />

        <MetaSection
          title="Required Permission"
          value={requiredPermissionLabel}
        />
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
