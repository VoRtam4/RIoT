import { useQuery } from "@apollo/client/react";
import {
  SdInstancesByTypeDocument,
  type SdInstancesByTypeQuery,
  type SdInstancesByTypeQueryVariables,
} from "../../../generated/graphql";

export const useSdInstancesByType = (sdTypeID: string | null) => {
  const { data, loading, error } = useQuery<
    SdInstancesByTypeQuery,
    SdInstancesByTypeQueryVariables
  >(SdInstancesByTypeDocument, {
    skip: !sdTypeID,
    variables: {
      id: sdTypeID ?? "",
    },
    fetchPolicy: "cache-and-network",
  });

  const sdInstances = data?.sdInstancesByType ?? [];

  return {
    sdInstances,
    loading,
    error,
  };
};
