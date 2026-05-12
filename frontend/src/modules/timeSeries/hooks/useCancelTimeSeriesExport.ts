import { useMutation } from "@apollo/client/react";
import {
  CancelTimeSeriesExportDocument,
  type CancelTimeSeriesExportMutation,
  type CancelTimeSeriesExportMutationVariables,
} from "../../../generated/graphql";

export const useCancelTimeSeriesExport = () => {
  const [cancelExportMutation, { loading, error }] = useMutation<
    CancelTimeSeriesExportMutation,
    CancelTimeSeriesExportMutationVariables
  >(CancelTimeSeriesExportDocument);

  const cancelTimeSeriesExport = async (
    id: CancelTimeSeriesExportMutationVariables["id"],
  ) => {
    const res = await cancelExportMutation({
      variables: { id },
    });

    return res.data?.cancelTimeSeriesExport ?? null;
  };

  return {
    cancelTimeSeriesExport,
    loading,
    error,
  };
};
