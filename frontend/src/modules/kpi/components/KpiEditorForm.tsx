/**
 * @file KpiEditorForm.tsx
 * @brief Hlavní formulář editoru KPI včetně výběru typu, instance a podmínek.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useState, useMemo, useEffect } from "react";
import Select, { type MultiValue } from "react-select";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";

import { useSdTypes } from "../../sdTypes/hooks/useSdTypes";
import { useSdInstancesByType } from "../../sdInstances/hooks/useSdInstancesByType";

import { useCreateKpi } from "../../kpi/hooks/useCreateKpi";
import { useUpdateKpi } from "../../kpi/hooks/useUpdateKpi";
import { useDeleteKpi } from "../../kpi/hooks/useDeleteKpi";
import { useKpiDefinition } from "../hooks/useKpiDefinition";

import KpiQueryBuilder from "./KpiQueryBuilder";
import { mapQueryToKPINodes } from "../utils/mapQueryToKPINodes";
import { mapKPINodesToQuery } from "../utils/mapKPINodesToQuery";
import {
  buildOptionSearchText,
  filterSelectOption,
} from "../../../utils/reactSelectSearch";
import { virtualizedSelectProps } from "../../../utils/reactSelectVirtualized";

type Props = {
  kpiUID?: string;
};

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

function extractScopedUIDSuffix(
  uid: string | null | undefined,
  scopeUID: string | null | undefined,
  childPrefix: string,
) {
  const trimmedUID = String(uid ?? "").trim();
  const trimmedScopeUID = String(scopeUID ?? "").trim();
  const prefix = `${trimmedScopeUID}.${childPrefix}:`;
  if (trimmedScopeUID && trimmedUID.startsWith(prefix)) {
    return trimmedUID.slice(prefix.length);
  }
  const childOnlyPrefix = `${childPrefix}:`;
  if (trimmedUID.startsWith(childOnlyPrefix)) {
    return trimmedUID.slice(childOnlyPrefix.length);
  }
  return trimmedUID;
}

export default function KpiEditorForm({ kpiUID }: Props) {
  const navigate = useNavigate();
  const isEdit = !!kpiUID;

  const { sdTypes, loading: sdTypesLoading } = useSdTypes();
  const { createKpi } = useCreateKpi();
  const { updateKpi } = useUpdateKpi();
  const { deleteKpi } = useDeleteKpi();
  const { kpi, loading: kpiLoading } = useKpiDefinition(kpiUID);

  const [label, setLabel] = useState("");
  const [uidSuffix, setUidSuffix] = useState("");
  const [sdTypeUID, setSdTypeUID] = useState<string | null>(null);
  const [mode, setMode] = useState<"all" | "selected">("all");
  const [selectedInstanceUIDs, setSelectedInstanceUIDs] = useState<string[]>(
    [],
  );
  const [query, setQuery] = useState<any>(null);
  const [initialQuery, setInitialQuery] = useState<any>(null);

  const selectedType = useMemo(
    () => sdTypes.find((t) => String(t.uid) === String(sdTypeUID)),
    [sdTypes, sdTypeUID],
  );

  const { entry: instancesEntry, loading: instancesLoading } =
    useSdInstancesByType(selectedType?.uid ?? null);

  const parameters = useMemo(() => {
    if (!selectedType) return [];

    return selectedType.parameters.map((p) => ({
      id: String(p.denotation),
      denotation: p.denotation,
      label: p.label,
      role: String(p.role).toLowerCase() === "tag" ? "tag" : "field",
      type:
        String(p.type).toLowerCase() === "number"
          ? "number"
          : String(p.type).toLowerCase() === "boolean"
            ? "boolean"
            : "string",
    })) as {
      id: string;
      denotation: string;
      label?: string;
      role?: "field" | "tag";
      type: "string" | "number" | "boolean";
    }[];
  }, [selectedType]);

  const sdTypeOptions: Option[] = useMemo(
    () =>
      sdTypes.map((t) => ({
        value: String(t.uid),
        label: t.label || t.uid || "",
        searchText: buildOptionSearchText(t.label, t.uid),
      })),
    [sdTypes],
  );

  const modeOptions: Option[] = [
    { value: "all", label: "ALL" },
    { value: "selected", label: "SELECTED" },
  ];

  const instanceOptions: Option[] = useMemo(
    () =>
      (instancesEntry?.rawSortedAsc ?? []).map((i) => ({
        value: String(i.uid),
        label: i.label || i.uid,
        searchText: buildOptionSearchText(i.label, i.uid),
      })),
    [instancesEntry],
  );

  const selectedInstanceOptions = useMemo(() => {
    return selectedInstanceUIDs.map((uid) => {
      const found = instanceOptions.find((o) => o.value === uid);
      return found ?? { value: uid, label: uid };
    });
  }, [selectedInstanceUIDs, instanceOptions]);

  useEffect(() => {
    if (!kpi) return;
    if (!sdTypes.length) return;

    const type = sdTypes.find((t) => String(t.uid) === String(kpi.sdTypeUID));
    if (!type) return;
    const nextSdTypeUID = String(type.uid);

    setLabel(kpi.label ?? "");
    setUidSuffix(extractScopedUIDSuffix(kpi.uid, nextSdTypeUID, "kpi"));
    setSdTypeUID(nextSdTypeUID);

    const backendMode = String(kpi.sdInstanceMode).toLowerCase();
    setMode(backendMode === "selected" ? "selected" : "all");

    setSelectedInstanceUIDs(
      (kpi.selectedSDInstanceUIDs ?? []).map((uid) => String(uid)),
    );

    const mappedQuery = mapKPINodesToQuery(
      kpi.nodes ?? [],
      type.parameters ?? [],
    );
    setInitialQuery(mappedQuery);
    setQuery(mappedQuery);
  }, [kpi, sdTypes]);

  useEffect(() => {
    if (isEdit) return;
    if (!sdTypes.length) return;
    if (sdTypeUID) return;

    setSdTypeUID(String(sdTypes[0].uid));
  }, [sdTypes, isEdit, sdTypeUID]);

  const hasRules = (q: any): boolean => {
    if (!q || !Array.isArray(q.rules)) return false;

    return q.rules.some((r: any) => {
      if (r?.rules) return hasRules(r);
      return !!r?.field;
    });
  };

  if (sdTypesLoading || kpiLoading || (sdTypeUID && instancesLoading)) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  const handleSubmit = async () => {
    if (!sdTypeUID || !selectedType) {
      toast.error("Select an SD type");
      return;
    }

    const finalQuery = query && hasRules(query) ? query : initialQuery;

    if (!finalQuery || !hasRules(finalQuery)) {
      toast.error("Fill in the KPI conditions");
      return;
    }

    if (uidSuffix.includes(".") || uidSuffix.includes(":")) {
      toast.error("KPI UID suffix cannot contain dot or colon");
      return;
    }

    if (mode === "selected" && selectedInstanceUIDs.length === 0) {
      toast.error("Select at least one instance");
      return;
    }

    const nodes = mapQueryToKPINodes(finalQuery, selectedType.parameters);

    const input = {
      uid: uidSuffix.trim() || undefined,
      label,
      sdTypeUID: selectedType.uid,
      userIdentifier: "",
      nodes,
      sdInstanceMode: mode,
      selectedSDInstanceUIDs: mode === "selected" ? selectedInstanceUIDs : [],
    };

    try {
      let result;

      if (isEdit && kpiUID) {
        result = await updateKpi(kpiUID, kpi?.sdTypeUID ?? "", input);
        toast.success("KPI updated");
      } else {
        result = await createKpi(input);
        toast.success("KPI created");
      }

      if (!result?.uid) throw new Error("Missing UID");

      navigate(`/kpi/${result.uid}`);
    } catch (e) {
      console.error(e);
      toast.error("Failed to save KPI");
    }
  };

  const handleDelete = async () => {
    if (!kpiUID) return;

    const confirmed = window.confirm("Do you really want to delete KPI?");
    if (!confirmed) return;

    try {
      const ok = await deleteKpi(kpiUID, kpi?.sdTypeUID ?? "");

      if (!ok) throw new Error("Delete failed");

      toast.success("KPI deleted");
      navigate("/kpi");
    } catch (e) {
      console.error(e);
      toast.error("Delete failed");
    }
  };

  return (
    <div className="card p-3">
      <div className="row g-3">
        <div className="col-md-6">
          <label className="form-label">Name</label>
          <input
            className="form-control"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
          />
        </div>

        <div className="col-md-6">
          <label className="form-label">Model</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={sdTypeOptions}
            value={
              sdTypeOptions.find((o) => o.value === String(sdTypeUID)) ?? null
            }
            onChange={(v) => {
              const value = v ? v.value : null;
              setSdTypeUID(value);
              setUidSuffix((current) =>
                extractScopedUIDSuffix(current, sdTypeUID, "kpi"),
              );
              setSelectedInstanceUIDs([]);
              setQuery(null);
              setInitialQuery(null);
            }}
            isClearable={false}
            isDisabled={isEdit}
            filterOption={filterSelectOption}
            {...virtualizedSelectProps}
          />
        </div>

        <div className="col-md-6">
          <label className="form-label">UID</label>
          <div className="input-group">
            <span className="input-group-text">
              {sdTypeUID ? `${sdTypeUID}.kpi:` : "sdt:type.kpi:"}
            </span>
            <input
              className="form-control"
              value={uidSuffix}
              onChange={(e) => setUidSuffix(e.target.value)}
              placeholder="uid-suffix"
            />
          </div>
        </div>

        <div className="col-md-6">
          <label className="form-label">Mode</label>
          <Select<Option, false>
            classNamePrefix="react-select"
            options={modeOptions}
            value={modeOptions.find((o) => o.value === mode) ?? null}
            onChange={(v) => {
              if (!v) return;
              setMode(v.value as "all" | "selected");
            }}
            isClearable={false}
            {...virtualizedSelectProps}
          />
        </div>

        {mode === "selected" && (
          <div className="col-md-6">
            <label className="form-label">Devices</label>
            <Select<Option, true>
              classNamePrefix="react-select"
              isMulti
              options={instanceOptions}
              value={selectedInstanceOptions}
              onChange={(v: MultiValue<Option>) => {
                setSelectedInstanceUIDs(v.map((item) => item.value));
              }}
              closeMenuOnSelect={false}
              filterOption={filterSelectOption}
              {...virtualizedSelectProps}
            />
          </div>
        )}
      </div>

      {sdTypeUID && parameters.length > 0 && (
        <div className="mt-4">
          <label className="form-label">Rules</label>
          <KpiQueryBuilder
            parameters={parameters}
            onChange={setQuery}
            initialQuery={initialQuery}
          />
        </div>
      )}

      <div className="d-flex justify-content-end gap-2 mt-4">
        {isEdit && (
          <button className="btn btn-outline-danger" onClick={handleDelete}>
            Delete
          </button>
        )}

        <button className="btn btn-primary" onClick={handleSubmit}>
          {isEdit ? "Update" : "Create"}
        </button>
      </div>
    </div>
  );
}
