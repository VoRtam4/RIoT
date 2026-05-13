/**
 * @file injectCssVariables.ts
 * @brief Vložení barevných a layoutových CSS proměnných do dokumentu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { colors } from "./colors";

const toKebabCase = (str: string) =>
  str.replace(/[A-Z]/g, (letter) => `-${letter.toLowerCase()}`);

export const injectCssVariables = () => {
  const root = document.documentElement;

  Object.entries(colors).forEach(([key, value]) => {
    const cssKey = toKebabCase(key);
    root.style.setProperty(`--${cssKey}`, value);
  });
};
