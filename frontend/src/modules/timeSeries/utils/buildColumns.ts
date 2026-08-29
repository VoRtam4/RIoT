/**
 * @file buildColumns.ts
 * @brief Sestavení sloupců tabulky historických dat podle vybraných parametrů.
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

export const buildColumns = (parameters: TimeSeriesParameterLike[]) => {
  if (!parameters?.length) return [];

  return parameters
    .filter((p) => {
      const role = p.role?.toLowerCase();

      if (role === "meta") {
        return (
          p.denotation === "sdInstanceUID" ||
          p.denotation === "kpiDefinitionUID"
        );
      }

      return true;
    })
    .map((p) => ({
      accessorKey: p.denotation,
      header: p.label || p.denotation,
    }));
};
