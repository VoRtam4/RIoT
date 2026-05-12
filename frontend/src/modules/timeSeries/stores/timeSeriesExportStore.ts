import { create } from "zustand";
import type { TimeSeriesExport } from "../../../generated/graphql";

type Store = {
  activeExport: TimeSeriesExport | null;
  setActiveExport: (activeExport: TimeSeriesExport | null) => void;
};

export const useTimeSeriesExportStore = create<Store>((set) => ({
  activeExport: null,
  setActiveExport: (activeExport) => set({ activeExport }),
}));
