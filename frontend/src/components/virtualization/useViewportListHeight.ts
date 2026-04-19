import { useEffect, useState, type RefObject } from "react";

type Options = {
  bottomOffset?: number;
  minHeight?: number;
};

function measureHeight(
  element: HTMLElement | null,
  { bottomOffset = 24, minHeight = 320 }: Options,
) {
  if (!element || typeof window === "undefined") {
    return minHeight;
  }

  const top = element.getBoundingClientRect().top;
  return Math.max(minHeight, window.innerHeight - top - bottomOffset);
}

export function useViewportListHeight<T extends HTMLElement>(
  ref: RefObject<T | null>,
  options: Options = {},
) {
  const [height, setHeight] = useState(options.minHeight ?? 320);

  useEffect(() => {
    const element = ref.current;
    if (!element) {
      return;
    }

    const updateHeight = () => {
      setHeight(measureHeight(element, options));
    };

    updateHeight();

    const resizeObserver = new ResizeObserver(updateHeight);
    resizeObserver.observe(element);
    window.addEventListener("resize", updateHeight);

    return () => {
      resizeObserver.disconnect();
      window.removeEventListener("resize", updateHeight);
    };
  }, [options.bottomOffset, options.minHeight, ref]);

  return height;
}
