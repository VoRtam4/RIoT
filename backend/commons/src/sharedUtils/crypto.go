/**
 * @file crypto.go
 * @brief Pomocné kryptografické funkce pro tvorbu hashů.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHexHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
