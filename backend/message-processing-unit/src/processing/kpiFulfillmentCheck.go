/**
 * @file kpiFulfillmentCheck.go
 * @brief Vyhodnocení stromu KPI podmínek nad aktuálními parametry zdroje dat.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní implementace vyhodnocování KPI podmínek.
 * - Vojtěch Hubáček: doplnění operací !=, exist, not exist a logické negace NOT.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type KPIEvaluationContext struct {
	CurrentValues  map[string]interface{}
	PreviousValues map[string]interface{}
}

func CheckKPIFulfillment(kpiDefinition sharedModel.KPIDefinitionMPU, context KPIEvaluationContext) sharedUtils.Result[bool] {
	return checkKPINodeFulfillment(kpiDefinition.RootNode, context)
}

func checkKPINodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	return getCheckerForKPINode(kpiNode).checkNodeFulfillment(kpiNode, context)
}

type kpiNodeFulfillmentChecker interface {
	checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool]
}

func getCheckerForKPINode(kpiNode sharedModel.KPINode) kpiNodeFulfillmentChecker {
	switch kpiNode.(type) {
	case *sharedModel.StringEQAtomKPINode:
		return &stringEQKPINodeFulfillmentChecker{}
	case *sharedModel.StringNEQAtomKPINode:
		return &stringNEQKPINodeFulfillmentChecker{}
	case *sharedModel.StringExistsAtomKPINode:
		return &stringExistsKPINodeFulfillmentChecker{}
	case *sharedModel.StringNotExistsAtomKPINode:
		return &stringNotExistsKPINodeFulfillmentChecker{}

	case *sharedModel.BooleanEQAtomKPINode:
		return &booleanEQKPINodeFulfillmentChecker{}
	case *sharedModel.BooleanNEQAtomKPINode:
		return &booleanNEQKPINodeFulfillmentChecker{}
	case *sharedModel.BooleanExistsAtomKPINode:
		return &booleanExistsKPINodeFulfillmentChecker{}
	case *sharedModel.BooleanNotExistsAtomKPINode:
		return &booleanNotExistsKPINodeFulfillmentChecker{}

	case *sharedModel.NumericEQAtomKPINode:
		return &numericEQKPINodeFulfillmentChecker{}
	case *sharedModel.NumericNEQAtomKPINode:
		return &numericNEQKPINodeFulfillmentChecker{}
	case *sharedModel.NumericGTAtomKPINode:
		return &numericGTKPINodeFulfillmentChecker{}
	case *sharedModel.NumericGEQAtomKPINode:
		return &numericGEQKPINodeFulfillmentChecker{}
	case *sharedModel.NumericLTAtomKPINode:
		return &numericLTKPINodeFulfillmentChecker{}
	case *sharedModel.NumericLEQAtomKPINode:
		return &numericLEQKPINodeFulfillmentChecker{}
	case *sharedModel.LogicalOperationKPINode:
		return &logicalOperationKPINodeFulfillmentChecker{}
	case *sharedModel.NumericExistsAtomKPINode:
		return &numericExistsKPINodeFulfillmentChecker{}
	case *sharedModel.NumericNotExistsAtomKPINode:
		return &numericNotExistsKPINodeFulfillmentChecker{}
	}
	panic(fmt.Errorf("unsupported type of KPI node: %T", kpiNode))
}

type stringEQKPINodeFulfillmentChecker struct{}
type stringNEQKPINodeFulfillmentChecker struct{}
type stringExistsKPINodeFulfillmentChecker struct{}
type stringNotExistsKPINodeFulfillmentChecker struct{}

type booleanEQKPINodeFulfillmentChecker struct{}
type booleanNEQKPINodeFulfillmentChecker struct{}
type booleanExistsKPINodeFulfillmentChecker struct{}
type booleanNotExistsKPINodeFulfillmentChecker struct{}

type numericEQKPINodeFulfillmentChecker struct{}
type numericNEQKPINodeFulfillmentChecker struct{}
type numericGTKPINodeFulfillmentChecker struct{}
type numericGEQKPINodeFulfillmentChecker struct{}
type numericLTKPINodeFulfillmentChecker struct{}
type numericLEQKPINodeFulfillmentChecker struct{}
type numericExistsKPINodeFulfillmentChecker struct{}
type numericNotExistsKPINodeFulfillmentChecker struct{}

type logicalOperationKPINodeFulfillmentChecker struct{}

func (_ *stringEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	stringEQAtomKPINode := kpiNode.(*sharedModel.StringEQAtomKPINode)
	return compareAtomValues[string](context, stringEQAtomKPINode.SDParameterSpecification, stringEQAtomKPINode.ReferenceMode, stringEQAtomKPINode.ReferenceValue, stringEQAtomKPINode.ComparedSDParameterSpecification, stringEQAtomKPINode.ComparedRecordOffset, func(actual, reference string) bool {
		return actual == reference
	})
}

func (_ *stringNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	stringNEQAtomKPINode := kpiNode.(*sharedModel.StringNEQAtomKPINode)
	return compareAtomValues[string](context, stringNEQAtomKPINode.SDParameterSpecification, stringNEQAtomKPINode.ReferenceMode, stringNEQAtomKPINode.ReferenceValue, stringNEQAtomKPINode.ComparedSDParameterSpecification, stringNEQAtomKPINode.ComparedRecordOffset, func(actual, reference string) bool {
		return actual != reference
	})
}

func (_ *stringExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	stringExistsAtomKPINode := kpiNode.(*sharedModel.StringExistsAtomKPINode)
	val, exists := context.CurrentValues[stringExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *stringNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	stringNotExistsAtomKPINode := kpiNode.(*sharedModel.StringNotExistsAtomKPINode)
	val, exists := context.CurrentValues[stringNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func (_ *booleanEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	booleanEQAtomKPINode := kpiNode.(*sharedModel.BooleanEQAtomKPINode)
	return compareAtomValues[bool](context, booleanEQAtomKPINode.SDParameterSpecification, booleanEQAtomKPINode.ReferenceMode, booleanEQAtomKPINode.ReferenceValue, booleanEQAtomKPINode.ComparedSDParameterSpecification, booleanEQAtomKPINode.ComparedRecordOffset, func(actual, reference bool) bool {
		return actual == reference
	})
}

func (_ *booleanNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	booleanEQAtomKPINode := kpiNode.(*sharedModel.BooleanNEQAtomKPINode)
	return compareAtomValues[bool](context, booleanEQAtomKPINode.SDParameterSpecification, booleanEQAtomKPINode.ReferenceMode, booleanEQAtomKPINode.ReferenceValue, booleanEQAtomKPINode.ComparedSDParameterSpecification, booleanEQAtomKPINode.ComparedRecordOffset, func(actual, reference bool) bool {
		return actual != reference
	})
}

func (_ *booleanExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	booleanExistsAtomKPINode := kpiNode.(*sharedModel.BooleanExistsAtomKPINode)
	val, exists := context.CurrentValues[booleanExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *booleanNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	booleanNotExistsAtomKPINode := kpiNode.(*sharedModel.BooleanNotExistsAtomKPINode)
	val, exists := context.CurrentValues[booleanNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func (_ *numericEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericEQAtomKPINode := kpiNode.(*sharedModel.NumericEQAtomKPINode)
	return compareAtomValues[float64](context, numericEQAtomKPINode.SDParameterSpecification, numericEQAtomKPINode.ReferenceMode, numericEQAtomKPINode.ReferenceValue, numericEQAtomKPINode.ComparedSDParameterSpecification, numericEQAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual == reference
	})
}

func (_ *numericNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericNEQAtomKPINode := kpiNode.(*sharedModel.NumericNEQAtomKPINode)
	return compareAtomValues[float64](context, numericNEQAtomKPINode.SDParameterSpecification, numericNEQAtomKPINode.ReferenceMode, numericNEQAtomKPINode.ReferenceValue, numericNEQAtomKPINode.ComparedSDParameterSpecification, numericNEQAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual != reference
	})
}

func (_ *numericGTKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericGTAtomKPINode := kpiNode.(*sharedModel.NumericGTAtomKPINode)
	return compareAtomValues[float64](context, numericGTAtomKPINode.SDParameterSpecification, numericGTAtomKPINode.ReferenceMode, numericGTAtomKPINode.ReferenceValue, numericGTAtomKPINode.ComparedSDParameterSpecification, numericGTAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual > reference
	})
}

func (_ *numericGEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericGEQAtomKPINode := kpiNode.(*sharedModel.NumericGEQAtomKPINode)
	return compareAtomValues[float64](context, numericGEQAtomKPINode.SDParameterSpecification, numericGEQAtomKPINode.ReferenceMode, numericGEQAtomKPINode.ReferenceValue, numericGEQAtomKPINode.ComparedSDParameterSpecification, numericGEQAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual >= reference
	})
}

func (_ *numericLTKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericLTAtomKPINode := kpiNode.(*sharedModel.NumericLTAtomKPINode)
	return compareAtomValues[float64](context, numericLTAtomKPINode.SDParameterSpecification, numericLTAtomKPINode.ReferenceMode, numericLTAtomKPINode.ReferenceValue, numericLTAtomKPINode.ComparedSDParameterSpecification, numericLTAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual < reference
	})
}

func (_ *numericLEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericLEQAtomKPINode := kpiNode.(*sharedModel.NumericLEQAtomKPINode)
	return compareAtomValues[float64](context, numericLEQAtomKPINode.SDParameterSpecification, numericLEQAtomKPINode.ReferenceMode, numericLEQAtomKPINode.ReferenceValue, numericLEQAtomKPINode.ComparedSDParameterSpecification, numericLEQAtomKPINode.ComparedRecordOffset, func(actual, reference float64) bool {
		return actual <= reference
	})
}

func (_ *numericExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericExistsAtomKPINode := kpiNode.(*sharedModel.NumericExistsAtomKPINode)
	val, exists := context.CurrentValues[numericExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *numericNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	numericNotExistsAtomKPINode := kpiNode.(*sharedModel.NumericNotExistsAtomKPINode)
	val, exists := context.CurrentValues[numericNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func compareAtomValues[T any](context KPIEvaluationContext, actualSDParameterSpecification string, referenceMode sharedModel.KPIReferenceMode, literalReferenceValue T, comparedSDParameterSpecification string, comparedRecordOffset *int, comparison func(T, T) bool) sharedUtils.Result[bool] {
	actualSDParameterValueResult := getSDParameterValue[T](context.CurrentValues, actualSDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	referenceSDParameterValueResult := getReferenceValue[T](context, referenceMode, literalReferenceValue, comparedSDParameterSpecification, comparedRecordOffset)
	if referenceSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](referenceSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	referenceSDParameterValuePointer := referenceSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil || referenceSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](comparison(*actualSDParameterValuePointer, *referenceSDParameterValuePointer))
}

func getReferenceValue[T any](context KPIEvaluationContext, referenceMode sharedModel.KPIReferenceMode, literalReferenceValue T, comparedSDParameterSpecification string, comparedRecordOffset *int) sharedUtils.Result[*T] {
	if referenceMode == "" || referenceMode == sharedModel.KPIReferenceModeLiteral {
		return sharedUtils.NewSuccessResult(&literalReferenceValue)
	}
	if referenceMode != sharedModel.KPIReferenceModeParameter {
		return sharedUtils.NewFailureResult[*T](fmt.Errorf("unsupported KPI reference mode: %s", referenceMode))
	}
	if comparedRecordOffset == nil {
		return sharedUtils.NewFailureResult[*T](fmt.Errorf("missing compared record offset"))
	}
	values, err := valuesByComparedRecordOffset(*comparedRecordOffset, context)
	if err != nil {
		return sharedUtils.NewFailureResult[*T](err)
	}
	return getSDParameterValue[T](values, comparedSDParameterSpecification)
}

func valuesByComparedRecordOffset(comparedRecordOffset int, context KPIEvaluationContext) (map[string]interface{}, error) {
	switch comparedRecordOffset {
	case 0:
		return context.CurrentValues, nil
	case -1:
		return context.PreviousValues, nil
	default:
		return nil, fmt.Errorf("unsupported compared record offset: %d", comparedRecordOffset)
	}
}

func getSDParameterValue[T any](sdParameterValueMap map[string]interface{}, sdParameterSpecification string) sharedUtils.Result[*T] {
	sdParameterValue, exists := sdParameterValueMap[sdParameterSpecification]
	if !exists || sdParameterValue == nil {
		return sharedUtils.NewSuccessResult[*T](nil)
	}
	if !sharedUtils.TypeIs[T](sdParameterValue) {
		return sharedUtils.NewFailureResult[*T](fmt.Errorf("the actual SD parameter '%s' is of different type that the one specified", sdParameterSpecification))
	}
	value := sdParameterValue.(T)
	return sharedUtils.NewSuccessResult[*T](&value)
}

func (_ *logicalOperationKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, context KPIEvaluationContext) sharedUtils.Result[bool] {
	logicalOperationKPINode := kpiNode.(*sharedModel.LogicalOperationKPINode)
	switch logicalOperationKPINode.Type {
	case sharedModel.AND:
		for _, node := range logicalOperationKPINode.ChildNodes {
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, context)
			if kpiNodeFulfillmentResult.IsFailure() {
				return sharedUtils.NewFailureResult[bool](kpiNodeFulfillmentResult.GetError())
			}
			if kpiNodeFulfillment := kpiNodeFulfillmentResult.GetPayload(); !kpiNodeFulfillment {
				return sharedUtils.NewSuccessResult[bool](false)
			}
		}
		return sharedUtils.NewSuccessResult[bool](true)
	case sharedModel.OR:
		for _, node := range logicalOperationKPINode.ChildNodes {
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, context)
			if kpiNodeFulfillmentResult.IsFailure() {
				return sharedUtils.NewFailureResult[bool](kpiNodeFulfillmentResult.GetError())
			}
			if kpiNodeFulfillment := kpiNodeFulfillmentResult.GetPayload(); kpiNodeFulfillment {
				return sharedUtils.NewSuccessResult[bool](true)
			}
		}
		return sharedUtils.NewSuccessResult[bool](false)
	case sharedModel.NOR:
		for _, node := range logicalOperationKPINode.ChildNodes {
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, context)
			if kpiNodeFulfillmentResult.IsFailure() {
				return sharedUtils.NewFailureResult[bool](kpiNodeFulfillmentResult.GetError())
			}
			if kpiNodeFulfillment := kpiNodeFulfillmentResult.GetPayload(); kpiNodeFulfillment {
				return sharedUtils.NewSuccessResult[bool](false)
			}
		}
		return sharedUtils.NewSuccessResult[bool](true)
	case sharedModel.NOT:
		if len(logicalOperationKPINode.ChildNodes) != 1 {
			return sharedUtils.NewFailureResult[bool](fmt.Errorf("NOT operation must have exactly one child"))
		}
		childResult := checkKPINodeFulfillment(logicalOperationKPINode.ChildNodes[0], context)
		if childResult.IsFailure() {
			return sharedUtils.NewFailureResult[bool](childResult.GetError())
		}
		return sharedUtils.NewSuccessResult[bool](!childResult.GetPayload())
	default:
		panic(fmt.Errorf("unsupported type of logical operation: %s", logicalOperationKPINode.Type))
	}
}
