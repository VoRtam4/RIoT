package dbClient

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbUtil"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/db2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2db"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var (
	rdbClientInstance RelationalDatabaseClient
	once              sync.Once
)

type RelationalDatabaseClient interface {
	setup()
	PerformOnStartupOperations(permissions map[string]map[string]bool, formatPermissionLabel func(string) string) error
	PersistKPIDefinition(userID uint32, kpiDefinition sharedModel.KPIDefinition) sharedUtils.Result[uint32]
	LoadAllKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition]
	LoadKPIDefinition(userID uint32, id uint32) sharedUtils.Result[sharedModel.KPIDefinition]
	LoadKPIDefinitions(userID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition]
	LoadKPIDefinitionsBySDType(userID uint32, sdTypeID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition]
	LoadKPIDefinitionsBySDInstance(userID uint32, sdInstanceID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition]
	DeleteKPIDefinition(id uint32) error
	PersistSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType]
	UpsertSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType]
	LoadSDType(id uint32) sharedUtils.Result[dllModel.SDType]
	LoadSDInstancesByType(sdTypeID uint32) sharedUtils.Result[[]dllModel.SDInstance]
	LoadSDInstancesByKpiDefinition(kpiDefinitionID uint32) sharedUtils.Result[[]dllModel.SDInstance]
	LoadSDTypeBasedOnUID(uid string) sharedUtils.Result[dllModel.SDType]
	LoadSDTypes() sharedUtils.Result[[]dllModel.SDType]
	DeleteSDType(id uint32) error
	PersistSDInstance(sdInstance dllModel.SDInstance) sharedUtils.Result[uint32]
	UpsertSDInstance(uid string, sdTypeSpecification string, label string) sharedUtils.Result[dllModel.SDInstance]
	LoadSDInstance(id uint32) sharedUtils.Result[dllModel.SDInstance]
	LoadSDInstanceBasedOnUID(uid string) sharedUtils.Result[sharedUtils.Optional[dllModel.SDInstance]]
	LoadSDInstances() sharedUtils.Result[[]dllModel.SDInstance]
	PersistRawDataPoints(points []dllModel.RawDataPoint) sharedUtils.Result[[]dllModel.RawDataPoint]
	LoadRawDataPointsBySDType(sdTypeID uint32) sharedUtils.Result[[]dllModel.RawDataPoint]
	LoadRawDataPoint(sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.RawDataPoint]]
	PersistKPIFulfillmentCheckResults(points []dllModel.KPIFulfillmentCheckResult, reprocess bool) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
	LoadKPIFulfillmentCheckResults(userID uint32) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
	LoadKPIFulfillmentCheckResultsByKPI(userID uint32, kpiDefinitionID uint32) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
	LoadKPIFulfillmentCheckResult(userID uint32, input dllModel.KPIFulfillmentCheckResultRequest) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]]
	LoadSDInstanceGroups() sharedUtils.Result[[]dllModel.SDInstanceGroup]
	LoadSDInstanceGroup(id uint32) sharedUtils.Result[dllModel.SDInstanceGroup]
	PersistSDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) sharedUtils.Result[uint32]
	DeleteSDInstanceGroup(id uint32) error
	PersistUser(user dllModel.User) sharedUtils.Result[uint]
	LoadUserBasedOnOAuth2ProviderIssuedID(oauth2ProviderIssuedID string) sharedUtils.Result[sharedUtils.Optional[dllModel.User]]
	LoadUser(id uint) sharedUtils.Result[dllModel.User]
	LoadUserSessionBasedOnRefreshTokenHash(refreshTokenHash string) sharedUtils.Result[sharedUtils.Optional[dllModel.UserSession]]
	PersistUserSession(userSession dllModel.UserSession) sharedUtils.Result[uint]
	PersistUserConfig(userConfig dllModel.UserConfig) sharedUtils.Result[uint32]
	LoadUserConfig(userId uint32) sharedUtils.Result[dllModel.UserConfig]
	DeleteUserConfig(userId uint32) error
	LoadAPIKeyByHash(hash string) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]]
	LoadAPIKeyByID(id uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]]
	LoadAPIKeysForUser(userID uint32) sharedUtils.Result[[]dllModel.APIKey]
	CreateAPIKey(userID uint32, k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey]
	UpdateAPIKey(k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey]
	DeleteAPIKey(id uint32) error
	LoadRoles() sharedUtils.Result[[]dllModel.Role]
	LoadRoleByUID(uid string) sharedUtils.Result[dllModel.Role]
	GetUserRole(userID uint32) sharedUtils.Result[dllModel.Role]
	SetUserRole(userID uint32, roleID uint32) error
}

var ErrOperationWouldLeadToForeignKeyIntegrityBreach = errors.New("operation would lead to foreign key integrity breach")

type relationalDatabaseClientImpl struct {
	db *gorm.DB
	mu sync.Mutex
}

func GetRelationalDatabaseClientInstance() RelationalDatabaseClient {
	once.Do(func() {
		rdbClientInstance = new(relationalDatabaseClientImpl)
		rdbClientInstance.setup()
	})
	return rdbClientInstance
}

func (r *relationalDatabaseClientImpl) setup() {
	rawPostgresURL := sharedUtils.GetEnvironmentVariableValue("POSTGRES_URL").GetPayloadOrDefault("postgres://admin:password@postgres:5432/postgres-db")
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(postgres.Open(rawPostgresURL), new(gorm.Config))
		if err != nil && strings.Contains(err.Error(), "SQLSTATE 57P03") {
			fmt.Println("[RDB client (GORM)]: Database seems to be starting up, retrying connection in 5 seconds...")
			time.Sleep(time.Second * 5)
			continue
		}
		sharedUtils.TerminateOnError(err, "[RDB client (GORM)]: couldn't connect to the database")
		break
	}
	session := new(gorm.Session)
	session.Logger = logger.Default.LogMode(logger.Warn)
	r.db = db.Session(session)
	sharedUtils.TerminateOnError(r.db.AutoMigrate(
		new(dbModel.KPIDefinitionEntity),
		new(dbModel.KPINodeEntity),
		new(dbModel.LogicalOperationKPINodeEntity),
		new(dbModel.AtomKPINodeEntity),
		new(dbModel.SDTypeEntity),
		new(dbModel.SDParameterEntity),
		new(dbModel.SDInstanceEntity),
		new(dbModel.RawDataPointEntity),
		new(dbModel.KPIFulfillmentCheckResultEntity),
		new(dbModel.SDInstanceGroupEntity),
		new(dbModel.SDInstanceGroupMembershipEntity),
		new(dbModel.SDInstanceKPIDefinitionRelationshipEntity),
		new(dbModel.UserEntity),
		new(dbModel.UserSessionEntity),
		new(dbModel.UserConfigEntity),
		new(dbModel.GraphQLOperationEntity),
		new(dbModel.RoleEntity),
		new(dbModel.RolesPermissionsMappingEntity),
		new(dbModel.APIKeysPermissionsMappingEntity),
		new(dbModel.PermissionEntity),
		new(dbModel.OperationTypeAccessPermissionEntity),
		new(dbModel.SingleOperationPermissionEntity),
		new(dbModel.APIKeyEntity),
		new(dbModel.APIKeyIPRestrictionEntity),
	), "[RDB client (GORM)]: auto-migration failed")
}

