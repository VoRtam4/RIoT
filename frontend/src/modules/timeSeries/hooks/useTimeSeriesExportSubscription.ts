/**
 * @file useTimeSeriesExportSubscription.ts
 * @brief Subscription hook pro živé aktualizace stavu exportů časových řad.
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
  OnTimeSeriesExportUpdatedDocument,
  type TimeSeriesExport,
} from "../../../generated/graphql";

export const useTimeSeriesExportSubscription = (id?: string | null) => {
  const { data } = useSubscription(OnTimeSeriesExportUpdatedDocument, {
    variables: {
      filter: {
        ids: id ? [id] : undefined,
      },
    },
    skip: !id,
  });

  const latest: TimeSeriesExport | undefined = data?.onTimeSeriesExportUpdated;

  return {
    latest,
  };
};
