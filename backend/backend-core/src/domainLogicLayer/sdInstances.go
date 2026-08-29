/**
 * @file sdInstances.go
 * @brief Doménová logika pro správu sledovaných instancí.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ práce se sledovanými instancemi.
 * - Vojtěch Hubáček: doplnění čtení instancí podle sledovaného typu a KPI definice.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func GetSDInstance(uid string) sharedUtils.Result[graphQLModel.SDInstance] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(strings.TrimSpace(uid))
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstance](loadResult.GetError())
	}
	sdInstanceOptional := loadResult.GetPayload()
	if sdInstanceOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstance](fmt.Errorf("couldn't find SD instance for UID: %s", uid))
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelSDInstance(sdInstanceOptional.GetPayload()))
}

func GetSDInstances() sharedUtils.Result[[]graphQLModel.SDInstance] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstances()
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.SDInstance](sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelSDInstance))
}

func GetSDInstancesByType(sdTypeUID string) sharedUtils.Result[[]graphQLModel.SDInstance] {
	normalizedSDTypeUID, normalizeErr := normalizeSDTypeUID(sdTypeUID)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](normalizeErr)
	}
	sdTypeResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(normalizedSDTypeUID)
	if sdTypeResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](sdTypeResult.GetError())
	}
	sdTypeIDOptional := sdTypeResult.GetPayload().ID
	if sdTypeIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](fmt.Errorf("SD type loaded by UID has no internal ID: %s", sdTypeUID))
	}
	sdTypeID := sdTypeIDOptional.GetPayload()
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstancesByType(sdTypeID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelSDInstance))
}

func GetSDInstancesByKpiDefinition(userID uint32, kpiDefinitionUID string) sharedUtils.Result[[]graphQLModel.SDInstance] {
	kpiDefinitionResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, strings.TrimSpace(kpiDefinitionUID))
	if kpiDefinitionResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](kpiDefinitionResult.GetError())
	}
	kpiDefinitionID := kpiDefinitionResult.GetPayload().ID
	if kpiDefinitionID == nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", kpiDefinitionUID))
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstancesByKpiDefinition(*kpiDefinitionID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstance](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), dll2gql.ToGraphQLModelSDInstance))
}

func UpdateSDInstance(uid string, sdInstanceUpdateInput graphQLModel.SDInstanceUpdateInput) sharedUtils.Result[graphQLModel.SDInstance] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(strings.TrimSpace(uid))
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstance](loadResult.GetError())
	}
	sdInstanceOptional := loadResult.GetPayload()
	if sdInstanceOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstance](fmt.Errorf("couldn't find SD instance for UID: %s", uid))
	}
	sdInstance := sdInstanceOptional.GetPayload()
	sharedUtils.NewOptionalFromPointer(sdInstanceUpdateInput.UserIdentifier).DoIfPresent(func(userIdentifier string) {
		sdInstance.UserIdentifier = userIdentifier
	})
	sharedUtils.NewOptionalFromPointer(sdInstanceUpdateInput.ConfirmedByUser).DoIfPresent(func(confirmedByUser bool) {
		sdInstance.ConfirmedByUser = confirmedByUser
	})
	if persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistSDInstance(sdInstance); persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstance](persistResult.GetError())
	}
	isc.EnqueueMessageRepresentingCurrentSDInstanceConfiguration(getDLLRabbitMQClient())
	return sharedUtils.NewSuccessResult[graphQLModel.SDInstance](dll2gql.ToGraphQLModelSDInstance(sdInstance))
}
