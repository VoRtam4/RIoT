package domainLogicLayer

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func HistoryRead(sdTypeInput graphQLModel.SDTypeInput) sharedUtils.Result[graphQLModel.SDType] {
	sdType := gql2dll.ToDLLModelSDType(sdTypeInput)
	persistResult := dbClient.GetRelationalDatabaseClientInstance().PersistSDType(sdType)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDType](persistResult.GetError())
	}
	isc.EnqueueMessageRepresentingCurrentSDTypeConfiguration(getDLLRabbitMQClient())
	return sharedUtils.NewSuccessResult[graphQLModel.SDType](dll2gql.ToGraphQLModelSDType(persistResult.GetPayload()))
}
