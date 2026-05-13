/**
 * @file mapToRows.ts
 * @brief Mapování odpovědi časových řad do tabulkových řádků.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
type TimeSeriesParameterLike = {
  denotation: string;
  label: string;
  role: string;
};

type TimeSeriesDataPointLike = {
  time: string;
  tags?: unknown;
  data?: unknown;
};

export const mapToRows = (
  data: TimeSeriesDataPointLike[],
  parameters: TimeSeriesParameterLike[],
) => {
  if (!data?.length) return [];

  return data.map((point) => {
    let parsedData: Record<string, unknown> = {};
    let parsedTags: Record<string, unknown> = {};

    if (typeof point.data === "string") {
      try {
        parsedData = JSON.parse(point.data);
      } catch {
        parsedData = {};
      }
    } else if (typeof point.data === "object" && point.data !== null) {
      parsedData = point.data as Record<string, unknown>;
    }

    if (typeof point.tags === "string") {
      try {
        parsedTags = JSON.parse(point.tags);
      } catch {
        parsedTags = {};
      }
    } else if (typeof point.tags === "object" && point.tags !== null) {
      parsedTags = point.tags as Record<string, unknown>;
    }

    const row: Record<string, unknown> = Object.fromEntries(
      parameters.map((p) => [p.denotation, ""]),
    );

    for (const p of parameters) {
      const role = p.role?.toLowerCase();

      if (
        role === "meta" &&
        p.denotation !== "sdInstanceUID" &&
        p.denotation !== "kpiDefinitionID"
      ) {
        continue;
      }

      let value: unknown;

      switch (role) {
        case "time":
          value = point.time;
          break;

        case "tag":
          value = parsedTags[p.denotation];
          break;

        case "field":
          value = parsedData[p.denotation];
          break;

        case "meta":
          if (p.denotation === "sdInstanceUID") {
            value = parsedTags["sdInstanceUID"] ?? parsedData["sdInstanceUID"];
          } else {
            value = parsedData[p.denotation] ?? parsedTags[p.denotation];
          }
          break;

        default:
          value =
            parsedData[p.denotation] ??
            parsedTags[p.denotation] ??
            (p.denotation === "time" ? point.time : undefined);
      }

      if (typeof value === "boolean") {
        row[p.denotation] = value ? "true" : "false";
      } else if (value === null || value === undefined) {
        row[p.denotation] = "";
      } else if (typeof value === "object") {
        row[p.denotation] = JSON.stringify(value);
      } else {
        row[p.denotation] = String(value);
      }
    }

    return row;
  });
};
