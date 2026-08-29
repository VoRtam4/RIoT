/**
 * @file kpiDefinitions.go
 * @brief Doménová logika pro správu KPI definic a spuštění navazujícího reprocessingu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ správy KPI definic.
 * - Vojtěch Hubáček: doplnění userID vazeb, By operací, aktivace reprocessingu a kontroly reálné změny definice před přepočtem.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
)

func CreateKPIDefinition(userID uint32, kpiDefinitionInput graphQLModel.KPIDefinitionInput) sharedUtils.Result[graphQLModel.KPIDefinition] {
	sdType, selectedSDInstanceIDs, selectedSDInstanceUIDs, prepareErr := prepareKPIDefinitionInputContext(kpiDefinitionInput)
	if prepareErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](prepareErr)
	}
	kpiDefinitionInput.SdTypeUID = sdType.UID
	kpiDefinitionInput.SelectedSDInstanceUIDs = selectedSDInstanceUIDs
	toDLLModelTransformResult := gql2dll.ToDLLModelKPIDefinition(kpiDefinitionInput, sdType, selectedSDInstanceIDs)
	if toDLLModelTransformResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](toDLLModelTransformResult.GetError())
	}
	kpiDefinition := toDLLModelTransformResult.GetPayload()
	normalizedUIDResult := normalizeKPIUID(kpiDefinitionInput.UID, kpiDefinition.SDTypeSpecification, kpiDefinition.Label)
	if normalizedUIDResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](normalizedUIDResult.GetError())
	}
	kpiDefinition.UID = normalizedUIDResult.GetPayload()
	persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistKPIDefinition(userID, kpiDefinition)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](persistResult.GetError())
	}
	jobID := uuid.New().String()
	isc.EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(getDLLRabbitMQClient(), jobID)
	id := persistResult.GetPayload()
	kpiDefinition.ID = &id
	if err := isc.EnqueueKPIReprocessRequest(getDLLRabbitMQClient(), kpiDefinition, time.Now().UTC(), jobID, false); err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](err)
	}
	return sharedUtils.NewSuccessResult[graphQLModel.KPIDefinition](dll2gql.ToGraphQLModelKPIDefinition(kpiDefinition))
}

func UpdateKPIDefinition(userID uint32, uid string, kpiDefinitionInput graphQLModel.KPIDefinitionInput) sharedUtils.Result[graphQLModel.KPIDefinition] {
	existingKPIDefinitionResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, strings.TrimSpace(uid))
	if existingKPIDefinitionResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](existingKPIDefinitionResult.GetError())
	}
	existingKPIDefinition := existingKPIDefinitionResult.GetPayload()
	if existingKPIDefinition.ID == nil {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", uid))
	}
	sdType, selectedSDInstanceIDs, selectedSDInstanceUIDs, prepareErr := prepareKPIDefinitionInputContext(kpiDefinitionInput)
	if prepareErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](prepareErr)
	}
	kpiDefinitionInput.SdTypeUID = sdType.UID
	kpiDefinitionInput.SelectedSDInstanceUIDs = selectedSDInstanceUIDs
	toDLLModelTransformResult := gql2dll.ToDLLModelKPIDefinition(kpiDefinitionInput, sdType, selectedSDInstanceIDs)
	if toDLLModelTransformResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](toDLLModelTransformResult.GetError())
	}
	kpiDefinition := toDLLModelTransformResult.GetPayload()
	id := *existingKPIDefinition.ID
	kpiDefinition.ID = &id
	fallbackUIDSuffix := kpiDefinition.Label
	if existingKPIDefinition.UID != nil {
		fallbackUIDSuffix = extractKPIUIDSuffix(*existingKPIDefinition.UID, existingKPIDefinition.SDTypeSpecification)
	}
	normalizedUIDResult := normalizeKPIUID(kpiDefinitionInput.UID, kpiDefinition.SDTypeSpecification, fallbackUIDSuffix)
	if normalizedUIDResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](normalizedUIDResult.GetError())
	}
	kpiDefinition.UID = normalizedUIDResult.GetPayload()
	shouldReprocess := ShouldReprocessKPIDefinition(existingKPIDefinition, kpiDefinition)
	persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistKPIDefinition(userID, kpiDefinition)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](persistResult.GetError())
	}
	jobID := uuid.New().String()
	if shouldReprocess {
		if err := isc.EnqueueKPIDeleteRequest(getDLLRabbitMQClient(), kpiDefinition, jobID); err != nil {
			return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](err)
		}
		if err := isc.EnqueueKPIReprocessRequest(getDLLRabbitMQClient(), kpiDefinition, time.Now().UTC(), jobID, true); err != nil {
			return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](err)
		}
	}
	isc.EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(getDLLRabbitMQClient(), jobID)
	return sharedUtils.NewSuccessResult[graphQLModel.KPIDefinition](dll2gql.ToGraphQLModelKPIDefinition(kpiDefinition))
}

func prepareKPIDefinitionInputContext(kpiDefinitionInput graphQLModel.KPIDefinitionInput) (dllModel.SDType, []uint32, []string, error) {
	sdTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(kpiDefinitionInput.SdTypeUID, "sdt", "SD type")
	if normalizeErr != nil {
		return dllModel.SDType{}, nil, nil, normalizeErr
	}
	sdTypeResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(sdTypeUID)
	if sdTypeResult.IsFailure() {
		return dllModel.SDType{}, nil, nil, sdTypeResult.GetError()
	}
	sdType := sdTypeResult.GetPayload()
	selectedSDInstanceIDs := make([]uint32, 0, len(kpiDefinitionInput.SelectedSDInstanceUIDs))
	selectedSDInstanceUIDs := make([]string, 0, len(kpiDefinitionInput.SelectedSDInstanceUIDs))
	for _, sdInstanceUID := range kpiDefinitionInput.SelectedSDInstanceUIDs {
		normalizedUID, _, normalizeErr := sharedUtils.NormalizeScopedUID(sdInstanceUID, sdTypeUID, "sdi", "selected SD instance")
		if normalizeErr != nil {
			return dllModel.SDType{}, nil, nil, normalizeErr
		}
		sdInstanceResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(normalizedUID)
		if sdInstanceResult.IsFailure() {
			return dllModel.SDType{}, nil, nil, sdInstanceResult.GetError()
		}
		sdInstanceOptional := sdInstanceResult.GetPayload()
		if sdInstanceOptional.IsEmpty() {
			return dllModel.SDType{}, nil, nil, fmt.Errorf("couldn't find SD instance for UID: %s", normalizedUID)
		}
		sdInstance := sdInstanceOptional.GetPayload()
		if sdInstance.SDType.UID != "" && sdInstance.SDType.UID != sdType.UID {
			return dllModel.SDType{}, nil, nil, fmt.Errorf("selected SD instance %s doesn't belong to SD type %s", normalizedUID, sdType.UID)
		}
		if sdInstance.ID.IsEmpty() {
			return dllModel.SDType{}, nil, nil, fmt.Errorf("SD instance loaded by UID has no internal ID: %s", normalizedUID)
		}
		selectedSDInstanceIDs = append(selectedSDInstanceIDs, sdInstance.ID.GetPayload())
		selectedSDInstanceUIDs = append(selectedSDInstanceUIDs, normalizedUID)
	}
	return sdType, selectedSDInstanceIDs, selectedSDInstanceUIDs, nil
}

func DeleteKPIDefinition(userID uint32, uid string) error {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, strings.TrimSpace(uid))
	if result.IsFailure() {
		return result.GetError()
	}
	kpiDefinition := result.GetPayload()
	if kpiDefinition.ID == nil {
		return fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", uid)
	}
	id := *kpiDefinition.ID
	if err := dbClient.GetRelationalDatabaseClientInstance().DeleteKPIDefinition(id); err != nil {
		return err
	}
	if err := isc.EnqueueKPIDeleteRequest(getDLLRabbitMQClient(), kpiDefinition, ""); err != nil {
		return err
	}
	isc.EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(getDLLRabbitMQClient(), "")
	return nil
}

func GetKPIDefinitions(userID uint32) sharedUtils.Result[[]graphQLModel.KPIDefinition] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitions(userID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]graphQLModel.KPIDefinition](sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelKPIDefinition))
}

func GetKPIDefinition(userID uint32, uid string) sharedUtils.Result[graphQLModel.KPIDefinition] {
	normalizedUID := strings.TrimSpace(uid)
	if normalizedUID == "" {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](fmt.Errorf("KPI UID must not be empty"))
	}
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionByUID(userID, normalizedUID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[graphQLModel.KPIDefinition](dll2gql.ToGraphQLModelKPIDefinition(loadResult.GetPayload()))
}

func normalizeKPIUID(inputUID *string, sdTypeUID string, fallbackUIDSuffix string) sharedUtils.Result[*string] {
	normalizedSDTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(sdTypeUID, "sdt", "SD type")
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[*string](normalizeErr)
	}
	uidSuffix := strings.TrimSpace(fallbackUIDSuffix)
	if inputUID != nil {
		uidSuffix = strings.TrimSpace(*inputUID)
	}
	uid, _, err := sharedUtils.NormalizeScopedUID(uidSuffix, normalizedSDTypeUID, "kpi", "KPI definition")
	if err != nil {
		return sharedUtils.NewFailureResult[*string](err)
	}
	return sharedUtils.NewSuccessResult(&uid)
}

func extractKPIUIDSuffix(uid string, sdTypeUID string) string {
	prefix := sdTypeUID + ".kpi:"
	if strings.HasPrefix(uid, prefix) {
		return strings.TrimPrefix(uid, prefix)
	}
	return uid
}

func GetKPIDefinitionsBySDType(userID uint32, sdTypeUID string) sharedUtils.Result[[]graphQLModel.KPIDefinition] {
	normalizedSDTypeUID, normalizeErr := normalizeSDTypeUID(sdTypeUID)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](normalizeErr)
	}
	sdTypeResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypeBasedOnUID(normalizedSDTypeUID)
	if sdTypeResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](sdTypeResult.GetError())
	}
	sdTypeIDOptional := sdTypeResult.GetPayload().ID
	if sdTypeIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](fmt.Errorf("SD type loaded by UID has no internal ID: %s", sdTypeUID))
	}
	sdTypeID := sdTypeIDOptional.GetPayload()
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionsBySDType(userID, sdTypeID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelKPIDefinition))
}

func GetKPIDefinitionsBySDInstance(userID uint32, sdInstanceUID string) sharedUtils.Result[[]graphQLModel.KPIDefinition] {
	sdInstanceResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(strings.TrimSpace(sdInstanceUID))
	if sdInstanceResult.IsFailure() {
		log.Printf("[DLL] failed: %s", sdInstanceResult.GetError().Error())
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](sdInstanceResult.GetError())
	}
	sdInstanceOptional := sdInstanceResult.GetPayload()
	if sdInstanceOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](fmt.Errorf("couldn't find SD instance for UID: %s", sdInstanceUID))
	}
	sdInstanceIDOptional := sdInstanceOptional.GetPayload().ID
	if sdInstanceIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](fmt.Errorf("SD instance loaded by UID has no internal ID: %s", sdInstanceUID))
	}
	sdInstanceID := sdInstanceIDOptional.GetPayload()
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionsBySDInstance(userID, sdInstanceID)
	if loadResult.IsFailure() {
		log.Printf("[DLL] failed: %s", loadResult.GetError().Error())
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelKPIDefinition))
}

func ShouldReprocessKPIDefinition(oldDef sharedModel.KPIDefinition, newDef sharedModel.KPIDefinition) bool {
	if oldDef.SDTypeID != newDef.SDTypeID {
		return true
	}
	if oldDef.SDInstanceMode != newDef.SDInstanceMode {
		return true
	}
	oldSelected := append([]uint32(nil), oldDef.SelectedSDInstanceIDs...)
	newSelected := append([]uint32(nil), newDef.SelectedSDInstanceIDs...)
	slices.Sort(oldSelected)
	slices.Sort(newSelected)
	if !reflect.DeepEqual(oldSelected, newSelected) {
		return true
	}
	if !AreKPINodesEqual(oldDef.RootNode, newDef.RootNode) {
		return true
	}
	return false
}

func normalizedReferenceMode(mode sharedModel.KPIReferenceMode) sharedModel.KPIReferenceMode {
	if mode == "" {
		return sharedModel.KPIReferenceModeLiteral
	}
	return mode
}

func referenceConfigEqual(aMode sharedModel.KPIReferenceMode, bMode sharedModel.KPIReferenceMode, aComparedSpecification string, bComparedSpecification string, aComparedOffset *int, bComparedOffset *int) bool {
	aNormalizedMode := normalizedReferenceMode(aMode)
	bNormalizedMode := normalizedReferenceMode(bMode)
	if aNormalizedMode != bNormalizedMode {
		return false
	}
	if aNormalizedMode == sharedModel.KPIReferenceModeLiteral {
		return true
	}
	return aComparedSpecification == bComparedSpecification && reflect.DeepEqual(aComparedOffset, bComparedOffset)
}

func AreKPINodesEqual(a sharedModel.KPINode, b sharedModel.KPINode) bool {
	switch ta := a.(type) {
	case *sharedModel.LogicalOperationKPINode:
		tb, ok := b.(*sharedModel.LogicalOperationKPINode)
		if !ok {
			return false
		}
		if ta.Type != tb.Type {
			return false
		}
		if len(ta.ChildNodes) != len(tb.ChildNodes) {
			return false
		}
		for i := range ta.ChildNodes {
			if !AreKPINodesEqual(ta.ChildNodes[i], tb.ChildNodes[i]) {
				return false
			}
		}
		return true
	case *sharedModel.StringEQAtomKPINode:
		tb, ok := b.(*sharedModel.StringEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.StringNEQAtomKPINode:
		tb, ok := b.(*sharedModel.StringNEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.StringExistsAtomKPINode:
		tb, ok := b.(*sharedModel.StringExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	case *sharedModel.StringNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.StringNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	case *sharedModel.BooleanEQAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.BooleanNEQAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanNEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.BooleanExistsAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	case *sharedModel.BooleanNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	case *sharedModel.NumericEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericNEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericNEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericLTAtomKPINode:
		tb, ok := b.(*sharedModel.NumericLTAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericLEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericLEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericGTAtomKPINode:
		tb, ok := b.(*sharedModel.NumericGTAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericGEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericGEQAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification &&
			ta.ReferenceValue == tb.ReferenceValue &&
			referenceConfigEqual(ta.ReferenceMode, tb.ReferenceMode, ta.ComparedSDParameterSpecification, tb.ComparedSDParameterSpecification, ta.ComparedRecordOffset, tb.ComparedRecordOffset)
	case *sharedModel.NumericExistsAtomKPINode:
		tb, ok := b.(*sharedModel.NumericExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	case *sharedModel.NumericNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.NumericNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterSpecification == tb.SDParameterSpecification
	default:
		return false
	}
}
