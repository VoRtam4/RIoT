package dbModel

import (
	"time"

	"gorm.io/gorm"
)

type KPIDefinitionEntity struct {
	ID                                         uint32                                      `gorm:"column:id;primaryKey"`
	UserIdentifier                             string                                      `gorm:"column:user_identifier;not null"`
	RootNodeID                                 *uint32                                     `gorm:"column:root_node_id;not null"`
	RootNode                                   *KPINodeEntity                              `gorm:"foreignKey:RootNodeID"`
	SDTypeID                                   uint32                                      `gorm:"column:sd_type_id;not null"`
	SDType                                     SDTypeEntity                                `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
	SDInstanceMode                             string                                      `gorm:"column:sd_instance_mode;not null"`
	KPIFulfillmentCheckResults                 []KPIFulfillmentCheckResultEntity           `gorm:"foreignKey:KPIDefinitionID;constraint:OnDelete:CASCADE"`
	SDInstanceKPIDefinitionRelationshipRecords []SDInstanceKPIDefinitionRelationshipEntity `gorm:"foreignKey:KPIDefinitionID;constraint:OnDelete:CASCADE"`
}

func (KPIDefinitionEntity) TableName() string {
	return "kpi_definitions"
}

type KPINodeEntity struct {
	ID           uint32         `gorm:"column:id;primaryKey;not null"`
	ParentNodeID *uint32        `gorm:"column:parent_node_id"`
	ParentNode   *KPINodeEntity `gorm:"foreignKey:ParentNodeID;constraint:OnDelete:CASCADE"`
}

func (KPINodeEntity) TableName() string {
	return "kpi_nodes"
}

type LogicalOperationKPINodeEntity struct {
	NodeID *uint32        `gorm:"column:node_id;primaryKey;not null"`
	Node   *KPINodeEntity `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE"`
	Type   string         `gorm:"column:type;not null"`
}

func (LogicalOperationKPINodeEntity) TableName() string {
	return "logical_operation_kpi_nodes"
}

type AtomKPINodeEntity struct {
	NodeID                *uint32           `gorm:"column:node_id;primaryKey;not null"`
	Node                  *KPINodeEntity    `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE"`
	SDParameterID         uint32            `gorm:"column:sd_parameter_id;not null"`
	SDParameter           SDParameterEntity `gorm:"foreignKey:SDParameterID;constraint:OnDelete:CASCADE"`
	Type                  string            `gorm:"column:type;not null"`
	StringReferenceValue  *string           `gorm:"column:string_reference_value"`
	BooleanReferenceValue *bool             `gorm:"column:boolean_reference_value"`
	NumericReferenceValue *float64          `gorm:"column:numeric_reference_value"`
}

func (AtomKPINodeEntity) TableName() string {
	return "atom_kpi_nodes"
}

type SDTypeEntity struct {
	ID         uint32              `gorm:"column:id;primaryKey;not null"`
	Denotation string              `gorm:"column:denotation;not null;index"` // Denotation is a separately indexed field
	Parameters []SDParameterEntity `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
	Commands   []SDCommandEntity   `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
}

func (SDTypeEntity) TableName() string {
	return "sd_types"
}

type SDParameterEntity struct {
	ID         uint32 `gorm:"column:id;primaryKey;not null"`
	SDTypeID   uint32 `gorm:"column:sd_type_id;not null"`
	Denotation string `gorm:"column:denotation;not null"`
	Type       string `gorm:"column:type;not null"`
}

func (SDParameterEntity) TableName() string {
	return "sd_parameters"
}

type SDInstanceEntity struct {
	ID                                         uint32                                      `gorm:"column:id;primaryKey;not null"`
	UID                                        string                                      `gorm:"column:uid;not null;index"` // UID is a separately indexed field
	ConfirmedByUser                            bool                                        `gorm:"column:confirmed_by_user;not null"`
	UserIdentifier                             string                                      `gorm:"column:user_identifier;not null"`
	SDTypeID                                   uint32                                      `gorm:"column:sd_type_id"`
	SDType                                     SDTypeEntity                                `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
	GroupMembershipRecords                     []SDInstanceGroupMembershipEntity           `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
	KPIFulfillmentCheckResults                 []KPIFulfillmentCheckResultEntity           `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
	SDInstanceKPIDefinitionRelationshipRecords []SDInstanceKPIDefinitionRelationshipEntity `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
}

func (SDInstanceEntity) TableName() string {
	return "sd_instances"
}

