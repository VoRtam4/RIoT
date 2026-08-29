/**
 * @file dateTimeUtils.ts
 * @brief Pomocné funkce pro formátování a převody datumu a času ve KPI modulech.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
export function toLocalInputValue(date: Date) {
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60000);
  return local.toISOString().slice(0, 19);
}

export function toUTCString(localValue: string) {
  return new Date(localValue).toISOString();
}

export function nowLocal() {
  return toLocalInputValue(new Date());
}

export function minusDaysLocal(days: number) {
  const d = new Date();
  d.setDate(d.getDate() - days);
  return toLocalInputValue(d);
}
