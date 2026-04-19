package processing

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func CheckKPIFulfillment(kpiDefinition sharedModel.KPIDefinitionMPU, sdParameterValueMap *any) sharedUtils.Result[bool] {
	return checkKPINodeFulfillment(kpiDefinition.RootNode, sdParameterValueMap)
}

func checkKPINodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	return getCheckerForKPINode(kpiNode).checkNodeFulfillment(kpiNode, sdParameterValueMap)
}

type kpiNodeFulfillmentChecker interface {
	checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool]
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

func (_ *stringEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	stringEQAtomKPINode := kpiNode.(*sharedModel.StringEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[string](sdParameterValueMap, stringEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer == stringEQAtomKPINode.ReferenceValue)
}

func (_ *stringNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	stringNEQAtomKPINode := kpiNode.(*sharedModel.StringNEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[string](sdParameterValueMap, stringNEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer != stringNEQAtomKPINode.ReferenceValue)
}

func (_ *stringExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	stringExistsAtomKPINode := kpiNode.(*sharedModel.StringExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[stringExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *stringNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	stringNotExistsAtomKPINode := kpiNode.(*sharedModel.StringNotExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[stringNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func (_ *booleanEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	booleanEQAtomKPINode := kpiNode.(*sharedModel.BooleanEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[bool](sdParameterValueMap, booleanEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer == booleanEQAtomKPINode.ReferenceValue)
}

func (_ *booleanNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	booleanEQAtomKPINode := kpiNode.(*sharedModel.BooleanNEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[bool](sdParameterValueMap, booleanEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer != booleanEQAtomKPINode.ReferenceValue)
}

func (_ *booleanExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	booleanExistsAtomKPINode := kpiNode.(*sharedModel.BooleanExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[booleanExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *booleanNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	booleanNotExistsAtomKPINode := kpiNode.(*sharedModel.BooleanNotExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[booleanNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func (_ *numericEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericEQAtomKPINode := kpiNode.(*sharedModel.NumericEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer == numericEQAtomKPINode.ReferenceValue)
}

func (_ *numericNEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericNEQAtomKPINode := kpiNode.(*sharedModel.NumericNEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericNEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer != numericNEQAtomKPINode.ReferenceValue)
}

func (_ *numericGTKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericGTAtomKPINode := kpiNode.(*sharedModel.NumericGTAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericGTAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer > numericGTAtomKPINode.ReferenceValue)
}

func (_ *numericGEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericGEQAtomKPINode := kpiNode.(*sharedModel.NumericGEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericGEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer >= numericGEQAtomKPINode.ReferenceValue)
}

func (_ *numericLTKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericLTAtomKPINode := kpiNode.(*sharedModel.NumericLTAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericLTAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer < numericLTAtomKPINode.ReferenceValue)
}

func (_ *numericLEQKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericLEQAtomKPINode := kpiNode.(*sharedModel.NumericLEQAtomKPINode)
	actualSDParameterValueResult := getSDParameterValue[float64](sdParameterValueMap, numericLEQAtomKPINode.SDParameterSpecification)
	if actualSDParameterValueResult.IsFailure() {
		return sharedUtils.NewFailureResult[bool](actualSDParameterValueResult.GetError())
	}
	actualSDParameterValuePointer := actualSDParameterValueResult.GetPayload()
	if actualSDParameterValuePointer == nil {
		return sharedUtils.NewSuccessResult[bool](false)
	}
	return sharedUtils.NewSuccessResult[bool](*actualSDParameterValuePointer <= numericLEQAtomKPINode.ReferenceValue)
}

func (_ *numericExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericExistsAtomKPINode := kpiNode.(*sharedModel.NumericExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[numericExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](exists && val != nil)
}

func (_ *numericNotExistsKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	numericNotExistsAtomKPINode := kpiNode.(*sharedModel.NumericNotExistsAtomKPINode)
	val, exists := (*sdParameterValueMap).(map[string]any)[numericNotExistsAtomKPINode.SDParameterSpecification]
	return sharedUtils.NewSuccessResult[bool](!exists || val == nil)
}

func getSDParameterValue[T any](sdParameterValueMap *any, sdParameterSpecification string) sharedUtils.Result[*T] {
	sdParameterValue, exists := (*sdParameterValueMap).(map[string]any)[sdParameterSpecification]
	if !exists || sdParameterValue == nil {
		return sharedUtils.NewSuccessResult[*T](nil)
	}
	if !sharedUtils.TypeIs[T](sdParameterValue) {
		return sharedUtils.NewFailureResult[*T](fmt.Errorf("the actual SD parameter '%s' is of different type that the one specified", sdParameterSpecification))
	}
	value := sdParameterValue.(T)
	return sharedUtils.NewSuccessResult[*T](&value)
}

func (_ *logicalOperationKPINodeFulfillmentChecker) checkNodeFulfillment(kpiNode sharedModel.KPINode, sdParameterValueMap *any) sharedUtils.Result[bool] {
	logicalOperationKPINode := kpiNode.(*sharedModel.LogicalOperationKPINode)
	switch logicalOperationKPINode.Type {
	case sharedModel.AND:
		for _, node := range logicalOperationKPINode.ChildNodes {
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, sdParameterValueMap)
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
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, sdParameterValueMap)
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
			kpiNodeFulfillmentResult := checkKPINodeFulfillment(node, sdParameterValueMap)
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
		childResult := checkKPINodeFulfillment(logicalOperationKPINode.ChildNodes[0], sdParameterValueMap)
		if childResult.IsFailure() {
			return sharedUtils.NewFailureResult[bool](childResult.GetError())
		}
		return sharedUtils.NewSuccessResult[bool](!childResult.GetPayload())
	default:
		panic(fmt.Errorf("unsupported type of logical operation: %s", logicalOperationKPINode.Type))
	}
}
