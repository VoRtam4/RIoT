/**
 * @file typeSystem.go
 * @brief Pomocné funkce pro práci s typovou kontrolou generických hodnot.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

func TypeIs[T any](subject any) bool {
	_, ok := subject.(T)
	return ok
}
