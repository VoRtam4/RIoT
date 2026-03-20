package graphql

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/graphql/gsc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func authorizeOperation(ctx context.Context, operation string, opType string) (*auth.Principal, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	if !auth.CanAccessOperation(principal, operation, opType) {
		return nil, fmt.Errorf("forbidden")
	}
	return principal, nil
}

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
	if _, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationCreate); err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	createKPIDefinitionResult := domainLogicLayer.CreateKPIDefinition(input)
	if createKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (create KPI definition): %s\n", createKPIDefinitionResult.GetError().Error())
	}
	return createKPIDefinitionResult.Unwrap()
}

func (r *mutationResolver) UpdateKPIDefinition(ctx context.Context, id uint32, input graphQLModel.KPIDefinitionInput) (graphQLModel.KPIDefinition, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationUpdate); err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	updateKPIDefinitionResult := domainLogicLayer.UpdateKPIDefinition(id, input)
	if updateKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (update KPI definition): %s\n", updateKPIDefinitionResult.GetError().Error())
	}
	return updateKPIDefinitionResult.Unwrap()
}

func (r *mutationResolver) DeleteKPIDefinition(ctx context.Context, id uint32) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteKPIDefinition(id); err != nil {
		log.Printf("Error occurred (delete KPI definition): %s\n", err.Error())
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateSDInstanceGroup(ctx context.Context, input graphQLModel.SDInstanceGroupInput) (graphQLModel.SDInstanceGroup, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
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
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationUpdate); err != nil {
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

func (r *mutationResolver) UpdateUserConfig(ctx context.Context, userID uint32, input graphQLModel.UserConfigInput) (graphQLModel.UserConfig, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationUpdate); err != nil {
		return graphQLModel.UserConfig{}, err
	}
	updateUserConfigResult := domainLogicLayer.UpdateUserConfig(userID, input)
	if updateUserConfigResult.IsFailure() {
		log.Printf("Error occurred (update user config): %s\n", updateUserConfigResult.GetError().Error())
	}
	return updateUserConfigResult.Unwrap()
}

func (r *mutationResolver) DeleteUserConfig(ctx context.Context, userID uint32) (bool, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationDelete); err != nil {
		return false, err
	}
	if err := domainLogicLayer.DeleteUserConfig(userID); err != nil {
		log.Printf("Error occurred (delete user configuration): %s\n", err.Error())
		return false, err
	}
	return true, nil
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

func (r *queryResolver) KpiDefinition(ctx context.Context, id uint32) (graphQLModel.KPIDefinition, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead); err != nil {
		return graphQLModel.KPIDefinition{}, err
	}
	getKPIDefinitionResult := domainLogicLayer.GetKPIDefinition(id)
	if getKPIDefinitionResult.IsFailure() {
		log.Printf("Error occurred (get KPI definition): %s\n", getKPIDefinitionResult.GetError().Error())
	}
	return getKPIDefinitionResult.Unwrap()
}

func (r *queryResolver) KpiDefinitions(ctx context.Context) ([]graphQLModel.KPIDefinition, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceKPIDefinitions, auth.OperationRead); err != nil {
		return nil, err
	}
	getKPIDefinitionsResult := domainLogicLayer.GetKPIDefinitions()
	if getKPIDefinitionsResult.IsFailure() {
		log.Printf("Error occurred (get KPI definitions): %s\n", getKPIDefinitionsResult.GetError().Error())
	}
	return getKPIDefinitionsResult.Unwrap()
}

func (r *queryResolver) KpiFulfillmentCheckResults(ctx context.Context) ([]graphQLModel.KPIFulfillmentCheckResult, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead); err != nil {
		return nil, err
	}
	getKPIFulfillmentCheckResultsResult := domainLogicLayer.GetKPIFulfillmentCheckResults()
	if getKPIFulfillmentCheckResultsResult.IsFailure() {
		log.Printf("Error occurred (get KPI fulfillment check results): %s\n", getKPIFulfillmentCheckResultsResult.GetError().Error())
	}
	return getKPIFulfillmentCheckResultsResult.Unwrap()
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

