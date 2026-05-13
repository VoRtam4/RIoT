/**
 * @file rawDataPoint.go
 * @brief Mapování raw datových bodů z doménového modelu do GraphQL modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelRawDataPoint(rawDataPoint dllModel.RawDataPoint) graphQLModel.RawDataPoint {
	return graphQLModel.RawDataPoint{
		SdTypeID:     rawDataPoint.SDTypeID,
		SdInstanceID: rawDataPoint.SDInstanceID,
		Payload:      string(rawDataPoint.Payload),
		EventTime:    rawDataPoint.EventTime.Format(time.RFC3339Nano),
	}
}
