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
