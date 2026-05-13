/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Doménový model výsledků vyhodnocení KPI.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní model výsledků vyhodnocení KPI.
 * - Vojtěch Hubáček: doplnění času události EventTime.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import "time"

type KPIFulfillmentCheckResult struct {
	SDTypeID        uint32
	KPIDefinitionID uint32
	SDInstanceID    uint32
	Fulfilled       bool
	EventTime       time.Time
}

type KPIFulfillmentCheckResultRequest struct {
	KPIDefinitionID uint32
	SDInstanceID    uint32
}
