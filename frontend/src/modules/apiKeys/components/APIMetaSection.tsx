/**
 * @file APIMetaSection.tsx
 * @brief Meta informace API dokumentace, například oprávnění, metody a cílové adresy.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
type Props = {
  title: string;
  value: string;
};

export default function MetaSection({ title, value }: Props) {
  return (
    <div className="api-docs-meta-card">
      <div className="form-label mb-2">{title}</div>

      <code className="api-docs-inline-code">{value}</code>
    </div>
  );
}
