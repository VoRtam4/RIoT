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
