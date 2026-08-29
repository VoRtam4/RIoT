/**
 * @file sdInstanceGroups.go
 * @brief Doménový model skupin instancí zdrojů dat.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní model skupin instancí zdrojů dat.
 * - Vojtěch Hubáček: doplnění labelů.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import "github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"

type SDInstanceGroup struct {
	ID             sharedUtils.Optional[uint32]
	UID            string
	Label          string
	UserIdentifier string
	SDInstanceIDs  []uint32
	SDInstanceUIDs []string
}
