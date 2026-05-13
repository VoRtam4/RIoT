/**
 * @file ipValidator.ts
 * @brief Validace IP adres a CIDR rozsahů používaných u omezení API klíčů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
export const isValidCIDR = (v: string) => {
  const value = v.trim();
  return /^(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\/(?:3[0-2]|[12]?\d))?$/.test(
    value,
  );
};
