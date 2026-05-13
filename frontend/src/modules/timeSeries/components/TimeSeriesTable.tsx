/**
 * @file TimeSeriesTable.tsx
 * @brief Tabulka historických časových řad s dynamickými sloupci.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import {
  MaterialReactTable,
  useMaterialReactTable,
} from "material-react-table";
import { useMemo, useRef, useCallback, useEffect } from "react";

import { mapToRows } from "../utils/mapToRows";
import { buildColumns } from "../utils/buildColumns";

type Props = {
  query: any;
};

export default function TimeSeriesTable({ query }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);

  const allData = useMemo(() => {
    return query.data ?? [];
  }, [query.data]);

  const parameters = useMemo(() => {
    return query.parameters ?? [];
  }, [query.parameters]);

  const rows = useMemo(() => {
    return mapToRows(allData, parameters);
  }, [allData, parameters]);

  const columns = useMemo(() => {
    return buildColumns(parameters);
  }, [parameters]);

  const fetchMore = useCallback(
    (el?: HTMLDivElement | null) => {
      if (!el) return;
      if (query.isFetchingNextPage) return;
      if (!query.hasNextPage) return;

      const { scrollHeight, scrollTop, clientHeight } = el;

      if (scrollHeight - scrollTop - clientHeight < 200) {
        query.fetchNextPage();
      }
    },
    [query],
  );

  useEffect(() => {
    fetchMore(containerRef.current);
  }, [fetchMore]);

  const table = useMaterialReactTable({
    columns,
    data: rows,

    enablePagination: false,
    enableRowVirtualization: true,
    enableRowNumbers: true,

    muiTableContainerProps: {
      ref: containerRef,
      sx: {
        maxHeight: "600px",
        backgroundColor: "var(--bg-card)",
      },
      onScroll: (e: any) => fetchMore(e.target as HTMLDivElement),
    },

    muiTableHeadCellProps: {
      sx: {
        backgroundColor: "var(--bg-navbar)",
        color: "var(--text-main)",
      },
    },

    muiTableBodyCellProps: {
      sx: {
        backgroundColor: "var(--bg-card)",
        color: "var(--text-main)",
      },
    },

    muiTableBodyRowProps: {
      sx: {
        "&:hover": {
          backgroundColor: "var(--bg-input)",
        },
      },
    },

    state: {
      isLoading: query.loading,
      showProgressBars: query.isFetchingNextPage,
    },
  });

  return (
    <div>
      <style>
        {`
          .MuiPaper-root {
            background-color: var(--bg-card);
            box-shadow: none;
          }
          .MuiBox-root {
            box-shadow: none;
          }
          .MuiPaper-root > .MuiBox-root {
            background-color: var(--bg-card);
          }
          .MuiPaper-root > .MuiBox-root:last-child {
            max-height: 15px;
          }
        `}
      </style>
      <MaterialReactTable table={table} />
    </div>
  );
}
