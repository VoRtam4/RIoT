/**
 * @file kpiDefinitions.go
 * @brief Mapování doménové reprezentace KPI definic do databázového stromu entit.
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
package dll2db

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func ToDBModelEntitiesKPIDefinition(kpiDefinition sharedModel.KPIDefinition, sdParameterIDsBySpecification map[string]uint32) (*dbModel.KPINodeEntity, []*dbModel.KPINodeEntity, []dbModel.LogicalOperationKPINodeEntity, []dbModel.AtomKPINodeEntity, error) {
	return transformKPIDefinitionTree(kpiDefinition.RootNode, nil, make([]*dbModel.KPINodeEntity, 0), make([]dbModel.LogicalOperationKPINodeEntity, 0), make([]dbModel.AtomKPINodeEntity, 0), sdParameterIDsBySpecification)
}

func atomNodeEntity(node *dbModel.KPINodeEntity, sdParameterID uint32, atomType string, referenceModeValue sharedModel.KPIReferenceMode, comparedSDParameterID *uint32, comparedRecordOffset *int) dbModel.AtomKPINodeEntity {
	referenceMode := string(referenceModeValue)
	if referenceMode == "" {
		referenceMode = string(sharedModel.KPIReferenceModeLiteral)
	}
	return dbModel.AtomKPINodeEntity{
		Node:                  node,
		SDParameterID:         sdParameterID,
		Type:                  atomType,
		ReferenceMode:         referenceMode,
		ComparedSDParameterID: comparedSDParameterID,
		ComparedRecordOffset:  comparedRecordOffset,
	}
}

func resolveSDParameterID(sdParameterSpecification string, sdParameterIDsBySpecification map[string]uint32) (uint32, error) {
	sdParameterID, exists := sdParameterIDsBySpecification[sdParameterSpecification]
	if !exists {
		return 0, fmt.Errorf("unknown SD parameter specification in KPI definition: %s", sdParameterSpecification)
	}
	return sdParameterID, nil
}

func resolveComparedSDParameterID(referenceMode sharedModel.KPIReferenceMode, comparedSDParameterSpecification string, sdParameterIDsBySpecification map[string]uint32) (*uint32, error) {
	if referenceMode == "" || referenceMode == sharedModel.KPIReferenceModeLiteral {
		return nil, nil
	}
	if referenceMode != sharedModel.KPIReferenceModeParameter {
		return nil, fmt.Errorf("unsupported KPI reference mode: %s", referenceMode)
	}
	sdParameterID, err := resolveSDParameterID(comparedSDParameterSpecification, sdParameterIDsBySpecification)
	if err != nil {
		return nil, err
	}
	return &sdParameterID, nil
}

func atomNodeEntityFromSpecification(node *dbModel.KPINodeEntity, sdParameterSpecification string, atomType string, referenceModeValue sharedModel.KPIReferenceMode, comparedSDParameterSpecification string, comparedRecordOffset *int, sdParameterIDsBySpecification map[string]uint32) (dbModel.AtomKPINodeEntity, error) {
	sdParameterID, err := resolveSDParameterID(sdParameterSpecification, sdParameterIDsBySpecification)
	if err != nil {
		return dbModel.AtomKPINodeEntity{}, err
	}
	comparedSDParameterID, err := resolveComparedSDParameterID(referenceModeValue, comparedSDParameterSpecification, sdParameterIDsBySpecification)
	if err != nil {
		return dbModel.AtomKPINodeEntity{}, err
	}
	return atomNodeEntity(node, sdParameterID, atomType, referenceModeValue, comparedSDParameterID, comparedRecordOffset), nil
}

func transformKPIDefinitionTree(node sharedModel.KPINode, parentKPINodeEntity *dbModel.KPINodeEntity, kpiNodeEntities []*dbModel.KPINodeEntity, logicalOperationNodeEntities []dbModel.LogicalOperationKPINodeEntity, atomNodeEntities []dbModel.AtomKPINodeEntity, sdParameterIDsBySpecification map[string]uint32) (*dbModel.KPINodeEntity, []*dbModel.KPINodeEntity, []dbModel.LogicalOperationKPINodeEntity, []dbModel.AtomKPINodeEntity, error) {
	currentNodeEntity := &dbModel.KPINodeEntity{
		ParentNode: parentKPINodeEntity,
	}
	kpiNodeEntities = append(kpiNodeEntities, currentNodeEntity)
	switch typedNode := node.(type) {
	case *sharedModel.LogicalOperationKPINode:
		logicalOperationNodeEntities = append(logicalOperationNodeEntities, dbModel.LogicalOperationKPINodeEntity{
			Node: currentNodeEntity,
			Type: string(typedNode.Type),
		})
		for _, childNode := range typedNode.ChildNodes {
			var err error
			_, kpiNodeEntities, logicalOperationNodeEntities, atomNodeEntities, err = transformKPIDefinitionTree(childNode, currentNodeEntity, kpiNodeEntities, logicalOperationNodeEntities, atomNodeEntities, sdParameterIDsBySpecification)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		}
	case *sharedModel.StringEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "string_eq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.StringReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.StringNEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "string_neq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.StringReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.StringExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "string_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.StringNotExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "string_not_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.BooleanEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "boolean_eq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.BooleanReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.BooleanNEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "boolean_neq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.BooleanReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.BooleanExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "boolean_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.BooleanNotExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "boolean_not_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_eq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericNEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_neq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericLTAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_lt", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericLEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_leq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericGTAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_gt", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericGEQAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_geq", typedNode.ReferenceMode, typedNode.ComparedSDParameterSpecification, typedNode.ComparedRecordOffset, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		entity.NumericReferenceValue = &typedNode.ReferenceValue
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	case *sharedModel.NumericNotExistsAtomKPINode:
		entity, err := atomNodeEntityFromSpecification(currentNodeEntity, typedNode.SDParameterSpecification, "numeric_not_exists", sharedModel.KPIReferenceModeLiteral, "", nil, sdParameterIDsBySpecification)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		atomNodeEntities = append(atomNodeEntities, entity)
	}
	return currentNodeEntity, kpiNodeEntities, logicalOperationNodeEntities, atomNodeEntities, nil
}
