type Props = {
  value: string;
};

export default function APIKeyCreatedBox({ value }: Props) {
  return (
    <div className="d-flex flex-column justify-content-center align-items-center h-100">
      <div className="card p-3 text-center" style={{ maxWidth: 500 }}>
        <div className="mb-3 text-white">Generated key</div>
        <code style={{ wordBreak: "break-all" }}>{value}</code>
        <button
          className="btn btn-primary mt-3"
          onClick={() => navigator.clipboard.writeText(value)}
        >
          Copy
        </button>
      </div>
    </div>
  );
}
