/**
 * @file mapQueryToFilterNode.ts
 * @brief Převod dotazové struktury historických dat na strom filtrů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
export function mapToFilterNode(query: any): any {
  if (!query || !query.rules) return undefined;

  const mapGroup = (group: any): any => {
    return {
      type: "logical",
      operator: group.combinator,
      nodes: group.rules.map((r: any) => {
        if (r.rules) {
          return mapGroup(r);
        }

        return {
          type: "rule",
          rule: {
            tag: r.field,
            operator: r.operator,
            value: String(r.value),
          },
        };
      }),
    };
  };

  return mapGroup(query);
}
