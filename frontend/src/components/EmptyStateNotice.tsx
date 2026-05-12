type Props = {
  title: string;
  description?: string;
  compact?: boolean;
};

export default function EmptyStateNotice({
  title,
  description,
  compact = false,
}: Props) {
  return (
    <div
      className={`empty-state-notice${compact ? " is-compact" : ""}`}
      role="status"
      aria-live="polite"
    >
      <div className="empty-state-notice__title">{title}</div>
      {description ? (
        <div className="empty-state-notice__description">{description}</div>
      ) : null}
    </div>
  );
}
