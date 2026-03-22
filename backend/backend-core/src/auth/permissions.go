package auth

const (
	RoleAdmin = "Admin"
	RoleUser  = "User"
)

const (
	OperationRead   = "read"
	OperationCreate = "create"
	OperationUpdate = "update"
	OperationDelete = "delete"

	OperationSubscribe = "subscribe"
)

const (
	ResourceSDTypes        = "sd_types"
	ResourceSDInstances    = "sd_instances"
	ResourceKPIDefinitions = "kpi_definitions"
	ResourceKPIResults     = "kpi_results"
	ResourceUserConfig     = "user_config"

	ResourceEvents = "events"

	ResourceStatistics = "statistics"

	ResourceTimeSeries = "time_series"

	ResourceAPIKeys = "api_keys"
)

var RolePermissions = map[string]map[string]bool{
	RoleAdmin: {
		ResourceSDTypes + "." + OperationRead:   true,
		ResourceSDTypes + "." + OperationCreate: true,
		ResourceSDTypes + "." + OperationDelete: true,

		ResourceSDInstances + "." + OperationRead:   true,
		ResourceSDInstances + "." + OperationUpdate: true,

		ResourceKPIDefinitions + "." + OperationRead:   true,
		ResourceKPIDefinitions + "." + OperationCreate: true,
		ResourceKPIDefinitions + "." + OperationUpdate: true,
		ResourceKPIDefinitions + "." + OperationDelete: true,

		ResourceKPIResults + "." + OperationRead: true,

		ResourceUserConfig + "." + OperationRead:   true,
		ResourceUserConfig + "." + OperationUpdate: true,
		ResourceUserConfig + "." + OperationDelete: true,

		ResourceEvents + "." + OperationSubscribe: true,

		ResourceTimeSeries + "." + OperationRead: true,

		ResourceAPIKeys + "." + OperationRead:   true,
		ResourceAPIKeys + "." + OperationCreate: true,
		ResourceAPIKeys + "." + OperationUpdate: true,
		ResourceAPIKeys + "." + OperationDelete: true,
	},

	RoleUser: {
		ResourceSDTypes + "." + OperationRead: true,

		ResourceSDInstances + "." + OperationRead: true,

		ResourceKPIDefinitions + "." + OperationRead: true,
		ResourceKPIResults + "." + OperationRead:     true,

		ResourceUserConfig + "." + OperationRead:   true,
		ResourceUserConfig + "." + OperationUpdate: true,

		ResourceEvents + "." + OperationSubscribe: true,

		ResourceTimeSeries + "." + OperationRead: true,

		ResourceAPIKeys + "." + OperationRead:   true,
		ResourceAPIKeys + "." + OperationCreate: true,
		ResourceAPIKeys + "." + OperationUpdate: true,
		ResourceAPIKeys + "." + OperationDelete: true,
	},
}

func GetAllRoleLabels() []string {
	labels := make([]string, 0, len(RolePermissions))
	for label := range RolePermissions {
		labels = append(labels, label)
	}
	return labels
}
