/**
 * @file generators.go
 * @brief Pomocné generátory hodnot pro sdílené použití v modulech.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

import "sync/atomic"

const PublicUIDRandomPartLength = 16

func SequentialNumberGenerator() func() uint32 {
	counter := uint32(0)
	return func() uint32 {
		return atomic.AddUint32(&counter, 1)
	}
}

func GeneratePublicUID(prefix string) string {
	return prefix + ":" + GenerateRandomAlphanumericString(PublicUIDRandomPartLength)
}
