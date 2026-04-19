import CreatableSelect from "react-select/creatable";
import { isValidCIDR } from "../utils/ipValidator";

type Props = {
  value: string[];
  onChange: (v: string[]) => void;
};

export default function APIKeyIPInput({ value, onChange }: Props) {
  return (
    <CreatableSelect
      classNamePrefix="react-select"
      isMulti
      options={[]}
      value={value.map((v) => ({ label: v, value: v }))}
      onChange={(items) => onChange(items.map((i) => i.value))}
      onCreateOption={(input) => {
        if (!isValidCIDR(input)) return;
        if (value.includes(input)) return;
        onChange([...value, input]);
      }}
    />
  );
}
