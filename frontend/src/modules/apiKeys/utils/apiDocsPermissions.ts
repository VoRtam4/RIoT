/**
 * @file apiDocsPermissions.ts
 * @brief Mapování API dokumentace na oprávnění vyžadovaná jednotlivými operacemi.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
const featurePermissionResource: Record<string, string> = {
  "sd-types": "sd_types",
  "sd-instances": "sd_instances",
  "raw-data-points": "raw_data",
  "kpi-definitions": "kpi_definitions",
  "kpi-results": "kpi_results",
  "device-groups": "sd_instances",
  "time-series": "time_series",
  "api-keys": "api_keys",
  "user-config": "user_config",
  roles: "roles",
};

export function getRequiredPermissionLabel(featureId: string, actionId: string) {
  const resource = featurePermissionResource[featureId];

  if (!resource) {
    return "Unknown Permission";
  }

  return formatPermissionLabel(`${resource}.${inferPermissionOperation(actionId)}`);
}

function inferPermissionOperation(actionId: string) {
  if (actionId.includes("stream")) {
    return "subscribe";
  }

  if (actionId.startsWith("create-")) {
    return "create";
  }

  if (actionId.startsWith("update-") || actionId === "assign-role") {
    return "update";
  }

  if (actionId.startsWith("delete-")) {
    return "delete";
  }

  return "read";
}

function formatPermissionLabel(permission: string) {
  const parts = permission.split(".");

  if (parts.length !== 2) {
    return permission;
  }

  const result = `${parts[1]} ${parts[0]}`.replaceAll("_", " ");

  return result
    .split(" ")
    .filter(Boolean)
    .map((word) => word[0].toUpperCase() + word.slice(1))
    .join(" ");
}