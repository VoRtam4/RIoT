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
	"gorm.io/gorm/logger"
)

var (
	rdbClientInstance RelationalDatabaseClient
	once              sync.Once
)

type RelationalDatabaseClient interface {
	setup()
	PerformOnStartupOperations(permissions map[string]map[string]bool) error
	PersistKPIDefinition(userID uint32, kpiDefinition sharedModel.KPIDefinition) sharedUtils.Result[uint32]
	LoadAllKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition]
	LoadKPIDefinition(userID uint32, id uint32) sharedUtils.Result[sharedModel.KPIDefinition]
	LoadKPIDefinitions(userID uint32) sharedUtils.Result[[]sharedModel.KPIDefinition]
	DeleteKPIDefinition(id uint32) error
	PersistSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType]
	UpsertSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType]
	LoadSDType(id uint32) sharedUtils.Result[dllModel.SDType]
	LoadSDInstancesByType(sdTypeID uint32) sharedUtils.Result[[]dllModel.SDInstance]
	LoadSDInstancesByKpiDefinition(kpiDefinitionID uint32) sharedUtils.Result[[]dllModel.SDInstance]
	LoadSDTypeBasedOnDenotation(denotation string) sharedUtils.Result[dllModel.SDType]
	LoadSDTypes() sharedUtils.Result[[]dllModel.SDType]
	DeleteSDType(id uint32) error
	PersistSDInstance(sdInstance dllModel.SDInstance) sharedUtils.Result[uint32]
	PersistNewSDInstance(uid string, sdTypeSpecification string) sharedUtils.Result[dllModel.SDInstance]
	LoadSDInstance(id uint32) sharedUtils.Result[dllModel.SDInstance]
	LoadSDInstanceBasedOnUID(uid string) sharedUtils.Result[sharedUtils.Optional[dllModel.SDInstance]]
	LoadSDInstances() sharedUtils.Result[[]dllModel.SDInstance]
	PersistKPIFulFulfillmentCheckResultTuple(sdInstanceUID string, kpiDefinitionIDs []uint32, fulfillmentStatuses []bool, eventTimes []time.Time) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
	LoadKPIFulFulfillmentCheckResult(kpiDefinitionID uint32, sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]]
	LoadKPIFulFulfillmentCheckResults() sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
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
	GetRoleIDByLabel(label string) sharedUtils.Result[dllModel.Role]
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
		new(dbModel.PermissionEntity),
		new(dbModel.OperationTypeAccessPermissionEntity),
		new(dbModel.SingleOperationPermissionEntity),
		new(dbModel.APIKeyEntity),
		new(dbModel.APIKeyIPRestrictionEntity),
	), "[RDB client (GORM)]: auto-migration failed")
}

