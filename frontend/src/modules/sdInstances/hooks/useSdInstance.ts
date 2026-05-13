/**
 * @file useSdInstance.ts
 * @brief Hook pro načtení detailu jedné sledované instance.
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
  SdInstanceDocument,
  type SdInstanceQuery,
  type SdInstanceQueryVariables,
} from "../../../generated/graphql";

export const useSdInstance = (id?: string | null) => {
  const { data, loading, error } = useQuery<
    SdInstanceQuery,
    SdInstanceQueryVariables
  >(SdInstanceDocument, {
    skip: !id,
    variables: { id: id ?? "" },
    fetchPolicy: "cache-and-network",
  });

  return {
    sdInstance: data?.sdInstance ?? null,
    loading,
    error,
  };
};