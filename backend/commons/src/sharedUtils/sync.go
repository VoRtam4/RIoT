/**
 * @file sync.go
 * @brief Pomocná synchronizační funkce pro paralelní spuštění více operací a čekání na jejich dokončení.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedUtils

import "sync"

func WaitForAll(functions ...func()) {
	wg := new(sync.WaitGroup)
	wg.Add(len(functions))
	for _, function := range functions {
		go func(function func()) {
			defer wg.Done()
			function()
		}(function)
	}
	wg.Wait()
}
