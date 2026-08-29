/**
 * @file kpiDefinitions.go
 * @brief Mapování databázové reprezentace KPI definic do doménového stromu KPI podmínek.
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
package db2dll

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func referenceModeFromAtomEntity(atomKPINodeEntity dbModel.AtomKPINodeEntity) sharedModel.KPIReferenceMode {
	referenceMode := sharedModel.KPIReferenceMode(atomKPINodeEntity.ReferenceMode)
	if referenceMode == "" {
		referenceMode = sharedModel.KPIReferenceModeLiteral
	}
	return referenceMode
}

func comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity dbModel.AtomKPINodeEntity) string {
	if atomKPINodeEntity.ComparedSDParameter != nil {
		return atomKPINodeEntity.ComparedSDParameter.Denotation
	}
	return ""
}

func reconstructKPINodeTree(currentKPINodeID uint32, kpiNodeParentChildrenMap map[uint32][]uint32, logicalOperationKPINodeEntities []dbModel.LogicalOperationKPINodeEntity, atomKPINodeEntities []dbModel.AtomKPINodeEntity) sharedModel.KPINode {
	logicalOperationKPINodeEntityOptional := sharedUtils.FindFirst(logicalOperationKPINodeEntities, func(logicalOperationKPINodeEntity dbModel.LogicalOperationKPINodeEntity) bool {
		return *logicalOperationKPINodeEntity.NodeID == currentKPINodeID
	})
	if logicalOperationKPINodeEntityOptional.IsPresent() {
		logicalOperationKPINodeEntity := logicalOperationKPINodeEntityOptional.GetPayload()
		childNodeIDs := kpiNodeParentChildrenMap[*logicalOperationKPINodeEntity.NodeID]
		childNodes := make([]sharedModel.KPINode, 0)
		sharedUtils.ForEach(childNodeIDs, func(childNodeID uint32) {
			childNodes = append(childNodes, reconstructKPINodeTree(childNodeID, kpiNodeParentChildrenMap, logicalOperationKPINodeEntities, atomKPINodeEntities))
		})
		return &sharedModel.LogicalOperationKPINode{
			Type:       sharedModel.LogicalOperationNodeType(logicalOperationKPINodeEntity.Type),
			ChildNodes: childNodes,
		}
	}
	atomKPINodeEntityOptional := sharedUtils.FindFirst(atomKPINodeEntities, func(atomKPINodeEntity dbModel.AtomKPINodeEntity) bool {
		return *atomKPINodeEntity.NodeID == currentKPINodeID
	})
	if atomKPINodeEntityOptional.IsPresent() {
		atomKPINodeEntity := atomKPINodeEntityOptional.GetPayload()
		sdParameterDenotation := atomKPINodeEntity.SDParameter.Denotation
		switch atomKPINodeEntity.Type {
		case "string_eq":
			return &sharedModel.StringEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.StringReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "string_neq":
			return &sharedModel.StringNEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.StringReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "string_exists":
			return &sharedModel.StringExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		case "string_not_exists":
			return &sharedModel.StringNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		case "boolean_eq":
			return &sharedModel.BooleanEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.BooleanReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "boolean_neq":
			return &sharedModel.BooleanNEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.BooleanReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "boolean_exists":
			return &sharedModel.BooleanExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		case "boolean_not_exists":
			return &sharedModel.BooleanNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		case "numeric_eq":
			return &sharedModel.NumericEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_neq":
			return &sharedModel.NumericNEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_lt":
			return &sharedModel.NumericLTAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_leq":
			return &sharedModel.NumericLEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_gt":
			return &sharedModel.NumericGTAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_geq":
			return &sharedModel.NumericGEQAtomKPINode{
				SDParameterSpecification:         sdParameterDenotation,
				ReferenceValue:                   sharedUtils.NewOptionalFromPointer(atomKPINodeEntity.NumericReferenceValue).GetPayload(),
				ReferenceMode:                    referenceModeFromAtomEntity(atomKPINodeEntity),
				ComparedSDParameterSpecification: comparedSDParameterSpecificationFromAtomEntity(atomKPINodeEntity),
				ComparedRecordOffset:             atomKPINodeEntity.ComparedRecordOffset,
			}
		case "numeric_exists":
			return &sharedModel.NumericExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		case "numeric_not_exists":
			return &sharedModel.NumericNotExistsAtomKPINode{
				SDParameterSpecification: sdParameterDenotation,
			}
		}
	}
	panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
}

func prepareKPINodeParentChildrenMap(kpiNodeEntities []dbModel.KPINodeEntity) map[uint32][]uint32 {
	kpiNodeParentChildrenMap := make(map[uint32][]uint32)
	for _, kpiNodeEntity := range kpiNodeEntities {
		sharedUtils.NewOptionalFromPointer(kpiNodeEntity.ParentNodeID).DoIfPresent(func(parentNodeID uint32) {
			kpiNodeParentChildrenMap[parentNodeID] = append(kpiNodeParentChildrenMap[parentNodeID], kpiNodeEntity.ID)
		})
	}
	return kpiNodeParentChildrenMap
}

func ToDLLModelKPIDefinition(kpiDefinitionEntity dbModel.KPIDefinitionEntity, kpiNodeEntities []dbModel.KPINodeEntity, logicalOperationKPINodeEntities []dbModel.LogicalOperationKPINodeEntity, atomKPINodeEntities []dbModel.AtomKPINodeEntity) sharedModel.KPIDefinition {
	kpiDefinitionRootOptional := sharedUtils.NewOptionalOf(reconstructKPINodeTree(*kpiDefinitionEntity.RootNodeID, prepareKPINodeParentChildrenMap(kpiNodeEntities), logicalOperationKPINodeEntities, atomKPINodeEntities))
	return sharedModel.KPIDefinition{
		ID:                  &kpiDefinitionEntity.ID,
		UID:                 kpiDefinitionEntity.UID,
		Label:               kpiDefinitionEntity.Label,
		SDTypeID:            kpiDefinitionEntity.SDTypeID,
		SDTypeSpecification: kpiDefinitionEntity.SDType.UID,
		UserIdentifier:      kpiDefinitionEntity.UserIdentifier,
		RootNode:            kpiDefinitionRootOptional.GetPayload(),
		SDInstanceMode:      sharedModel.SDInstanceMode(kpiDefinitionEntity.SDInstanceMode),
		SelectedSDInstanceIDs: sharedUtils.Map(kpiDefinitionEntity.SDInstanceKPIDefinitionRelationshipRecords, func(sdInstanceKPIDefinitionRelationshipEntity dbModel.SDInstanceKPIDefinitionRelationshipEntity) uint32 {
			return sdInstanceKPIDefinitionRelationshipEntity.SDInstanceID
		}),
		SelectedSDInstanceUIDs: sharedUtils.Map(kpiDefinitionEntity.SDInstanceKPIDefinitionRelationshipRecords, func(sdInstanceKPIDefinitionRelationshipEntity dbModel.SDInstanceKPIDefinitionRelationshipEntity) string {
			return sdInstanceKPIDefinitionRelationshipEntity.SDInstanceUID
		}),
	}
}
