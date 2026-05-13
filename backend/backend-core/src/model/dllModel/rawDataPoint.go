/**
 * @file rawDataPoint.go
 * @brief Doménový model surového datového bodu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import "time"

type RawDataPoint struct {
	SDTypeID     uint32
	SDInstanceID uint32
	EventTime    time.Time
	Payload      []byte
}
