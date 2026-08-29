/**
 * @file rawDataPoint.go
 * @brief Doménová logika pro čtení aktuálních raw datových bodů sledovaných instancí.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality raw datových bodů v doménové vrstvě.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func GetRawDataPointsBySDType(sdTypeID uint32) sharedUtils.Result[[]graphQLModel.RawDataPoint] {
	if sdTypeID == 0 {
		return sharedUtils.NewFailureResult[[]graphQLModel.RawDataPoint](fmt.Errorf("invalid sdTypeID"))
	}
	dbResult := dbClient.GetRelationalDatabaseClientInstance().LoadRawDataPointsBySDType(sdTypeID)
	if dbResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.RawDataPoint](dbResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.RawDataPoint](sharedUtils.Map(dbResult.GetPayload(), dll2gql.ToGraphQLModelRawDataPoint))
}

func GetRawDataPointsBySDTypeUID(sdTypeUID string) sharedUtils.Result[[]graphQLModel.RawDataPoint] {
	normalizedUID, normalizeErr := normalizeSDTypeUID(sdTypeUID)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.RawDataPoint](normalizeErr)
	}
	sdTypeResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(normalizedUID)
	if sdTypeResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.RawDataPoint](sdTypeResult.GetError())
	}
	sdType := sdTypeResult.GetPayload()
	if sdType.ID.IsEmpty() {
		return sharedUtils.NewFailureResult[[]graphQLModel.RawDataPoint](fmt.Errorf("SD type loaded by UID has no internal ID: %s", normalizedUID))
	}
	return GetRawDataPointsBySDType(sdType.ID.GetPayload())
}

func GetRawDataPoint(sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[graphQLModel.RawDataPoint]] {
	if sdInstanceID == 0 {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](fmt.Errorf("invalid sdInstanceID"))
	}
	dbResult := dbClient.GetRelationalDatabaseClientInstance().LoadRawDataPoint(sdInstanceID)
	if dbResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](dbResult.GetError())
	}
	opt := dbResult.GetPayload()
	if opt.IsEmpty() {
		return sharedUtils.NewSuccessResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](sharedUtils.NewEmptyOptional[graphQLModel.RawDataPoint]())
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](
		sharedUtils.NewOptionalOf(dll2gql.ToGraphQLModelRawDataPoint(opt.GetPayload())))
}

func GetRawDataPointBySDInstanceUID(sdInstanceUID string) sharedUtils.Result[sharedUtils.Optional[graphQLModel.RawDataPoint]] {
	trimmedUID := strings.TrimSpace(sdInstanceUID)
	if trimmedUID == "" {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](fmt.Errorf("invalid sdInstanceUID"))
	}
	sdInstanceResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(trimmedUID)
	if sdInstanceResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](sdInstanceResult.GetError())
	}
	sdInstanceOptional := sdInstanceResult.GetPayload()
	if sdInstanceOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](fmt.Errorf("couldn't find SD instance for UID: %s", trimmedUID))
	}
	sdInstanceIDOptional := sdInstanceOptional.GetPayload().ID
	if sdInstanceIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.RawDataPoint]](fmt.Errorf("SD instance loaded by UID has no internal ID: %s", trimmedUID))
	}
	return GetRawDataPoint(sdInstanceIDOptional.GetPayload())
}
