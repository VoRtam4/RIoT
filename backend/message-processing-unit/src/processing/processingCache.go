package processing

import (
	"log"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

var (
	lastRaw sync.Map
	lastKPI sync.Map
)

func StartCacheCleanup(interval time.Duration, ttl time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()
			lastRaw.Range(func(key, value any) bool {
				state := value.(sharedModel.RawState)
				if now.Sub(state.SynchronizedAt) > ttl {
					lastRaw.Delete(key)
					log.Printf("[RAW][CLEANUP] Removed stale entry | uid=%v", key)
				}
				return true
			})
			lastKPI.Range(func(key, value any) bool {
				state := value.(sharedModel.KPIState)

				if now.Sub(state.SynchronizedAt) > ttl {
					lastKPI.Delete(key)
					log.Printf("[KPI][CLEANUP] Removed stale entry | key=%v", key)
				}
				return true
			})
		}
	}()
}
