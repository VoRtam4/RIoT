/**
 * @file slice.go
 * @brief Pomocné funkce pro konstrukci prázdných nebo předvyplněných slice hodnot.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

func SliceOf[T any](items ...T) []T {
	return items
}

func EmptySlice[T any]() []T {
	return make([]T, 0)
}
