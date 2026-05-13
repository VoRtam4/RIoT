/**
 * @file muiTheme.ts
 * @brief Definice Material UI tématu navázaná na barevný systém aplikace.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { createTheme } from "@mui/material";
import { colors } from "./colors";

export const muiTheme = createTheme({
  palette: {
    mode: "dark",
    primary: { main: colors.primary },
    background: {
      default: colors.bgMain,
      paper: colors.bgCard,
    },
    text: {
      primary: colors.textMain,
    },
    divider: colors.border,
  },
});
