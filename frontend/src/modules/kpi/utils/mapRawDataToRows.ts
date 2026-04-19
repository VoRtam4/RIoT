type Row = {
  key: string;
  value: any;
};

export function mapRawDataToRows(payload: any): Row[] {
  if (!payload || typeof payload !== "object") return [];

  return Object.entries(payload).map(([key, value]) => ({
    key,
    value: typeof value === "object" ? JSON.stringify(value) : value,
  }));
}
