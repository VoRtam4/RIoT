/**
 * @file useRawDataSubscription.ts
 * @brief Subscription hook pro živý příjem raw datových bodů sledovaných instancí.
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
  OnRawDataPointArrivedDocument,
  type RawDataPoint,
} from "../../../generated/graphql";

export const useRawDataSubscription = (
  sdTypeID?: string,
  sdInstanceID?: string | null,
) => {
  const { data } = useSubscription(OnRawDataPointArrivedDocument, {
    variables: {
      filter: {
        sdTypeIDs: sdTypeID ? [sdTypeID] : undefined,
        sdInstanceIDs: sdInstanceID ? [sdInstanceID] : undefined,
      },
    },
    skip: !sdTypeID || !sdInstanceID,
  });

  const latest: RawDataPoint | undefined = data?.onRawDataPointArrived?.[0];

  return {
    latest,
  };
};
