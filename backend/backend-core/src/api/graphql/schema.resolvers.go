package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/graphql/gsc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func (r *mutationResolver) CreateSDType(ctx context.Context, input graphQLModel.SDTypeInput) (graphQLModel.SDType, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationCreate); err != nil {
		return graphQLModel.SDType{}, err
	}
	createSDTypeResult := domainLogicLayer.CreateSDType(input)
	if createSDTypeResult.IsFailure() {
		log.Printf("Error occurred (create SD type): %s\n", createSDTypeResult.GetError().Error())
	}
	return createSDTypeResult.Unwrap()
}

func (r *mutationResolver) DeleteSDType(ctx context.Context, uid string) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteSDType(uid); err != nil {
		log.Printf("Error occurred (delete SD type): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) UpdateSDInstance(ctx context.Context, uid string, input graphQLModel.SDInstanceUpdateInput) (graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
		return graphQLModel.SDInstance{}, err
	}
	updateSDInstanceResult := domainLogicLayer.UpdateSDInstance(uid, input)
	if updateSDInstanceResult.IsFailure() {
		log.Printf("Error occurred (update SD instance): %s\n", updateSDInstanceResult.GetError().Error())
	}
	return updateSDInstanceResult.Unwrap()
}

func (r *mutationResolver) CreateKPIDefinition(ctx context.Context, input graphQLModel.KPIDefinitionInput) (graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationCreate)
	if err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	createKPIDefinitionResult := domainLogicLayer.CreateKPIDefinition(principal.UserID, input)
	if createKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (create KPI definition): %s\n", createKPIDefinitionResult.GetError().Error())
	}
	return createKPIDefinitionResult.Unwrap()
}

func (r *mutationResolver) UpdateKPIDefinition(ctx context.Context, uid string, input graphQLModel.KPIDefinitionInput) (graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	updateKPIDefinitionResult := domainLogicLayer.UpdateKPIDefinition(principal.UserID, uid, input)
	if updateKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (update KPI definition): %s\n", updateKPIDefinitionResult.GetError().Error())
	}
	return updateKPIDefinitionResult.Unwrap()
}

func (r *mutationResolver) DeleteKPIDefinition(ctx context.Context, uid string) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationDelete)
	if err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteKPIDefinition(principal.UserID, uid); err != nil {
		log.Printf("Error occurred (delete KPI definition): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateSDInstanceGroup(ctx context.Context, input graphQLModel.SDInstanceGroupInput) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationCreate); err != nil {
		return graphQLModel.SDInstanceGroup{}, err
	}
	createSDInstanceGroupResult := domainLogicLayer.CreateSDInstanceGroup(input)
	if createSDInstanceGroupResult.IsFailure() {
		log.Printf("Error occurred (create SD instance group): %s\n", createSDInstanceGroupResult.GetError().Error())
	}
	return createSDInstanceGroupResult.Unwrap()
}

func (r *mutationResolver) UpdateSDInstanceGroup(ctx context.Context, uid string, input graphQLModel.SDInstanceGroupInput) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
		return graphQLModel.SDInstanceGroup{}, err
	}
	updateSDInstanceGroupResult := domainLogicLayer.UpdateSDInstanceGroup(uid, input)
	if updateSDInstanceGroupResult.IsFailure() {
		log.Printf("Error occurred (update SD instance group): %s\n", updateSDInstanceGroupResult.GetError().Error())
	}
	return updateSDInstanceGroupResult.Unwrap()
}

func (r *mutationResolver) DeleteSDInstanceGroup(ctx context.Context, uid string) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteSDInstanceGroup(uid); err != nil {
		log.Printf("Error occurred (delete SD instance group): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) StatisticsMutate(ctx context.Context, inputData graphQLModel.InputData) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceStatistics, auth.OperationRead); err != nil {
		return false, err
	}
	return domainLogicLayer.Save(inputData).Unwrap()
}

