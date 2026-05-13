/**
 * @file ternaryOperator.go
 * @brief Generická pomocná funkce pro ternární výběr hodnoty.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

func Ternary[T any](cond bool, r1 T, r2 T) T {
	if cond {
		return r1
	} else {
		return r2
	}
}
