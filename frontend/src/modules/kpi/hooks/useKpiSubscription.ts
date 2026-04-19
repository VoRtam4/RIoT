import { useSubscription } from "@apollo/client/react";
import {
  OnKpiFulfillmentCheckedDocument,
  type KpiFulfillmentCheckResult,
} from "../../../generated/graphql";

export const useKpiSubscription = (
  kpiDefinitionID?: string,
  sdInstanceID?: string | null,
) => {
  const { data } = useSubscription(OnKpiFulfillmentCheckedDocument, {
    variables: {
      filter: {
        kpiDefinitions: kpiDefinitionID ? [kpiDefinitionID] : undefined,
        sdInstanceIDs: sdInstanceID ? [sdInstanceID] : undefined,
      },
    },
    skip: !kpiDefinitionID || !sdInstanceID,
  });

  const latest: KpiFulfillmentCheckResult | undefined =
    data?.onKPIFulfillmentChecked?.[0];

  return {
    latest,
  };
};
