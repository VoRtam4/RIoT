/**
 * @file reprocessSync.go
 * @brief Synchronizace aktivních přepočtů KPI s požadavky na rušení a mazání výsledků.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"fmt"
	"log"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func ProcessDelete(req sharedModel.KPIDeleteResultsRequestISCMessage) error {
	if req.KPIDefinitionID != 0 {
		key := makeKey(req.SDTypeUID, req.KPIDefinitionID)
		activeReprocessJobs.Delete(key)
		log.Printf("[MPU][REPROCESS] Delete request received -> cancelling active reprocess | kpiID=%d sdType=%s key=%s", req.KPIDefinitionID, req.SDTypeUID, key)
	}
	return nil
}

func makeKey(sdType string, kpiID uint32) string {
	return fmt.Sprintf("%s|%d", sdType, kpiID)
}

func isActive(jobKey string, jobID string) bool {
	v, ok := activeReprocessJobs.Load(jobKey)
	if !ok {
		return false
	}
	job, ok := v.(string)
	return ok && job == jobID
}

func waitForConfig(jobID string) {
	if jobID == "" {
		return
	}
	chIface, _ := jobSyncMap.LoadOrStore(jobID, make(chan struct{}))
	ch := chIface.(chan struct{})
	log.Printf("[MPU][SYNC] Waiting for config | jobID=%s", jobID)
	select {
	case <-ch:
		log.Printf("[MPU][SYNC] Config ready | jobID=%s", jobID)
	case <-time.After(5 * time.Second):
		log.Printf("[MPU][SYNC] Config wait timeout, proceeding with current cache | jobID=%s", jobID)
	}
}

func signalConfigReady(jobID string) {
	if jobID == "" {
		return
	}
	chIface, _ := jobSyncMap.LoadOrStore(jobID, make(chan struct{}))
	ch := chIface.(chan struct{})
	select {
	case <-ch:
	default:
		close(ch)
	}
	log.Printf("[MPU][SYNC] Config signaled | jobID=%s", jobID)
}
