package dll2gql

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func toGraphQLModelKPINode(kpiNode sharedModel.KPINode, id uint32, parentNodeID *uint32) graphQLModel.KPINode {
	switch typedKPINode := kpiNode.(type) {
	case *sharedModel.StringEQAtomKPINode:
		return graphQLModel.StringEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			StringReferenceValue:     typedKPINode.ReferenceValue,
		}
	case *sharedModel.StringNEQAtomKPINode:
		return graphQLModel.StringNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringNEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			StringReferenceValue:     typedKPINode.ReferenceValue,
		}
	case *sharedModel.StringExistsAtomKPINode:
		return graphQLModel.StringExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
		}
	case *sharedModel.StringNotExistsAtomKPINode:
		return graphQLModel.StringNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeStringNotExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
		}
	case *sharedModel.BooleanEQAtomKPINode:
		return graphQLModel.BooleanEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			BooleanReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.BooleanNEQAtomKPINode:
		return graphQLModel.BooleanNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanNEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			BooleanReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.BooleanExistsAtomKPINode:
		return graphQLModel.BooleanExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
		}
	case *sharedModel.BooleanNotExistsAtomKPINode:
		return graphQLModel.BooleanNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeBooleanNotExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
		}
	case *sharedModel.NumericEQAtomKPINode:
		return graphQLModel.NumericEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericNEQAtomKPINode:
		return graphQLModel.NumericNEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericNEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericLTAtomKPINode:
		return graphQLModel.NumericLTAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericLTAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericLEQAtomKPINode:
		return graphQLModel.NumericLEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericLEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericGTAtomKPINode:
		return graphQLModel.NumericGTAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericGTAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericGEQAtomKPINode:
		return graphQLModel.NumericGEQAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericGEQAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
			NumericReferenceValue:    typedKPINode.ReferenceValue,
		}
	case *sharedModel.NumericExistsAtomKPINode:
		return graphQLModel.NumericExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
		}
	case *sharedModel.NumericNotExistsAtomKPINode:
		return graphQLModel.NumericNotExistsAtomKPINode{
			ID:                       id,
			ParentNodeID:             parentNodeID,
			NodeType:                 graphQLModel.KPINodeTypeNumericNotExistsAtom,
			SdParameterID:            typedKPINode.SDParameterID,
			SdParameterSpecification: typedKPINode.SDParameterSpecification,
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
		ID:                    sharedUtils.NewOptionalFromPointer(kpiDefinition.ID).GetPayload(),
		Label:                 kpiDefinition.Label,
		SdTypeID:              kpiDefinition.SDTypeID,
		SdTypeUID:             kpiDefinition.SDTypeSpecification,
		UserIdentifier:        kpiDefinition.UserIdentifier,
		Nodes:                 nodes,
		SdInstanceMode:        graphQLModel.SDInstanceMode(strings.ToUpper(string(kpiDefinition.SDInstanceMode))),
		SelectedSDInstanceIDs: kpiDefinition.SelectedSDInstanceIDs,
	}
}
