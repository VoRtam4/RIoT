/**
 * @file TimeSeriesPage.tsx
 * @brief Stránka pro dotazování, agregaci a export historických časových řad.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useMemo, useRef, useState } from "react";
import toast from "react-hot-toast";

import TimeSeriesFilters from "../modules/timeSeries/components/TimeSeriesFilters";
import TimeSeriesTable from "../modules/timeSeries/components/TimeSeriesTable";

import { useTimeSeries } from "../modules/timeSeries/hooks/useTimeSeries";
import { useStartTimeSeriesExport } from "../modules/timeSeries/hooks/useStartTimeSeriesExport";
import { useTimeSeriesExportStore } from "../modules/timeSeries/stores/timeSeriesExportStore";
import { type ExportStatus, type TimeSeriesReadInput } from "../generated/graphql";
import { usePageState } from "../app/navigation/usePageState";
import {
  draftToTimeSeriesInput,
  mergeTimeSeriesDraft,
  timeSeriesPageStateCodec,
  type TimeSeriesDraftState,
  type TimeSeriesEntryState,
} from "../modules/timeSeries/state/timeSeriesPageState";

function isRunningStatus(status: ExportStatus) {
  return status === "pending" || status === "processing";
}

export default function TimeSeriesPage() {
  const activeExport = useTimeSeriesExportStore((s) => s.activeExport);
  const setActiveExport = useTimeSeriesExportStore((s) => s.setActiveExport);
  const { startTimeSeriesExport } = useStartTimeSeriesExport();
  const { query: queryState, entry, setPageState } = usePageState(
    timeSeriesPageStateCodec,
  );
  const draft = useMemo(
    () => mergeTimeSeriesDraft(queryState, entry),
    [entry, queryState],
  );
  const [applied, setApplied] = useState<TimeSeriesReadInput | null>(
    queryState.autoRun ? draftToTimeSeriesInput(draft) : null,
  );
  const didAutoRunRef = useRef(false);

  useEffect(() => {
    if (!draft.autoRun || didAutoRunRef.current) {
      return;
    }

    setApplied(draftToTimeSeriesInput(draft));
    didAutoRunRef.current = true;
    updateDraft({
      ...draft,
      autoRun: false,
    });
  }, [draft]);

  const query = useTimeSeries(applied as TimeSeriesReadInput, !!applied);

  const updateDraft = (nextDraft: TimeSeriesDraftState) => {
    const { sdInstanceIDs, kpiDefinitionIDs, ...nextQuery } = nextDraft;
    const nextEntry: TimeSeriesEntryState = {
      v: 1,
      sdInstanceIDs,
      kpiDefinitionIDs,
    };

    setPageState(
      {
        query: {
          ...nextQuery,
          autoRun: false,
        },
        entry: nextEntry,
      },
      { replace: true },
    );
  };

  const handleExport = async (input: TimeSeriesReadInput) => {
    try {
      if (activeExport && isRunningStatus(activeExport.status)) {
        toast.error("An export is already in progress.");
        return;
      }

      const exportJob = await startTimeSeriesExport(input);

      if (!exportJob) {
        throw new Error("Export failed");
      }

      setActiveExport(exportJob);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Export failed.");
    }
  };

  return (
    <div className="container mt-3">
      <TimeSeriesFilters
        value={draft}
        onChange={updateDraft}
        onSubmit={() => {
          setApplied(draftToTimeSeriesInput(draft));
        }}
        onExport={() => handleExport(draftToTimeSeriesInput(draft))}
      />

      {applied && <TimeSeriesTable query={query} />}
    </div>
  );
}
