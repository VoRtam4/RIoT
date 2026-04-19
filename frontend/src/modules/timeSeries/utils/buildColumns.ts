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
          p.denotation === "sdInstanceUID" || p.denotation === "kpiDefinitionID"
        );
      }

      return true;
    })
    .map((p) => ({
      accessorKey: p.denotation,
      header: p.label || p.denotation,
    }));
};
