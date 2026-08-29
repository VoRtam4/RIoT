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
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
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

func GetKPIFulfillmentCheckResultsByKPIUID(userID uint32, kpiDefinitionUID string) sharedUtils.Result[[]graphQLModel.KPIFulfillmentCheckResult] {
	normalizedUID := strings.TrimSpace(kpiDefinitionUID)
	if userID == 0 || normalizedUID == "" {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](fmt.Errorf("invalid input"))
	}
	kpiDefinitionResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, normalizedUID)
	if kpiDefinitionResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](kpiDefinitionResult.GetError())
	}
	kpiDefinition := kpiDefinitionResult.GetPayload()
	if kpiDefinition.ID == nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIFulfillmentCheckResult](fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", normalizedUID))
	}
	return GetKPIFulfillmentCheckResultsByKPI(userID, *kpiDefinition.ID)
}

func GetKPIFulfillmentCheckResult(userID uint32, input graphQLModel.KPIFulfillmentCheckResultRequest) sharedUtils.Result[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]] {
	kpiDefinitionUID := strings.TrimSpace(input.KpiDefinitionUID)
	sdInstanceUID := strings.TrimSpace(input.SdInstanceUID)
	if userID == 0 || kpiDefinitionUID == "" || sdInstanceUID == "" {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](fmt.Errorf("invalid input"))
	}
	kpiDefinitionResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, kpiDefinitionUID)
	if kpiDefinitionResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](kpiDefinitionResult.GetError())
	}
	kpiDefinition := kpiDefinitionResult.GetPayload()
	if kpiDefinition.ID == nil {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", kpiDefinitionUID))
	}
	sdInstanceResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(sdInstanceUID)
	if sdInstanceResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](sdInstanceResult.GetError())
	}
	sdInstanceOptional := sdInstanceResult.GetPayload()
	if sdInstanceOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](fmt.Errorf("couldn't find SD instance for UID: %s", sdInstanceUID))
	}
	sdInstance := sdInstanceOptional.GetPayload()
	if sdInstance.ID.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[graphQLModel.KPIFulfillmentCheckResult]](fmt.Errorf("SD instance loaded by UID has no internal ID: %s", sdInstanceUID))
	}
	dllInput := dllModel.KPIFulfillmentCheckResultRequest{
		KPIDefinitionID: *kpiDefinition.ID,
		SDInstanceID:    sdInstance.ID.GetPayload(),
	}
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