func (r *relationalDatabaseClientImpl) PerformOnStartupOperations(permissions map[string]map[string]bool, formatPermissionLabel func(string) string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		for roleLabel, perms := range permissions {

			var role dbModel.RoleEntity

			err := tx.Where("label = ?", roleLabel).First(&role).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				role = dbModel.RoleEntity{Label: roleLabel}
				if err := tx.Create(&role).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			for permLabel, allowed := range perms {
				if !allowed {
					continue
				}

				var perm dbModel.PermissionEntity
				err := tx.Where("uid = ?", permLabel).First(&perm).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					perm = dbModel.PermissionEntity{UID: permLabel, Label: formatPermissionLabel(permLabel)}
					if err := tx.Create(&perm).Error; err != nil {
						return err
					}
				} else if err != nil {
					return err
				}

				if err := tx.Model(&role).Association("Permissions").Append(&perm); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *relationalDatabaseClientImpl) PersistKPIDefinition(userID uint32, kpiDefinition sharedModel.KPIDefinition) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiDefinitionID := sharedUtils.NewOptionalFromPointer(kpiDefinition.ID).GetPayloadOrDefault(0)
	idsOfKPINodeEntitiesFormingTheKPIDefinition := sharedUtils.EmptySlice[uint32]()
	if kpiDefinitionID != 0 {
		getIDsResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, kpiDefinitionID)
		if getIDsResult.IsFailure() {
			return sharedUtils.NewFailureResult[uint32](getIDsResult.GetError())
		}
		idsOfKPINodeEntitiesFormingTheKPIDefinition = getIDsResult.GetPayload()
	}
	kpiNodeEntity, kpiNodeEntities, logicalNodes, atomNodes := dll2db.ToDBModelEntitiesKPIDefinition(kpiDefinition)
	referencedSDInstances := make([]dbModel.SDInstanceEntity, 0)
	if len(kpiDefinition.SelectedSDInstanceIDs) > 0 {
		referencedSDInstancesResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Where("id IN (?)", kpiDefinition.SelectedSDInstanceIDs))
		if referencedSDInstancesResult.IsFailure() {
			return sharedUtils.NewFailureResult[uint32](referencedSDInstancesResult.GetError())
		}
		referencedSDInstances = referencedSDInstancesResult.GetPayload()
	}
	kpiDefinitionEntity := dbModel.KPIDefinitionEntity{
		ID:             kpiDefinitionID,
		Label:          kpiDefinition.Label,
		UserID:         userID,
		SDTypeID:       kpiDefinition.SDTypeID,
		UserIdentifier: kpiDefinition.UserIdentifier,
		RootNode:       kpiNodeEntity,
		SDInstanceMode: string(kpiDefinition.SDInstanceMode),
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range kpiNodeEntities {
			if err := dbUtil.PersistEntityIntoDB(tx, entity); err != nil {
				return err
			}
		}
		if err := dbUtil.PersistEntityIntoDB(tx, &kpiDefinitionEntity); err != nil {
			return err
		}
		for _, entity := range logicalNodes {
			if err := dbUtil.PersistEntityIntoDB(tx, &entity); err != nil {
				return err
			}
		}
		for _, entity := range atomNodes {
			if err := dbUtil.PersistEntityIntoDB(tx, &entity); err != nil {
				return err
			}
		}
		if kpiDefinitionID != 0 {
			if err := dbUtil.DeleteEntitiesBasedOnWhereClauses[dbModel.SDInstanceKPIDefinitionRelationshipEntity](
				tx,
				dbUtil.Where("kpi_definition_id = ?", kpiDefinitionID),
			); err != nil {
				return err
			}
		}
		for _, sdInstance := range referencedSDInstances {
			entity := dbModel.SDInstanceKPIDefinitionRelationshipEntity{
				KPIDefinitionID: kpiDefinitionEntity.ID,
				SDInstanceID:    sdInstance.ID,
				SDInstanceUID:   sdInstance.UID,
			}
			if err := dbUtil.PersistEntityIntoDB(tx, &entity); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	if len(idsOfKPINodeEntitiesFormingTheKPIDefinition) > 0 {
		if err := dbUtil.DeleteEntitiesBasedOnSliceOfIds[dbModel.KPINodeEntity](
			r.db,
			idsOfKPINodeEntitiesFormingTheKPIDefinition,
		); err != nil {
			log.Printf("cleanup failed: %s\n", err.Error())
		}
	}
	return sharedUtils.NewSuccessResult[uint32](kpiDefinitionEntity.ID)
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinition(userID uint32, id uint32) sharedUtils.Result[sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Where("id = ? AND user_id = ?", id, userID), dbUtil.Preload("SDType"), dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](result.GetError())
	}
	entities := result.GetPayload()
	if len(entities) == 0 {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](errors.New("not found"))
	}
	entity := entities[0]
	idsResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, entity.ID)
	if idsResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](idsResult.GetError())
	}
	nodeIDs := idsResult.GetPayload()
	kpiNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](
		r.db,
		dbUtil.Where("id IN (?)", nodeIDs),
	)
	if kpiNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](kpiNodesResult.GetError())
	}
	kpiNodes := kpiNodesResult.GetPayload()

	logicalNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](
		r.db,
		dbUtil.Where("node_id IN (?)", nodeIDs),
	)
	if logicalNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](logicalNodesResult.GetError())
	}
	logicalNodes := logicalNodesResult.GetPayload()

	atomNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](
		r.db,
		dbUtil.Where("node_id IN (?)", nodeIDs),
		dbUtil.Preload("SDParameter"),
	)
	if atomNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](atomNodesResult.GetError())
	}
	atomNodes := atomNodesResult.GetPayload()
	kpi := db2dll.ToDLLModelKPIDefinition(entity, kpiNodes, logicalNodes, atomNodes)
	return sharedUtils.NewSuccessResult(kpi)
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinitions(userID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](
		r.db,
		dbUtil.Where("user_id = ?", userID),
		dbUtil.Preload("SDType"),
		dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"),
	)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](result.GetError())
	}

	entities := result.GetPayload()
	if len(entities) == 0 {
		return sharedUtils.NewSuccessResult([]sharedModel.KPIDefinition{})
	}

	allNodeIDs := make([]uint32, 0)
	for _, entity := range entities {
		idsResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, entity.ID)
		if idsResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](idsResult.GetError())
		}
		allNodeIDs = append(allNodeIDs, idsResult.GetPayload()...)
	}

	kpiNodes := make([]dbModel.KPINodeEntity, 0)
	logicalNodes := make([]dbModel.LogicalOperationKPINodeEntity, 0)
	atomNodes := make([]dbModel.AtomKPINodeEntity, 0)

	if len(allNodeIDs) > 0 {
		kpiNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](
			r.db,
			dbUtil.Where("id IN (?)", allNodeIDs),
		)
		if kpiNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](kpiNodesResult.GetError())
		}
		kpiNodes = kpiNodesResult.GetPayload()

		logicalNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](
			r.db,
			dbUtil.Where("node_id IN (?)", allNodeIDs),
		)
		if logicalNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](logicalNodesResult.GetError())
		}
		logicalNodes = logicalNodesResult.GetPayload()

		atomNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](
			r.db,
			dbUtil.Where("node_id IN (?)", allNodeIDs),
			dbUtil.Preload("SDParameter"),
		)
		if atomNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](atomNodesResult.GetError())
		}
		atomNodes = atomNodesResult.GetPayload()
	}

	kpis := make([]sharedModel.KPIDefinition, 0, len(entities))
	for _, entity := range entities {
		kpi := db2dll.ToDLLModelKPIDefinition(entity, kpiNodes, logicalNodes, atomNodes)
		kpis = append(kpis, kpi)
	}

	return sharedUtils.NewSuccessResult(kpis)
}

