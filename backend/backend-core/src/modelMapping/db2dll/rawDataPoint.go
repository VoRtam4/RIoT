/**
 * @file rawDataPoint.go
 * @brief Mapování raw datových bodů z databázového modelu do doménového modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDLLModelRawDataPoint(rawDataPointEntity dbModel.RawDataPointEntity) dllModel.RawDataPoint {
	return dllModel.RawDataPoint{
		SDTypeID:     rawDataPointEntity.SDInstance.SDTypeID,
		SDInstanceID: rawDataPointEntity.SDInstanceID,
		EventTime:    rawDataPointEntity.EventTime,
		Payload:      rawDataPointEntity.Payload,
	}
}
