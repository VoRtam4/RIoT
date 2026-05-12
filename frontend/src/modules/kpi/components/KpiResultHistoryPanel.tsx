import { useEffect, useMemo, useState } from "react";
import { LineChart } from "@mui/x-charts/LineChart";

import { useKpiResult } from "../hooks/useKpiResult";
import { useKpiSubscription } from "../hooks/useKpiSubscription";
import { useTimeSeriesAggregateKpi } from "../../timeSeries/hooks/useTimeSeriesAggregateKpi";
import { toUTCString } from "../utils/dateTimeUtils";

import { colors } from "../../../theme/colors";
import EmptyStateNotice from "../../../components/EmptyStateNotice";

type Props = {
  kpiDefinitionID?: string;
  sdInstanceID?: string | null;
  sdTypeID?: string | null;
  from: string;
  to: string;
  onFromChange: (value: string) => void;
  onToChange: (value: string) => void;
};

function computeAggregateSeconds(from: string, to: string) {
  const diff =
    (new Date(to).getTime() - new Date(from).getTime()) / 1000;

  if (diff <= 50) return 1;

  return Math.floor(diff / 50);
}

export default function KpiResultHistoryPanel({
  kpiDefinitionID,
  sdInstanceID,
  sdTypeID,
  from,
  to,
  onFromChange,
  onToChange,
}: Props) {
  const [debouncedFrom, setDebouncedFrom] = useState(from);
  const [debouncedTo, setDebouncedTo] = useState(to);

  useEffect(() => {
    const t = setTimeout(() => {
      setDebouncedFrom(from);
      setDebouncedTo(to);
    }, 800);

    return () => clearTimeout(t);
  }, [from, to]);

  const safeInstanceID = sdInstanceID ?? undefined;

  const { result } = useKpiResult(kpiDefinitionID, safeInstanceID);
  const { latest } = useKpiSubscription(kpiDefinitionID, safeInstanceID);

  const liveResult = latest ?? result;

  const aggregateSeconds = computeAggregateSeconds(
    debouncedFrom,
    debouncedTo,
  );

  const tsQuery = useTimeSeriesAggregateKpi(
    {
      from: toUTCString(debouncedFrom),
      to: toUTCString(debouncedTo),
      aggregateSeconds,

      kpiDefinitionIDs: kpiDefinitionID ? [kpiDefinitionID] : [],
      sdInstanceIDs: sdInstanceID ? [String(sdInstanceID)] : [],
      sdTypeID: sdTypeID ? String(sdTypeID) : undefined,
    },
    !!kpiDefinitionID && !!sdInstanceID,
  );

  const chartData = useMemo(() => {
    const raw = tsQuery.data ?? [];
    if (!raw.length) return [];

    const points = raw
      .map((d: any) => {
        let parsed = d.data;

        if (typeof parsed === "string") {
          try {
            parsed = JSON.parse(parsed);
          } catch {
            parsed = {};
          }
        }

        return {
          time: new Date(d.time).getTime(),
          value: parsed?.fulfilled === true ? 1 : 0,
        };
      })
      .sort((a: any, b: any) => a.time - b.time);

    const fromTs = new Date(toUTCString(debouncedFrom)).getTime();
    const toTs = new Date(toUTCString(debouncedTo)).getTime();

    const result: { time: number; value: number }[] = [];

    let currentValue = 0;

    for (let i = 0; i < points.length; i++) {
      if (points[i].time <= fromTs) {
        currentValue = points[i].value;
      }
    }

    result.push({
      time: fromTs,
      value: currentValue,
    });

    for (let i = 0; i < points.length; i++) {
      const p = points[i];

      if (p.time < fromTs) continue;
      if (p.time > toTs) break;

      result.push({
        time: p.time,
        value: currentValue,
      });

      currentValue = p.value;

      result.push({
        time: p.time,
        value: currentValue,
      });
    }

    result.push({
      time: toTs,
      value: currentValue,
    });

    return result;
  }, [tsQuery.data, debouncedFrom, debouncedTo]);

  const x = chartData.map((d) => d.time);
  const y = chartData.map((d) => d.value);

  if (!sdInstanceID || ! kpiDefinitionID) {
    return (
      <div className="h-100 d-flex align-items-center">
        <EmptyStateNotice
          title="Select a device"
          description="Choose a device from the sidebar to view the KPI evaluation history."
        />
      </div>
    );
  }

  return (
    <div className="h-100 d-flex flex-column">
      {/* TOP BAR */}
      <div className="row g-3 mb-3 align-items-center">

        {/* STATUS */}
        <div className="col-md-4 d-flex align-items-center gap-2 form-label">
          <div
            style={{
              width: 12,
              height: 12,
              borderRadius: "50%",
              backgroundColor: liveResult?.fulfilled
                ? colors.success
                : colors.error,
            }}
          />

          <span className="form-label">
            {liveResult?.fulfilled ? "True" : "False"}
          </span>

          {liveResult?.eventTime && (
            <small className="form-label" style={{ opacity: 0.6 }}>
              ({new Date(liveResult.eventTime).toLocaleString("cs-CZ")})
            </small>
          )}
        </div>

        {/* FROM */}
        <div className="col-md-4">
          <label className="form-label">From</label>
          <input
            type="datetime-local"
            step="1"
            className="form-control"
            value={from}
            onChange={(e) => onFromChange(e.target.value)}
          />
        </div>

        {/* TO */}
        <div className="col-md-4">
          <label className="form-label">To</label>
          <input
            type="datetime-local"
            step="1"
            className="form-control"
            value={to}
            onChange={(e) => onToChange(e.target.value)}
          />
        </div>
      </div>

      {/* GRAPH */}
      <div style={{ flex: 1, color: colors.textMain }}>
        <LineChart
          xAxis={[
            {
              data: x,
              scaleType: "time",
            },
          ]}
          series={[
            {
              data: y,
              curve: "stepAfter",
              area: true,
            },
          ]}
          height={300}
        />
      </div>
    </div>
  );
}
