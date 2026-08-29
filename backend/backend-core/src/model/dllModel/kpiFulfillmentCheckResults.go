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
	SDTypeID         uint32
	SDTypeUID        string
	KPIDefinitionID  uint32
	KPIDefinitionUID string
	SDInstanceID     uint32
	SDInstanceUID    string
	Fulfilled        bool
	EventTime        time.Time
}

type KPIFulfillmentCheckResultRequest struct {
	KPIDefinitionID uint32
	SDInstanceID    uint32
}
