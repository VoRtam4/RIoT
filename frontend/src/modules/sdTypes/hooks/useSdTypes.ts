import { useMemo } from "react";
import { useQuery } from "@apollo/client/react";
import { SdTypesDocument, type SdTypesQuery } from "../../../generated/graphql";

export const useSdTypes = () => {
  const { data, loading, error, refetch } = useQuery<SdTypesQuery>(
    SdTypesDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const sdTypes = data?.sdTypes ?? [];

  const sdTypeMap = useMemo(() => {
    return new Map(
      sdTypes.map((t) => [String(t.id), t.label ?? t.uid ?? String(t.id)]),
    );
  }, [sdTypes]);

  return {
    sdTypes,
    sdTypeMap,
    loading,
    error,
    refetch,
  };
};