func (r *relationalDatabaseClientImpl) LoadAllKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()

	kpiDefinitionEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](
		r.db,
		dbUtil.Preload("SDType"),
		dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"),
	)
	if kpiDefinitionEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI definition entities from the database: %w", kpiDefinitionEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	kpiDefinitionEntities := kpiDefinitionEntitiesLoadResult.GetPayload()

	if len(kpiDefinitionEntities) == 0 {
		return sharedUtils.NewSuccessResult([]sharedModel.KPIDefinition{})
	}

	kpiNodeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](r.db)
	if kpiNodeEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI node entities: %w", kpiNodeEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	kpiNodeEntities := kpiNodeEntitiesLoadResult.GetPayload()

	logicalNodesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](r.db)
	if logicalNodesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load logical KPI nodes: %w", logicalNodesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	logicalNodes := logicalNodesLoadResult.GetPayload()

	atomNodesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](
		r.db,
		dbUtil.Preload("SDParameter"),
	)
	if atomNodesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load atom KPI nodes: %w", atomNodesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	atomNodes := atomNodesLoadResult.GetPayload()

	kpiDefinitions := make([]sharedModel.KPIDefinition, 0, len(kpiDefinitionEntities))
	for _, entity := range kpiDefinitionEntities {
		kpiDefinition := db2dll.ToDLLModelKPIDefinition(entity, kpiNodeEntities, logicalNodes, atomNodes)
		kpiDefinitions = append(kpiDefinitions, kpiDefinition)
	}

	return sharedUtils.NewSuccessResult(kpiDefinitions)
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinitionsBySDType(userID uint32, sdTypeID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Where("user_id = ? AND sd_type_id = ?", userID, sdTypeID), dbUtil.Preload("SDType"), dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](result.GetError())
	}
	entities := result.GetPayload()
	allNodeIDs := make([]uint32, 0)
	for _, entity := range entities {
		idsResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, entity.ID)
		if idsResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](idsResult.GetError())
		}
		allNodeIDs = append(allNodeIDs, idsResult.GetPayload()...)
	}
	if len(allNodeIDs) == 0 {
		return sharedUtils.NewSuccessResult([]sharedModel.KPIDefinition{})
	}
	kpiNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](r.db, dbUtil.Where("id IN (?)", allNodeIDs))
	if kpiNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](kpiNodesResult.GetError())
	}
	kpiNodes := kpiNodesResult.GetPayload()
	logicalNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](r.db, dbUtil.Where("node_id IN (?)", allNodeIDs))
	if logicalNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](logicalNodesResult.GetError())
	}
	logicalNodes := logicalNodesResult.GetPayload()
	atomNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](r.db, dbUtil.Where("node_id IN (?)", allNodeIDs), dbUtil.Preload("SDParameter"))
	if atomNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](atomNodesResult.GetError())
	}
	atomNodes := atomNodesResult.GetPayload()
	kpis := make([]sharedModel.KPIDefinition, 0, len(entities))
	for _, entity := range entities {
		kpi := db2dll.ToDLLModelKPIDefinition(entity, kpiNodes, logicalNodes, atomNodes)
		kpis = append(kpis, kpi)
	}
	return sharedUtils.NewSuccessResult(kpis)
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinitionsBySDInstance(userID uint32, sdInstanceID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Where("id = ?", sdInstanceID))
	if sdInstanceResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](sdInstanceResult.GetError())
	}
	sdInstance := sdInstanceResult.GetPayload()
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Where("user_id = ?", userID), dbUtil.Preload("SDType"), dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](result.GetError())
	}
	entities := result.GetPayload()
	filtered := make([]dbModel.KPIDefinitionEntity, 0)
	for _, e := range entities {
		if e.SDTypeID != sdInstance.SDTypeID {
			continue
		}
		if e.SDInstanceMode == string(sharedModel.ALL) {
			filtered = append(filtered, e)
			continue
		}
		for _, rel := range e.SDInstanceKPIDefinitionRelationshipRecords {
			if rel.SDInstanceID == sdInstanceID {
				filtered = append(filtered, e)
				break
			}
		}
	}
	if len(filtered) == 0 {
		return sharedUtils.NewSuccessResult([]sharedModel.KPIDefinition{})
	}
	allNodeIDs := make([]uint32, 0)
	for _, entity := range filtered {
		idsResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, entity.ID)
		if idsResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](idsResult.GetError())
		}
		allNodeIDs = append(allNodeIDs, idsResult.GetPayload()...)
	}
	kpiNodes := make([]dbModel.KPINodeEntity, 0)
	logicalNodes := make([]dbModel.LogicalOperationKPINodeEntity, 0)
	atomNodes := make([]dbModel.AtomKPINodeEntity, 0)
	if len(allNodeIDs) > 0 {
		kpiNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](r.db, dbUtil.Where("id IN (?)", allNodeIDs))
		if kpiNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](kpiNodesResult.GetError())
		}
		kpiNodes = kpiNodesResult.GetPayload()
		logicalNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](r.db, dbUtil.Where("node_id IN (?)", allNodeIDs))
		if logicalNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](logicalNodesResult.GetError())
		}
		logicalNodes = logicalNodesResult.GetPayload()
		atomNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](r.db, dbUtil.Where("node_id IN (?)", allNodeIDs), dbUtil.Preload("SDParameter"))
		if atomNodesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](atomNodesResult.GetError())
		}
		atomNodes = atomNodesResult.GetPayload()
	}
	kpis := make([]sharedModel.KPIDefinition, 0, len(filtered))
	for _, entity := range filtered {
		kpi := db2dll.ToDLLModelKPIDefinition(entity, kpiNodes, logicalNodes, atomNodes)
		kpis = append(kpis, kpi)
	}
	return sharedUtils.NewSuccessResult(kpis)
}

