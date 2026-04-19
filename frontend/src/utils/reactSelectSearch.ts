type SearchableOption = {
  label?: string | null;
  value?: string | null;
  searchText?: string | null;
};

function normalizeSearchValue(value: string | null | undefined): string {
  return (value ?? "").trim().toLowerCase();
}

export function buildOptionSearchText(...values: Array<string | null | undefined>) {
  return values
    .map((value) => value?.trim())
    .filter((value): value is string => !!value)
    .join(" ");
}

export function matchesSearchText(
  inputValue: string,
  ...values: Array<string | null | undefined>
) {
  const terms = normalizeSearchValue(inputValue).split(/\s+/).filter(Boolean);
  if (terms.length === 0) {
    return true;
  }

  const haystack = normalizeSearchValue(buildOptionSearchText(...values));
  return terms.every((term) => haystack.includes(term));
}

export function filterSelectOption(
  option: { data: SearchableOption },
  inputValue: string,
) {
  const needle = normalizeSearchValue(inputValue);
  if (!needle) {
    return true;
  }

  const haystack = normalizeSearchValue(
    option.data.searchText ?? `${option.data.label ?? ""} ${option.data.value ?? ""}`,
  );

  return haystack.includes(needle);
}
