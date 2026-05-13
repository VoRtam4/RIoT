/**
 * @file debug.go
 * @brief Pomocné ladicí funkce pro výpis hodnot v běhu aplikace.
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
	"bytes"
	"log"

	"github.com/davecgh/go-spew/spew"
)

func Dump(a ...any) {
	buffer := new(bytes.Buffer)
	spew.Fdump(buffer, a)
	log.Println(buffer.String())
}
