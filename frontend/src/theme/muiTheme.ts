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
