/**
 * @file pair.go
 * @brief Generická dvojice hodnot pro sdílené předávání jednoduchých dvojic dat.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

type Pair[T any, U any] struct {
	e1 T
	e2 U
}

func NewPairOf[T any, U any](e1 T, e2 U) Pair[T, U] {
	return Pair[T, U]{
		e1: e1,
		e2: e2,
	}
}

func (p Pair[T, U]) GetFirst() T {
	return p.e1
}

func (p Pair[T, U]) GetSecond() U {
	return p.e2
}
