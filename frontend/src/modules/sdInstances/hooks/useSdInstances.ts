/**
 * @file useSdInstances.ts
 * @brief Hook pro načítání sledovaných instancí dostupných aktuálnímu uživateli.
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
  SdInstancesDocument,
  type SdInstancesQuery,
} from "../../../generated/graphql";

export const useSdInstances = () => {
  const { data, loading, error, refetch } = useQuery<SdInstancesQuery>(
    SdInstancesDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const sdInstances = data?.sdInstances ?? [];

  return {
    sdInstances,
    loading,
    error,
    refetch,
  };
};
