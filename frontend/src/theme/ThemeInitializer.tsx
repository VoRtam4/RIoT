/**
 * @file ThemeInitializer.tsx
 * @brief Inicializace CSS proměnných a Material UI tématu při startu frontendu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { injectCssVariables } from "./injectCssVariables";

export default function ThemeInitializer() {
  useEffect(() => {
    injectCssVariables();
  }, []);

  return null;
}
