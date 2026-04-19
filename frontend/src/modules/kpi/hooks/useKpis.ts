import { useQuery } from "@apollo/client/react";
import {
  KpiDefinitionsDocument,
  type KpiDefinitionsQuery,
} from "../../../generated/graphql";

export const useKpis = () => {
  const { data, loading, error, refetch } = useQuery<KpiDefinitionsQuery>(
    KpiDefinitionsDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const kpis = data?.kpiDefinitions ?? [];

  return {
    kpis,
    loading,
    error,
    refetch,
  };
};
