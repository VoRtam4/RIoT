import { useEffect, useRef, useState } from "react";
import { searchModes, type SearchMode } from "../utils/reactSelectSearch";

type Props = {
  value: string;
  onChange: (value: string) => void;
  className?: string;
  placeholder?: string;
  mode?: SearchMode;
  onModeChange?: (mode: SearchMode) => void;
};

const searchModeLabels: Record<SearchMode, string> = {
  or: "OR",
  and: "AND",
  exact: "EXACT",
};

const searchModeTitles: Record<SearchMode, string> = {
  or: "Matches at least one separated term",
  and: "Matches all separated terms",
  exact: "Matches the exact entered text including spaces",
};

function nextSearchMode(mode: SearchMode): SearchMode {
  const index = searchModes.indexOf(mode);
  return searchModes[(index + 1) % searchModes.length] ?? "or";
}

export default function UrlSyncedSearchInput({
  value,
  onChange,
  className = "form-control",
  placeholder,
  mode,
  onModeChange,
}: Props) {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [localValue, setLocalValue] = useState(value);

  useEffect(() => {
    if (document.activeElement === inputRef.current) {
      return;
    }

    queueMicrotask(() => {
      setLocalValue(value);
    });
  }, [value]);

  const input = (
    <input
      ref={inputRef}
      className={className}
      placeholder={placeholder}
      value={localValue}
      onChange={(event) => {
        const nextValue = event.target.value;
        setLocalValue(nextValue);
        onChange(nextValue);
      }}
      onBlur={() => {
        if (localValue !== value) {
          setLocalValue(value);
        }
      }}
    />
  );

  if (!mode || !onModeChange) {
    return input;
  }

  return (
    <div className="search-mode-input">
      <button
        type="button"
        className="btn btn-outline-light search-mode-input__toggle"
        aria-label={`Search mode: ${searchModeLabels[mode]}`}
        title={searchModeTitles[mode]}
        onClick={() => onModeChange(nextSearchMode(mode))}
      >
        {searchModeLabels[mode]}
      </button>
      {input}
    </div>
  );
}
