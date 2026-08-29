/**
 * @file reactSelectSearch.ts
 * @brief Pomocné funkce pro textové vyhledávání a filtrování položek v react-select.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
type SearchableOption = {
  label?: string | number | null;
  value?: string | number | null;
  searchText?: string | null;
};

export type SearchMode = "or" | "and" | "exact";

export const searchModes: SearchMode[] = ["or", "and", "exact"];

export function isSearchMode(value: string | null): value is SearchMode {
  return searchModes.includes(value as SearchMode);
}

function toSearchString(value: string | number | null | undefined): string {
  if (value == null) {
    return "";
  }

  return String(value);
}

function normalizeSearchValue(
  value: string | number | null | undefined,
): string {
  return toSearchString(value).trim().toLowerCase();
}

function sanitizeSearchToken(token: string): string {
  return token
    .replace(/^"+|"+$/g, "")
    .trim()
    .toLowerCase();
}

export function parseSearchInput(inputValue: string): string[] {
  if (!inputValue) {
    return [];
  }

  if (
    inputValue.startsWith('"') &&
    inputValue.endsWith('"') &&
    inputValue.length >= 2
  ) {
    const phrase = sanitizeSearchToken(inputValue.slice(1, -1));
    return phrase ? [phrase] : [];
  }

  return inputValue.split(/\s+/).map(sanitizeSearchToken).filter(Boolean);
}

function parseExactSearchInput(inputValue: string): string {
  return sanitizeSearchToken(inputValue);
}

export function buildOptionSearchText(
  ...values: Array<string | number | null | undefined>
) {
  return values
    .map((value) => toSearchString(value).trim())
    .filter((value): value is string => !!value)
    .join(" ");
}

export function matchesSearchText(
  inputValue: string,
  modeOrFirstValue?: SearchMode | string | number | null,
  ...values: Array<string | number | null | undefined>
) {
  const isExplicitMode =
    modeOrFirstValue === "and" ||
    modeOrFirstValue === "or" ||
    modeOrFirstValue === "exact";
  const mode: SearchMode = isExplicitMode ? modeOrFirstValue : "or";
  const haystackValues = isExplicitMode
    ? values
    : [modeOrFirstValue, ...values];
  const haystack = normalizeSearchValue(
    buildOptionSearchText(...haystackValues),
  );

  if (mode === "exact") {
    const phrase = parseExactSearchInput(inputValue);
    return phrase ? haystack.includes(phrase) : true;
  }

  const terms = parseSearchInput(inputValue);
  if (terms.length === 0) {
    return true;
  }

  const matcher = (term: string) => haystack.includes(term);

  return mode === "and" ? terms.every(matcher) : terms.some(matcher);
}

export function filterSelectOption(
  option: { data: SearchableOption },
  inputValue: string,
) {
  const terms = parseSearchInput(inputValue);
  if (terms.length === 0) {
    return true;
  }

  const haystack = normalizeSearchValue(
    option.data.searchText ??
      `${toSearchString(option.data.label)} ${toSearchString(option.data.value)}`,
  );

  return terms.some((term) => haystack.includes(term));
}
