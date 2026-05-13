/**
 * @file mapRawDataToRows.ts
 * @brief Mapování raw dat sledované instance do tabulkových řádků.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
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
