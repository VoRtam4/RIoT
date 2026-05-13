/**
 * @file useRawDataPoint.ts
 * @brief Hook pro čtení aktuálního raw datového bodu sledované instance.
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
  RawDataPointDocument,
  type RawDataPointQuery,
  type RawDataPointQueryVariables,
} from "../../../generated/graphql";

export const useRawDataPoint = (sdInstanceID?: string | null) => {
  const { data, loading, error, refetch } = useQuery<
    RawDataPointQuery,
    RawDataPointQueryVariables
  >(RawDataPointDocument, {
    variables: {
      id: sdInstanceID as string,
    },
    skip: sdInstanceID == null,
    fetchPolicy: "cache-and-network",
  });

  return {
    rawDataPoint: data?.rawDataPoint ?? null,
    loading,
    error,
    refetch,
  };
};
