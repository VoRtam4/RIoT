/**
 * @file useKpiDefinition.ts
 * @brief Hook pro načtení jedné KPI definice a jejích detailů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useQuery } from "@apollo/client/react";
import {
  KpiDefinitionDocument,
  type KpiDefinitionQuery,
  type KpiDefinitionQueryVariables,
} from "../../../generated/graphql";

export const useKpiDefinition = (id?: string) => {
  const { data, loading, error, refetch } = useQuery<
    KpiDefinitionQuery,
    KpiDefinitionQueryVariables
  >(KpiDefinitionDocument, {
    variables: {
      id: id ?? "",
    },
    skip: !id,
    fetchPolicy: "cache-and-network",
  });

  return {
    kpi: data?.kpiDefinition ?? null,
    loading,
    error,
    refetch,
  };
};
