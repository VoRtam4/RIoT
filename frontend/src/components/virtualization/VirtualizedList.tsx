import { useRef, type CSSProperties, type ReactNode } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";

type Props<T> = {
  items: T[];
  rowHeight: number;
  renderItem: (item: T, index: number) => ReactNode;
  emptyState?: ReactNode;
  overscan?: number;
  threshold?: number;
  style?: CSSProperties;
  className?: string;
};

export default function VirtualizedList<T>({
  items,
  rowHeight,
  renderItem,
  emptyState,
  overscan = 4,
  threshold = 40,
  style,
  className,
}: Props<T>) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const shouldVirtualize = items.length >= threshold;

  const virtualizer = useVirtualizer({
    count: shouldVirtualize ? items.length : 0,
    getScrollElement: () => containerRef.current,
    estimateSize: () => rowHeight,
    overscan,
  });

  const virtualItems = virtualizer.getVirtualItems();

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
          ? items.map((item, index) => <div key={index}>{renderItem(item, index)}</div>)
          : (
              <div
                style={{
                  height: virtualizer.getTotalSize(),
                  position: "relative",
                  width: "100%",
                }}
              >
                {virtualItems.map((virtualItem) => (
                  <div
                    key={virtualItem.key}
                    style={{
                      position: "absolute",
                      top: 0,
                      left: 0,
                      right: 0,
                      height: virtualItem.size,
                      transform: `translateY(${virtualItem.start}px)`,
                    }}
                  >
                    {renderItem(items[virtualItem.index], virtualItem.index)}
                  </div>
                ))}
              </div>
            )}
    </div>
  );
}
