import {
  useRef,
  useState,
  useLayoutEffect,
} from "react";

type Props<T> = {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
  rowHeight: number;
  minColumnWidth: number;
  gap?: number;
  overscan?: number;
  style?: React.CSSProperties;
  className?: string;
  emptyState?: React.ReactNode;
};

export default function VirtualizedCardGrid<T>({
  items,
  renderItem,
  rowHeight,
  minColumnWidth,
  gap = 12,
  overscan = 2,
  style,
  className,
  emptyState,
}: Props<T>) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  const [size, setSize] = useState({
    width: 0,
    height: 0,
    scrollTop: 0,
  });

  useLayoutEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const update = () => {
      setSize((prev) => ({
        ...prev,
        width: el.clientWidth,
        height: el.clientHeight,
      }));
    };

    update();

    const observer = new ResizeObserver(update);
    observer.observe(el);

    return () => observer.disconnect();
  }, []);

  const onScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.currentTarget;
    setSize((prev) => ({
      ...prev,
      scrollTop: target.scrollTop,
    }));
  };

  if (!items.length) {
    return (
      <div ref={containerRef} style={{ ...style }}>
        {emptyState ?? null}
      </div>
    );
  }

  if (size.width === 0 || size.height === 0) {
    return (
      <div
        ref={containerRef}
        className={className}
        style={{
          overflowY: "auto",
          overflowX: "hidden",
          minHeight: 0,
          width: "100%",
          ...style,
        }}
      />
    );
  }

  const columnCount = Math.max(
    1,
    Math.floor((size.width + gap) / (minColumnWidth + gap)),
  );

  const rowCount = Math.ceil(items.length / columnCount);

  const visibleRowStart = Math.floor(size.scrollTop / rowHeight);
  const visibleRowEnd = Math.ceil(
    (size.scrollTop + size.height) / rowHeight,
  );

  const startRow = Math.max(0, visibleRowStart - overscan);
  const endRow = Math.min(rowCount, visibleRowEnd + overscan);

  const startIndex = startRow * columnCount;
  const endIndex = Math.min(items.length, endRow * columnCount);

  const visibleItems = items.slice(startIndex, endIndex);

  const offsetY = startRow * rowHeight;
  const totalHeight = rowCount * rowHeight;

  return (
    <div
      ref={containerRef}
      onScroll={onScroll}
      className={className}
      style={{
        overflowY: "auto",
        overflowX: "hidden",
        willChange: "transform",
        contain: "strict",
        width: "100%",
        ...style,
      }}
    >
      <div
        style={{
          height: totalHeight,
          position: "relative",
          width: "100%",
        }}
      >
        <div
          style={{
            transform: `translateY(${offsetY}px)`,
            display: "grid",

            gridTemplateColumns: `repeat(${columnCount}, minmax(0, 1fr))`,

            gap,
            width: "100%",
          }}
        >
          {visibleItems.map((item, index) => (
            <div
              key={startIndex + index}
              style={{
                minWidth: 0,
              }}
            >
              {renderItem(item)}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}