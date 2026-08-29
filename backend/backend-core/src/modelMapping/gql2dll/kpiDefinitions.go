/**
 * @file kpiDefinitions.go
 * @brief Mapování GraphQL vstupů KPI definic do doménového stromu KPI podmínek.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování základních KPI operací.
 * - Vojtěch Hubáček: rozšíření mapování o nově přidané KPI operace a labely KPI definic.
 *
 * @ingroup riot_backend_core
 */
package gql2dll

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func failureResultDueToMissingInputProperty() sharedUtils.Result[sharedModel.KPINode] {
	return sharedUtils.NewFailureResult[sharedModel.KPINode](fmt.Errorf("model mapping failure – a necessary input property is missing"))
}

func referenceModeFromInput(kpiNodeInput graphQLModel.KPINodeInput) sharedModel.KPIReferenceMode {
	referenceMode := sharedModel.KPIReferenceModeLiteral
	if kpiNodeInput.ReferenceMode != nil {
		referenceMode = sharedModel.KPIReferenceMode(strings.ToLower(string(*kpiNodeInput.ReferenceMode)))
	}
	return referenceMode
}

func comparedSDParameterSpecificationFromInput(kpiNodeInput graphQLModel.KPINodeInput) string {
	if kpiNodeInput.ComparedSDParameterSpecification == nil {
		return ""
	}
	return *kpiNodeInput.ComparedSDParameterSpecification
}

func sdParameterIDBySpecification(sdType dllModel.SDType) sharedUtils.Result[map[string]uint32] {
	idsBySpecification := make(map[string]uint32)
	for _, sdParameter := range sdType.Parameters {
		if sdParameter.ID.IsEmpty() {
			return sharedUtils.NewFailureResult[map[string]uint32](fmt.Errorf("SD parameter loaded for SD type %s has no internal ID: %s", sdType.UID, sdParameter.Denotation))
		}
		idsBySpecification[sdParameter.Denotation] = sdParameter.ID.GetPayload()
	}
	return sharedUtils.NewSuccessResult(idsBySpecification)
}

func resolveSDParameterID(sdParameterSpecification string, idsBySpecification map[string]uint32) sharedUtils.Result[uint32] {
	sdParameterID, ok := idsBySpecification[sdParameterSpecification]
	if !ok {
		return sharedUtils.NewFailureResult[uint32](fmt.Errorf("unknown SD parameter specification in KPI input: %s", sdParameterSpecification))
	}
	return sharedUtils.NewSuccessResult(sdParameterID)
}

func literalValueRequired(kpiNodeInput graphQLModel.KPINodeInput) bool {
	if kpiNodeInput.ReferenceMode == nil {
		return true
	}
	return sharedModel.KPIReferenceMode(strings.ToLower(string(*kpiNodeInput.ReferenceMode))) == sharedModel.KPIReferenceModeLiteral
}

