import { useMemo } from "react";
import {
  MaterialReactTable,
  useMaterialReactTable,
} from "material-react-table";

import { useRawDataPoint } from "../../sdInstances/hooks/useRawDataPoint";
import { useRawDataSubscription } from "../../sdInstances/hooks/useRawDataSubscription";

type SDParameter = {
  id: string;
  label: string;
  denotation: string;
  type: string;
};

type SDType = {
  id: string;
  label: string;
  uid: string;
  parameters: SDParameter[];
};

type Props = {
  sdInstanceID?: string | null;
  sdType?: SDType | null;
};

type Row = {
  parameter: string;
  value: string;
  type: string;
};

function isDate(value: any): boolean {
  if (typeof value !== "string") return false;
  const trimmed = value.trim();
  if (!trimmed) return false;
  if (!/^\d{4}-\d{2}-\d{2}(?:[T\s].*)?$/.test(trimmed)) return false;
  return !isNaN(Date.parse(value));
}

function formatDate(value: string): string {
  try {
    return new Date(value).toLocaleString("cs-CZ", {
      timeZone: "Europe/Prague",
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  } catch {
    return value;
  }
}

export default function KpiSdTypeDataPanel({ sdInstanceID, sdType }: Props) {
  const { rawDataPoint } = useRawDataPoint(sdInstanceID);

  const { latest } = useRawDataSubscription(sdType?.id, sdInstanceID);

  const effectiveData = latest ?? rawDataPoint;
  const effectiveEventTime =
    latest?.eventTime ??
    rawDataPoint?.eventTime ??
    (typeof effectiveData?.payload === "object" &&
    effectiveData?.payload !== null &&
    "eventTime" in effectiveData.payload
      ? String((effectiveData.payload as Record<string, unknown>).eventTime ?? "")
      : "");
  
  const parsed = useMemo(() => {
    if (!effectiveData?.payload) return {};

    try {
      return typeof effectiveData.payload === "string"
        ? JSON.parse(effectiveData.payload)
        : effectiveData.payload;
    } catch {
      return {};
    }
  }, [effectiveData]);
  
  const rows: Row[] = useMemo(() => {
    const dataRows = (sdType?.parameters ?? []).map((p) => {
      const val = parsed[p.denotation];

      let formattedValue = "";

      if (val !== undefined && val !== null) {
        if (typeof val === "object") {
          formattedValue = JSON.stringify(val);
        } else if (isDate(val)) {
          formattedValue = formatDate(val);
        } else {
          formattedValue = String(val);
        }
      }

      return {
        parameter: p.label ?? p.denotation,
        value: formattedValue,
        type: p.type,
      };
    });

    if (!effectiveEventTime) {
      return dataRows;
    }

    return [
      {
        parameter: "Time",
        value: formatDate(effectiveEventTime),
        type: "string",
      },
      ...dataRows,
    ];
  }, [effectiveEventTime, parsed, sdType]);

  const columns = useMemo(
    () => [
      {
        accessorKey: "parameter",
        header: "Parameter",
      },
      {
        accessorKey: "value",
        header: "Value",
      },
      {
        accessorKey: "type",
        header: "Type",
      },
    ],
    [],
  );

  const table = useMaterialReactTable({
    columns,
    data: rows,

    enablePagination: false,
    enableTopToolbar: false,
    enableBottomToolbar: false,
    enableStickyHeader: true,

    muiTableContainerProps: {
      sx: {
        height: "100%",
        backgroundColor: "var(--bg-card)",
      },
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
  });

  return (
    <div style={{ height: "100%" }}>
      <MaterialReactTable table={table} />
    </div>
  );
}
