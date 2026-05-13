/**
 * @file SdInstanceList.tsx
 * @brief Seznam sledovaných instancí s virtuálním vykreslováním a výběrem detailu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useRef } from "react";
import SdInstanceCard from "./SdInstanceCard";
import VirtualizedCardGrid from "../../../components/virtualization/VirtualizedCardGrid";
import EmptyStateNotice from "../../../components/EmptyStateNotice";

type Instance = {
  id: string;
  label?: string | null;
  uid?: string | null;
};

type Props = {
  instances?: Instance[];
  loading?: any;
  onOpen: (id: string) => void;
};

export default function SdInstanceList({ instances = [], loading, }: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  
  if (loading) {
    return (
      <div className="d-flex h-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  return (
    <div ref={containerRef}>
      <VirtualizedCardGrid
        items={instances}
        rowHeight={136}
        minColumnWidth={300}
        style={{ height: "calc(100vh - 220px)" }}
        emptyState={
          <EmptyStateNotice
            title="No devices match the current filter"
            description="Try a different search phrase or model."
          />
        }
        renderItem={(inst) => (
          <SdInstanceCard
            key={inst.id}
            instance={inst}
          />
        )}
      />
    </div>
  );
}
