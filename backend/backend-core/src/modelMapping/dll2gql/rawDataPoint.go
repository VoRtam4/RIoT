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