func (r *relationalDatabaseClientImpl) DeleteKPIDefinition(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	idsOfKPINodeEntitiesFormingTheDefinitionResult := dbModel.GetIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, id)
	if idsOfKPINodeEntitiesFormingTheDefinitionResult.IsFailure() {
		return idsOfKPINodeEntitiesFormingTheDefinitionResult.GetError()
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := dbUtil.DeleteCertainEntityBasedOnId[dbModel.KPIDefinitionEntity](tx, id); err != nil {
			return err
		}
		if err := dbUtil.DeleteEntitiesBasedOnSliceOfIds[dbModel.KPINodeEntity](tx, idsOfKPINodeEntitiesFormingTheDefinitionResult.GetPayload()); err != nil {
			return err
		}
		return nil
	})
}

func (r *relationalDatabaseClientImpl) PersistSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdTypeEntity := dll2db.ToDBModelEntitySDType(sdType)
	if err := dbUtil.PersistEntityIntoDB[dbModel.SDTypeEntity](r.db, &sdTypeEntity); err != nil {
		return sharedUtils.NewFailureResult[dllModel.SDType](err)
	}
	return sharedUtils.NewSuccessResult[dllModel.SDType](db2dll.ToDLLModelSDType(sdTypeEntity))
}

func (r *relationalDatabaseClientImpl) UpsertSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Preload("Parameters"), dbUtil.Where("uid = ?", sdType.UID))
	var entity dbModel.SDTypeEntity
	if existingResult.IsFailure() {
		if errors.Is(existingResult.GetError(), gorm.ErrRecordNotFound) {
			entity = dll2db.ToDBModelEntitySDType(sdType)
			if err := dbUtil.PersistEntityIntoDB(r.db, &entity); err != nil {
				return sharedUtils.NewFailureResult[dllModel.SDType](err)
			}
			return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDType(entity))
		}
		return sharedUtils.NewFailureResult[dllModel.SDType](existingResult.GetError())
	}
	entity = existingResult.GetPayload()
	if sdType.Label != "" && entity.Label != sdType.Label {
		entity.Label = sdType.Label
		if err := r.db.Save(&entity).Error; err != nil {
			return sharedUtils.NewFailureResult[dllModel.SDType](err)
		}
	}
	existingMap := make(map[string]dbModel.SDParameterEntity)
	for _, p := range entity.Parameters {
		existingMap[p.Denotation] = p
	}
	incomingMap := make(map[string]dllModel.SDParameter)
	for _, p := range sdType.Parameters {
		incomingMap[p.Denotation] = p
	}
	for _, newParam := range sdType.Parameters {
		if existingParam, ok := existingMap[newParam.Denotation]; ok {
			if existingParam.Label != newParam.Label ||
				existingParam.Type != string(newParam.Type) ||
				existingParam.Role != string(newParam.Role) {
				existingParam.Label = newParam.Label
				existingParam.Type = string(newParam.Type)
				existingParam.Role = string(newParam.Role)

				if err := r.db.Save(&existingParam).Error; err != nil {
					return sharedUtils.NewFailureResult[dllModel.SDType](err)
				}
			}
			continue
		}
		paramEntity := dbModel.SDParameterEntity{
			Denotation: newParam.Denotation,
			Label:      newParam.Label,
			Type:       string(newParam.Type),
			Role:       string(newParam.Role),
			SDTypeID:   entity.ID,
		}
		if err := dbUtil.PersistEntityIntoDB(r.db, &paramEntity); err != nil {
			return sharedUtils.NewFailureResult[dllModel.SDType](err)
		}
	}
	for denotation, existingParam := range existingMap {
		if _, ok := incomingMap[denotation]; !ok {
			if err := r.db.Delete(&existingParam).Error; err != nil {
				return sharedUtils.NewFailureResult[dllModel.SDType](err)
			}
		}
	}
	updated := loadSDType(r.db, dbUtil.Where("id = ?", entity.ID))
	if updated.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDType](updated.GetError())
	}

	return updated
}

func loadSDType(g *gorm.DB, whereClause dbUtil.WhereClause) sharedUtils.Result[dllModel.SDType] {
	sdTypeEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](g, dbUtil.Preload("Parameters"), whereClause)
	if sdTypeEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDType](sdTypeEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[dllModel.SDType](db2dll.ToDLLModelSDType(sdTypeEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) LoadSDType(id uint32) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return loadSDType(r.db, dbUtil.Where("id = ?", id))
}

func (r *relationalDatabaseClientImpl) LoadSDTypeBasedOnUID(uid string) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return loadSDType(r.db, dbUtil.Where("uid = ?", uid))
}

