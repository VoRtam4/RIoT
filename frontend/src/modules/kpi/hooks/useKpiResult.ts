import { useQuery } from "@apollo/client/react";
import {
  KpiResultDocument,
  type KpiResultQuery,
  type KpiResultQueryVariables,
} from "../../../generated/graphql";

export const useKpiResult = (
  kpiDefinitionID?: string,
  sdInstanceID?: string,
) => {
  const { data, loading, error, refetch } = useQuery<
    KpiResultQuery,
    KpiResultQueryVariables
  >(KpiResultDocument, {
    variables: {
      request: {
        kpiDefinitionID: kpiDefinitionID ?? "",
        sdInstanceID: sdInstanceID ?? "",
      },
    },
    skip: !kpiDefinitionID || !sdInstanceID,
    fetchPolicy: "cache-and-network",
  });

  return {
    result: data?.kpiResult ?? null,
    loading,
    error,
    refetch,
  };
};
