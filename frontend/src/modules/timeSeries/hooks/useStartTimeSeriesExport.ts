import { useMutation } from "@apollo/client/react";
import {
  StartTimeSeriesExportDocument,
  type StartTimeSeriesExportMutation,
  type StartTimeSeriesExportMutationVariables,
} from "../../../generated/graphql";

export const useStartTimeSeriesExport = () => {
  const [startExportMutation, { loading, error }] = useMutation<
    StartTimeSeriesExportMutation,
    StartTimeSeriesExportMutationVariables
  >(StartTimeSeriesExportDocument);

  const startTimeSeriesExport = async (
    input: StartTimeSeriesExportMutationVariables["input"],
  ) => {
    const res = await startExportMutation({
      variables: { input },
    });

    return res.data?.startTimeSeriesExport ?? null;
  };

  return {
    startTimeSeriesExport,
    loading,
    error,
  };
};
