/**
 * @file useKpiResult.ts
 * @brief Hook pro čtení aktuálního nebo historického výsledku KPI.
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
  KpiResultDocument,
  type KpiResultQuery,
  type KpiResultQueryVariables,
} from "../../../generated/graphql";

export const useKpiResult = (
  kpiDefinitionUID?: string,
  sdInstanceUID?: string,
) => {
  const { data, loading, error, refetch } = useQuery<
    KpiResultQuery,
    KpiResultQueryVariables
  >(KpiResultDocument, {
    variables: {
      request: {
        kpiDefinitionUID: kpiDefinitionUID ?? "",
        sdInstanceUID: sdInstanceUID ?? "",
      },
    },
    skip: !kpiDefinitionUID || !sdInstanceUID,
    fetchPolicy: "cache-and-network",
  });

  return {
    result: data?.kpiResult ?? null,
    loading,
    error,
    refetch,
  };
};
