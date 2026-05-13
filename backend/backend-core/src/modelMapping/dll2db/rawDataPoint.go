/**
 * @file rawDataPoint.go
 * @brief Mapování raw datových bodů z doménového modelu do databázového modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelRawDataPoint(raw dllModel.RawDataPoint) dbModel.RawDataPointEntity {
	return dbModel.RawDataPointEntity{
		SDInstanceID: raw.SDInstanceID,
		EventTime:    raw.EventTime,
		Payload:      raw.Payload,
	}
}