func (r *relationalDatabaseClientImpl) PerformOnStartupOperations(permissions map[string]map[string]bool) error {
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
				err := tx.Where("label = ?", permLabel).First(&perm).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					perm = dbModel.PermissionEntity{Label: permLabel}
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

	referencedSDInstancesResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](
		r.db,
		dbUtil.Where("uid IN (?)", kpiDefinition.SelectedSDInstanceUIDs),
	)
	if referencedSDInstancesResult.IsFailure() {
		return sharedUtils.NewFailureResult[uint32](referencedSDInstancesResult.GetError())
	}
	referencedSDInstances := referencedSDInstancesResult.GetPayload()

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

		return nil
	})

	if err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}

	if kpiDefinitionID != 0 {
		if err := dbUtil.DeleteEntitiesBasedOnWhereClauses[dbModel.SDInstanceKPIDefinitionRelationshipEntity](r.db, dbUtil.Where("kpi_definition_id = ?", kpiDefinitionID)); err != nil {
			return sharedUtils.NewFailureResult[uint32](err)
		}
	}

	sdInstanceRelations := sharedUtils.Map(referencedSDInstances, func(sdInstance dbModel.SDInstanceEntity) dbModel.SDInstanceKPIDefinitionRelationshipEntity {
		return dbModel.SDInstanceKPIDefinitionRelationshipEntity{
			KPIDefinitionID: kpiDefinitionEntity.ID,
			SDInstanceID:    sdInstance.ID,
			SDInstanceUID:   sdInstance.UID,
		}
	})

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range sdInstanceRelations {
			if err := dbUtil.PersistEntityIntoDB(tx, &entity); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}

	if err := dbUtil.DeleteEntitiesBasedOnSliceOfIds[dbModel.KPINodeEntity](r.db, idsOfKPINodeEntitiesFormingTheKPIDefinition); err != nil {
		log.Printf("cleanup failed: %s\n", err.Error())
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
	result := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Where("user_id = ?", userID), dbUtil.Preload("SDType"), dbUtil.Preload("SDInstanceKPIDefinitionRelationshipRecords"))
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
	kpiNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](
		r.db,
		dbUtil.Where("id IN (?)", allNodeIDs),
	)
	if kpiNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](kpiNodesResult.GetError())
	}
	kpiNodes := kpiNodesResult.GetPayload()

	logicalNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](
		r.db,
		dbUtil.Where("node_id IN (?)", allNodeIDs),
	)
	if logicalNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](logicalNodesResult.GetError())
	}
	logicalNodes := logicalNodesResult.GetPayload()

	atomNodesResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](r.db, dbUtil.Where("node_id IN (?)", allNodeIDs), dbUtil.Preload("SDParameter"))
	if atomNodesResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](atomNodesResult.GetError())
	}
	atomNodes := atomNodesResult.GetPayload()
	kpis := make([]sharedModel.KPIDefinition, 0)
	for _, entity := range entities {
		kpi := db2dll.ToDLLModelKPIDefinition(entity, kpiNodes, logicalNodes, atomNodes)
		kpis = append(kpis, kpi)
	}
	return sharedUtils.NewSuccessResult(kpis)
}

func (r *relationalDatabaseClientImpl) LoadAllKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()

	kpiDefinitionEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Preload("SDType", "SDInstanceKPIDefinitionRelationshipRecords"))
	if kpiDefinitionEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI definition entities from the database: %w", kpiDefinitionEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	kpiDefinitionEntities := kpiDefinitionEntitiesLoadResult.GetPayload()

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
	existingResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Preload("Parameters"), dbUtil.Where("denotation = ?", sdType.Denotation))
	if existingResult.IsFailure() {
		if errors.Is(existingResult.GetError(), gorm.ErrRecordNotFound) {
			sdTypeEntity := dll2db.ToDBModelEntitySDType(sdType)
			if err := dbUtil.PersistEntityIntoDB(r.db, &sdTypeEntity); err != nil {
				return sharedUtils.NewFailureResult[dllModel.SDType](err)
			}
			return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDType(sdTypeEntity))
		}
		return sharedUtils.NewFailureResult[dllModel.SDType](existingResult.GetError())
	}
	existing := existingResult.GetPayload()
	existingParams := map[string]bool{}
	for _, p := range existing.Parameters {
		existingParams[p.Denotation] = true
	}
	for _, newParam := range sdType.Parameters {
		if existingParams[newParam.Denotation] {
			continue
		}
		paramEntity := dbModel.SDParameterEntity{
			Denotation: newParam.Denotation,
			Type:       string(newParam.Type),
			SDTypeID:   existing.ID,
		}
		if err := dbUtil.PersistEntityIntoDB(r.db, &paramEntity); err != nil {
			return sharedUtils.NewFailureResult[dllModel.SDType](err)
		}
		log.Printf("SDType %s rozšířen o parametr %s", sdType.Denotation, newParam.Denotation)
	}
	updated := loadSDType(r.db, dbUtil.Where("id = ?", existing.ID))
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

func (r *relationalDatabaseClientImpl) LoadSDTypeBasedOnDenotation(denotation string) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return loadSDType(r.db, dbUtil.Where("denotation = ?", denotation))
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

