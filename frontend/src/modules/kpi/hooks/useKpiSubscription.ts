/**
 * @file useKpiSubscription.ts
 * @brief Subscription hook pro příjem živých změn výsledků KPI.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useSubscription } from "@apollo/client/react";
import {
  OnKpiFulfillmentCheckedDocument,
  type KpiFulfillmentCheckResult,
} from "../../../generated/graphql";

export const useKpiSubscription = (
  kpiDefinitionUID?: string,
  sdInstanceUID?: string | null,
) => {
  const { data } = useSubscription(OnKpiFulfillmentCheckedDocument, {
    variables: {
      filter: {
        kpiDefinitionUIDs: kpiDefinitionUID ? [kpiDefinitionUID] : undefined,
        sdInstanceUIDs: sdInstanceUID ? [sdInstanceUID] : undefined,
      },
    },
    skip: !kpiDefinitionUID || !sdInstanceUID,
  });

  const latest: KpiFulfillmentCheckResult | undefined =
    data?.onKPIFulfillmentChecked?.[0];

  return {
    latest,
  };
};
