import { useEffect, useMemo, useRef, useState } from "react";
import { LineChart } from "@mui/x-charts/LineChart";

import { useKpiResult } from "../hooks/useKpiResult";
import { useKpiSubscription } from "../hooks/useKpiSubscription";
import { useTimeSeriesAggregateKpi } from "../../timeSeries/hooks/useTimeSeriesAggregateKpi";

import { colors } from "../../../theme/colors";

type Props = {
  kpiDefinitionID?: string;
  sdInstanceID?: string | null;
  sdTypeID?: string | null;
};

function toLocalInputValue(date: Date) {
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60000);
  return local.toISOString().slice(0, 19);
}

function toUTCString(localValue: string) {
  return new Date(localValue).toISOString();
}

function nowLocal() {
  return toLocalInputValue(new Date());
}

function minus7DaysLocal() {
  const d = new Date();
  d.setDate(d.getDate() - 7);
  return toLocalInputValue(d);
}

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
}: Props) {
  const [from, setFrom] = useState(minus7DaysLocal());
  const [to, setTo] = useState(nowLocal());

  const [debouncedFrom, setDebouncedFrom] = useState(from);
  const [debouncedTo, setDebouncedTo] = useState(to);

  const fromRef = useRef<HTMLInputElement | null>(null);
  const toRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    const t = setTimeout(() => {
      const isTyping =
        document.activeElement === fromRef.current ||
        document.activeElement === toRef.current;

      if (!isTyping) {
        setDebouncedFrom(from);
        setDebouncedTo(to);
      }
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
    return <div className="d-flex justify-content-center align-items-center">
      <span className="form-label">Select the device item to view the KPI definition evaluation historyVyberte položku pro zobrazení historie</span>
    </div>;
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
          <label className="form-label">Od</label>
          <input
            ref={fromRef}
            type="datetime-local"
            step="1"
            className="form-control"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
          />
        </div>

        {/* TO */}
        <div className="col-md-4">
          <label className="form-label">Do</label>
          <input
            ref={toRef}
            type="datetime-local"
            step="1"
            className="form-control"
            value={to}
            onChange={(e) => setTo(e.target.value)}
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