func (r *relationalDatabaseClientImpl) LoadSDTypes() sharedUtils.Result[[]dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdTypeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Preload("Parameters"))
	if sdTypeEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDType](sdTypeEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.SDType](sharedUtils.Map(sdTypeEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDType))
}

func (r *relationalDatabaseClientImpl) DeleteSDType(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Transaction(func(tx *gorm.DB) error {
		relatedKPIDefinitionEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](tx, dbUtil.Where("sd_type_id = ?", id))
		if relatedKPIDefinitionEntitiesLoadResult.IsFailure() {
			return fmt.Errorf("failed to load KPI definition entities related to the SD type with ID = %d from the database: %w", id, relatedKPIDefinitionEntitiesLoadResult.GetError())
		}
		if err := dbUtil.DeleteCertainEntityBasedOnId[dbModel.SDTypeEntity](tx, id); err != nil {
			return fmt.Errorf("failed to delete SD type entity with ID = %d from the database: %w", id, err)
		}
		relatedKPIDefinitionEntities := relatedKPIDefinitionEntitiesLoadResult.GetPayload()
		for _, relatedKPIDefinitionEntity := range relatedKPIDefinitionEntities {
			rootNodeID := sharedUtils.NewOptionalFromPointer(relatedKPIDefinitionEntity.RootNodeID).GetPayload()
			if err := dbUtil.DeleteCertainEntityBasedOnId[dbModel.KPINodeEntity](tx, rootNodeID); err != nil {
				return fmt.Errorf("failed to delete KPI node entity with ID = %d from the database: %w", rootNodeID, err)
			}
		}
		return nil
	})
}

func (r *relationalDatabaseClientImpl) PersistSDInstance(sdInstance dllModel.SDInstance) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntity := dll2db.ToDBModelEntitySDInstance(sdInstance)
	if err := dbUtil.PersistEntityIntoDB(r.db, &sdInstanceEntity); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	return sharedUtils.NewSuccessResult[uint32](sdInstanceEntity.ID)
}

func (r *relationalDatabaseClientImpl) UpsertSDInstance(uid string, sdTypeSpecification string, label string) sharedUtils.Result[dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("uid = ?", uid))
	if existingResult.IsSuccess() {
		entity := existingResult.GetPayload()
		if label != "" && entity.Label != label {
			entity.Label = label
			if err := r.db.Save(&entity).Error; err != nil {
				return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
			}
		}
		return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDInstance(entity))
	}
	if err := existingResult.GetError(); !errors.Is(err, gorm.ErrRecordNotFound) {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
	}
	sdTypeResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Where("uid = ?", sdTypeSpecification))
	var sdTypeEntity dbModel.SDTypeEntity
	if sdTypeResult.IsFailure() {
		sdTypeErr := sdTypeResult.GetError()
		if errors.Is(sdTypeErr, gorm.ErrRecordNotFound) {
			sdTypeEntity = dbModel.SDTypeEntity{
				UID:   sdTypeSpecification,
				Label: sharedUtils.SafeLabel("", sdTypeSpecification),
			}
			if err := dbUtil.PersistEntityIntoDB(r.db, &sdTypeEntity); err != nil {
				return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
			}
		} else {
			return sharedUtils.NewFailureResult[dllModel.SDInstance](sdTypeErr)
		}
	} else {
		sdTypeEntity = sdTypeResult.GetPayload()
	}
	entity := dbModel.SDInstanceEntity{
		UID:             uid,
		Label:           sharedUtils.SafeLabel(label, uid),
		ConfirmedByUser: false,
		UserIdentifier:  uid,
		SDTypeID:        sdTypeEntity.ID,
		SDType:          sdTypeEntity,
	}
	if err := dbUtil.PersistEntityIntoDB(r.db, &entity); err != nil {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDInstance(entity))
}

func (r *relationalDatabaseClientImpl) LoadSDInstance(id uint32) sharedUtils.Result[dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("id = ?", id))
	if sdInstanceEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](sdInstanceEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[dllModel.SDInstance](db2dll.ToDLLModelSDInstance(sdInstanceEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceBasedOnUID(uid string) sharedUtils.Result[sharedUtils.Optional[dllModel.SDInstance]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("uid = ?", uid))
	if sdInstanceEntityLoadResult.IsFailure() {
		err := sdInstanceEntityLoadResult.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.SDInstance]](sharedUtils.NewEmptyOptional[dllModel.SDInstance]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.SDInstance]](err)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.SDInstance]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelSDInstance(sdInstanceEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) LoadSDInstances() sharedUtils.Result[[]dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"))
	if sdInstanceEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstance](sdInstanceEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.SDInstance](sharedUtils.Map(sdInstanceEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDInstance))
}

func (r *relationalDatabaseClientImpl) LoadSDInstancesByType(sdTypeID uint32) sharedUtils.Result[[]dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("sd_type_id = ?", sdTypeID))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstance](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), db2dll.ToDLLModelSDInstance))
}

func (r *relationalDatabaseClientImpl) LoadSDInstancesByKpiDefinition(kpiDefinitionID uint32) sharedUtils.Result[[]dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiResult := dbUtil.LoadEntityFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"), dbUtil.Where("id = ?", kpiDefinitionID))
	if kpiResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstance](kpiResult.GetError())
	}
	kpi := kpiResult.GetPayload()
	if kpi.SDInstanceMode == string(sharedModel.ALL) {
		instancesResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("sd_type_id = ?", kpi.SDTypeID))
		if instancesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]dllModel.SDInstance](instancesResult.GetError())
		}
		return sharedUtils.NewSuccessResult(sharedUtils.Map(instancesResult.GetPayload(), db2dll.ToDLLModelSDInstance))
	}
	if kpi.SDInstanceMode == string(sharedModel.SELECTED) {
		relations := kpi.SDInstanceKPIDefinitionRelationshipRecords
		if len(relations) == 0 {
			return sharedUtils.NewSuccessResult([]dllModel.SDInstance{})
		}
		uids := sharedUtils.Map(relations, func(r dbModel.SDInstanceKPIDefinitionRelationshipEntity) string {
			return r.SDInstanceUID
		})
		instancesResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("uid IN (?)", uids))
		if instancesResult.IsFailure() {
			return sharedUtils.NewFailureResult[[]dllModel.SDInstance](instancesResult.GetError())
		}
		return sharedUtils.NewSuccessResult(sharedUtils.Map(instancesResult.GetPayload(), db2dll.ToDLLModelSDInstance))
	}
	return sharedUtils.NewSuccessResult([]dllModel.SDInstance{})
}

