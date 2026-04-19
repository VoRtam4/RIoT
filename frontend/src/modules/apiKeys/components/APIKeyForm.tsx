import { useEffect, useState } from "react";
import Select from "react-select";
import APIKeyIPInput from "./APIKeyIPInput";
import { usePermissions } from "../../auth/hooks/usePermission";

type APIKeyFormData = {
  label: string;
  expiresAt?: string;
  permissions: string[];
  ipRestrictions: string[];
  revoked?: boolean;
};

type Props = {
  mode: "create" | "detail";
  apiKey?: {
    label: string;
    expiresAt?: string | null;
    permissions: string[];
    ipRestrictions: string[];
    revoked: boolean;
  };
  onSubmit: (data: APIKeyFormData) => void;
  onDelete?: () => void;
};

export default function APIKeyForm({
  mode,
  apiKey,
  onSubmit,
  onDelete,
}: Props) {
  const { permissions } = usePermissions();

  const [form, setForm] = useState<APIKeyFormData>({
    label: "",
    expiresAt: "",
    permissions: [],
    ipRestrictions: [],
    revoked: false,
  });

  useEffect(() => {
    if (apiKey) {
      setForm({
        label: apiKey.label ?? "",
        expiresAt: apiKey.expiresAt
          ? String(apiKey.expiresAt).slice(0, 10)
          : "",
        permissions: apiKey.permissions ?? [],
        ipRestrictions: apiKey.ipRestrictions ?? [],
        revoked: apiKey.revoked ?? false,
      });
    }
  }, [apiKey]);

  const permissionOptions = permissions.map((p: any) => ({
    label: p.label,
    value: p.uid,
  }));

  return (
    <div className="d-flex flex-column">
      {/* FORM */}
      <div className="d-flex flex-column justify-content-center flex-grow-1">
        <div className="mb-3">
          <label className="form-label">Name</label>
          <input
            className="form-control"
            value={form.label}
            onChange={(e) => setForm({ ...form, label: e.target.value })}
          />
        </div>

        <div className="mb-3">
          <label className="form-label">Expiration</label>
          <input
            type="date"
            className="form-control"
            value={form.expiresAt || ""}
            onChange={(e) => setForm({ ...form, expiresAt: e.target.value })}
          />
        </div>

        <div className="mb-3">
          <label className="form-label">Permissions</label>
          <Select
            classNamePrefix="react-select"
            isMulti
            closeMenuOnSelect={false}
            options={permissionOptions}
            value={permissionOptions.filter((o) =>
              form.permissions.includes(o.value),
            )}
            onChange={(items) =>
              setForm({
                ...form,
                permissions: (items ?? []).map((i) => i.value),
              })
            }
          />
        </div>

        <div className="mb-3">
          <label className="form-label">IPs</label>
          <APIKeyIPInput
            value={form.ipRestrictions}
            onChange={(v) => setForm({ ...form, ipRestrictions: v })}
          />
        </div>

        {mode === "detail" && (
          <div className="form-check mt-2">
            <input
              type="checkbox"
              className="form-check-input"
              checked={!!form.revoked}
              onChange={(e) => setForm({ ...form, revoked: e.target.checked })}
            />
            <label className="form-check-label">Revoked</label>
          </div>
        )}
      </div>

      {/* BUTTONS */}
      <div className="d-flex justify-content-end gap-2 mt-3">
        {mode === "detail" && (
          <button className="btn btn-outline-danger" onClick={onDelete}>
            Smazat
          </button>
        )}

        <button className="btn btn-primary" onClick={() => onSubmit(form)}>
          {mode === "create" ? "Create" : "Save"}
        </button>
      </div>
    </div>
  );
}
