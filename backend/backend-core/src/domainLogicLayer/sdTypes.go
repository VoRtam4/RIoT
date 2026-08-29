/**
 * @file sdTypes.go
 * @brief Doménová logika pro správu sledovaných typů a jejich konfigurace.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func CreateSDType(sdTypeInput graphQLModel.SDTypeInput) sharedUtils.Result[graphQLModel.SDType] {
	sdType := gql2dll.ToDLLModelSDType(sdTypeInput)
	persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistSDType(sdType)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDType](persistResult.GetError())
	}
	isc.EnqueueMessageRepresentingCurrentSDTypeConfiguration(getDLLRabbitMQClient())
	return sharedUtils.NewSuccessResult[graphQLModel.SDType](dll2gql.ToGraphQLModelSDType(persistResult.GetPayload()))
}

func DeleteSDType(uid string) error {
	normalizedUID, normalizeErr := normalizeSDTypeUID(uid)
	if normalizeErr != nil {
		return normalizeErr
	}
	sdTypeResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(normalizedUID)
	if sdTypeResult.IsFailure() {
		return sdTypeResult.GetError()
	}
	sdTypeIDOptional := sdTypeResult.GetPayload().ID
	if sdTypeIDOptional.IsEmpty() {
		return fmt.Errorf("SD type loaded by UID has no internal ID: %s", uid)
	}
	id := sdTypeIDOptional.GetPayload()
	if err := dbClient.GetRelationalDatabaseClientInstance().DeleteSDType(id); err != nil {
		return err
	}
	isc.EnqueueMessagesRepresentingCurrentSystemConfiguration(getDLLRabbitMQClient())
	return nil
}

func GetSDType(uid string) sharedUtils.Result[graphQLModel.SDType] {
	normalizedUID, normalizeErr := normalizeSDTypeUID(uid)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDType](normalizeErr)
	}
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(normalizedUID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDType](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[graphQLModel.SDType](dll2gql.ToGraphQLModelSDType(loadResult.GetPayload()))
}

func GetSDTypes() sharedUtils.Result[[]graphQLModel.SDType] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypes()
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDType](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.SDType](sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelSDType))
}