func (r *relationalDatabaseClientImpl) PersistRawDataPoints(points []dllModel.RawDataPoint) sharedUtils.Result[[]dllModel.RawDataPoint] {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(points) == 0 {
		return sharedUtils.NewSuccessResult([]dllModel.RawDataPoint{})
	}
	points = dedupeRawDataPoints(points)
	entities := sharedUtils.Map(points, dll2db.ToDBModelRawDataPoint)
	tx := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "sd_instance_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"event_time": gorm.Expr("EXCLUDED.event_time"),
			"payload":    gorm.Expr("EXCLUDED.payload"),
		}),
		Where: clause.Where{
			Exprs: []clause.Expression{
				gorm.Expr("EXCLUDED.event_time > raw_data_point_entities.event_time"),
			},
		},
	}).Create(&entities)
	if tx.Error != nil {
		return sharedUtils.NewFailureResult[[]dllModel.RawDataPoint](tx.Error)
	}
	return sharedUtils.NewSuccessResult(points)
}

func (r *relationalDatabaseClientImpl) LoadRawDataPointsBySDType(sdTypeID uint32) sharedUtils.Result[[]dllModel.RawDataPoint] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.RawDataPointEntity](r.db, dbUtil.Preload("SDInstance"), dbUtil.Where("sd_instance_id IN (?)", r.db.Table("sd_instances").Select("id").Where("sd_type_id = ?", sdTypeID)))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.RawDataPoint](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), db2dll.ToDLLModelRawDataPoint))
}

func (r *relationalDatabaseClientImpl) LoadRawDataPoint(sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.RawDataPoint]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntityFromDB[dbModel.RawDataPointEntity](r.db, dbUtil.Preload("SDInstance"), dbUtil.Where("sd_instance_id = ?", sdInstanceID))
	if result.IsFailure() {
		err := result.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.RawDataPoint]())
		}
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.RawDataPoint]](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(db2dll.ToDLLModelRawDataPoint(result.GetPayload())))
}

func (r *relationalDatabaseClientImpl) PersistKPIFulfillmentCheckResults(points []dllModel.KPIFulfillmentCheckResult, reprocess bool) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(points) == 0 {
		return sharedUtils.NewSuccessResult([]dllModel.KPIFulfillmentCheckResult{})
	}
	points = dedupeKPIFulfillmentCheckResults(points)
	entities := sharedUtils.Map(points, dll2db.ToDBModelEntityKPIFulfillmentCheckResult)
	tx := &gorm.DB{}
	if reprocess {
		tx = r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "kpi_definition_id"},
				{Name: "sd_instance_id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"event_time": gorm.Expr("EXCLUDED.event_time"),
				"fulfilled":  gorm.Expr("EXCLUDED.fulfilled"),
			}),
		}).Create(&entities)
	} else {
		tx = r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "kpi_definition_id"},
				{Name: "sd_instance_id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"event_time": gorm.Expr("EXCLUDED.event_time"),
				"fulfilled":  gorm.Expr("EXCLUDED.fulfilled"),
			}),
			Where: clause.Where{
				Exprs: []clause.Expression{
					gorm.Expr("EXCLUDED.event_time > kpi_fulfillment_check_results.event_time"),
				},
			},
		}).Create(&entities)
	}
	if tx.Error != nil {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](tx.Error)
	}
	return sharedUtils.NewSuccessResult(points)
}

func dedupeRawDataPoints(points []dllModel.RawDataPoint) []dllModel.RawDataPoint {
	latestByInstance := make(map[uint32]dllModel.RawDataPoint, len(points))
	for _, point := range points {
		current, exists := latestByInstance[point.SDInstanceID]
		if !exists || point.EventTime.After(current.EventTime) {
			latestByInstance[point.SDInstanceID] = point
		}
	}
	deduped := make([]dllModel.RawDataPoint, 0, len(latestByInstance))
	for _, point := range latestByInstance {
		deduped = append(deduped, point)
	}
	return deduped
}

func dedupeKPIFulfillmentCheckResults(points []dllModel.KPIFulfillmentCheckResult) []dllModel.KPIFulfillmentCheckResult {
	type key struct {
		KPIDefinitionID uint32
		SDInstanceID    uint32
	}
	latestByKey := make(map[key]dllModel.KPIFulfillmentCheckResult, len(points))
	for _, point := range points {
		k := key{
			KPIDefinitionID: point.KPIDefinitionID,
			SDInstanceID:    point.SDInstanceID,
		}
		current, exists := latestByKey[k]
		if !exists || point.EventTime.After(current.EventTime) {
			latestByKey[k] = point
		}
	}
	deduped := make([]dllModel.KPIFulfillmentCheckResult, 0, len(latestByKey))
	for _, point := range latestByKey {
		deduped = append(deduped, point)
	}
	return deduped
}

func (r *relationalDatabaseClientImpl) LoadKPIFulfillmentCheckResults(userID uint32) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db, dbUtil.Preload("SDInstance"), dbUtil.Where("kpi_definition_id IN (?)", r.db.Table("kpi_definitions").Select("id").Where("user_id = ?", userID)))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), db2dll.ToDLLModelKPIFulfillmentCheckResult))
}

func (r *relationalDatabaseClientImpl) LoadKPIFulfillmentCheckResultsByKPI(userID uint32, kpiDefinitionID uint32) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db, dbUtil.Preload("SDInstance"), dbUtil.Where(`kpi_definition_id = ? AND kpi_definition_id IN (?)`, kpiDefinitionID, r.db.Table("kpi_definitions").Select("id").Where("user_id = ?", userID)))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), db2dll.ToDLLModelKPIFulfillmentCheckResult))
}

