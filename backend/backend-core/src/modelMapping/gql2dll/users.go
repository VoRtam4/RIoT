/**
 * @file users.go
 * @brief Mapování uživatelské konfigurace z GraphQL modelu do doménového modelu.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToDLLModelUserConfig(userConfigInput graphQLModel.UserConfigInput) dllModel.UserConfig {
	return dllModel.UserConfig{
		UserID: 0,
		Config: userConfigInput.Config,
	}
}
