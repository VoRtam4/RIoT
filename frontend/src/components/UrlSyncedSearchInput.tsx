import { useEffect, useRef, useState } from "react";

type Props = {
  value: string;
  onChange: (value: string) => void;
  className?: string;
  placeholder?: string;
};

export default function UrlSyncedSearchInput({
  value,
  onChange,
  className = "form-control",
  placeholder,
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

  return (
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
}
