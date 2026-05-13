/**
 * @file readJobs.go
 * @brief Evidence aktivních čtecích úloh a jejich rušení během dlouhých exportů dat.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_time_series_store
 */
package internal

import (
	"context"
	"sync"
)

var (
	activeReadJobs    sync.Map
	cancelledReadJobs sync.Map
)

func RegisterReadJob(jobID uint32, cancel context.CancelFunc) {
	if jobID == 0 {
		return
	}
	activeReadJobs.Store(jobID, cancel)
	if _, cancelled := cancelledReadJobs.Load(jobID); cancelled {
		cancel()
	}
}

func UnregisterReadJob(jobID uint32) {
	if jobID == 0 {
		return
	}
	activeReadJobs.Delete(jobID)
	cancelledReadJobs.Delete(jobID)
}

func CancelReadJob(jobID uint32) {
	if jobID == 0 {
		return
	}
	cancelledReadJobs.Store(jobID, struct{}{})
	if cancelFn, ok := activeReadJobs.Load(jobID); ok {
		cancelFn.(context.CancelFunc)()
	}
}
