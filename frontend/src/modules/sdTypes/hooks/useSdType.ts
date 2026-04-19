import { useQuery } from "@apollo/client/react";
import {
  SdTypeByIdDocument,
  type SdTypeByIdQuery,
  type SdTypeByIdQueryVariables,
} from "../../../generated/graphql";

export const useSdType = (id?: string) => {
  const { data, loading, error, refetch } = useQuery<
    SdTypeByIdQuery,
    SdTypeByIdQueryVariables
  >(SdTypeByIdDocument, {
    variables: {
      id: id ?? "",
    },
    skip: !id,
    fetchPolicy: "cache-and-network",
  });

  return {
    sdType: data?.sdType ?? null,
    loading,
    error,
    refetch,
  };
};
