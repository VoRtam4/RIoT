import { useQuery } from "@apollo/client/react";
import {
  KpiDefinitionsBySdInstanceDocument,
  type KpiDefinitionsBySdInstanceQuery,
  type KpiDefinitionsBySdInstanceQueryVariables,
} from "../../../generated/graphql";

export const useKpiDefinitionsBySdInstance = (sdInstanceID: string | null) => {
  const { data, loading, error, refetch } = useQuery<
    KpiDefinitionsBySdInstanceQuery,
    KpiDefinitionsBySdInstanceQueryVariables
  >(KpiDefinitionsBySdInstanceDocument, {
    skip: !sdInstanceID,
    variables: {
      id: sdInstanceID ?? "",
    },
    fetchPolicy: "cache-and-network",
    notifyOnNetworkStatusChange: true,
  });

  const kpiDefinitions = data?.kpiDefinitionsBySdInstance ?? [];

  return {
    kpiDefinitions,
    loading,
    error,
    refetch,
  };
};
