/**
 * @file schema.resolvers.go
 * @brief Implementace resolverů GraphQL API Backend Core.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní resolverová logika pro základní GraphQL API.
 * - Vojtěch Hubáček: doplnění autorizace operací, per-user operací, API klíčů, IP restrikcí, raw dat, time-series rozhraní, rolí, tag/field parametrů, labelů, nových KPI operací, By dotazů a sjednocených subscription událostí pro všechna rozhraní.
 *
 * @ingroup riot_backend_core
 */
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

func (r *mutationResolver) DeleteSDType(ctx context.Context, id uint32) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteSDType(id); err != nil {
		log.Printf("Error occurred (delete SD type): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) UpdateSDInstance(ctx context.Context, id uint32, input graphQLModel.SDInstanceUpdateInput) (graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
		return graphQLModel.SDInstance{}, err
	}
	updateSDInstanceResult := domainLogicLayer.UpdateSDInstance(id, input)
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

func (r *mutationResolver) UpdateKPIDefinition(ctx context.Context, id uint32, input graphQLModel.KPIDefinitionInput) (graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationUpdate)
	if err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	updateKPIDefinitionResult := domainLogicLayer.UpdateKPIDefinition(principal.UserID, id, input)
	if updateKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (update KPI definition): %s\n", updateKPIDefinitionResult.GetError().Error())
	}
	return updateKPIDefinitionResult.Unwrap()
}

func (r *mutationResolver) DeleteKPIDefinition(ctx context.Context, id uint32) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationDelete)
	if err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteKPIDefinition(principal.UserID, id); err != nil {
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

func (r *mutationResolver) UpdateSDInstanceGroup(ctx context.Context, id uint32, input graphQLModel.SDInstanceGroupInput) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
		return graphQLModel.SDInstanceGroup{}, err
	}
	updateSDInstanceGroupResult := domainLogicLayer.UpdateSDInstanceGroup(id, input)
	if updateSDInstanceGroupResult.IsFailure() {
		log.Printf("Error occurred (update SD instance group): %s\n", updateSDInstanceGroupResult.GetError().Error())
	}
	return updateSDInstanceGroupResult.Unwrap()
}

func (r *mutationResolver) DeleteSDInstanceGroup(ctx context.Context, id uint32) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteSDInstanceGroup(id); err != nil {
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

func (r *mutationResolver) AssignRoleToUser(ctx context.Context, input graphQLModel.AssignRoleInput) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationUpdate); err != nil {
		return false, err
	}
	err := domainLogicLayer.AssignRoleToUser(input.UserID, input.RoleID)
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

func (r *mutationResolver) UpdateAPIKey(ctx context.Context, id uint32, input graphQLModel.APIKeyInput) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return false, err
	}
	err = domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, id, input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DeleteAPIKey(ctx context.Context, id uint32) (bool, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationDelete)
	if err != nil {
		return false, err
	}
	err = domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, id)
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

func (r *mutationResolver) CancelTimeSeriesExport(ctx context.Context, id uint32) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.CancelTimeSeriesExport(principal.UserID, id)
}

func (r *queryResolver) SdType(ctx context.Context, id uint32) (graphQLModel.SDType, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDTypes, auth.OperationRead); err != nil {
		return graphQLModel.SDType{}, err
	}
	getSDTypeResult := domainLogicLayer.GetSDType(id)
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

func (r *queryResolver) SdInstance(ctx context.Context, id uint32) (graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return graphQLModel.SDInstance{}, err
	}
	getSDInstanceResult := domainLogicLayer.GetSDInstance(id)
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

func (r *queryResolver) SdInstancesByType(ctx context.Context, id uint32) ([]graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetSDInstancesByType(id)
	if result.IsFailure() {
		log.Printf("Error occurred (get SD instances by type): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) SdInstancesByKpiDefinition(ctx context.Context, id uint32) ([]graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetSDInstancesByKpiDefinition(id)
	if result.IsFailure() {
		log.Printf("Error occurred (get SD instances by KPI): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiDefinition(ctx context.Context, id uint32) (graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	getKPIDefinitionResult := domainLogicLayer.GetKPIDefinition(principal.UserID, id)
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

func (r *queryResolver) KpiDefinitionsBySdType(ctx context.Context, id uint32) ([]graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDType(principal.UserID, id)
	if result.IsFailure() {
		log.Printf("Error occurred (get KPI definitions by SDType): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) KpiDefinitionsBySdInstance(ctx context.Context, id uint32) ([]graphQLModel.KPIDefinition, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDInstance(principal.UserID, id)
	if result.IsFailure() {
		log.Printf("Error occurred (get KPI definitions by SDInstance): %s\n", result.GetError().Error())
	}
	return result.Unwrap()
}

func (r *queryResolver) RawDataPointsBySDType(ctx context.Context, id uint32) ([]graphQLModel.RawDataPoint, error) {
	_, err := authorizeOperation(ctx, auth.ResourceRawData, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetRawDataPointsBySDType(id)
	if result.IsFailure() {
		log.Printf("Error occurred (get raw data points by sd type): %s\n", result.GetError().Error())
		return nil, result.GetError()
	}
	return result.Unwrap()
}

func (r *queryResolver) RawDataPoint(ctx context.Context, id uint32) (graphQLModel.RawDataPoint, error) {
	_, err := authorizeOperation(ctx, auth.ResourceRawData, auth.OperationRead)
	if err != nil {
		return graphQLModel.RawDataPoint{}, err
	}
	opt, err := domainLogicLayer.GetRawDataPoint(id).Unwrap()
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

func (r *queryResolver) KpiResultsByKpi(ctx context.Context, id uint32) ([]graphQLModel.KPIFulfillmentCheckResult, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.GetKPIFulfillmentCheckResultsByKPI(principal.UserID, id)
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

func (r *queryResolver) SdInstanceGroup(ctx context.Context, id uint32) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return graphQLModel.SDInstanceGroup{}, err
	}
	getSDInstanceGroupResult := domainLogicLayer.GetSDInstanceGroup(id)
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

func (r *queryResolver) UserRole(ctx context.Context, id uint32) (*graphQLModel.Role, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead); err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUserRole(id)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	role := result.GetPayload()
	return &role, nil
}

func (r *queryResolver) Role(ctx context.Context) (*graphQLModel.Role, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceRoles, auth.OperationRead)
	if err != nil {
		return nil, err
	}
	result := domainLogicLayer.LoadUserRole(principal.UserID)
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

func (r *queryResolver) APIKey(ctx context.Context, id uint32) (graphQLModel.APIKey, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationRead)
	if err != nil {
		return graphQLModel.APIKey{}, err
	}
	result := domainLogicLayer.LoadAPIKeyByID(principal.UserID, id)
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

func (r *queryResolver) TimeSeriesExport(ctx context.Context, id uint32) (graphQLModel.TimeSeriesExport, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceTimeSeries, auth.OperationRead)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return domainLogicLayer.GetTimeSeriesExport(principal.UserID, id)
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
	for _, id := range filter.Ids {
		if _, err := domainLogicLayer.GetTimeSeriesExport(principal.UserID, id); err != nil {
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