func kpiNodeInputToKPINode(kpiNodeInput graphQLModel.KPINodeInput, idsBySpecification map[string]uint32) sharedUtils.Result[sharedModel.KPINode] {
	nodeType := kpiNodeInput.Type
	if nodeType == graphQLModel.KPINodeTypeLogicalOperation {
		logicalOperationTypeOptional := sharedUtils.NewOptionalFromPointer[graphQLModel.LogicalOperationType](kpiNodeInput.LogicalOperationType)
		if logicalOperationTypeOptional.IsEmpty() {
			return failureResultDueToMissingInputProperty()
		}
		return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.LogicalOperationKPINode{
			Type:       sharedModel.LogicalOperationNodeType(logicalOperationTypeOptional.GetPayload()),
			ChildNodes: make([]sharedModel.KPINode, 0),
		})
	} else {
		sdParameterSpecificationOptional := sharedUtils.NewOptionalFromPointer[string](kpiNodeInput.SdParameterSpecification)
		if sdParameterSpecificationOptional.IsEmpty() {
			return failureResultDueToMissingInputProperty()
		}
		sdParameterSpecification := sdParameterSpecificationOptional.GetPayload()
		if sdParameterIDResult := resolveSDParameterID(sdParameterSpecification, idsBySpecification); sdParameterIDResult.IsFailure() {
			return sharedUtils.NewFailureResult[sharedModel.KPINode](sdParameterIDResult.GetError())
		}
		referenceMode := referenceModeFromInput(kpiNodeInput)
		comparedSDParameterSpecification := comparedSDParameterSpecificationFromInput(kpiNodeInput)
		if referenceMode == sharedModel.KPIReferenceModeParameter {
			if kpiNodeInput.ComparedSDParameterSpecification == nil || kpiNodeInput.ComparedRecordOffset == nil {
				return failureResultDueToMissingInputProperty()
			}
			comparedSDParameterIDResult := resolveSDParameterID(comparedSDParameterSpecification, idsBySpecification)
			if comparedSDParameterIDResult.IsFailure() {
				return sharedUtils.NewFailureResult[sharedModel.KPINode](comparedSDParameterIDResult.GetError())
			}
		}
		switch nodeType {
		case graphQLModel.KPINodeTypeStringEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[string](kpiNodeInput.StringReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.StringEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(""),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeStringNEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[string](kpiNodeInput.StringReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.StringNEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(""),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeStringExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.StringExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		case graphQLModel.KPINodeTypeStringNotExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.StringNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		case graphQLModel.KPINodeTypeBooleanEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[bool](kpiNodeInput.BooleanReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.BooleanEQAtomKPINode{
				SDParameterSpecification:         sharedUtils.NewOptionalFromPointer[string](kpiNodeInput.SdParameterSpecification).GetPayload(),
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(false),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeBooleanNEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[bool](kpiNodeInput.BooleanReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.BooleanNEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(false),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericNEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericNEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		case graphQLModel.KPINodeTypeNumericNotExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		case graphQLModel.KPINodeTypeNumericLTAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericLTAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericLEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericLEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericGTAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericGTAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeNumericGEQAtom:
			referenceValueOptional := sharedUtils.NewOptionalFromPointer[float64](kpiNodeInput.NumericReferenceValue)
			if referenceValueOptional.IsEmpty() && literalValueRequired(kpiNodeInput) {
				return failureResultDueToMissingInputProperty()
			}
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.NumericGEQAtomKPINode{
				SDParameterSpecification:         sdParameterSpecification,
				ReferenceValue:                   referenceValueOptional.GetPayloadOrDefault(0),
				ReferenceMode:                    referenceMode,
				ComparedSDParameterSpecification: comparedSDParameterSpecification,
				ComparedRecordOffset:             kpiNodeInput.ComparedRecordOffset,
			})
		case graphQLModel.KPINodeTypeBooleanExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.BooleanExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		case graphQLModel.KPINodeTypeBooleanNotExistsAtom:
			return sharedUtils.NewSuccessResult[sharedModel.KPINode](&sharedModel.BooleanNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterSpecification,
			})
		}
	}
	return sharedUtils.NewFailureResult[sharedModel.KPINode](fmt.Errorf("unexpected model mapping failure for KPI node type %q", nodeType))
}

func constructKPINodeByIDMap(kpiNodeInputs []graphQLModel.KPINodeInput, idsBySpecification map[string]uint32) sharedUtils.Result[map[uint32]sharedModel.KPINode] {
	kpiNodeByIDMap := make(map[uint32]sharedModel.KPINode)
	for _, kpiNodeInput := range kpiNodeInputs {
		kpiNodeResult := kpiNodeInputToKPINode(kpiNodeInput, idsBySpecification)
		if kpiNodeResult.IsFailure() {
			return sharedUtils.NewFailureResult[map[uint32]sharedModel.KPINode](kpiNodeResult.GetError())
		}
		kpiNodeByIDMap[kpiNodeInput.ID] = kpiNodeResult.GetPayload()
	}
	return sharedUtils.NewSuccessResult[map[uint32]sharedModel.KPINode](kpiNodeByIDMap)
}

func ToDLLModelKPIDefinition(kpiDefinitionInput graphQLModel.KPIDefinitionInput, sdType dllModel.SDType, selectedSDInstanceIDs []uint32) sharedUtils.Result[sharedModel.KPIDefinition] {
	sdTypeIDOptional := sdType.ID
	if sdTypeIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](fmt.Errorf("SD type loaded by UID has no internal ID: %s", kpiDefinitionInput.SdTypeUID))
	}
	idsBySpecificationResult := sdParameterIDBySpecification(sdType)
	if idsBySpecificationResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](idsBySpecificationResult.GetError())
	}
	kpiNodeInputs := kpiDefinitionInput.Nodes
	kpiNodeByIDMapConstructionResult := constructKPINodeByIDMap(kpiNodeInputs, idsBySpecificationResult.GetPayload())
	if kpiNodeByIDMapConstructionResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](kpiNodeByIDMapConstructionResult.GetError())
	}
	kpiNodeByIDMap := kpiNodeByIDMapConstructionResult.GetPayload()
	rootNodeIDOptional := sharedUtils.NewEmptyOptional[uint32]()
	for _, kpiNodeInput := range kpiNodeInputs {
		parentNodeIDOptional := sharedUtils.NewOptionalFromPointer(kpiNodeInput.ParentNodeID)
		if parentNodeIDOptional.IsEmpty() {
			rootNodeIDOptional = sharedUtils.NewOptionalOf(kpiNodeInput.ID)
			continue
		}
		parentNodeID := parentNodeIDOptional.GetPayload()
		parentNode := kpiNodeByIDMap[parentNodeID]
		if !sharedUtils.TypeIs[*sharedModel.LogicalOperationKPINode](parentNode) {
			return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](fmt.Errorf("model mapping failure – detected a parent node that's not a logical operation node"))
		}
		logicalOperationNode := parentNode.(*sharedModel.LogicalOperationKPINode)
		logicalOperationNode.ChildNodes = append(logicalOperationNode.ChildNodes, kpiNodeByIDMap[kpiNodeInput.ID])
		kpiNodeByIDMap[parentNodeID] = logicalOperationNode
	}
	if rootNodeIDOptional.IsEmpty() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](fmt.Errorf("model mapping failure – couldn't find root node ID"))
	}
	return sharedUtils.NewSuccessResult[sharedModel.KPIDefinition](sharedModel.KPIDefinition{
		ID:                     nil,
		UID:                    kpiDefinitionInput.UID,
		Label:                  kpiDefinitionInput.Label,
		SDTypeID:               sdTypeIDOptional.GetPayload(),
		SDTypeSpecification:    kpiDefinitionInput.SdTypeUID,
		UserIdentifier:         kpiDefinitionInput.UserIdentifier,
		RootNode:               kpiNodeByIDMap[rootNodeIDOptional.GetPayload()],
		SDInstanceMode:         sharedModel.SDInstanceMode(strings.ToLower(string(kpiDefinitionInput.SdInstanceMode))),
		SelectedSDInstanceIDs:  selectedSDInstanceIDs,
		SelectedSDInstanceUIDs: kpiDefinitionInput.SelectedSDInstanceUIDs,
	})
}
