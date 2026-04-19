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