type KPIFulfillmentCheckResultEntity struct {
	KPIDefinitionID uint32    `gorm:"column:kpi_definition_id;primaryKey;not null"`
	SDInstanceID    uint32    `gorm:"column:sd_instance_id;primaryKey;not null"`
	Fulfilled       bool      `gorm:"column:fulfilled;not null"`
	EventTime       time.Time `gorm:"column:event_time;not null;index"`
}

func (KPIFulfillmentCheckResultEntity) TableName() string {
	return "kpi_fulfillment_check_results"
}

type SDInstanceGroupEntity struct {
	ID                     uint32                            `gorm:"column:id;primaryKey;not null"`
	UserIdentifier         string                            `gorm:"column:user_identifier;not null"`
	GroupMembershipRecords []SDInstanceGroupMembershipEntity `gorm:"foreignKey:SDInstanceGroupID;constraint:OnDelete:CASCADE"`
}

func (SDInstanceGroupEntity) TableName() string {
	return "sd_instance_groups"
}

type SDInstanceGroupMembershipEntity struct {
	SDInstanceGroupID uint32 `gorm:"column:sd_instance_group_id;primaryKey;not null;index"` // SDInstanceGroupID is a separately indexed field
	SDInstanceID      uint32 `gorm:"column:sd_instance_id;primaryKey;not null"`
}

func (SDInstanceGroupMembershipEntity) TableName() string {
	return "sd_instance_group_membership"
}

type SDInstanceKPIDefinitionRelationshipEntity struct { // TODO: Missing 'TableName' function, handle this across schema
	KPIDefinitionID uint32 `gorm:"column:kpi_definition_id;primaryKey;not null"`
	SDInstanceID    uint32 `gorm:"column:sd_instance_id;primaryKey;not null"`
	SDInstanceUID   string `gorm:"column:sd_instance_uid;not null"`
}

// UserEntity represents a user of the application who can log in using various OAuth providers.
type UserEntity struct {
	gorm.Model
	Username               string                      `gorm:"column:username;uniqueIndex"`
	Email                  string                      `gorm:"column:email;uniqueIndex"`
	Name                   *string                     `gorm:"column:name"`
	ProfileImageURL        *string                     `gorm:"column:profile_image_url"`
	OAuth2Provider         *string                     `gorm:"column:oauth2_provider;uniqueIndex:idx_oauth,priority:1"`
	OAuth2ProviderIssuedID *string                     `gorm:"column:oauth2_provider_issued_id;uniqueIndex:idx_oauth,priority:2"`
	LastLoginAt            *time.Time                  `gorm:"column:last_login_at"`
	Sessions               []UserSessionEntity         `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Invocations            []SDCommandInvocationEntity `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	UserConfig             UserConfigEntity            `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	/*
	   ID        uint           `gorm:"primaryKey"` // Primary key for the user
	   CreatedAt time.Time      // Timestamp of creation
	   UpdatedAt time.Time      // Timestamp of the last update
	   DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete field

	   // Basic User Information
	   Username     string `gorm:"uniqueIndex;size:100"` // Unique username for the user
	   Email        string `gorm:"uniqueIndex;size:255"` // User's email address (unique)
	   Name         string `gorm:"size:255"`             // Full name of the user
	   ProfileImage string `gorm:"size:500"`             // URL to the user's profile image

	   // OAuth Information
	   Provider     string    `gorm:"size:50"`        // OAuth provider name (e.g., google, github)
	   ProviderID   string    `gorm:"size:255;index"` // Unique ID provided by the OAuth provider
	   OAuthToken   string    `gorm:"size:500"`       // OAuth access token
	   RefreshToken string    `gorm:"size:500"`       // OAuth refresh token, if available
	   TokenExpiry  time.Time // Expiration time of the OAuth token

	   // Additional Metadata
	   LastLoginAt time.Time                   // Timestamp of the last login
	   IsActive    bool                        `gorm:"default:true"` // Whether the user's account is active
	   Invocations []SDCommandInvocationEntity `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	   UserConfig  UserConfigEntity            `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	*/
}

func (UserEntity) TableName() string { // TODO: Standardize table names, e.g. 'user' × 'users'
	return "user"
}

type UserSessionEntity struct {
	gorm.Model                 // TODO: Standardize 'gorm.Model' usage: either use it everywhere, or not at all
	UserID           uint      `gorm:"column:user_id"`
	RefreshTokenHash string    `gorm:"column:refresh_token_hash;not null;uniqueIndex"`
	ExpiresAt        time.Time `gorm:"column:expires_at"` // TODO: Are the 'column:...' entries necessary? If not, get rid of them across entire schema
	Revoked          bool      `gorm:"column:revoked;not null;default:false"`
	IPAddress        string    `gorm:"column:ip_address"`
	UserAgent        string    `gorm:"column:user_agent"`
}

func (UserSessionEntity) TableName() string {
	return "user_sessions"
}

type UserConfigEntity struct { // TODO: Consider embedding this inside 'UserEntity' for global user config and or inside 'UserSessionEntity' for per-session user config
	UserID uint32 `gorm:"primaryKey;column:user_id;not null"`
	Config string `gorm:"column:config;type:jsonb;not null"` // Store JSON as a string
}

func (UserConfigEntity) TableName() string {
	return "user_config"
}

type SDCommandEntity struct {
	ID          uint32                      `gorm:"column:id;primaryKey;not null"`
	SDTypeID    uint32                      `gorm:"column:sd_type_id;not null"`
	Denotation  string                      `gorm:"column:denotation;not null"`
	Type        string                      `gorm:"column:type;not null"`
	Payload     string                      `gorm:"column:payload;not null"`
	Invocations []SDCommandInvocationEntity `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
}