func (r *mutationResolver) UpdateUserConfig(ctx context.Context, input graphQLModel.UserConfigInput) (graphQLModel.UserConfig, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.UserConfig{}, err
	}
	updateUserConfigResult := domainLogicLayer.UpdateUserConfig(principal.UserID, input)
	if updateUserConfigResult.IsFailure() {
		log.Printf("Error occurred (update user config): %s\n", updateUserConfigResult.GetError().Error())
	}
	return updateUserConfigResult.Unwrap()
}

func (r *mutationResolver) DeleteUserConfig(ctx context.Context) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationDelete)
	if err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteUserConfig(principal.UserID); err != nil {
		log.Printf("Error occurred (delete user configuration): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, uid string, input graphQLModel.UserUpdateInput) (graphQLModel.User, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationUpdate); err != nil {
		return graphQLModel.User{}, err
	}
	result := domainLogicLayer.UpdateUser(uid, input)
	if result.IsFailure() {
		return graphQLModel.User{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) DisableUser(ctx context.Context, uid string, reason *string) (graphQLModel.User, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.User{}, err
	}
	result := domainLogicLayer.DisableUser(principal.UserID, uid, reason)
	if result.IsFailure() {
		return graphQLModel.User{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) EnableUser(ctx context.Context, uid string) (graphQLModel.User, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationUpdate); err != nil {
		return graphQLModel.User{}, err
	}
	result := domainLogicLayer.EnableUser(uid)
	if result.IsFailure() {
		return graphQLModel.User{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) RevokeUserSessions(ctx context.Context, uid string) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationUpdate); err != nil {
		return false, err
	}
	if err := domainLogicLayer.RevokeUserSessions(uid); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RevokeSession(ctx context.Context, uid string) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceSessions, auth.OperationUpdate)
	if err != nil {
		return false, err
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	if err := domainLogicLayer.RevokeSession(principal.UserID, uid, allowForeign); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RevokeAllSessionsForUser(ctx context.Context, userUID string) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationUpdate); err != nil {
		return false, err
	}
	if err := domainLogicLayer.RevokeUserSessions(userUID); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RevokeOwnOtherSessions(ctx context.Context) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceSessions, auth.OperationUpdate)
	if err != nil {
		return false, err
	}
	if principal.SessionID == nil {
		return false, fmt.Errorf("current session not available")
	}
	if err := domainLogicLayer.RevokeOwnOtherSessions(principal.UserID, *principal.SessionID); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateRole(ctx context.Context, input graphQLModel.RoleInput) (graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationCreate); err != nil {
		return graphQLModel.Role{}, err
	}
	result := domainLogicLayer.CreateRole(input)
	if result.IsFailure() {
		return graphQLModel.Role{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdateRole(ctx context.Context, uid string, input graphQLModel.RoleInput) (graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationUpdate); err != nil {
		return graphQLModel.Role{}, err
	}
	result := domainLogicLayer.UpdateRole(uid, input)
	if result.IsFailure() {
		return graphQLModel.Role{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, uid string) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteRole(uid); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CloneRole(ctx context.Context, uid string, label string) (graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationCreate); err != nil {
		return graphQLModel.Role{}, err
	}
	result := domainLogicLayer.CloneRole(uid, label)
	if result.IsFailure() {
		return graphQLModel.Role{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdatePermissionLabel(ctx context.Context, uid string, label string) (graphQLModel.Permission, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationUpdate); err != nil {
		return graphQLModel.Permission{}, err
	}
	result := domainLogicLayer.UpdatePermissionLabel(uid, label)
	if result.IsFailure() {
		return graphQLModel.Permission{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) AssignRoleToUser(ctx context.Context, input graphQLModel.AssignRoleInput) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationUpdate); err != nil {
		return false, err
	}
	err := domainLogicLayer.AssignRoleToUser(input.UserUID, input.RoleUID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateAPIKey(ctx context.Context, input graphQLModel.APIKeyInput) (string, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationCreate)
	if err != nil {
		return "", err
	}
	result := domainLogicLayer.CreateAPIKey(principal.UserID, input)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdateAPIKey(ctx context.Context, uid string, input graphQLModel.APIKeyInput) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return false, err
	}
	err = domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, uid, input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RevokeAPIKey(ctx context.Context, uid string) (graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.APIKey{}, err
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RevokeAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		return graphQLModel.APIKey{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) RotateAPIKey(ctx context.Context, uid string) (string, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return "", err
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RotateAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdateAPIKeyPermissions(ctx context.Context, uid string, permissionUIDs []string) (graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.APIKey{}, err
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyPermissions(principal.UserID, uid, permissionUIDs, allowForeign)
	if result.IsFailure() {
		return graphQLModel.APIKey{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdateAPIKeyRestrictions(ctx context.Context, uid string, input graphQLModel.APIKeyRestrictionsInput) (graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.APIKey{}, err
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyRestrictions(principal.UserID, uid, input, allowForeign)
	if result.IsFailure() {
		return graphQLModel.APIKey{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) DeleteAPIKey(ctx context.Context, uid string) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationDelete)
	if err != nil {
		return false, err
	}
	err = domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, uid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) StartTimeSeriesExport(ctx context.Context, input graphQLModel.TimeSeriesReadInput) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.StartTimeSeriesExport(principal.UserID, input)
}

func (r *mutationResolver) StartTimeSeriesExportAggregateKpi(ctx context.Context, input graphQLModel.TimeSeriesReadAggregateKPIInput) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.StartTimeSeriesExportAggregateKPI(principal.UserID, input)
}

func (r *mutationResolver) CancelTimeSeriesExport(ctx context.Context, uid string) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.CancelTimeSeriesExport(principal.UserID, uid)
}

func (r *queryResolver) SdType(ctx context.Context, uid string) (graphQLModel.SDType, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationRead); err != nil {
		return graphQLModel.SDType{}, err
	}
	getSDTypeResult := domainLogicLayer.GetSDType(uid)
	if getSDTypeResult.IsFailure() {
		log.Printf("Error occurred (get SD type): %s\n", getSDTypeResult.GetError().Error())
	}
	return getSDTypeResult.Unwrap()
}

func (r *queryResolver) SdTypes(ctx context.Context) ([]graphQLModel.SDType, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationRead); err != nil {
		return nil, err
	}
	getSDTypesResult := domainLogicLayer.GetSDTypes()
	if getSDTypesResult.IsFailure() {
		log.Printf("Error occurred (get SD types): %s\n", getSDTypesResult.GetError().Error())
	}
	return getSDTypesResult.Unwrap()
}

func (r *queryResolver) SdInstance(ctx context.Context, uid string) (graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return graphQLModel.SDInstance{}, err
	}
	getSDInstanceResult := domainLogicLayer.GetSDInstance(uid)
	if getSDInstanceResult.IsFailure() {
		log.Printf("Error occurred (get SD instance): %s\n", getSDInstanceResult.GetError().Error())
	}
	return getSDInstanceResult.Unwrap()
}

func (r *queryResolver) SdInstances(ctx context.Context) ([]graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	getSDInstancesResult := domainLogicLayer.GetSDInstances()
	if getSDInstancesResult.IsFailure() {
		log.Printf("Error occurred (get SD instances): %s\n", getSDInstancesResult.GetError().Error())
	}
	return getSDInstancesResult.Unwrap()
}

func (r *queryResolver) SdInstancesByType(ctx context.Context, uid string) ([]graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetSDInstancesByType(uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get SD instances by type): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) SdInstancesByKpiDefinition(ctx context.Context, uid string) ([]graphQLModel.SDInstance, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetSDInstancesByKpiDefinition(principal.UserID, uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get SD instances by KPI): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiDefinition(ctx context.Context, uid string) (graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	getKPIDefinitionResult := domainLogicLayer.GetKPIDefinition(principal.UserID, uid)
	if getKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (get KPI definition): %s\n", getKPIDefinitionResult.GetError().Error())
	}
	return getKPIDefinitionResult.Unwrap()
}

func (r *queryResolver) KpiDefinitions(ctx context.Context) ([]graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	getKPIDefinitionsResult := domainLogicLayer.GetKPIDefinitions(principal.UserID)
	if getKPIDefinitionsResult.IsFailure() {
		log.Printf("Error occurred (get KPI definitions): %s\n", getKPIDefinitionsResult.GetError().Error())
	}
	return getKPIDefinitionsResult.Unwrap()
}

func (r *queryResolver) KpiDefinitionsBySdType(ctx context.Context, uid string) ([]graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDType(principal.UserID, uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get KPI definitions by SDType): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiDefinitionsBySdInstance(ctx context.Context, uid string) ([]graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDInstance(principal.UserID, uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get KPI definitions by SDInstance): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) RawDataPointsBySDType(ctx context.Context, uid string) ([]graphQLModel.RawDataPoint, error) {
	_, err := authorizeOperation(ctx, auth.ResourceRawData, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetRawDataPointsBySDTypeUID(uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get raw data points by sd type): %s\n", result.GetError().Error())
		return nil, result.GetError()
	}
	return result.Unwrap()
}

func (r *queryResolver) RawDataPoint(ctx context.Context, uid string) (graphQLModel.RawDataPoint, error) {
	_, err := authorizeOperation(ctx, auth.ResourceRawData, auth.OperationRead)
	if err != nil {
		return graphQLModel.RawDataPoint{}, err
	}
	opt, err := domainLogicLayer.GetRawDataPointBySDInstanceUID(uid).Unwrap()
	if err != nil {
		log.Printf("Error occurred (get raw data point): %s\n", err.Error())
		return graphQLModel.RawDataPoint{}, err
	}
	if opt.IsEmpty() {
		return graphQLModel.RawDataPoint{}, fmt.Errorf("raw data point not found")
	}
	return opt.GetPayload(), nil
}

func (r *queryResolver) KpiResults(ctx context.Context) ([]graphQLModel.KPIFulfillmentCheckResult, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIFulfillmentCheckResults(principal.UserID)
	if result.IsFailure() {
		log.Printf("Error occurred (get kpi results): %s\n", result.GetError().Error())
		return nil, result.GetError()
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiResultsByKpi(ctx context.Context, uid string) ([]graphQLModel.KPIFulfillmentCheckResult, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIFulfillmentCheckResultsByKPIUID(principal.UserID, uid)
	if result.IsFailure() {
		log.Printf("Error occurred (get kpi results by kpi definition): %s\n", result.GetError().Error())
		return nil, result.GetError()
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiResult(ctx context.Context, request graphQLModel.KPIFulfillmentCheckResultRequest) (graphQLModel.KPIFulfillmentCheckResult, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead)
	if err != nil {
		return graphQLModel.KPIFulfillmentCheckResult{}, err
	}
	opt, err := domainLogicLayer.GetKPIFulfillmentCheckResult(principal.UserID, request).Unwrap()
	if err != nil {
		log.Printf("Error occurred (get kpi result): %s\n", err.Error())
		return graphQLModel.KPIFulfillmentCheckResult{}, err
	}
	if opt.IsEmpty() {
		return graphQLModel.KPIFulfillmentCheckResult{}, fmt.Errorf("kpi result not found")
	}
	return opt.GetPayload(), nil
}

func (r *queryResolver) SdInstanceGroup(ctx context.Context, uid string) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return graphQLModel.SDInstanceGroup{}, err
	}
	getSDInstanceGroupResult := domainLogicLayer.GetSDInstanceGroupByUID(uid)
	if getSDInstanceGroupResult.IsFailure() {
		log.Printf("Error occurred (get SD instance group): %s\n", getSDInstanceGroupResult.GetError().Error())
	}
	return getSDInstanceGroupResult.Unwrap()
}

func (r *queryResolver) SdInstanceGroups(ctx context.Context) ([]graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	getSDInstanceGroupsResult := domainLogicLayer.GetSDInstanceGroups()
	if getSDInstanceGroupsResult.IsFailure() {
		log.Printf("Error occurred (get SD instance groups): %s\n", getSDInstanceGroupsResult.GetError().Error())
	}
	return getSDInstanceGroupsResult.Unwrap()
}

func (r *queryResolver) StatisticsQuerySimpleSensors(ctx context.Context, request *graphQLModel.StatisticsInput, sensors graphQLModel.SimpleSensors) ([]graphQLModel.OutputData, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceStatistics, auth.OperationRead); err != nil {
		return nil, err
	}
	convertedRequest, err := domainLogicLayer.MapStatisticsInputToReadRequestBody(request, &sensors, nil)
	if err != nil {
		return nil, err
	}
	data := domainLogicLayer.Query(*convertedRequest)
	return data.Unwrap()
}

func (r *queryResolver) StatisticsQuerySensorsWithFields(ctx context.Context, request *graphQLModel.StatisticsInput, sensors graphQLModel.SensorsWithFields) ([]graphQLModel.OutputData, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceStatistics, auth.OperationRead); err != nil {
		return nil, err
	}
	convertedRequest, err := domainLogicLayer.MapStatisticsInputToReadRequestBody(request, nil, &sensors)
	if err != nil {
		return nil, err
	}
	data := domainLogicLayer.Query(*convertedRequest)
	return data.Unwrap()
}

func (r *queryResolver) UserConfig(ctx context.Context) (graphQLModel.UserConfig, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationRead)
	if err != nil {
		return graphQLModel.UserConfig{}, err
	}
	getUserConfigResult := domainLogicLayer.GetUserConfig(principal.UserID)
	if getUserConfigResult.IsFailure() {
		log.Printf("Error occurred (get user config results): %s\n", getUserConfigResult.GetError().Error())
	}
	return getUserConfigResult.Unwrap()
}

func (r *queryResolver) Users(ctx context.Context) ([]graphQLModel.User, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUsers()
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) User(ctx context.Context, uid string) (graphQLModel.User, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationRead); err != nil {
		return graphQLModel.User{}, err
	}
	result := domainLogicLayer.LoadUserByUID(uid)
	if result.IsFailure() {
		return graphQLModel.User{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) Me(ctx context.Context) (graphQLModel.User, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return graphQLModel.User{}, fmt.Errorf("unauthorized")
	}
	result := domainLogicLayer.LoadUserByID(principal.UserID)
	if result.IsFailure() {
		return graphQLModel.User{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) Sessions(ctx context.Context) ([]graphQLModel.UserSession, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceSessions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUserSessions(principal.UserID)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) SessionsByUser(ctx context.Context, userUID string) ([]graphQLModel.UserSession, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUserSessionsByUserUID(userUID)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) Roles(ctx context.Context) ([]graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadRoles()
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) Permissions(ctx context.Context) ([]graphQLModel.Permission, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadPermissions()
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) UserRole(ctx context.Context, uid string) (*graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUserRoleByUID(uid)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	role := result.GetPayload()
	return &role, nil
}

func (r *queryResolver) Role(ctx context.Context, uid *string) (*graphQLModel.Role, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	var result sharedUtils.Result[graphQLModel.Role]
	if uid != nil && *uid != "" {
		result = domainLogicLayer.LoadRoleByUID(*uid)
	} else {
		result = domainLogicLayer.LoadUserRole(principal.UserID)
	}
	if result.IsFailure() {
		return nil, result.GetError()
	}
	role := result.GetPayload()
	return &role, nil
}

func (r *queryResolver) APIKeys(ctx context.Context) ([]graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadAPIKeysForUser(principal.UserID)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) APIKeysByUser(ctx context.Context, userUID string) ([]graphQLModel.APIKey, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUsers, auth.OperationRead); err != nil {
		return nil, err
	}
	if _, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadAPIKeysByUserUID(userUID)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) APIKey(ctx context.Context, uid string) (graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationRead)
	if err != nil {
		return graphQLModel.APIKey{}, err
	}
	result := domainLogicLayer.LoadAPIKeyByUID(principal.UserID, uid)
	if result.IsFailure() {
		return graphQLModel.APIKey{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) TimeSeriesRead(ctx context.Context, request graphQLModel.TimeSeriesReadInput) (graphQLModel.TimeSeriesReadResponse, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesReadResponse{}, err
	}
	result := domainLogicLayer.ReadTimeSeries(principal.UserID, request)
	if result.IsFailure() {
		return graphQLModel.TimeSeriesReadResponse{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) TimeSeriesReadAggregateKpi(ctx context.Context, request graphQLModel.TimeSeriesReadAggregateKPIInput) (graphQLModel.TimeSeriesReadResponse, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesReadResponse{}, err
	}
	result := domainLogicLayer.ReadTimeSeriesAggregateKPI(principal.UserID, request)
	if result.IsFailure() {
		return graphQLModel.TimeSeriesReadResponse{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) TimeSeriesDistinctTagValues(ctx context.Context, request graphQLModel.TimeSeriesDistinctTagValuesInput) (graphQLModel.TimeSeriesDistinctTagValuesResponse, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	result := domainLogicLayer.DistinctTimeSeriesTagValues(principal.UserID, request)
	if result.IsFailure() {
		return graphQLModel.TimeSeriesDistinctTagValuesResponse{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *queryResolver) TimeSeriesExport(ctx context.Context, uid string) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.GetTimeSeriesExport(principal.UserID, uid)
}

func (r *subscriptionResolver) OnSDInstanceRegistered(ctx context.Context, filter *graphQLModel.SDInstanceRegisteredFilter) (<-chan graphQLModel.SDInstance, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationSubscribe)
	if err != nil {
		return nil, err
	}
	sub, err := events.SubscribeSDInstanceRegistered(ctx, filter, principal.UserID)
	if err != nil {
		return nil, err
	}
	return sub.Channel, nil
}

func (r *subscriptionResolver) OnRawDataPointArrived(ctx context.Context, filter *graphQLModel.RawDataPointArrivedFilter) (<-chan []graphQLModel.RawDataPoint, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceRawData, auth.OperationSubscribe)
	if err != nil {
		return nil, err
	}
	sub, err := events.SubscribeRawDataPointArrived(ctx, filter, principal.UserID)
	if err != nil {
		log.Println("SUB ERROR:", err)
		return nil, err
	}
	return sub.Channel, nil
}

func (r *subscriptionResolver) OnKPIFulfillmentChecked(ctx context.Context, filter *graphQLModel.KPIFulfillmentCheckedFilter) (<-chan []graphQLModel.KPIFulfillmentCheckResult, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationSubscribe)
	if err != nil {
		return nil, err
	}
	sub, err := events.SubscribeKPIFulfillmentChecked(ctx, filter, principal.UserID)
	if err != nil {
		log.Println("SUB ERROR:", err)
		return nil, err
	}
	return sub.Channel, nil
}

func (r *subscriptionResolver) OnTimeSeriesExportUpdated(ctx context.Context, filter graphQLModel.TimeSeriesExportFilter) (<-chan graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationSubscribe)
	if err != nil {
		return nil, err
	}
	for _, uid := range filter.Uids {
		if _, err := domainLogicLayer.GetTimeSeriesExport(principal.UserID, uid); err != nil {
			return nil, err
		}
	}
	sub, err := events.SubscribeTimeSeriesExportUpdated(ctx, &filter)
	if err != nil {
		log.Println("SUB ERROR:", err)
		return nil, err
	}
	return sub.Channel, nil
}

func (r *Resolver) Mutation() gsc.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() gsc.QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() gsc.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
