import { useQuery } from "@apollo/client/react";
import {
  KpiDefinitionsBySdTypeDocument,
  type KpiDefinitionsBySdTypeQuery,
  type KpiDefinitionsBySdTypeQueryVariables,
} from "../../../generated/graphql";

export const useKpiDefinitionsBySdType = (sdTypeID: string | null) => {
  const { data, loading, error, refetch } = useQuery<
    KpiDefinitionsBySdTypeQuery,
    KpiDefinitionsBySdTypeQueryVariables
  >(KpiDefinitionsBySdTypeDocument, {
    skip: !sdTypeID,
    variables: {
      id: sdTypeID ?? "",
    },
    fetchPolicy: "cache-and-network",
  });

  const kpiDefinitions = data?.kpiDefinitionsBySdType ?? [];

  return {
    kpiDefinitions,
    loading,
    error,
    refetch,
  };
};
