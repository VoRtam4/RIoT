/**
 * @file useSdInstancesByKpiDefinition.ts
 * @brief Hook pro načítání instancí navázaných na konkrétní KPI definici.
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
  SdInstancesByKpiDefinitionDocument,
  type SdInstancesByKpiDefinitionQuery,
  type SdInstancesByKpiDefinitionQueryVariables,
} from "../../../generated/graphql";

export const useSdInstancesByKpiDefinition = (
  kpiDefinitionID: string | null,
  enabled: boolean,
) => {
  const { data, loading, error, refetch } = useQuery<
    SdInstancesByKpiDefinitionQuery,
    SdInstancesByKpiDefinitionQueryVariables
  >(SdInstancesByKpiDefinitionDocument, {
    skip: !enabled || !kpiDefinitionID,
    variables: {
      id: kpiDefinitionID ?? "",
    },
    fetchPolicy: "cache-and-network",
  });

  const sdInstances = data?.sdInstancesByKpiDefinition ?? [];

  return {
    sdInstances,
    loading,
    error,
    refetch,
  };
};