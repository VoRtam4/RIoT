package events

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func BuildSDInstanceRegisteredFilter(_ uint32, filter *graphQLModel.SDInstanceRegisteredFilter) (func(graphQLModel.SDInstance) bool, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	if filter == nil {
		return func(_ graphQLModel.SDInstance) bool { return true }, nil
	}
	var allowedInstanceIDs map[uint32]struct{}
	var allowedTypeIDs map[uint32]struct{}
	if len(filter.SdInstanceIDs) > 0 {
		allowedInstanceIDs = make(map[uint32]struct{})
		instanceTypeIDs := make(map[uint32]struct{})
		for _, id := range filter.SdInstanceIDs {
			res := db.LoadSDInstance(id)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdInstance not found: %d", id)
			}
			instance := res.GetPayload()
			allowedInstanceIDs[id] = struct{}{}
			instanceTypeIDs[instance.SDType.ID.GetPayload()] = struct{}{}
		}
		if len(filter.SdTypeIDs) > 0 {
			allowedTypeIDs = make(map[uint32]struct{})
			for _, typeID := range filter.SdTypeIDs {
				res := db.LoadSDType(typeID)
				if res.IsFailure() {
					return nil, fmt.Errorf("sdType not found: %d", typeID)
				}
				allowedTypeIDs[typeID] = struct{}{}
			}
			for t := range instanceTypeIDs {
				if _, ok := allowedTypeIDs[t]; !ok {
					return nil, fmt.Errorf("sdInstances do not belong to provided sdTypeIDs")
				}
			}
		} else {
			allowedTypeIDs = instanceTypeIDs
		}
	}
	if len(filter.SdTypeIDs) > 0 && allowedTypeIDs == nil {
		allowedTypeIDs = make(map[uint32]struct{})
		for _, typeID := range filter.SdTypeIDs {
			res := db.LoadSDType(typeID)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdType not found: %d", typeID)
			}
			allowedTypeIDs[typeID] = struct{}{}
		}
	}
	if allowedInstanceIDs == nil && allowedTypeIDs == nil {
		return func(_ graphQLModel.SDInstance) bool { return true }, nil
	}
	if allowedInstanceIDs != nil && allowedTypeIDs == nil {
		return func(inst graphQLModel.SDInstance) bool {
			_, ok := allowedInstanceIDs[inst.Type.ID]
			return ok
		}, nil
	}
	if allowedInstanceIDs == nil && allowedTypeIDs != nil {
		return func(inst graphQLModel.SDInstance) bool {
			_, ok := allowedTypeIDs[inst.Type.ID]
			return ok
		}, nil
	}
	return func(inst graphQLModel.SDInstance) bool {
		if _, ok := allowedTypeIDs[inst.Type.ID]; !ok {
			return false
		}
		_, ok := allowedInstanceIDs[inst.ID]
		return ok
	}, nil
}

