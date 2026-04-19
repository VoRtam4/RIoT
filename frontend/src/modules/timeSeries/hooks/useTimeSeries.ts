import { useInfiniteQuery } from "@tanstack/react-query";
import { useApolloClient } from "@apollo/client/react";
import {
  TimeSeriesReadDocument,
  type TimeSeriesCursorInput,
  type TimeSeriesReadQuery,
  type TimeSeriesReadQueryVariables,
  type TimeSeriesReadResponse,
  type TimeSeriesReadInput,
} from "../../../generated/graphql";

export const useTimeSeries = (input: TimeSeriesReadInput, enabled: boolean) => {
  const client = useApolloClient();

  const query = useInfiniteQuery<
    TimeSeriesReadResponse,
    Error,
    { pages: TimeSeriesReadResponse[] },
    [string, TimeSeriesReadInput],
    TimeSeriesCursorInput | null
  >({
    queryKey: ["timeSeries", input],
    enabled,
    initialPageParam: null,
    refetchOnWindowFocus: false,

    queryFn: async ({ pageParam }) => {
      const cursor: TimeSeriesCursorInput | undefined = pageParam
        ? {
            time: pageParam.time,
            sdInstanceUID: pageParam.sdInstanceUID,
            ...(pageParam.kpiDefinitionID
              ? { kpiDefinitionID: pageParam.kpiDefinitionID }
              : {}),
          }
        : undefined;

      const request: TimeSeriesReadInput = {
        ...input,
        cursor,
      };

      const res = await client.query<
        TimeSeriesReadQuery,
        TimeSeriesReadQueryVariables
      >({
        query: TimeSeriesReadDocument,
        variables: { request },
        fetchPolicy: "no-cache",
      });

      if (!res.data?.timeSeriesRead) {
        throw new Error("No data");
      }

      return res.data.timeSeriesRead;
    },

    getNextPageParam: (lastPage) =>
      lastPage.hasMoreData ? (lastPage.nextCursor ?? null) : undefined,
  });

  const data = query.data?.pages.flatMap((p) => p.data) ?? [];
  const parameters = query.data?.pages[0]?.parameters ?? [];
  const base = query.data?.pages[0]?.base ?? null;

  return {
    data,
    parameters,
    base,
    loading: query.isLoading,
    error: query.error,
    fetchNextPage: query.fetchNextPage,
    hasNextPage: query.hasNextPage,
    isFetchingNextPage: query.isFetchingNextPage,
    refetch: query.refetch,
  };
};
