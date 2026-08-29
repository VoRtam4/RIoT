/**
 * @file sdInstances.go
 * @brief Mapování instancí zdrojů dat z doménového modelu do GraphQL modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování instancí zdrojů dat.
 * - Vojtěch Hubáček: doplnění mapování labelů.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelSDInstance(sdInstance dllModel.SDInstance) graphQLModel.SDInstance {
	return graphQLModel.SDInstance{
		UID:             sdInstance.UID,
		Label:           sdInstance.Label,
		ConfirmedByUser: sdInstance.ConfirmedByUser,
		UserIdentifier:  sdInstance.UserIdentifier,
		Type:            ToGraphQLModelSDType(sdInstance.SDType),
	}
}
