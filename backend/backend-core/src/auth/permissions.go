/**
 * @file permissions.go
 * @brief Definice rolí, zdrojů, operací a výchozích oprávnění pro autorizaci Backend Core.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package auth

import "strings"

const (
	RoleAdmin = "Admin"
	RoleUser  = "User"
	RoleGuest = "Guest"
)

const (
	OperationRead      = "read"
	OperationCreate    = "create"
	OperationUpdate    = "update"
	OperationDelete    = "delete"
	OperationSubscribe = "subscribe"
)

const (
	ResourceSDTypes        = "sd_types"
	ResourceSDInstances    = "sd_instances"
	ResourceKPIDefinitions = "kpi_definitions"
	ResourceKPIResults     = "kpi_results"
	ResourceUserConfig     = "user_config"

	ResourceRawData = "raw_data"

	ResourceStatistics = "statistics"

	ResourceTimeSeries = "time_series"

	ResourceAPIKeys = "api_keys"
	ResourceRoles   = "roles"
)

var RolePermissions = map[string]map[string]bool{
	RoleAdmin: {
		ResourceSDTypes + "." + OperationRead:      true,
		ResourceSDTypes + "." + OperationCreate:    true,
		ResourceSDTypes + "." + OperationDelete:    true,
		ResourceSDTypes + "." + OperationSubscribe: true,

		ResourceSDInstances + "." + OperationRead:      true,
		ResourceSDInstances + "." + OperationSubscribe: true,
		ResourceSDInstances + "." + OperationUpdate:    true,

		ResourceKPIDefinitions + "." + OperationRead:   true,
		ResourceKPIDefinitions + "." + OperationCreate: true,
		ResourceKPIDefinitions + "." + OperationUpdate: true,
		ResourceKPIDefinitions + "." + OperationDelete: true,

		ResourceRawData + "." + OperationRead:      true,
		ResourceRawData + "." + OperationSubscribe: true,

		ResourceKPIResults + "." + OperationRead:      true,
		ResourceKPIResults + "." + OperationSubscribe: true,

		ResourceUserConfig + "." + OperationRead:   true,
		ResourceUserConfig + "." + OperationUpdate: true,
		ResourceUserConfig + "." + OperationDelete: true,

		ResourceTimeSeries + "." + OperationRead:      true,
		ResourceTimeSeries + "." + OperationSubscribe: true,

		ResourceAPIKeys + "." + OperationRead:   true,
		ResourceAPIKeys + "." + OperationCreate: true,
		ResourceAPIKeys + "." + OperationUpdate: true,
		ResourceAPIKeys + "." + OperationDelete: true,

		ResourceRoles + "." + OperationRead:   true,
		ResourceRoles + "." + OperationUpdate: true,
	},

	RoleUser: {
		ResourceSDTypes + "." + OperationRead:      true,
		ResourceSDTypes + "." + OperationSubscribe: true,

		ResourceSDInstances + "." + OperationRead:      true,
		ResourceSDInstances + "." + OperationSubscribe: true,

		ResourceKPIDefinitions + "." + OperationRead: true,

		ResourceRawData + "." + OperationRead:      true,
		ResourceRawData + "." + OperationSubscribe: true,

		ResourceKPIResults + "." + OperationRead:      true,
		ResourceKPIResults + "." + OperationSubscribe: true,

		ResourceUserConfig + "." + OperationRead:   true,
		ResourceUserConfig + "." + OperationUpdate: true,

		ResourceTimeSeries + "." + OperationRead:      true,
		ResourceTimeSeries + "." + OperationSubscribe: true,

		ResourceAPIKeys + "." + OperationRead:   true,
		ResourceAPIKeys + "." + OperationCreate: true,
		ResourceAPIKeys + "." + OperationUpdate: true,
		ResourceAPIKeys + "." + OperationDelete: true,

		ResourceRoles + "." + OperationRead: true,
	},
	RoleGuest: {},
}

func GetAllRoleUIDs() []string {
	uids := make([]string, 0, len(RolePermissions))
	for uid := range RolePermissions {
		uids = append(uids, uid)
	}
	return uids
}

func FormatPermissionLabel(permission string) string {
	parts := strings.Split(permission, ".")
	if len(parts) != 2 {
		return permission
	}
	result := parts[1] + " " + parts[0]
	result = strings.ReplaceAll(result, "_", " ")
	words := strings.Split(result, " ")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
