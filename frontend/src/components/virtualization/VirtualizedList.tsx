/**
 * @file VirtualizedList.tsx
 * @brief Sdílený virtualizovaný seznam pro výkonné vykreslení většího počtu položek.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import {
  useRef,
  useCallback,
  useEffect,
  useState,
  type CSSProperties,
  type ReactNode,
} from "react";
import { useVirtualizer } from "@tanstack/react-virtual";

type Props<T> = {
  items: T[];
  rowHeight: number;
  itemSpacing?: number;
  renderItem: (item: T, index: number) => ReactNode;
  emptyState?: ReactNode;
  overscan?: number;
  threshold?: number;
  style?: CSSProperties;
  className?: string;
  scrollToIndex?: number;
  pinnedIndex?: number;
};

export default function VirtualizedList<T>({
  items,
  rowHeight,
  itemSpacing = 0,
  renderItem,
  emptyState,
  overscan = 4,
  threshold = 40,
  style,
  className,
  scrollToIndex,
  pinnedIndex,
}: Props<T>) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const pinnedRightInset = 8;
  const pinnedVerticalInset = 5;
  const [{ scrollTop, viewportHeight }, setViewportState] = useState({
    scrollTop: 0,
    viewportHeight: 0,
  });

  const shouldVirtualize = items.length >= threshold;

  const estimateSize = useCallback(
    (index: number) => rowHeight + (index === items.length - 1 ? 0 : itemSpacing),
    [itemSpacing, items.length, rowHeight],
  );
  const getScrollElement = useCallback(() => containerRef.current, []);

  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement,
    estimateSize,
    overscan,
  });

  useEffect(() => {
    const element = containerRef.current;
    if (!element) return;

    const syncViewportState = () => {
      setViewportState((current) => {
        const next = {
          scrollTop: element.scrollTop,
          viewportHeight: element.clientHeight,
        };

        if (
          current.scrollTop === next.scrollTop &&
          current.viewportHeight === next.viewportHeight
        ) {
          return current;
        }

        return next;
      });
    };

    syncViewportState();

    const resizeObserver = new ResizeObserver(syncViewportState);
    resizeObserver.observe(element);

    return () => {
      resizeObserver.disconnect();
    };
  }, [items.length, rowHeight, shouldVirtualize]);

  useEffect(() => {
    if (scrollToIndex == null) return;
    if (scrollToIndex < 0 || scrollToIndex >= items.length) return;

    virtualizer.scrollToIndex(scrollToIndex, {
      align: "auto",
    });
  }, [virtualizer, scrollToIndex, items.length]);

  const virtualItems = virtualizer?.getVirtualItems() ?? [];
  const shouldEvaluatePinnedItem =
    pinnedIndex != null &&
    pinnedIndex >= 0 &&
    pinnedIndex < items.length &&
    viewportHeight > 0;
  const pinnedItemStart = shouldEvaluatePinnedItem
    ? pinnedIndex! * (rowHeight + itemSpacing)
    : 0;
  const pinnedItemEnd = pinnedItemStart + rowHeight;
  const viewportEnd = scrollTop + viewportHeight;
  const pinToTop =
    shouldEvaluatePinnedItem && pinnedItemStart < scrollTop;
  const pinToBottom =
    shouldEvaluatePinnedItem && pinnedItemEnd > viewportEnd;
  const shouldRenderPinnedItem = pinToTop || pinToBottom;

  const scrollPinnedItemIntoView = useCallback(() => {
    if (pinnedIndex == null) return;
    if (pinnedIndex < 0 || pinnedIndex >= items.length) return;

    virtualizer.scrollToIndex(pinnedIndex, {
      align: "center",
    });
  }, [virtualizer, pinnedIndex, items.length]);

  return (
    <div
      className={className}
      style={{
        minHeight: 0,
        position: "relative",
        ...style,
      }}
    >
      {shouldRenderPinnedItem && (
        <div
          style={{
            position: "absolute",
            top: pinToTop ? -pinnedVerticalInset : undefined,
            bottom: pinToBottom ? -pinnedVerticalInset : undefined,
            left: 0,
            right: pinnedRightInset,
            height: rowHeight + pinnedVerticalInset * 2,
            zIndex: 3,
            pointerEvents: "none",
            margin: 0,
            boxSizing: "border-box",
          }}
        >
          <div
            style={{
              pointerEvents: "auto",
              position: "absolute",
              top: pinnedVerticalInset,
              bottom: pinnedVerticalInset,
              left: 0,
              right: 0,
              boxSizing: "border-box",
            }}
          >
            {renderItem(items[pinnedIndex!], pinnedIndex!)}

            <button
              type="button"
              aria-label="Scroll to selected item"
              onClick={scrollPinnedItemIntoView}
              style={{
                position: "absolute",
                inset: 0,
                border: 0,
                padding: 0,
                margin: 0,
                background: "transparent",
                cursor: "pointer",
              }}
            />
          </div>
        </div>
      )}

      <div
        ref={containerRef}
        onScroll={(event) => {
          const target = event.currentTarget;
          setViewportState((current) => {
            const next = {
              scrollTop: target.scrollTop,
              viewportHeight: target.clientHeight,
            };

            if (
              current.scrollTop === next.scrollTop &&
              current.viewportHeight === next.viewportHeight
            ) {
              return current;
            }

            return next;
          });
        }}
        style={{
          overflowY: "auto",
          minHeight: 0,
          height: "100%",
        }}
      >
        {!items.length
          ? (emptyState ?? null)
          : !shouldVirtualize
            ? items.map((item: any, index) => (
                <div
                  key={item.id ?? index}
                  style={{
                    height: rowHeight + (index === items.length - 1 ? 0 : itemSpacing),
                  }}
                >
                  <div
                    style={{
                      height: rowHeight,
                      marginBottom: index === items.length - 1 ? 0 : itemSpacing,
                    }}
                  >
                    {renderItem(item, index)}
                  </div>
                </div>
              ))
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
                      <div
                        style={{
                          height: rowHeight,
                          marginBottom:
                            virtualItem.index === items.length - 1 ? 0 : itemSpacing,
                        }}
                      >
                        {renderItem(items[virtualItem.index], virtualItem.index)}
                      </div>
                    </div>
                  ))}
                </div>
              )}
      </div>
    </div>
  );
}
