/**
 * @file kpiDefinitionModel.go
 * @brief Sdílený model stromové definice KPI a jejích logických a porovnávacích operátorů.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní model stromové KPI definice a základních operátorů.
 * - Vojtěch Hubáček: rozšíření modelu o nové KPI operace !=, exist, not exist a NOT.
 *
 * @ingroup riot_commons
 */
package sharedModel

type LogicalOperationNodeType string

const (
	AND LogicalOperationNodeType = "and"
	OR  LogicalOperationNodeType = "or"
	NOR LogicalOperationNodeType = "nor"
	NOT LogicalOperationNodeType = "not"
)

type KPINodeType string

const (
	StringEQAtom        KPINodeType = "string_eq_atom"
	StringNEQAtom       KPINodeType = "string_neq_atom"
	StringExistsAtom    KPINodeType = "string_exists_atom"
	StringNotExistsAtom KPINodeType = "string_not_exists_atom"

	BooleanEQAtom        KPINodeType = "boolean_eq_atom"
	BooleanNEQAtom       KPINodeType = "boolean_neq_atom"
	BooleanExistsAtom    KPINodeType = "boolean_exists_atom"
	BooleanNotExistsAtom KPINodeType = "boolean_not_exists_atom"

	NumericEQAtom        KPINodeType = "numeric_eq_atom"
	NumericNEQAtom       KPINodeType = "numeric_neq_atom"
	NumericGTAtom        KPINodeType = "numeric_gt_atom"
	NumericGEQAtom       KPINodeType = "numeric_geq_atom"
	NumericLTAtom        KPINodeType = "numeric_lt_atom"
	NumericLEQAtom       KPINodeType = "numeric_leq_atom"
	NumericExistsAtom    KPINodeType = "numeric_exists_atom"
	NumericNotExistsAtom KPINodeType = "numeric_not_exists_atom"

	LogicalOperation KPINodeType = "logical_operation"
)

type SDInstanceMode string

const (
	ALL      SDInstanceMode = "all"
	SELECTED SDInstanceMode = "selected"
)

type KPIReferenceMode string

const (
	KPIReferenceModeLiteral   KPIReferenceMode = "literal"
	KPIReferenceModeParameter KPIReferenceMode = "parameter"
)

type KPIDefinition struct {
	ID                     *uint32        `json:"id,omitempty"`
	UID                    *string        `json:"uid,omitempty"`
	Label                  string         `json:"label"`
	UserID                 *uint32        `json:"userID,omitempty"`
	SDTypeID               uint32         `json:"sdTypeID"`
	SDTypeSpecification    string         `json:"sdTypeSpecification"`
	UserIdentifier         string         `json:"userIdentifier"`
	RootNode               KPINode        `json:"rootNode"`
	SDInstanceMode         SDInstanceMode `json:"sdInstanceMode"`
	SelectedSDInstanceIDs  []uint32       `json:"selectedSDInstanceIDs"`
	SelectedSDInstanceUIDs []string       `json:"selectedSDInstanceUIDs"`
}

type KPIDefinitionMPU struct {
	UID                    string         `json:"uid"`
	SDTypeUID              string         `json:"sdTypeUID"`
	RootNode               KPINode        `json:"rootNode"`
	SDInstanceMode         SDInstanceMode `json:"sdInstanceMode"`
	SelectedSDInstanceUIDs []string       `json:"selectedSDInstanceUIDs"`
}

type KPINode interface {
	GetType() KPINodeType
}

type StringEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   string           `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*StringEQAtomKPINode) GetType() KPINodeType {
	return StringEQAtom
}

type StringNEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   string           `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*StringNEQAtomKPINode) GetType() KPINodeType {
	return StringNEQAtom
}

type StringExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*StringExistsAtomKPINode) GetType() KPINodeType {
	return StringExistsAtom
}

type StringNotExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*StringNotExistsAtomKPINode) GetType() KPINodeType {
	return StringNotExistsAtom
}

type BooleanEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   bool             `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*BooleanEQAtomKPINode) GetType() KPINodeType {
	return BooleanEQAtom
}

type BooleanNEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   bool             `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*BooleanNEQAtomKPINode) GetType() KPINodeType {
	return BooleanNEQAtom
}

type BooleanExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*BooleanExistsAtomKPINode) GetType() KPINodeType {
	return BooleanExistsAtom
}

type BooleanNotExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*BooleanNotExistsAtomKPINode) GetType() KPINodeType {
	return BooleanNotExistsAtom
}

type NumericEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericEQAtomKPINode) GetType() KPINodeType {
	return NumericEQAtom
}

type NumericNEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericNEQAtomKPINode) GetType() KPINodeType {
	return NumericNEQAtom
}

type NumericGTAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericGTAtomKPINode) GetType() KPINodeType {
	return NumericGTAtom
}

type NumericGEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericGEQAtomKPINode) GetType() KPINodeType {
	return NumericGEQAtom
}

type NumericLTAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericLTAtomKPINode) GetType() KPINodeType {
	return NumericLTAtom
}

type NumericLEQAtomKPINode struct {
	SDParameterSpecification         string           `json:"sdParameterSpecification"`
	ReferenceValue                   float64          `json:"referenceValue"`
	ReferenceMode                    KPIReferenceMode `json:"referenceMode,omitempty"`
	ComparedSDParameterSpecification string           `json:"comparedSDParameterSpecification,omitempty"`
	ComparedRecordOffset             *int             `json:"comparedRecordOffset,omitempty"`
}

func (*NumericLEQAtomKPINode) GetType() KPINodeType {
	return NumericLEQAtom
}

type NumericExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*NumericExistsAtomKPINode) GetType() KPINodeType {
	return NumericExistsAtom
}

type NumericNotExistsAtomKPINode struct {
	SDParameterSpecification string `json:"sdParameterSpecification"`
}

func (*NumericNotExistsAtomKPINode) GetType() KPINodeType {
	return NumericNotExistsAtom
}

type LogicalOperationKPINode struct {
	Type       LogicalOperationNodeType `json:"type"`
	ChildNodes []KPINode                `json:"childNodes"`
}

func (*LogicalOperationKPINode) GetType() KPINodeType {
	return LogicalOperation
}