func (r *relationalDatabaseClientImpl) PersistNewSDInstance(uid string, sdTypeSpecification string) sharedUtils.Result[dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Where("uid = ?", uid))
	if sdInstanceEntityLoadResult.IsSuccess() {
		return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDInstance(sdInstanceEntityLoadResult.GetPayload()))
	} else if sdInstanceEntityLoadError := sdInstanceEntityLoadResult.GetError(); !errors.Is(sdInstanceEntityLoadError, gorm.ErrRecordNotFound) {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](sdInstanceEntityLoadError)
	}
	referencedSDTypeEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Where("denotation = ?", sdTypeSpecification))
	if referencedSDTypeEntityLoadResult.IsFailure() {
		referencedSDTypeEntityLoadError := referencedSDTypeEntityLoadResult.GetError()
		err := sharedUtils.Ternary(errors.Is(referencedSDTypeEntityLoadError, gorm.ErrRecordNotFound), ErrOperationWouldLeadToForeignKeyIntegrityBreach, referencedSDTypeEntityLoadError)
		return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
	}
	sdInstanceEntity := dbModel.SDInstanceEntity{
		UID:             uid,
		ConfirmedByUser: false,
		UserIdentifier:  uid,
		SDTypeID:        referencedSDTypeEntityLoadResult.GetPayload().ID,
	}
	if err := dbUtil.PersistEntityIntoDB(r.db, &sdInstanceEntity); err != nil {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](err)
	}
	return sharedUtils.NewSuccessResult[dllModel.SDInstance](db2dll.ToDLLModelSDInstance(sdInstanceEntity))
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

func (r *relationalDatabaseClientImpl) PersistKPIFulFulfillmentCheckResultTuple(sdInstanceUID string, kpiDefinitionIDs []uint32, fulfillmentStatuses []bool, eventTimes []time.Time) sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	referencedKPIDefinitionEntitiesExistCheckResult := dbUtil.DoIDsExist[dbModel.KPIDefinitionEntity](r.db, kpiDefinitionIDs)
	if referencedKPIDefinitionEntitiesExistCheckResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](referencedKPIDefinitionEntitiesExistCheckResult.GetError())
	}
	if referencedKPIDefinitionEntitiesExist := referencedKPIDefinitionEntitiesExistCheckResult.GetPayload(); !referencedKPIDefinitionEntitiesExist {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](ErrOperationWouldLeadToForeignKeyIntegrityBreach)
	}
	referencedSDInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Where("uid = ?", sdInstanceUID))
	if referencedSDInstanceEntityLoadResult.IsFailure() {
		referencedSDInstanceEntityLoadError := referencedSDInstanceEntityLoadResult.GetError()
		err := sharedUtils.Ternary(errors.Is(referencedSDInstanceEntityLoadError, gorm.ErrRecordNotFound), ErrOperationWouldLeadToForeignKeyIntegrityBreach, referencedSDInstanceEntityLoadError)
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](err)
	}
	referencedSDInstanceEntityID := referencedSDInstanceEntityLoadResult.GetPayload().ID
	kpiFulfillmentCheckResultEntities := make([]dbModel.KPIFulfillmentCheckResultEntity, 0)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for index, kpiDefinitionID := range kpiDefinitionIDs {
			eventTime := eventTimes[index]
			existingResult := dbModel.KPIFulfillmentCheckResultEntity{}
			err := tx.Where("kpi_definition_id = ? AND sd_instance_id = ?", kpiDefinitionID, referencedSDInstanceEntityID).First(&existingResult).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newResult := dbModel.KPIFulfillmentCheckResultEntity{
						KPIDefinitionID: kpiDefinitionID,
						SDInstanceID:    referencedSDInstanceEntityID,
						Fulfilled:       fulfillmentStatuses[index],
						EventTime:       eventTime,
					}
					if err := tx.Create(&newResult).Error; err != nil {
						return err
					}
					kpiFulfillmentCheckResultEntities = append(kpiFulfillmentCheckResultEntities, newResult)
					continue
				}
				return err
			}
			if eventTime.After(existingResult.EventTime) {
				existingResult.Fulfilled = fulfillmentStatuses[index]
				existingResult.EventTime = eventTime
				if err := tx.Save(&existingResult).Error; err != nil {
					return err
				}
				kpiFulfillmentCheckResultEntities = append(kpiFulfillmentCheckResultEntities, existingResult)
			}
		}
		return nil
	})
	if err != nil {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(kpiFulfillmentCheckResultEntities, db2dll.ToDLLModelKPIFulfillmentCheckResult))
}

