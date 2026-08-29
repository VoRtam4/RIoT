/**
 * @file useTimeSeriesAggregateKpi.ts
 * @brief Hook pro agregované čtení historických KPI výsledků.
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
  TimeSeriesReadAggregateKpiDocument,
  type TimeSeriesReadAggregateKpiQuery,
  type TimeSeriesReadAggregateKpiQueryVariables,
} from "../../../generated/graphql";

type Params = {
  from: string;
  to: string;

  aggregateSeconds: number;

  kpiDefinitionUIDs?: string[];
  sdInstanceUIDs?: string[];
  sdTypeUID?: string;
};

export const useTimeSeriesAggregateKpi = (params: Params, enabled: boolean) => {
  const { data, loading, error } = useQuery<
    TimeSeriesReadAggregateKpiQuery,
    TimeSeriesReadAggregateKpiQueryVariables
  >(TimeSeriesReadAggregateKpiDocument, {
    skip: !enabled,
    variables: {
      request: {
        from: params.from,
        to: params.to,
        aggregateSeconds: params.aggregateSeconds,

        kpiDefinitionUIDs: params.kpiDefinitionUIDs,
        sdInstanceUIDs: params.sdInstanceUIDs,
        sdTypeUID: params.sdTypeUID,
      },
    },
    fetchPolicy: "cache-and-network",
  });

  return {
    data: data?.timeSeriesReadAggregateKPI?.data ?? [],
    loading,
    error,
  };
};
