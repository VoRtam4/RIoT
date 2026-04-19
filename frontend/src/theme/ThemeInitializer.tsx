import { useEffect } from "react";
import { injectCssVariables } from "./injectCssVariables";

export default function ThemeInitializer() {
  useEffect(() => {
    injectCssVariables();
  }, []);

  return null;
}
