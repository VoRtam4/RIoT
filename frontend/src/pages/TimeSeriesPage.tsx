import { useMemo, useState } from "react";
import { useApolloClient } from "@apollo/client/react";
import { useSearchParams } from "react-router-dom";

import TimeSeriesFilters from "../modules/timeSeries/components/TimeSeriesFilters";
import TimeSeriesTable from "../modules/timeSeries/components/TimeSeriesTable";

import { useTimeSeries } from "../modules/timeSeries/hooks/useTimeSeries";

import {
  StartTimeSeriesExportDocument,
  type TimeSeriesReadInput,
} from "../generated/graphql";

export default function TimeSeriesPage() {
  const [searchParams] = useSearchParams();
  const client = useApolloClient();

  const initialInput: TimeSeriesReadInput | null = useMemo(() => {
    const type = searchParams.get("type") as "raw" | "kpi" | null;

    if (!type) return null;

    const sdInstanceID = searchParams.get("sdInstanceID");
    const sdTypeID = searchParams.get("sdTypeID");
    const kpiDefinitionID = searchParams.get("kpiDefinitionID");

    return {
      type,
      sortDesc: true,
      limit: 50,

      sdTypeID: sdTypeID ?? undefined,
      sdInstanceIDs: sdInstanceID ? [sdInstanceID] : [],

      kpiDefinitionIDs:
        type === "kpi" && kpiDefinitionID
          ? [kpiDefinitionID]
          : undefined,
    };
  }, [searchParams]);

  const autoRun = searchParams.get("auto") === "1";

  const [input, setInput] = useState<TimeSeriesReadInput | null>(
    autoRun ? initialInput : null,
  );

  const [submitted, setSubmitted] = useState(autoRun);

  const query = useTimeSeries(
    input as TimeSeriesReadInput,
    submitted && !!input,
  );

  const handleExport = async (input: TimeSeriesReadInput) => {
    const res = await client.mutate({
      mutation: StartTimeSeriesExportDocument,
      variables: { input },
    });

    const id = res.data?.startTimeSeriesExport;

    if (!id) {
      throw new Error("Export failed");
    }

    window.open(`rest/time-series/export/${id}`, "_blank");
  };

  return (
    <div className="container mt-3">
      <TimeSeriesFilters
        initialInput={initialInput}
        onSubmit={(v) => {
          setInput(v);
          setSubmitted(true);
        }}
        onExport={handleExport}
      />

      {submitted && input && <TimeSeriesTable query={query} />}
    </div>
  );
}