func (r *relationalDatabaseClientImpl) LoadKPIFulFulfillmentCheckResult(kpiDefinitionID uint32, sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiFulFulfillmentCheckResultEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db, dbUtil.Where("kpi_definition_id = ?", kpiDefinitionID), dbUtil.Where("sd_instance_id = ?", sdInstanceID))
	if kpiFulFulfillmentCheckResultEntityLoadResult.IsFailure() {
		err := kpiFulFulfillmentCheckResultEntityLoadResult.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](sharedUtils.NewEmptyOptional[dllModel.KPIFulfillmentCheckResult]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](err)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelKPIFulfillmentCheckResult(kpiFulFulfillmentCheckResultEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) LoadKPIFulFulfillmentCheckResults() sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiFulFulfillmentCheckResultEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db)
	if kpiFulFulfillmentCheckResultEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](kpiFulFulfillmentCheckResultEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.KPIFulfillmentCheckResult](sharedUtils.Map(kpiFulFulfillmentCheckResultEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelKPIFulfillmentCheckResult))
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
	result := dbUtil.LoadEntityFromDB[dbModel.APIKeyEntity](r.db, dbUtil.Preload("IPRestrictions"), dbUtil.Preload("Role.Permissions"), dbUtil.Where("key_hash = ?", hash))
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

func (c *relationalDatabaseClientImpl) LoadAPIKeyByID(id uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]] {
	var entity dbModel.APIKeyEntity
	err := c.db.Preload("Role.Permissions").Preload("IPRestrictions").First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.APIKey]())
		}
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](err)
	}
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(db2dll.ToDLLModelAPIKey(entity)))
}

func (c *relationalDatabaseClientImpl) LoadAPIKeysForUser(userID uint32) sharedUtils.Result[[]dllModel.APIKey] {
	var entities []dbModel.APIKeyEntity
	err := c.db.Preload("Role.Permissions").Preload("IPRestrictions").Where("user_id = ?", userID).Find(&entities).Error
	if err != nil {
		return sharedUtils.NewFailureResult[[]dllModel.APIKey](err)
	}
	result := sharedUtils.Map(entities, func(e dbModel.APIKeyEntity) dllModel.APIKey {
		return db2dll.ToDLLModelAPIKey(e)
	})
	return sharedUtils.NewSuccessResult(result)
}

func (c *relationalDatabaseClientImpl) CreateAPIKey(userID uint32, k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey] {
	perms, err := c.getOrCreatePermissions(k.Permissions)
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	role := dbModel.RoleEntity{
		Label: fmt.Sprintf("api_key_%d", time.Now().UnixNano()),
	}
	if err := c.db.Create(&role).Error; err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	if err := c.db.Model(&role).Association("Permissions").Append(perms); err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	entity := dbModel.APIKeyEntity{
		UserID:    userID,
		RoleID:    role.ID,
		Label:     k.Label,
		ExpiresAt: k.ExpiresAt,
		Revoked:   k.Revoked,
		RateLimit: k.RateLimit,
	}
	if k.KeyHash != nil {
		entity.KeyHash = *k.KeyHash
	}
	if err := c.db.Create(&entity).Error; err != nil {
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
		for cidr, r := range currentMap {
			if _, ok := newMap[cidr]; !ok {
				if err := c.db.Delete(&r).Error; err != nil {
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
				if err := c.db.Create(&ip).Error; err != nil {
					return sharedUtils.NewFailureResult[dllModel.APIKey](err)
				}
			}
		}
	}
	err = c.db.Preload("Role.Permissions").Preload("IPRestrictions").First(&entity, entity.ID).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelAPIKey(entity))
}

