/**
 * @file useCancelTimeSeriesExport.ts
 * @brief Hook pro zrušení běžícího exportu historických dat.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
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
