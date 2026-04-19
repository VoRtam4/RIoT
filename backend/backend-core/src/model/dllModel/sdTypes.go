package dllModel

import "github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"

type SDParameterType string

const (
	SDParameterTypeString  SDParameterType = "string"
	SDParameterTypeNumber  SDParameterType = "number"
	SDParameterTypeBoolean SDParameterType = "boolean"
)

type SDParameterRole string

const (
	SDParameterRoleField SDParameterRole = "field"
	SDParameterRoleTag   SDParameterRole = "tag"
)

type SDParameter struct {
	ID         sharedUtils.Optional[uint32]
	Label      string
	Denotation string
	Type       SDParameterType
	Role       SDParameterRole
}

type SDType struct {
	ID         sharedUtils.Optional[uint32]
	UID        string
	Label      string
	Parameters []SDParameter
}
