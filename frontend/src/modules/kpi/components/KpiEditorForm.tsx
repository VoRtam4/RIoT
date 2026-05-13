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
  kpiId?: string;
};

type Option = {
  value: string;
  label: string;
  searchText?: string;
};

export default function KpiEditorForm({ kpiId }: Props) {
  const navigate = useNavigate();
  const isEdit = !!kpiId;

  const { sdTypes, loading: sdTypesLoading } = useSdTypes();
  const { createKpi } = useCreateKpi();
  const { updateKpi } = useUpdateKpi();
  const { deleteKpi } = useDeleteKpi();
  const { kpi, loading: kpiLoading } = useKpiDefinition(kpiId);

  const [label, setLabel] = useState("");
  const [sdTypeID, setSdTypeID] = useState<string | null>(null);
  const [mode, setMode] = useState<"all" | "selected">("all");
  const [selectedInstanceIds, setSelectedInstanceIds] = useState<string[]>([]);
  const [query, setQuery] = useState<any>(null);
  const [initialQuery, setInitialQuery] = useState<any>(null);

  const selectedType = useMemo(
    () => sdTypes.find((t) => String(t.id) === String(sdTypeID)),
    [sdTypes, sdTypeID],
  );

  const { entry: instancesEntry, loading: instancesLoading } =
    useSdInstancesByType(sdTypeID);

  const parameters = useMemo(() => {
    if (!selectedType) return [];

    return selectedType.parameters.map((p) => ({
      id: String(p.id),
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
        value: String(t.id),
        label: t.label ?? "",
        searchText: buildOptionSearchText(t.label, t.uid, String(t.id)),
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
        value: String(i.id),
        label: i.label || i.uid,
        searchText: buildOptionSearchText(i.label, i.uid, String(i.id)),
      })),
    [instancesEntry],
  );

  const selectedInstanceOptions = useMemo(() => {
    return selectedInstanceIds.map((id) => {
      const found = instanceOptions.find((o) => o.value === id);
      return found ?? { value: id, label: id };
    });
  }, [selectedInstanceIds, instanceOptions]);

  useEffect(() => {
    if (!kpi) return;
    if (!sdTypes.length) return;

    const nextSdTypeID = String(kpi.sdTypeID);
    const type = sdTypes.find((t) => String(t.id) === nextSdTypeID);
    if (!type) return;

    setLabel(kpi.label ?? "");
    setSdTypeID(nextSdTypeID);

    const backendMode = String(kpi.sdInstanceMode).toLowerCase();
    setMode(backendMode === "selected" ? "selected" : "all");

    setSelectedInstanceIds(
      (kpi.selectedSDInstanceIDs ?? []).map((id) => String(id)),
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
    if (sdTypeID) return;

    setSdTypeID(String(sdTypes[0].id));
  }, [sdTypes, isEdit, sdTypeID]);

  const hasRules = (q: any): boolean => {
    if (!q || !Array.isArray(q.rules)) return false;

    return q.rules.some((r: any) => {
      if (r?.rules) return hasRules(r);
      return !!r?.field;
    });
  };

  if (sdTypesLoading || kpiLoading || (sdTypeID && instancesLoading)) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" />
      </div>
    );
  }

  const handleSubmit = async () => {
    if (!sdTypeID || !selectedType) {
      toast.error("Select an SD type");
      return;
    }

    const finalQuery = query && hasRules(query) ? query : initialQuery;

    if (!finalQuery || !hasRules(finalQuery)) {
      toast.error("Fill in the KPI conditions");
      return;
    }

    if (mode === "selected" && selectedInstanceIds.length === 0) {
      toast.error("Select at least one instance");
      return;
    }

    const nodes = mapQueryToKPINodes(finalQuery, selectedType.parameters);

    const input = {
      label,
      sdTypeID,
      sdTypeUID: selectedType.uid,
      userIdentifier: "",
      nodes,
      sdInstanceMode: mode,
      selectedSDInstanceIDs: mode === "selected" ? selectedInstanceIds : [],
    };

    try {
      let result;

      if (isEdit && kpiId) {
        result = await updateKpi(kpiId, kpi?.sdTypeID ?? "0", input);
        toast.success("KPI updated");
      } else {
        result = await createKpi(input);
        toast.success("KPI created");
      }

      if (!result?.id) throw new Error("Missing ID");

      navigate(`/kpi/${result.id}`);
    } catch (e) {
      console.error(e);
      toast.error("Failed to save KPI");
    }
  };

  const handleDelete = async () => {
    if (!kpiId) return;

    const confirmed = window.confirm("Do you really want to delete KPI?");
    if (!confirmed) return;

    try {
      const ok = await deleteKpi(kpiId, kpi?.sdTypeID ?? "0");

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
              sdTypeOptions.find((o) => o.value === String(sdTypeID)) ?? null
            }
            onChange={(v) => {
              const value = v ? v.value : null;
              setSdTypeID(value);
              setSelectedInstanceIds([]);
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
                setSelectedInstanceIds(v.map((item) => item.value));
              }}
              closeMenuOnSelect={false}
              filterOption={filterSelectOption}
              {...virtualizedSelectProps}
            />
          </div>
        )}
      </div>

      {sdTypeID && parameters.length > 0 && (
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
          <button
            className="btn btn-outline-danger"
            onClick={handleDelete}
          >
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
