import { Children, useMemo, useRef } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import { components, type GroupBase, type MenuListProps } from "react-select";

const OPTION_HEIGHT = 38;
const OVERSCAN = 4;
const VIRTUALIZE_FROM = 120;

function VirtualizedMenuList<
  Option,
  IsMulti extends boolean,
  Group extends GroupBase<Option>,
>(props: MenuListProps<Option, IsMulti, Group>) {
  const rows = Children.toArray(props.children);
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const shouldVirtualize = rows.length >= VIRTUALIZE_FROM;

  const rowVirtualizer = useVirtualizer({
    count: shouldVirtualize ? rows.length : 0,
    getScrollElement: () => scrollRef.current,
    estimateSize: () => OPTION_HEIGHT,
    overscan: OVERSCAN,
  });

  const content = useMemo(() => {
    if (!shouldVirtualize) {
      return props.children;
    }

    return (
      <div
        style={{
          height: rowVirtualizer.getTotalSize(),
          position: "relative",
        }}
      >
        {rowVirtualizer.getVirtualItems().map((virtualItem) => (
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
            {rows[virtualItem.index]}
          </div>
        ))}
      </div>
    );
  }, [props.children, rowVirtualizer, rows, shouldVirtualize]);

  return (
    <components.MenuList
      {...props}
      innerRef={(instance) => {
        scrollRef.current = instance;
        if (typeof props.innerRef === "function") {
          props.innerRef(instance);
          return;
        }

        if (props.innerRef) {
          props.innerRef.current = instance;
        }
      }}
    >
      {content}
    </components.MenuList>
  );
}

export const virtualizedSelectComponents = {
  MenuList: VirtualizedMenuList,
};

export const virtualizedSelectProps = {
  components: virtualizedSelectComponents as any,
  maxMenuHeight: 320,
} as const;