func (r *relationalDatabaseClientImpl) LoadKPIFulfillmentCheckResult(userID uint32, input dllModel.KPIFulfillmentCheckResultRequest) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntityFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db, dbUtil.Preload("SDInstance"), dbUtil.Where(`kpi_definition_id = ? AND sd_instance_id = ? AND kpi_definition_id IN (?)`, input.KPIDefinitionID, input.SDInstanceID, r.db.Table("kpi_definitions").Select("id").Where("user_id = ?", userID)))
	if result.IsFailure() {
		err := result.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.KPIFulfillmentCheckResult]())
		}
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(db2dll.ToDLLModelKPIFulfillmentCheckResult(result.GetPayload())))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceGroups() sharedUtils.Result[[]dllModel.SDInstanceGroup] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceGroupEntity](r.db, dbUtil.Preload("GroupMembershipRecords"))
	if sdInstanceGroupEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstanceGroup](sdInstanceGroupEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(sdInstanceGroupEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDInstanceGroup))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceGroup(id uint32) sharedUtils.Result[dllModel.SDInstanceGroup] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceGroupEntity](r.db, dbUtil.Preload("GroupMembershipRecords"), dbUtil.Where("id = ?", id))
	if sdInstanceGroupEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDInstanceGroup](sdInstanceGroupEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDInstanceGroup(sdInstanceGroupEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) PersistSDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntity := dll2db.ToDBModelEntitySDInstanceGroup(sdInstanceGroup)
	if err := dbUtil.PersistEntityIntoDB(r.db, &sdInstanceGroupEntity); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	return sharedUtils.NewSuccessResult[uint32](sdInstanceGroupEntity.ID)
}

func (r *relationalDatabaseClientImpl) DeleteSDInstanceGroup(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return dbUtil.DeleteCertainEntityBasedOnId[dbModel.SDInstanceGroupEntity](r.db, id)
}

func (r *relationalDatabaseClientImpl) PersistUser(user dllModel.User) sharedUtils.Result[uint] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userEntity := dll2db.ToDBModelEntityUser(user)
	if err := dbUtil.PersistEntityIntoDB(r.db, &userEntity); err != nil {
		return sharedUtils.NewFailureResult[uint](err)
	}
	return sharedUtils.NewSuccessResult[uint](userEntity.Model.ID)
}

func (r *relationalDatabaseClientImpl) LoadUserBasedOnOAuth2ProviderIssuedID(oauth2ProviderIssuedID string) sharedUtils.Result[sharedUtils.Optional[dllModel.User]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.UserEntity](r.db, dbUtil.Where("oauth2_provider_issued_id = ?", oauth2ProviderIssuedID))
	if userEntityLoadResult.IsFailure() {
		userEntityLoadError := userEntityLoadResult.GetError()
		if errors.Is(userEntityLoadError, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.User]](sharedUtils.NewEmptyOptional[dllModel.User]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.User]](userEntityLoadError)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.User]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelUser(userEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) LoadUser(id uint) sharedUtils.Result[dllModel.User] {
	r.mu.Lock()
	defer r.mu.Unlock()
	// TODO: Implement
	return sharedUtils.NewFailureResult[dllModel.User](errors.New("[RDB client (GORM)]: not implemented"))
}

func (r *relationalDatabaseClientImpl) LoadUserSessionBasedOnRefreshTokenHash(refreshTokenHash string) sharedUtils.Result[sharedUtils.Optional[dllModel.UserSession]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userSessionEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.UserSessionEntity](r.db, dbUtil.Where("refresh_token_hash = ?", refreshTokenHash))
	if userSessionEntityLoadResult.IsFailure() {
		userSessionEntityLoadError := userSessionEntityLoadResult.GetError()
		if errors.Is(userSessionEntityLoadError, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.UserSession]](sharedUtils.NewEmptyOptional[dllModel.UserSession]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.UserSession]](userSessionEntityLoadError)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.UserSession]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelUserSession(userSessionEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) PersistUserSession(userSession dllModel.UserSession) sharedUtils.Result[uint] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userSessionEntity := dll2db.ToDBModelEntityUserSession(userSession)
	if err := dbUtil.PersistEntityIntoDB(r.db, &userSessionEntity); err != nil {
		return sharedUtils.NewFailureResult[uint](err)
	}
	return sharedUtils.NewSuccessResult[uint](userSessionEntity.ID)
}

func (r *relationalDatabaseClientImpl) PersistUserConfig(userConfig dllModel.UserConfig) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userConfigEntity := dll2db.ToDBModelEntityUserConfig(userConfig)
	if err := dbUtil.PersistEntityIntoDB(r.db, &userConfigEntity); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	return sharedUtils.NewSuccessResult[uint32](userConfigEntity.UserID)
}

func (r *relationalDatabaseClientImpl) LoadUserConfig(userId uint32) sharedUtils.Result[dllModel.UserConfig] {
	r.mu.Lock()
	defer r.mu.Unlock()
	userConfigEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.UserConfigEntity](r.db, dbUtil.Where("user_id = ?", userId))
	if userConfigEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.UserConfig](userConfigEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelUserConfig(userConfigEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) DeleteUserConfig(userId uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return dbUtil.DeleteCertainEntityBasedOnId[dbModel.UserConfigEntity](r.db, userId)
}

func (r *relationalDatabaseClientImpl) LoadAPIKeyByHash(hash string) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := dbUtil.LoadEntityFromDB[dbModel.APIKeyEntity](r.db, dbUtil.Preload("IPRestrictions"), dbUtil.Preload("Permissions"), dbUtil.Where("key_hash = ?", hash))
	if result.IsFailure() {
		err := result.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.APIKey]())
		}
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](err)
	}
	entity := result.GetPayload()
	dll := db2dll.ToDLLModelAPIKey(entity)
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(dll))
}

func (r *relationalDatabaseClientImpl) LoadAPIKeyByID(id uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]] {
	var entity dbModel.APIKeyEntity
	err := r.db.Preload("Permissions").Preload("IPRestrictions").First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.APIKey]())
		}
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(db2dll.ToDLLModelAPIKey(entity)))
}

func (r *relationalDatabaseClientImpl) LoadAPIKeysForUser(userID uint32) sharedUtils.Result[[]dllModel.APIKey] {
	var entities []dbModel.APIKeyEntity
	err := r.db.Preload("Permissions").Preload("IPRestrictions").Where("user_id = ?", userID).Find(&entities).Error
	if err != nil {
		return sharedUtils.NewFailureResult[[]dllModel.APIKey](err)
	}
	result := sharedUtils.Map(entities, func(e dbModel.APIKeyEntity) dllModel.APIKey {
		return db2dll.ToDLLModelAPIKey(e)
	})
	return sharedUtils.NewSuccessResult(result)
}