func (c *relationalDatabaseClientImpl) getOrCreatePermissions(labels []string) ([]dbModel.PermissionEntity, error) {
	if len(labels) == 0 {
		return []dbModel.PermissionEntity{}, nil
	}
	var perms []dbModel.PermissionEntity
	err := c.db.Where("label IN ?", labels).Find(&perms).Error
	if err != nil {
		return nil, err
	}
	existing := map[string]bool{}
	for _, p := range perms {
		existing[p.Label] = true
	}
	for _, label := range labels {
		if !existing[label] {
			p := dbModel.PermissionEntity{Label: label}
			if err := c.db.Create(&p).Error; err != nil {
				return nil, err
			}
			perms = append(perms, p)
		}
	}
	return perms, nil
}

func (c *relationalDatabaseClientImpl) UpdateAPIKey(k dllModel.APIKey) sharedUtils.Result[dllModel.APIKey] {
	if k.ID.IsEmpty() {
		return sharedUtils.NewFailureResult[dllModel.APIKey](fmt.Errorf("missing id"))
	}
	var entity dbModel.APIKeyEntity
	err := c.db.Preload("Role.Permissions").Preload("IPRestrictions").First(&entity, k.ID.GetPayload()).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.APIKey](err)
	}
	entity.Label = k.Label
	entity.ExpiresAt = k.ExpiresAt
	entity.Revoked = k.Revoked
	entity.RateLimit = k.RateLimit
	if k.Permissions != nil {
		perms, err := c.getOrCreatePermissions(k.Permissions)
		if err != nil {
			return sharedUtils.NewFailureResult[dllModel.APIKey](err)
		}
		if err := c.db.Model(&entity.Role).Association("Permissions").Replace(perms); err != nil {
			return sharedUtils.NewFailureResult[dllModel.APIKey](err)
		}
	}
	if k.IPRestrictions != nil {
		if err := c.db.Where("api_key_id = ?", entity.ID).
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
			for cidr, r := range currentMap {
				if _, ok := newMap[cidr]; !ok {
					if err := c.db.Delete(&r).Error; err != nil {
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
					if err := c.db.Create(&ip).Error; err != nil {
						return sharedUtils.NewFailureResult[dllModel.APIKey](err)
					}
				}
			}
		}
	}
	if err := c.db.Save(&entity).Error; err != nil {
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

func (c *relationalDatabaseClientImpl) DeleteAPIKey(id uint32) error {
	var entity dbModel.APIKeyEntity
	if err := c.db.Preload("Role").First(&entity, id).Error; err != nil {
		return err
	}
	if err := c.db.Delete(&entity).Error; err != nil {
		return err
	}
	if err := c.db.Delete(&dbModel.RoleEntity{}, entity.RoleID).Error; err != nil {
		return err
	}
	return nil
}

func (c *relationalDatabaseClientImpl) GetRoleIDByLabel(label string) sharedUtils.Result[dllModel.Role] {
	var role dbModel.RoleEntity
	err := c.db.Preload("Permissions").Where("label = ?", label).First(&role).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.Role](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelRole(role))
}

func (c *relationalDatabaseClientImpl) GetUserRole(userID uint32) sharedUtils.Result[dllModel.Role] {
	var user dbModel.UserEntity
	err := c.db.Preload("Role.Permissions").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.Role](err)
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelRole(user.Role))
}

func (c *relationalDatabaseClientImpl) SetUserRole(userID uint32, roleID uint32) error {
	err := c.db.Model(&dbModel.UserEntity{}).Where("id = ?", userID).Update("role_id", roleID).Error
	if err != nil {
		return err
	}
	return nil
}