func (SDCommandEntity) TableName() string {
	return "command"
}

type SDCommandInvocationEntity struct {
	ID             uint32 `gorm:"column:id;primaryKey;not null"`
	InvocationTime time.Time
	Payload        string `gorm:"column:payload;not null"`
	UserId         uint32 `gorm:"column:user_id"`
}

func (SDCommandInvocationEntity) TableName() string {
	return "command_invocation"
}

/* ----- R.-B.-A.-C. ----- */

type GraphQLOperationEntity struct {
	ID            uint32 `gorm:"column:id;primaryKey;not null"`
	Identifier    string `gorm:"column:identifier;not null;uniqueIndex"`
	OperationType string `gorm:"column:operation_type;not null;check:operation_type IN ('query', 'mutation', 'subscription')"`
}

func (GraphQLOperationEntity) TableName() string {
	return "graphql_operations"
}

type RoleEntity struct {
	ID          uint32             `gorm:"column:id;primaryKey;not null"`
	Label       string             `gorm:"column:label;not null;uniqueIndex"`
	Permissions []PermissionEntity `gorm:"many2many:roles_permissions_mapping;joinForeignKey:RoleID;joinReferences:PermissionID"`
}

func (RoleEntity) TableName() string {
	return "roles"
}

type RolesPermissionsMappingEntity struct {
	RoleID       uint32 `gorm:"column:role_id;primaryKey;not null;index;constraint:OnDelete:CASCADE"`
	PermissionID uint32 `gorm:"column:permission_id;primaryKey;not null;constraint:OnDelete:CASCADE"`
}

func (RolesPermissionsMappingEntity) TableName() string {
	return "roles_permissions_mapping"
}

type PermissionEntity struct {
	ID                            uint32                               `gorm:"column:id;primaryKey;not null"`
	Label                         string                               `gorm:"column:label;not null;uniqueIndex"`
	Roles                         []RoleEntity                         `gorm:"many2many:roles_permissions_mapping;joinForeignKey:PermissionID;joinReferences:RoleID"`
	OperationTypeAccessPermission *OperationTypeAccessPermissionEntity `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
	SingleOperationPermission     *SingleOperationPermissionEntity     `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

func (PermissionEntity) TableName() string {
	return "permissions"
}

type OperationTypeAccessPermissionEntity struct {
	PermissionID  uint32 `gorm:"column:permission_id;primaryKey;not null"`
	OperationType string `gorm:"column:operation_type;not null;check:operation_type IN ('query', 'mutation', 'subscription')"`
}

func (OperationTypeAccessPermissionEntity) TableName() string {
	return "operation_type_access_permissions"
}

type SingleOperationPermissionEntity struct {
	PermissionID       uint32                 `gorm:"column:permission_id;primaryKey;not null"`
	GraphQLOperationID uint32                 `gorm:"column:graphql_operation_id;primaryKey;not null;uniqueIndex:idx_op_effect,priority:1"`
	GraphQLOperation   GraphQLOperationEntity `gorm:"foreignKey:GraphQLOperationID;references:ID;constraint:OnDelete:CASCADE"`
	Effect             string                 `gorm:"column:effect;not null;check:effect IN ('allow', 'deny');uniqueIndex:idx_op_effect,priority:2"`
}
