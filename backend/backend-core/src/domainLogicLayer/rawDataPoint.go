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