func BuildRawDataPointArrivedFilter(_ uint32, filter *graphQLModel.RawDataPointArrivedFilter) (func([]graphQLModel.RawDataPoint) bool, error) {

	db := dbClient.GetRelationalDatabaseClientInstance()
	if filter == nil {
		return func(_ []graphQLModel.RawDataPoint) bool { return true }, nil
	}
	var allowedInstanceIDs map[uint32]struct{}
	var allowedTypeIDs map[uint32]struct{}
	if len(filter.SdInstanceIDs) > 0 {
		allowedInstanceIDs = make(map[uint32]struct{})
		instanceTypeIDs := make(map[uint32]struct{})

		for _, id := range filter.SdInstanceIDs {
			res := db.LoadSDInstance(id)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdInstance not found: %d", id)
			}
			inst := res.GetPayload()
			allowedInstanceIDs[id] = struct{}{}
			instanceTypeIDs[inst.SDType.ID.GetPayload()] = struct{}{}
		}

		if len(filter.SdTypeIDs) > 0 {
			allowedTypeIDs = make(map[uint32]struct{})
			for _, t := range filter.SdTypeIDs {
				res := db.LoadSDType(t)
				if res.IsFailure() {
					return nil, fmt.Errorf("sdType not found: %d", t)
				}
				allowedTypeIDs[t] = struct{}{}
			}

			for t := range instanceTypeIDs {
				if _, ok := allowedTypeIDs[t]; !ok {
					return nil, fmt.Errorf("sdInstances do not belong to provided sdTypeIDs")
				}
			}
		} else {
			allowedTypeIDs = instanceTypeIDs
		}
	}
	if len(filter.SdTypeIDs) > 0 && allowedTypeIDs == nil {
		allowedTypeIDs = make(map[uint32]struct{})
		for _, t := range filter.SdTypeIDs {
			res := db.LoadSDType(t)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdType not found: %d", t)
			}
			allowedTypeIDs[t] = struct{}{}
		}
	}
	if allowedInstanceIDs == nil && allowedTypeIDs == nil {
		return func(_ []graphQLModel.RawDataPoint) bool { return true }, nil
	}
	if allowedInstanceIDs != nil && allowedTypeIDs == nil {
		return func(arr []graphQLModel.RawDataPoint) bool {
			for _, p := range arr {
				if _, ok := allowedInstanceIDs[p.SdInstanceID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	if allowedInstanceIDs == nil && allowedTypeIDs != nil {
		return func(arr []graphQLModel.RawDataPoint) bool {
			for _, p := range arr {
				if _, ok := allowedTypeIDs[p.SdTypeID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	return func(arr []graphQLModel.RawDataPoint) bool {
		for _, p := range arr {
			if _, ok := allowedTypeIDs[p.SdTypeID]; !ok {
				continue
			}
			if _, ok := allowedInstanceIDs[p.SdInstanceID]; ok {
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
	userKPISet := make(map[uint32]struct{})
	for _, k := range userKPIsResult.GetPayload() {
		if k.ID != nil {
			userKPISet[*k.ID] = struct{}{}
		}
	}
	var allowedInstanceIDs map[uint32]struct{}
	var allowedTypeIDs map[uint32]struct{}
	var allowedKPIIDs map[uint32]struct{}
	if filter != nil && len(filter.SdInstanceIDs) > 0 {
		allowedInstanceIDs = make(map[uint32]struct{})
		for _, id := range filter.SdInstanceIDs {
			res := db.LoadSDInstance(id)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdInstance not found: %d", id)
			}
			allowedInstanceIDs[id] = struct{}{}
		}
	}
	if filter != nil && len(filter.SdTypeIDs) > 0 {
		allowedTypeIDs = make(map[uint32]struct{})
		for _, t := range filter.SdTypeIDs {
			res := db.LoadSDType(t)
			if res.IsFailure() {
				return nil, fmt.Errorf("sdType not found: %d", t)
			}
			allowedTypeIDs[t] = struct{}{}
		}
	}
	if filter != nil && len(filter.KpiDefinitions) > 0 {
		allowedKPIIDs = make(map[uint32]struct{})

		for _, id := range filter.KpiDefinitions {
			if _, ok := userKPISet[id]; !ok {
				return nil, fmt.Errorf("kpiDefinition not accessible: %d", id)
			}
			allowedKPIIDs[id] = struct{}{}
		}
	}
	if filter == nil || (allowedInstanceIDs == nil && allowedTypeIDs == nil && allowedKPIIDs == nil) {
		return func(arr []graphQLModel.KPIFulfillmentCheckResult) bool {
			for _, p := range arr {
				if _, ok := userKPISet[p.KpiDefinitionID]; ok {
					return true
				}
			}
			return false
		}, nil
	}
	return func(arr []graphQLModel.KPIFulfillmentCheckResult) bool {
		for _, p := range arr {
			if _, ok := userKPISet[p.KpiDefinitionID]; !ok {
				continue
			}
			if allowedInstanceIDs != nil {
				if _, ok := allowedInstanceIDs[p.SdInstanceID]; !ok {
					continue
				}
			}
			if allowedTypeIDs != nil {
				if _, ok := allowedTypeIDs[p.SdTypeID]; !ok {
					continue
				}
			}
			if allowedKPIIDs != nil {
				if _, ok := allowedKPIIDs[p.KpiDefinitionID]; !ok {
					continue
				}
			}
			return true
		}
		return false
	}, nil
}
