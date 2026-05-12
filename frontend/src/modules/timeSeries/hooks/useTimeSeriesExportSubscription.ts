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
