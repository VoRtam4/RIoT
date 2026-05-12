import toast from "react-hot-toast";

type Props = {
  value: string;
  onNew: () => void;
};

export default function APIKeyCreatedBox({ value, onNew }: Props) {
  return (
    <div className="d-flex flex-column h-100">
      <div className="d-flex flex-column justify-content-center h-100">
        <div className="w-100 mx-auto" style={{ maxWidth: 720 }}>
          <div className="mb-2 form-label text-start">Generated key</div>

          <div className="d-flex align-items-center gap-3">
            <code
              className="api-docs-inline-code flex-grow-1"
              style={{
                wordBreak: "break-all",
                whiteSpace: "pre-wrap",
                fontSize: "14px",
                lineHeight: 1.45,
              }}
            >
              {value}
            </code>

            <button
              className="btn btn-outline-light flex-shrink-0"
              onClick={async () => {
                await navigator.clipboard.writeText(value);
                toast.success("Zkopírováno");
              }}
            >
              Copy
            </button>
          </div>
        </div>
      </div>

      <div className="d-flex justify-content-end pt-3">
        <button className="btn btn-primary" onClick={onNew}>
          New
        </button>
      </div>
    </div>
  );
}
