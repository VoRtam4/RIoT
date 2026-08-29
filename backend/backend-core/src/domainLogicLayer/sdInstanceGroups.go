/**
 * @file sdInstanceGroups.go
 * @brief Doménová logika pro správu skupin sledovaných instancí.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func GetSDInstanceGroup(id uint32) sharedUtils.Result[graphQLModel.SDInstanceGroup] {
	sdInstanceGroupLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceGroup(id)
	if sdInstanceGroupLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](sdInstanceGroupLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelSDInstanceGroup(sdInstanceGroupLoadResult.GetPayload()))
}

func GetSDInstanceGroupByUID(uid string) sharedUtils.Result[graphQLModel.SDInstanceGroup] {
	normalizedUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(uid, "grp", "SD instance group")
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](normalizeErr)
	}
	sdInstanceGroupLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceGroupByUID(normalizedUID)
	if sdInstanceGroupLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](sdInstanceGroupLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelSDInstanceGroup(sdInstanceGroupLoadResult.GetPayload()))
}

func GetSDInstanceGroups() sharedUtils.Result[[]graphQLModel.SDInstanceGroup] {
	sdInstanceGroupsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceGroups()
	if sdInstanceGroupsLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.SDInstanceGroup](sdInstanceGroupsLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(sdInstanceGroupsLoadResult.GetPayload(), dll2gql.ToGraphQLModelSDInstanceGroup))
}

func CreateSDInstanceGroup(input graphQLModel.SDInstanceGroupInput) sharedUtils.Result[graphQLModel.SDInstanceGroup] {
	sdInstanceGroup := gql2dll.ToDLLModelSDInstanceGroup(input)
	if err := prepareSDInstanceGroupForPersistence(&sdInstanceGroup); err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](err)
	}
	sdInstanceGroupPersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistSDInstanceGroup(sdInstanceGroup)
	if sdInstanceGroupPersistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](sdInstanceGroupPersistResult.GetError())
	}
	sdInstanceGroup.ID = sharedUtils.NewOptionalOf(sdInstanceGroupPersistResult.GetPayload())
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelSDInstanceGroup(sdInstanceGroup))
}

func UpdateSDInstanceGroup(uid string, input graphQLModel.SDInstanceGroupInput) sharedUtils.Result[graphQLModel.SDInstanceGroup] {
	normalizedUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(uid, "grp", "SD instance group")
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](normalizeErr)
	}
	inputUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(input.UID, "grp", "SD instance group")
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](normalizeErr)
	}
	if inputUID != normalizedUID {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](fmt.Errorf("sdInstanceGroup input UID must match path UID"))
	}
	existingGroupResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceGroupByUID(normalizedUID)
	if existingGroupResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](existingGroupResult.GetError())
	}
	sdInstanceGroup := gql2dll.ToDLLModelSDInstanceGroup(input)
	sdInstanceGroup.ID = existingGroupResult.GetPayload().ID
	if err := prepareSDInstanceGroupForPersistence(&sdInstanceGroup); err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](err)
	}
	if sdInstanceGroupPersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistSDInstanceGroup(sdInstanceGroup); sdInstanceGroupPersistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.SDInstanceGroup](sdInstanceGroupPersistResult.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelSDInstanceGroup(sdInstanceGroup))
}

func DeleteSDInstanceGroup(uid string) error {
	normalizedUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(uid, "grp", "SD instance group")
	if normalizeErr != nil {
		return normalizeErr
	}
	return dbClient.GetRelationalDatabaseClientInstance().DeleteSDInstanceGroupByUID(normalizedUID)
}

func prepareSDInstanceGroupForPersistence(sdInstanceGroup *dllModel.SDInstanceGroup) error {
	normalizedGroupUID, groupUIDSuffix, normalizeErr := sharedUtils.NormalizePrefixedUID(sdInstanceGroup.UID, "grp", "SD instance group")
	if normalizeErr != nil {
		return normalizeErr
	}
	sdInstanceGroup.UID = normalizedGroupUID
	sdInstanceGroup.UserIdentifier = sharedUtils.SafeLabel(sdInstanceGroup.UserIdentifier, groupUIDSuffix)
	sdInstanceGroup.SDInstanceIDs = make([]uint32, 0, len(sdInstanceGroup.SDInstanceUIDs))
	for _, uid := range sdInstanceGroup.SDInstanceUIDs {
		normalizedUID := strings.TrimSpace(uid)
		if normalizedUID == "" {
			return fmt.Errorf("invalid sdInstanceUID")
		}
		result := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstanceBasedOnUID(normalizedUID)
		if result.IsFailure() {
			return result.GetError()
		}
		optional := result.GetPayload()
		if optional.IsEmpty() {
			return fmt.Errorf("sdInstance not found: %s", normalizedUID)
		}
		instance := optional.GetPayload()
		if instance.ID.IsEmpty() {
			return fmt.Errorf("sdInstance loaded by UID has no internal ID: %s", normalizedUID)
		}
		sdInstanceGroup.SDInstanceIDs = append(sdInstanceGroup.SDInstanceIDs, instance.ID.GetPayload())
	}
	return nil
}
