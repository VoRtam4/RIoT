/**
 * @file timeSeriesExportStore.ts
 * @brief Lokální store běžících a dokončených exportů časových řad.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
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
