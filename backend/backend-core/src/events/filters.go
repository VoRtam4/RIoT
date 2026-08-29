/**
 * @file filters.go
 * @brief Sestavení filtrů pro odběr backendových událostí podle GraphQL vstupů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package events

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func BuildSDInstanceRegisteredFilter(_ uint32, filter *graphQLModel.SDInstanceRegisteredFilter) (func(graphQLModel.SDInstance) bool, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	if filter == nil {
		return func(_ graphQLModel.SDInstance) bool { return true }, nil
	}
	var allowedInstanceUIDs map[string]struct{}
	var allowedTypeUIDs map[string]struct{}
	instanceTypeUIDs := make(map[string]struct{})
	if len(filter.SdInstanceUIDs) > 0 {
		allowedInstanceUIDs = make(map[string]struct{})
		for _, uid := range filter.SdInstanceUIDs {
			res := db.LoadSDInstanceBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			instanceOptional := res.GetPayload()
			if instanceOptional.IsEmpty() {
				return nil, fmt.Errorf("sdInstance not found: %s", uid)
			}
			instance := instanceOptional.GetPayload()
			if instance.ID.IsEmpty() {
				return nil, fmt.Errorf("sdInstance loaded by UID has no internal ID: %s", uid)
			}
			allowedInstanceUIDs[instance.UID] = struct{}{}
			instanceTypeUIDs[instance.SDType.UID] = struct{}{}
		}
	}
	if len(filter.SdTypeUIDs) > 0 {
		allowedTypeUIDs = make(map[string]struct{})
		for _, uid := range filter.SdTypeUIDs {
			res := db.LoadSDTypeBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			sdType := res.GetPayload()
			if sdType.ID.IsEmpty() {
				return nil, fmt.Errorf("sdType loaded by UID has no internal ID: %s", uid)
			}
			allowedTypeUIDs[sdType.UID] = struct{}{}
		}
	}
	if len(instanceTypeUIDs) > 0 {
		if allowedTypeUIDs == nil {
			allowedTypeUIDs = instanceTypeUIDs
		} else {
			for t := range instanceTypeUIDs {
				if _, ok := allowedTypeUIDs[t]; !ok {
					return nil, fmt.Errorf("sdInstances do not belong to provided sdTypes")
				}
			}
		}
	}
	if allowedInstanceUIDs == nil && allowedTypeUIDs == nil {
		return func(_ graphQLModel.SDInstance) bool { return true }, nil
	}
	if allowedInstanceUIDs != nil && allowedTypeUIDs == nil {
		return func(inst graphQLModel.SDInstance) bool {
			_, ok := allowedInstanceUIDs[inst.UID]
			return ok
		}, nil
	}
	if allowedInstanceUIDs == nil && allowedTypeUIDs != nil {
		return func(inst graphQLModel.SDInstance) bool {
			_, ok := allowedTypeUIDs[inst.Type.UID]
			return ok
		}, nil
	}
	return func(inst graphQLModel.SDInstance) bool {
		if _, ok := allowedTypeUIDs[inst.Type.UID]; !ok {
			return false
		}
		_, ok := allowedInstanceUIDs[inst.UID]
		return ok
	}, nil
}

func BuildRawDataPointArrivedFilter(_ uint32, filter *graphQLModel.RawDataPointArrivedFilter) (func([]graphQLModel.RawDataPoint) bool, error) {

	db := dbClient.GetRelationalDatabaseClientInstance()
	if filter == nil {
		return func(_ []graphQLModel.RawDataPoint) bool { return true }, nil
	}
	var allowedInstanceUIDs map[string]struct{}
	var allowedTypeUIDs map[string]struct{}
	instanceTypeUIDs := make(map[string]struct{})
	if len(filter.SdInstanceUIDs) > 0 {
		if allowedInstanceUIDs == nil {
			allowedInstanceUIDs = make(map[string]struct{})
		}
		for _, uid := range filter.SdInstanceUIDs {
			res := db.LoadSDInstanceBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			instanceOptional := res.GetPayload()
			if instanceOptional.IsEmpty() {
				return nil, fmt.Errorf("sdInstance not found: %s", uid)
			}
			instance := instanceOptional.GetPayload()
			if instance.ID.IsEmpty() {
				return nil, fmt.Errorf("sdInstance loaded by UID has no internal ID: %s", uid)
			}
			allowedInstanceUIDs[instance.UID] = struct{}{}
			instanceTypeUIDs[instance.SDType.UID] = struct{}{}
		}
	}
	if len(filter.SdTypeUIDs) > 0 {
		if allowedTypeUIDs == nil {
			allowedTypeUIDs = make(map[string]struct{})
		}
		for _, uid := range filter.SdTypeUIDs {
			res := db.LoadSDTypeBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			sdType := res.GetPayload()
			if sdType.ID.IsEmpty() {
				return nil, fmt.Errorf("sdType loaded by UID has no internal ID: %s", uid)
			}
			allowedTypeUIDs[sdType.UID] = struct{}{}
		}
	}
	if len(instanceTypeUIDs) > 0 {
		if allowedTypeUIDs == nil {
			allowedTypeUIDs = instanceTypeUIDs
		} else {
			for t := range instanceTypeUIDs {
				if _, ok := allowedTypeUIDs[t]; !ok {
					return nil, fmt.Errorf("sdInstances do not belong to provided sdTypes")
				}
			}
		}
	}
	if allowedInstanceUIDs == nil && allowedTypeUIDs == nil {
		return func(_ []graphQLModel.RawDataPoint) bool { return true }, nil
	}
	if allowedInstanceUIDs != nil && allowedTypeUIDs == nil {
		return func(arr []graphQLModel.RawDataPoint) bool {
			for _, p := range arr {
				if _, ok := allowedInstanceUIDs[p.SdInstanceUID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	if allowedInstanceUIDs == nil && allowedTypeUIDs != nil {
		return func(arr []graphQLModel.RawDataPoint) bool {
			for _, p := range arr {
				if _, ok := allowedTypeUIDs[p.SdTypeUID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	return func(arr []graphQLModel.RawDataPoint) bool {
		for _, p := range arr {
			if _, ok := allowedTypeUIDs[p.SdTypeUID]; !ok {
				continue
			}
			if _, ok := allowedInstanceUIDs[p.SdInstanceUID]; ok {
				return true
			}
		}
		return false
	}, nil
}

func BuildKPIFulfillmentCheckedFilter(userID uint32, filter *graphQLModel.KPIFulfillmentCheckedFilter) (func([]graphQLModel.KPIFulfillmentCheckResult) bool, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	userKPIsResult := db.LoadKPIDefinitions(userID)
	if userKPIsResult.IsFailure() {
		return nil, userKPIsResult.GetError()
	}
	userKPISet := make(map[string]struct{})
	for _, k := range userKPIsResult.GetPayload() {
		if k.UID != nil {
			userKPISet[*k.UID] = struct{}{}
		}
	}
	var allowedInstanceUIDs map[string]struct{}
	var allowedTypeUIDs map[string]struct{}
	var allowedKPIUIDs map[string]struct{}
	if filter != nil && len(filter.SdInstanceUIDs) > 0 {
		if allowedInstanceUIDs == nil {
			allowedInstanceUIDs = make(map[string]struct{})
		}
		for _, uid := range filter.SdInstanceUIDs {
			res := db.LoadSDInstanceBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			instanceOptional := res.GetPayload()
			if instanceOptional.IsEmpty() {
				return nil, fmt.Errorf("sdInstance not found: %s", uid)
			}
			instance := instanceOptional.GetPayload()
			if instance.ID.IsEmpty() {
				return nil, fmt.Errorf("sdInstance loaded by UID has no internal ID: %s", uid)
			}
			allowedInstanceUIDs[instance.UID] = struct{}{}
		}
	}
	if filter != nil && len(filter.SdTypeUIDs) > 0 {
		if allowedTypeUIDs == nil {
			allowedTypeUIDs = make(map[string]struct{})
		}
		for _, uid := range filter.SdTypeUIDs {
			res := db.LoadSDTypeBasedOnUID(uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			sdType := res.GetPayload()
			if sdType.ID.IsEmpty() {
				return nil, fmt.Errorf("sdType loaded by UID has no internal ID: %s", uid)
			}
			allowedTypeUIDs[sdType.UID] = struct{}{}
		}
	}
	if filter != nil && len(filter.KpiDefinitionUIDs) > 0 {
		if allowedKPIUIDs == nil {
			allowedKPIUIDs = make(map[string]struct{})
		}
		for _, uid := range filter.KpiDefinitionUIDs {
			res := db.LoadKPIDefinitionByUID(userID, uid)
			if res.IsFailure() {
				return nil, res.GetError()
			}
			kpiDefinition := res.GetPayload()
			if kpiDefinition.ID == nil {
				return nil, fmt.Errorf("kpiDefinition loaded by UID has no internal ID: %s", uid)
			}
			if kpiDefinition.UID == nil {
				return nil, fmt.Errorf("kpiDefinition loaded by UID has no UID: %s", uid)
			}
			allowedKPIUIDs[*kpiDefinition.UID] = struct{}{}
		}
	}
	if filter == nil || (allowedInstanceUIDs == nil && allowedTypeUIDs == nil && allowedKPIUIDs == nil) {
		return func(arr []graphQLModel.KPIFulfillmentCheckResult) bool {
			for _, p := range arr {
				if _, ok := userKPISet[p.KpiDefinitionUID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	return func(arr []graphQLModel.KPIFulfillmentCheckResult) bool {
		for _, p := range arr {
			if _, ok := userKPISet[p.KpiDefinitionUID]; !ok {
				continue
			}
			if allowedInstanceUIDs != nil {
				if _, ok := allowedInstanceUIDs[p.SdInstanceUID]; !ok {
					continue
				}
			}
			if allowedTypeUIDs != nil {
				if _, ok := allowedTypeUIDs[p.SdTypeUID]; !ok {
					continue
				}
			}
			if allowedKPIUIDs != nil {
				if _, ok := allowedKPIUIDs[p.KpiDefinitionUID]; !ok {
					continue
				}
			}
			return true
		}
		return false
	}, nil
}

func BuildTimeSeriesExportFilter(filter *graphQLModel.TimeSeriesExportFilter) (func(graphQLModel.TimeSeriesExport) bool, error) {
	if filter == nil || len(filter.Uids) == 0 {
		return nil, fmt.Errorf("timeSeries export subscription requires at least one uid")
	}
	allowedUIDs := make(map[string]struct{}, len(filter.Uids))
	for _, uid := range filter.Uids {
		normalizedUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(uid, "exp", "time series export")
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		allowedUIDs[normalizedUID] = struct{}{}
	}
	return func(export graphQLModel.TimeSeriesExport) bool {
		_, ok := allowedUIDs[export.UID]
		return ok
	}, nil
}
