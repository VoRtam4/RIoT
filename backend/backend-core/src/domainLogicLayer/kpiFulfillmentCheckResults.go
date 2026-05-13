/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Doménová logika pro čtení výsledků vyhodnocení KPI.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ práce s výsledky vyhodnocení KPI.
 * - Vojtěch Hubáček: doplnění operací omezených podle userID a By operací pro čtení výsledků podle KPI, instance a dalších vazeb.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func GetKPIFulfillmentCheckResults(userID uint32) sharedUtils.Result[[]graphQLModel.KPIFulfillmentCheckResult] {
	if userID == 0 {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](fmt.Errorf("invalid userID"))
	}
	dbResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIFulfillmentCheckResults(userID)
	if dbResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](dbResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.KPIFulfillmentCheckResult](sharedUtils.Map(dbResult.GetPayload(), dll2gql.ToGraphQLModelKPIFulfillmentCheckResult))
}

func GetKPIFulfillmentCheckResultsByKPI(userID uint32, kpiDefinitionID uint32) sharedUtils.Result[[]graphQLModel.KPIFulfillmentCheckResult] {
	if userID == 0 || kpiDefinitionID == 0 {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](fmt.Errorf("invalid input"))
	}
	dbResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIFulfillmentCheckResultsByKPI(userID, kpiDefinitionID)
	if dbResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](dbResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.KPIFulfillmentCheckResult](sharedUtils.Map(dbResult.GetPayload(), dll2gql.ToGraphQLModelKPIFulfillmentCheckResult))
}

func GetKPIFulfillmentCheckResult(userID uint32, input graphQLModel.KPIFulfillmentCheckResultRequest) sharedUtils.Result[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]] {
	if userID == 0 || input.KpiDefinitionID == 0 || input.SdInstanceID == 0 {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](fmt.Errorf("invalid input"))
	}
	dllInput := gql2dll.ToDLLModelKPIFulfillmentCheckResultRequest(input)
	dbResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIFulfillmentCheckResult(userID, dllInput)
	if dbResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](dbResult.GetError())
	}
	opt := dbResult.GetPayload()
	if opt.IsEmpty() {
		return sharedUtils.NewSuccessResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](sharedUtils.NewEmptyOptional[graphQLModel.KPIFulfillmentCheckResult]())
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](sharedUtils.NewOptionalOf(dll2gql.ToGraphQLModelKPIFulfillmentCheckResult(opt.GetPayload())))
}
