import { useApolloClient } from "@apollo/client/react";
import {
  TimeSeriesExportDocument,
  type TimeSeriesExportQuery,
  type TimeSeriesExportQueryVariables,
} from "../../../generated/graphql";

export const useTimeSeriesExport = () => {
  const client = useApolloClient();

  const getTimeSeriesExport = async (
    id: TimeSeriesExportQueryVariables["id"],
  ) => {
    const res = await client.query<
      TimeSeriesExportQuery,
      TimeSeriesExportQueryVariables
    >({
      query: TimeSeriesExportDocument,
      variables: { id },
      fetchPolicy: "no-cache",
    });

    return res.data?.timeSeriesExport ?? null;
  };

  return {
    getTimeSeriesExport,
  };
};