func (r *queryResolver) UserConfig(ctx context.Context, id uint32) (graphQLModel.UserConfig, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationRead); err != nil {
		return graphQLModel.UserConfig{}, err
	}
	getUserConfigResult := domainLogicLayer.GetUserConfig(id)
	if getUserConfigResult.IsFailure() {
		log.Printf("Error occurred (get user config results): %s\n", getUserConfigResult.GetError().Error())
	}
	return getUserConfigResult.Unwrap()
}

func (r *queryResolver) MyUserConfig(ctx context.Context) (graphQLModel.UserConfig, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceUserConfig, auth.OperationRead)
	if err != nil {
		return graphQLModel.UserConfig{}, err
	}
	return domainLogicLayer.GetUserConfig(principal.UserID).Unwrap()
}

func (r *queryResolver) ApiKeys(ctx context.Context) ([]dbModel.APIKeyEntity, error) {
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

func (r *mutationResolver) CreateAPIKey(ctx context.Context, input graphQLModel.CreateAPIKeyInput) (string, error) {
	principal, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationCreate)
	if err != nil {
		return "", err
	}
	roleID, err := strconv.ParseUint(input.RoleID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid role id")
	}
	result := domainLogicLayer.CreateAPIKey(principal.UserID, uint32(roleID), input.Label, input.ExpiresAt)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.GetPayload(), nil
}

func (r *mutationResolver) UpdateAPIKey(ctx context.Context, id uint32, input graphQLModel.UpdateAPIKeyInput) (bool, error) {
	_, err := authorizeOperation(ctx, auth.ResourceAPIKeys, auth.OperationUpdate)
	if err != nil {
		return false, err
	}
	apiKeyResult := domainLogicLayer.LoadAPIKeyByID(id)
	if apiKeyResult.IsFailure() || apiKeyResult.GetPayload().IsEmpty() {
		return false, fmt.Errorf("api key not found")
	}
	apiKey := apiKeyResult.GetPayload().GetPayload()
	if input.Label != nil {
		apiKey.Label = *input.Label
	}
	if input.RoleID != nil {
		roleID, err := strconv.ParseUint(*input.RoleID, 10, 32)
		if err != nil {
			return false, err
		}
		apiKey.RoleID = uint32(roleID)
	}
	if input.ExpiresAt != nil {
		apiKey.ExpiresAt = input.ExpiresAt
	}
	if input.Revoked != nil {
		apiKey.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		apiKey.RateLimit = input.RateLimit
	}
	if err := domainLogicLayer.UpdateAPIKey(apiKey); err != nil {
		return false, err
	}
	return true, nil
}

func (r *subscriptionResolver) OnSDInstanceRegistered(ctx context.Context) (<-chan graphQLModel.SDInstance, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceEvents, auth.OperationSubscribe); err != nil {
		return nil, err
	}
	if _, err := authorizeOperation(ctx, auth.ResourceSDInstances, auth.OperationRead); err != nil {
		return nil, err
	}
	output := make(chan graphQLModel.SDInstance, 16)
	subscription := events.GetEventBus().Subscribe([]events.EventType{events.SDInstanceRegisteredEventType}, 16)
	go func() {
		defer close(output)
		defer events.GetEventBus().Unsubscribe(subscription.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-subscription.Channel:
				if !ok {
					return
				}
				payload, ok := event.Payload.(graphQLModel.SDInstance)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- payload:
				}
			}
		}
	}()
	return output, nil
}

func (r *subscriptionResolver) OnKPIFulfillmentChecked(ctx context.Context) (<-chan graphQLModel.KPIFulfillmentCheckResultTuple, error) {
	if _, err := authorizeOperation(ctx, auth.ResourceEvents, auth.OperationSubscribe); err != nil {
		return nil, err
	}
	if _, err := authorizeOperation(ctx, auth.ResourceKPIResults, auth.OperationRead); err != nil {
		return nil, err
	}
	output := make(chan graphQLModel.KPIFulfillmentCheckResultTuple, 16)
	subscription := events.GetEventBus().Subscribe([]events.EventType{events.KPIFulfillmentCheckedEventType}, 16)
	go func() {
		defer close(output)
		defer events.GetEventBus().Unsubscribe(subscription.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-subscription.Channel:
				if !ok {
					return
				}
				payload, ok := event.Payload.(graphQLModel.KPIFulfillmentCheckResultTuple)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- payload:
				}
			}
		}
	}()
	return output, nil
}

func (r *Resolver) Mutation() gsc.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() gsc.QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() gsc.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
