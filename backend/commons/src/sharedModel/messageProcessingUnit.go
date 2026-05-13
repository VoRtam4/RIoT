/**
 * @file messageProcessingUnit.go
 * @brief Sdílené modely runtime stavu a cache používané modulem Message Processing Unit.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedModel

import "time"

type KPIKey struct {
	SDInstanceUID   string
	KPIDefinitionID uint32
}

type RawState struct {
	Values         map[string]interface{}
	EventTime      time.Time
	SynchronizedAt time.Time
}

type KPIState struct {
	Value          bool
	EventTime      time.Time
	SynchronizedAt time.Time
}

type SDTypeDefinitionCache struct {
	Label      string
	Parameters map[string]SDParameter
}

type ProcessResult struct {
	Values        map[string]interface{}
	Changed       bool
	SDTypeChanged bool
}
