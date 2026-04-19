import { useQuery } from "@apollo/client/react";
import { ApiKeysDocument, type ApiKeysQuery } from "../../../generated/graphql";

export const useApiKeys = () => {
  const { data, loading, error, refetch } = useQuery<ApiKeysQuery>(
    ApiKeysDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const apiKeys = data?.apiKeys ?? [];

  return {
    apiKeys,
    loading,
    error,
    refetch,
  };
};