func (r *relationalDatabaseClientImpl) CreateAPIKey(userID uint32, k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey] {
	perms, err := r.loadPermissionsByUIDs(k.Permissions)
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	entity := dbModel.APIKeyEntity{
		UserID:    userID,
		Label:     k.Label,
		ExpiresAt: k.ExpiresAt,
		Revoked:   k.Revoked,
		RateLimit: k.RateLimit,
	}
	if k.KeyHash != nil {
		entity.KeyHash = *k.KeyHash
	}
	if err := r.db.Create(&entity).Error; err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	if err := r.db.Model(&entity).Association("Permissions").Replace(perms); err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	if k.IPRestrictions != nil {
		newList := uniqueStrings(k.IPRestrictions)
		currentMap := make(map[string]dbModel.APIKeyIPRestrictionEntity)
		for _, row := range entity.IPRestrictions {
			currentMap[row.CIDR] = row
		}
		newMap := make(map[string]struct{})
		for _, cidr := range newList {
			newMap[cidr] = struct{}{}
		}
		for cidr, row := range currentMap {
			if _, ok := newMap[cidr]; !ok {
				if err := r.db.Delete(&row).Error; err != nil {
					return sharedUtils.NewFailureResult[dllModel.APIKey](err)
				}
			}
		}
		for cidr := range newMap {
			if _, ok := currentMap[cidr]; !ok {
				ip := dbModel.APIKeyIPRestrictionEntity{
					APIKeyID: entity.ID,
					CIDR:     cidr,
				}
				if err := r.db.Create(&ip).Error; err != nil {
					return sharedUtils.NewFailureResult[dllModel.APIKey](err)
				}
			}
		}
	}
	err = r.db.Preload("Permissions").Preload("IPRestrictions").First(&entity, entity.ID).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelAPIKey(entity))
}

func (r *relationalDatabaseClientImpl) loadPermissionsByUIDs(uids []string) ([]dbModel.PermissionEntity, error) {
	uniqueUIDs := uniqueStrings(uids)
	if len(uniqueUIDs) == 0 {
		return []dbModel.PermissionEntity{}, nil
	}
	var perms []dbModel.PermissionEntity
	err := r.db.Where("uid IN ?", uniqueUIDs).Find(&perms).Error
	if err != nil {
		return nil, err
	}
	existing := make(map[string]bool, len(perms))
	for _, p := range perms {
		existing[p.UID] = true
	}
	for _, uid := range uniqueUIDs {
		if !existing[uid] {
			return nil, fmt.Errorf("permission '%s' not found", uid)
		}
	}
	return perms, nil
}

func (r *relationalDatabaseClientImpl) UpdateAPIKey(k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey] {
	if k.ID.IsEmpty() || k.UserID == nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](fmt.Errorf("missing id"))
	}
	var entity dbModel.APIKeyEntity
	err := r.db.Preload("Permissions").Preload("IPRestrictions").First(&entity, k.ID.GetPayload()).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	entity.Label = k.Label
	entity.ExpiresAt = k.ExpiresAt
	entity.Revoked = k.Revoked
	entity.RateLimit = k.RateLimit
	if k.Permissions != nil {
		perms, err := r.loadPermissionsByUIDs(k.Permissions)
		if err != nil {
			return sharedUtils.NewFailureResult[dllModel.APIKey](err)
		}
		if err := r.db.Model(&entity).Association("Permissions").Replace(perms); err != nil {
			return sharedUtils.NewFailureResult[dllModel.APIKey](err)
		}
	}
	if k.IPRestrictions != nil {
		if err := r.db.Where("api_key_id = ?", entity.ID).
			Delete(&dbModel.APIKeyIPRestrictionEntity{}).Error; err != nil {
			return sharedUtils.NewFailureResult[dllModel.APIKey](err)
		}
		if k.IPRestrictions != nil {
			newList := uniqueStrings(k.IPRestrictions)
			currentMap := make(map[string]dbModel.APIKeyIPRestrictionEntity)
			for _, r := range entity.IPRestrictions {
				currentMap[r.CIDR] = r
			}
			newMap := make(map[string]struct{})
			for _, cidr := range newList {
				newMap[cidr] = struct{}{}
			}
			for cidr, row := range currentMap {
				if _, ok := newMap[cidr]; !ok {
					if err := r.db.Delete(&row).Error; err != nil {
						return sharedUtils.NewFailureResult[dllModel.APIKey](err)
					}
				}
			}
			for cidr := range newMap {
				if _, ok := currentMap[cidr]; !ok {
					ip := dbModel.APIKeyIPRestrictionEntity{
						APIKeyID: entity.ID,
						CIDR:     cidr,
					}
					if err := r.db.Create(&ip).Error; err != nil {
						return sharedUtils.NewFailureResult[dllModel.APIKey](err)
					}
				}
			}
		}
	}
	if err := r.db.Save(&entity).Error; err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelAPIKey(entity))
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))
	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func (r *relationalDatabaseClientImpl) DeleteAPIKey(id uint32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var entity dbModel.APIKeyEntity
		if err := tx.First(&entity, id).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity).Association("Permissions").Clear(); err != nil {
			return err
		}
		if err := tx.Delete(&entity).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *relationalDatabaseClientImpl) LoadRoles() sharedUtils.Result[[]dllModel.Role] {
	var roles []dbModel.RoleEntity
	err := r.db.Preload("Permissions").Order("label ASC").Find(&roles).Error
	if err != nil {
		return sharedUtils.NewFailureResult[[]dllModel.Role](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(roles, db2dll.ToDLLModelRole))
}

func (r *relationalDatabaseClientImpl) LoadRoleByUID(uid string) sharedUtils.Result[dllModel.Role] {
	var role dbModel.RoleEntity
	err := r.db.Preload("Permissions").Where("label = ?", uid).First(&role).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.Role](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelRole(role))
}

func (r *relationalDatabaseClientImpl) GetUserRole(userID uint32) sharedUtils.Result[dllModel.Role] {
	var user dbModel.UserEntity
	err := r.db.Preload("Role.Permissions").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.Role](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelRole(user.Role))
}

func (r *relationalDatabaseClientImpl) SetUserRole(userID uint32, roleID uint32) error {
	err := r.db.Model(&dbModel.UserEntity{}).Where("id = ?", userID).Update("role_id", roleID).Error
	if err != nil {
		return err
	}
	return nil
}
