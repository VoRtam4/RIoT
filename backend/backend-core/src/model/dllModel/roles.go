/**
 * @file roles.go
 * @brief Doménový model uživatelských rolí a oprávnění.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dllModel

type Role struct {
	ID          uint32
	UID         string
	Label       string
	System      bool
	Permissions []Permission
}

type Permission struct {
	UID   string
	Label string
}
