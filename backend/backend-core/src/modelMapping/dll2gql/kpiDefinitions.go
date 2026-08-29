/**
 * @file kpiDefinitions.go
 * @brief Mapování doménové reprezentace KPI definic do GraphQL modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování základních KPI operací.
 * - Vojtěch Hubáček: rozšíření mapování o nově přidané KPI operace, labely KPI definic a zjednodušení převodu string enumů přetypováním na GraphQL typy.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func graphQLKPIReferenceMode(referenceMode sharedModel.KPIReferenceMode) graphQLModel.KPIReferenceMode {
	if referenceMode == "" {
		referenceMode = sharedModel.KPIReferenceModeLiteral
	}
	return graphQLModel.KPIReferenceMode(referenceMode)
}

func optionalStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toGraphQLModelKPINode(kpiNode sharedModel.KPINode, id uint32, parentNodeID *uint32) graphQLModel.KPINode {
	switch typedKPINode := kpiNode.(type) {
	case *sharedModel.StringEQAtomKPINode:
		return graphQLModel.StringEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset: typedKPINode.ComparedRecordOffset,
			StringReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.StringNEQAtomKPINode:
		return graphQLModel.StringNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringNEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset: typedKPINode.ComparedRecordOffset,
			StringReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.StringExistsAtomKPINode:
		return graphQLModel.StringExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.StringNotExistsAtomKPINode:
		return graphQLModel.StringNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringNotExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.BooleanEQAtomKPINode:
		return graphQLModel.BooleanEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			BooleanReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.BooleanNEQAtomKPINode:
		return graphQLModel.BooleanNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanNEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			BooleanReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.BooleanExistsAtomKPINode:
		return graphQLModel.BooleanExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.BooleanNotExistsAtomKPINode:
		return graphQLModel.BooleanNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanNotExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.NumericEQAtomKPINode:
		return graphQLModel.NumericEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericNEQAtomKPINode:
		return graphQLModel.NumericNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericNEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericLTAtomKPINode:
		return graphQLModel.NumericLTAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericLTAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericLEQAtomKPINode:
		return graphQLModel.NumericLEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericLEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericGTAtomKPINode:
		return graphQLModel.NumericGTAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericGTAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericGEQAtomKPINode:
		return graphQLModel.NumericGEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericGEQAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLKPIReferenceMode(typedKPINode.ReferenceMode),
			ComparedSDParameterSpecification: optionalStringPointer(
				typedKPINode.ComparedSDParameterSpecification,
			),
			ComparedRecordOffset:  typedKPINode.ComparedRecordOffset,
			NumericReferenceValue: typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericExistsAtomKPINode:
		return graphQLModel.NumericExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.NumericNotExistsAtomKPINode:
		return graphQLModel.NumericNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericNotExistsAtom,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			ReferenceMode:            graphQLModel.KPIReferenceModeLiteral,
		}
	case *sharedModel.LogicalOperationKPINode:
		return graphQLModel.LogicalOperationKPINode{
			ID:           id,
			ParentNodeID: parentNodeID,
			NodeType:     graphQLModel.KPINodeTypeLogicalOperation,
			Type:         graphQLModel.LogicalOperationType(typedKPINode.Type),
		}
	}
	panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
}

func processKPINode(node sharedModel.KPINode, generateNextNumber func() uint32, parentID *uint32) []graphQLModel.KPINode {
	nodeID := generateNextNumber()
	nodes := make([]graphQLModel.KPINode, 0)
	nodes = append(nodes, toGraphQLModelKPINode(node, nodeID, parentID))
	if sharedUtils.TypeIs[*sharedModel.LogicalOperationKPINode](node) {
		for _, childNode := range node.(*sharedModel.LogicalOperationKPINode).ChildNodes {
			nodes = append(nodes, processKPINode(childNode, generateNextNumber, &nodeID)...)
		}
	}
	return nodes
}

func ToGraphQLModelKPIDefinition(kpiDefinition sharedModel.KPIDefinition) graphQLModel.KPIDefinition {
	nodes := processKPINode(kpiDefinition.RootNode, sharedUtils.SequentialNumberGenerator(), nil)
	return graphQLModel.KPIDefinition{
		UID:                    kpiDefinition.UID,
		Label:                  kpiDefinition.Label,
		SdTypeUID:              kpiDefinition.SDTypeSpecification,
		UserIdentifier:         kpiDefinition.UserIdentifier,
		Nodes:                  nodes,
		SdInstanceMode:         graphQLModel.SDInstanceMode(strings.ToUpper(string(kpiDefinition.SDInstanceMode))),
		SelectedSDInstanceUIDs: kpiDefinition.SelectedSDInstanceUIDs,
	}
}
