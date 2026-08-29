/**
 * @file TimeSeriesExportManager.tsx
 * @brief Ovládání spuštění, sledování a stažení exportů historických dat.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect, useRef } from "react";
import toast from "react-hot-toast";
import { apiEndpoints } from "../../../app/apiEndpoints";
import type { ExportStatus } from "../../../generated/graphql";
import { useCancelTimeSeriesExport } from "../hooks/useCancelTimeSeriesExport";
import { useTimeSeriesExport } from "../hooks/useTimeSeriesExport";
import { useTimeSeriesExportSubscription } from "../hooks/useTimeSeriesExportSubscription";
import { useTimeSeriesExportStore } from "../stores/timeSeriesExportStore";
import TimeSeriesExportToast from "./TimeSeriesExportToast";

const EXPORT_TOAST_ID = "time-series-export";

function isRunningStatus(status: ExportStatus) {
  return status === "pending" || status === "processing";
}

function startDownload(downloadUrl: string) {
  const anchor = document.createElement("a");
  anchor.href = downloadUrl;
  anchor.rel = "noreferrer";
  anchor.style.display = "none";
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
}

export default function TimeSeriesExportManager() {
  const activeExport = useTimeSeriesExportStore((s) => s.activeExport);
  const setActiveExport = useTimeSeriesExportStore((s) => s.setActiveExport);
  const { cancelTimeSeriesExport } = useCancelTimeSeriesExport();
  const { getTimeSeriesExport } = useTimeSeriesExport();
  const { latest } = useTimeSeriesExportSubscription(activeExport?.uid);
  const handledTerminalStateRef = useRef<string | null>(null);
  const syncedAfterSubscribeRef = useRef<string | null>(null);

  const handleCancel = async () => {
    if (!activeExport) {
      return;
    }

    try {
      const cancelled = await cancelTimeSeriesExport(activeExport.uid);

      if (cancelled) {
        setActiveExport(cancelled);
      }
    } catch {
      toast.error("Failed to cancel export.");
    }
  };

  useEffect(() => {
    if (!latest) {
      return;
    }

    setActiveExport(latest);
  }, [latest, setActiveExport]);

  useEffect(() => {
    if (!activeExport || !isRunningStatus(activeExport.status)) {
      if (!activeExport) {
        syncedAfterSubscribeRef.current = null;
      }
      return;
    }

    const syncKey = activeExport.uid;
    if (syncedAfterSubscribeRef.current === syncKey) {
      return;
    }
    syncedAfterSubscribeRef.current = syncKey;

    void (async () => {
      try {
        const currentExport = await getTimeSeriesExport(activeExport.uid);
        if (currentExport) {
          setActiveExport(currentExport);
        }
      } catch {
        syncedAfterSubscribeRef.current = null;
      }
    })();
  }, [activeExport, getTimeSeriesExport, setActiveExport]);

  useEffect(() => {
    if (!activeExport || !isRunningStatus(activeExport.status)) {
      toast.dismiss(EXPORT_TOAST_ID);
      return;
    }

    toast.custom(
      () => (
        <TimeSeriesExportToast
          activeExport={activeExport}
          onCancel={handleCancel}
        />
      ),
      {
        id: EXPORT_TOAST_ID,
        duration: Number.POSITIVE_INFINITY,
      },
    );
  }, [activeExport, handleCancel]);

  useEffect(() => {
    if (!activeExport || isRunningStatus(activeExport.status)) {
      return;
    }

    const stateKey = `${activeExport.uid}:${activeExport.status}`;
    if (handledTerminalStateRef.current === stateKey) {
      return;
    }
    handledTerminalStateRef.current = stateKey;

    toast.dismiss(EXPORT_TOAST_ID);

    if (activeExport.status === "done") {
      if (activeExport.downloadUrl) {
        startDownload(activeExport.downloadUrl);
        toast.success("Export successful.");
      } else {
        toast.error("Export finished without a download URL.");
      }
      setActiveExport(null);
      return;
    }

    if (activeExport.status === "cancelled") {
      toast.error("Export cancelled.");
      setActiveExport(null);
      return;
    }

    if (activeExport.status === "failed") {
      toast.error(activeExport.error ?? "Export failed.");
      setActiveExport(null);
      return;
    }

    if (activeExport.status === "expired") {
      toast.error("Export expired before download.");
      setActiveExport(null);
    }
  }, [activeExport, setActiveExport]);

  useEffect(() => {
    if (!activeExport || !isRunningStatus(activeExport.status)) {
      return;
    }

    const handleBeforeUnload = () => {
      void fetch(
        `${apiEndpoints.rest}/time-series/export/${activeExport.uid}`,
        {
          method: "DELETE",
          credentials: "include",
          keepalive: true,
        },
      );
    };

    window.addEventListener("beforeunload", handleBeforeUnload);

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [activeExport]);

  return null;
}
