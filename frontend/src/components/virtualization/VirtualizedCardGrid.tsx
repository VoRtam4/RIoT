import { useEffect, useRef, useState, type CSSProperties, type ReactNode } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";

type Props<T> = {
  items: T[];
  renderItem: (item: T, index: number) => ReactNode;
  emptyState?: ReactNode;
  minColumnWidth?: number;
  rowHeight?: number;
  gap?: number;
  overscan?: number;
  threshold?: number;
  style?: CSSProperties;
  className?: string;
};

export default function VirtualizedCardGrid<T>({
  items,
  renderItem,
  emptyState,
  minColumnWidth = 300,
  rowHeight = 136,
  gap = 16,
  overscan = 2,
  threshold = 40,
  style,
  className,
}: Props<T>) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [containerSize, setContainerSize] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const element = containerRef.current;
    if (!element) {
      return;
    }

    const updateSize = () => {
      setContainerSize({
        width: element.clientWidth,
        height: element.clientHeight,
      });
    };

    updateSize();

    const resizeObserver = new ResizeObserver(updateSize);
    resizeObserver.observe(element);

    return () => resizeObserver.disconnect();
  }, []);

  const columns = Math.max(
    1,
    Math.floor((containerSize.width + gap) / (minColumnWidth + gap)),
  );
  const itemWidth = Math.max(
    minColumnWidth,
    (containerSize.width - gap * (columns - 1)) / columns,
  );
  const rowSpan = rowHeight + gap;
  const rowCount = Math.ceil(items.length / columns);
  const totalHeight = Math.max(0, rowCount * rowSpan - gap);
  const shouldVirtualize =
    items.length >= threshold &&
    containerSize.width > 0 &&
    containerSize.height > 0;

  const rowVirtualizer = useVirtualizer({
    count: shouldVirtualize ? rowCount : 0,
    getScrollElement: () => containerRef.current,
    estimateSize: () => rowSpan,
    overscan,
  });

  const virtualRows = rowVirtualizer.getVirtualItems();

  return (
    <div
      ref={containerRef}
      className={className}
      style={{
        overflowY: "auto",
        minHeight: 0,
        ...style,
      }}
    >
      {!items.length
        ? (emptyState ?? null)
        : !shouldVirtualize
          ? (
              <div
                style={{
                  display: "grid",
                  gridTemplateColumns: `repeat(auto-fill, minmax(${minColumnWidth}px, 1fr))`,
                  gap,
                }}
              >
                {items.map((item, index) => (
                  <div key={index}>{renderItem(item, index)}</div>
                ))}
              </div>
            )
          : (
              <div
                style={{
                  height: totalHeight,
                  position: "relative",
                  width: "100%",
                }}
              >
                {virtualRows.map((virtualRow) => {
                  const row = virtualRow.index;
                  const children: ReactNode[] = [];

                  for (let column = 0; column < columns; column += 1) {
                    const index = row * columns + column;
                    if (index >= items.length) {
                      break;
                    }

                    children.push(
                      <div
                        key={index}
                        style={{
                          position: "absolute",
                          top: virtualRow.start,
                          left: column * (itemWidth + gap),
                          width: itemWidth,
                          height: rowHeight,
                        }}
                      >
                        {renderItem(items[index], index)}
                      </div>,
                    );
                  }

                  return children;
                })}
              </div>
            )}
    </div>
  );
}
