type Props = {
  title: string;
  code: string;
};

export default function CodeSection({ title, code }: Props) {
  return (
    <div className="mb-3">
      <div className="form-label mb-2">{title}</div>

      <pre className="api-docs-code-block">
        <code>{code}</code>
      </pre>
    </div>
  );
}