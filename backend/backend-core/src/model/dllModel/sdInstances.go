/**
 * @file sdInstances.go
 * @brief Doménový model instancí zdrojů dat.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní model instancí zdrojů dat.
 * - Vojtěch Hubáček: doplnění labelů.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import "github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"

type SDInstance struct {
	ID              sharedUtils.Optional[uint32]
	UID             string
	Label           string
	ConfirmedByUser bool
	UserIdentifier  string
	SDType          SDType
}
