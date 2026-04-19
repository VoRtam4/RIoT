package domainLogicLayer

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
)

func CreateKPIDefinition(userID uint32, kpiDefinitionInput graphQLModel.KPIDefinitionInput) sharedUtils.Result[graphQLModel.KPIDefinition] {
	toDLLModelTransformResult := gql2dll.ToDLLModelKPIDefinition(kpiDefinitionInput)
	if toDLLModelTransformResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](toDLLModelTransformResult.GetError())
	}
	kpiDefinition := toDLLModelTransformResult.GetPayload()
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

func UpdateKPIDefinition(userID uint32, id uint32, kpiDefinitionInput graphQLModel.KPIDefinitionInput) sharedUtils.Result[graphQLModel.KPIDefinition] {
	existingKPIDefinitionResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinition(userID, id)
	if existingKPIDefinitionResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](existingKPIDefinitionResult.GetError())
	}
	existingKPIDefinition := existingKPIDefinitionResult.GetPayload()
	toDLLModelTransformResult := gql2dll.ToDLLModelKPIDefinition(kpiDefinitionInput)
	if toDLLModelTransformResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](toDLLModelTransformResult.GetError())
	}
	kpiDefinition := toDLLModelTransformResult.GetPayload()
	kpiDefinition.ID = &id
	shouldReprocess := ShouldReprocessKPIDefinition(existingKPIDefinition, kpiDefinition)
	persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistKPIDefinition(userID, kpiDefinition)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](persistResult.GetError())
	}
	jobID := uuid.New().String()
	if shouldReprocess {
		if err := isc.EnqueueKPIDeleteRequest(getDLLRabbitMQClient(), id, kpiDefinition.SDTypeID, jobID); err != nil {
			return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](err)
		}
		if err := isc.EnqueueKPIReprocessRequest(getDLLRabbitMQClient(), kpiDefinition, time.Now().UTC(), jobID, true); err != nil {
			return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](err)
		}
	}
	isc.EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(getDLLRabbitMQClient(), jobID)
	return sharedUtils.NewSuccessResult[graphQLModel.KPIDefinition](dll2gql.ToGraphQLModelKPIDefinition(kpiDefinition))
}

func DeleteKPIDefinition(userID uint32, id uint32) error {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinition(userID, id)
	if result.IsFailure() {
		return result.GetError()
	}
	if err := dbClient.GetRelationalDatabaseClientInstance().DeleteKPIDefinition(id); err != nil {
		return err
	}
	if err := isc.EnqueueKPIDeleteRequest(getDLLRabbitMQClient(), id, result.GetPayload().SDTypeID, ""); err != nil {
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

func GetKPIDefinition(userID uint32, id uint32) sharedUtils.Result[graphQLModel.KPIDefinition] {
	// FIXME: Searching for target KPI definition on domain-logic layer (suboptimal)
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitions(userID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](loadResult.GetError())
	}
	kpiDefinitions := loadResult.GetPayload()
	targetKPIDefinitionOptional := sharedUtils.FindFirst[sharedModel.KPIDefinition](kpiDefinitions, func(kpiDefinition sharedModel.KPIDefinition) bool {
		return *kpiDefinition.ID == id
	})
	if targetKPIDefinitionOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.KPIDefinition](fmt.Errorf("couldn't find KPI definition for id: %d", id))
	}
	return sharedUtils.NewSuccessResult[graphQLModel.KPIDefinition](dll2gql.ToGraphQLModelKPIDefinition(targetKPIDefinitionOptional.GetPayload()))
}

func GetKPIDefinitionsBySDType(userID uint32, sdTypeID uint32) sharedUtils.Result[[]graphQLModel.KPIDefinition] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitionsBySDType(userID, sdTypeID)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.KPIDefinition](loadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(loadResult.GetPayload(), dll2gql.ToGraphQLModelKPIDefinition))
}

func GetKPIDefinitionsBySDInstance(userID uint32, sdInstanceID uint32) sharedUtils.Result[[]graphQLModel.KPIDefinition] {
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
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.StringNEQAtomKPINode:
		tb, ok := b.(*sharedModel.StringNEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.StringExistsAtomKPINode:
		tb, ok := b.(*sharedModel.StringExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	case *sharedModel.StringNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.StringNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	case *sharedModel.BooleanEQAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.BooleanNEQAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanNEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.BooleanExistsAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	case *sharedModel.BooleanNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.BooleanNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	case *sharedModel.NumericEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericNEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericNEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericLTAtomKPINode:
		tb, ok := b.(*sharedModel.NumericLTAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericLEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericLEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericGTAtomKPINode:
		tb, ok := b.(*sharedModel.NumericGTAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericGEQAtomKPINode:
		tb, ok := b.(*sharedModel.NumericGEQAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID && ta.ReferenceValue == tb.ReferenceValue
	case *sharedModel.NumericExistsAtomKPINode:
		tb, ok := b.(*sharedModel.NumericExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	case *sharedModel.NumericNotExistsAtomKPINode:
		tb, ok := b.(*sharedModel.NumericNotExistsAtomKPINode)
		return ok &&
			ta.SDParameterID == tb.SDParameterID
	default:
		return false
	}
}
