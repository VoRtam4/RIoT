/**
 * @file useSdType.ts
 * @brief Hook pro načtení detailu jednoho sledovaného typu.
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
  SdTypeByUidDocument,
  type SdTypeByUidQuery,
  type SdTypeByUidQueryVariables,
} from "../../../generated/graphql";

export const useSdType = (uid?: string) => {
  const { data, loading, error, refetch } = useQuery<
    SdTypeByUidQuery,
    SdTypeByUidQueryVariables
  >(SdTypeByUidDocument, {
    variables: {
      uid: uid ?? "",
    },
    skip: !uid,
    fetchPolicy: "cache-and-network",
  });

  return {
    sdType: data?.sdType ?? null,
    loading,
    error,
    refetch,
  };
};
