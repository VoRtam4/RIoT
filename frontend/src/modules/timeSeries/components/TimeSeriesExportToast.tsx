import { useEffect, useState } from "react";
import type { TimeSeriesExport } from "../../../generated/graphql";

function formatElapsed(createdAt: string, now: number) {
  const diffMs = Math.max(0, now - new Date(createdAt).getTime());
  const totalSeconds = Math.floor(diffMs / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${String(seconds).padStart(2, "0")}`;
}

type Props = {
  activeExport: TimeSeriesExport;
  onCancel: () => void;
};

export default function TimeSeriesExportToast({
  activeExport,
  onCancel,
}: Props) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setNow(Date.now());
    }, 1000);
    return () => window.clearInterval(intervalId);
  }, [activeExport]);

  return (
    <div
      style={{
        minWidth: 320,
        background: "var(--bg-card)",
        color: "var(--text-main)",
        border: "1px solid var(--border)",
        borderRadius: "10px",
        padding: "10px 14px",
        fontSize: "14px",
      }}
    >
      <div className="d-flex align-items-center gap-3">
        <div className="spinner-border text-primary" role="status" />

        <div className="flex-grow-1">
          <div className="fw-semibold">Exporting dataset</div>
          <div style={{ fontSize: 13, opacity: 0.8 }}>
            Preparing download. Elapsed{" "}
            {formatElapsed(activeExport.createdAt, now)}
          </div>
        </div>

        <button
          type="button"
          className="btn btn-sm btn-outline-danger"
          onClick={() => {
            onCancel();
          }}
        >
          Cancel
        </button>
      </div>
    </div>
  );
}
