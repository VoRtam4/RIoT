/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from "@graphql-typed-document-node/core";
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  Date: { input: any; output: any };
  JSON: { input: any; output: any };
};

export type ApiKey = {
  __typename?: "APIKey";
  expiresAt?: Maybe<Scalars["Date"]["output"]>;
  ipRestrictions: Array<Scalars["String"]["output"]>;
  label: Scalars["String"]["output"];
  lastUsedAt?: Maybe<Scalars["Date"]["output"]>;
  permissions: Array<Scalars["String"]["output"]>;
  rateLimit?: Maybe<Scalars["ID"]["output"]>;
  revoked: Scalars["Boolean"]["output"];
  uid: Scalars["String"]["output"];
};

export type ApiKeyInput = {
  expiresAt?: InputMaybe<Scalars["Date"]["input"]>;
  ipRestrictions?: InputMaybe<Array<Scalars["String"]["input"]>>;
  label?: InputMaybe<Scalars["String"]["input"]>;
  permissions?: InputMaybe<Array<Scalars["String"]["input"]>>;
  rateLimit?: InputMaybe<Scalars["ID"]["input"]>;
  revoked?: InputMaybe<Scalars["Boolean"]["input"]>;
};

export type ApiKeyRestrictionsInput = {
  expiresAt?: InputMaybe<Scalars["Date"]["input"]>;
  ipRestrictions?: InputMaybe<Array<Scalars["String"]["input"]>>;
  rateLimit?: InputMaybe<Scalars["ID"]["input"]>;
  revoked?: InputMaybe<Scalars["Boolean"]["input"]>;
};

export type AssignRoleInput = {
  roleUID: Scalars["String"]["input"];
  userUID: Scalars["String"]["input"];
};

export type AtomKpiNode = {
  comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
  comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
  id: Scalars["ID"]["output"];
  nodeType: KpiNodeType;
  parentNodeID?: Maybe<Scalars["ID"]["output"]>;
  referenceMode: KpiReferenceMode;
  sdParameterSpecification: Scalars["String"]["output"];
};

export type BooleanEqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "BooleanEQAtomKPINode";
    booleanReferenceValue: Scalars["Boolean"]["output"];
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type BooleanExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "BooleanExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type BooleanNeqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "BooleanNEQAtomKPINode";
    booleanReferenceValue: Scalars["Boolean"]["output"];
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type BooleanNotExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "BooleanNotExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type ExportStatus =
  | "cancelled"
  | "done"
  | "expired"
  | "failed"
  | "pending"
  | "processing";

export type FilterNodeInput = {
  nodes?: InputMaybe<Array<FilterNodeInput>>;
  operator?: InputMaybe<LogicalOperator>;
  rule?: InputMaybe<FilterRuleInput>;
  type: FilterNodeType;
};

export type FilterNodeType = "logical" | "rule";

export type FilterOperator =
  | "contains"
  | "eq"
  | "in"
  | "neq"
  | "prefix"
  | "regex"
  | "suffix";

export type FilterRuleInput = {
  operator: FilterOperator;
  tag: Scalars["String"]["input"];
  value: Scalars["String"]["input"];
};

export type InputData = {
  data: Scalars["JSON"]["input"];
  deviceId: Scalars["String"]["input"];
  deviceType?: InputMaybe<Scalars["String"]["input"]>;
  time: Scalars["Date"]["input"];
};

export type KpiDefinition = {
  __typename?: "KPIDefinition";
  label: Scalars["String"]["output"];
  nodes: Array<KpiNode>;
  sdInstanceMode: SdInstanceMode;
  sdTypeUID: Scalars["String"]["output"];
  selectedSDInstanceUIDs: Array<Scalars["String"]["output"]>;
  uid?: Maybe<Scalars["String"]["output"]>;
  userIdentifier: Scalars["String"]["output"];
};

export type KpiDefinitionInput = {
  label: Scalars["String"]["input"];
  nodes: Array<KpiNodeInput>;
  sdInstanceMode: SdInstanceMode;
  sdTypeUID: Scalars["String"]["input"];
  selectedSDInstanceUIDs: Array<Scalars["String"]["input"]>;
  uid?: InputMaybe<Scalars["String"]["input"]>;
  userIdentifier: Scalars["String"]["input"];
};

export type KpiFulfillmentCheckResult = {
  __typename?: "KPIFulfillmentCheckResult";
  eventTime: Scalars["String"]["output"];
  fulfilled: Scalars["Boolean"]["output"];
  kpiDefinitionUID: Scalars["String"]["output"];
  sdInstanceUID: Scalars["String"]["output"];
  sdTypeUID: Scalars["String"]["output"];
};

export type KpiFulfillmentCheckResultRequest = {
  kpiDefinitionUID: Scalars["String"]["input"];
  sdInstanceUID: Scalars["String"]["input"];
};

export type KpiFulfillmentCheckedFilter = {
  kpiDefinitionUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
};

export type KpiNode = {
  id: Scalars["ID"]["output"];
  nodeType: KpiNodeType;
  parentNodeID?: Maybe<Scalars["ID"]["output"]>;
};

export type KpiNodeInput = {
  booleanReferenceValue?: InputMaybe<Scalars["Boolean"]["input"]>;
  comparedRecordOffset?: InputMaybe<Scalars["Int"]["input"]>;
  comparedSDParameterSpecification?: InputMaybe<Scalars["String"]["input"]>;
  id: Scalars["ID"]["input"];
  logicalOperationType?: InputMaybe<LogicalOperationType>;
  numericReferenceValue?: InputMaybe<Scalars["Float"]["input"]>;
  parentNodeID?: InputMaybe<Scalars["ID"]["input"]>;
  referenceMode?: InputMaybe<KpiReferenceMode>;
  sdParameterSpecification?: InputMaybe<Scalars["String"]["input"]>;
  stringReferenceValue?: InputMaybe<Scalars["String"]["input"]>;
  type: KpiNodeType;
};

export type KpiNodeType =
  | "BooleanEQAtom"
  | "BooleanExistsAtom"
  | "BooleanNEQAtom"
  | "BooleanNotExistsAtom"
  | "LogicalOperation"
  | "NumericEQAtom"
  | "NumericExistsAtom"
  | "NumericGEQAtom"
  | "NumericGTAtom"
  | "NumericLEQAtom"
  | "NumericLTAtom"
  | "NumericNEQAtom"
  | "NumericNotExistsAtom"
  | "StringEQAtom"
  | "StringExistsAtom"
  | "StringNEQAtom"
  | "StringNotExistsAtom";

export type KpiReferenceMode = "literal" | "parameter";

export type LogicalOperationKpiNode = KpiNode & {
  __typename?: "LogicalOperationKPINode";
  id: Scalars["ID"]["output"];
  nodeType: KpiNodeType;
  parentNodeID?: Maybe<Scalars["ID"]["output"]>;
  type: LogicalOperationType;
};

export type LogicalOperationType = "and" | "nor" | "not" | "or";

export type LogicalOperator = "and" | "not" | "or";

export type Mutation = {
  __typename?: "Mutation";
  assignRoleToUser: Scalars["Boolean"]["output"];
  cancelTimeSeriesExport: TimeSeriesExport;
  cloneRole: Role;
  createAPIKey: Scalars["String"]["output"];
  createKPIDefinition: KpiDefinition;
  createRole: Role;
  createSDInstanceGroup: SdInstanceGroup;
  createSDType: SdType;
  deleteAPIKey: Scalars["Boolean"]["output"];
  deleteKPIDefinition: Scalars["Boolean"]["output"];
  deleteRole: Scalars["Boolean"]["output"];
  deleteSDInstanceGroup: Scalars["Boolean"]["output"];
  deleteSDType: Scalars["Boolean"]["output"];
  deleteUserConfig: Scalars["Boolean"]["output"];
  disableUser: User;
  enableUser: User;
  revokeAPIKey: ApiKey;
  revokeAllSessionsForUser: Scalars["Boolean"]["output"];
  revokeOwnOtherSessions: Scalars["Boolean"]["output"];
  revokeSession: Scalars["Boolean"]["output"];
  revokeUserSessions: Scalars["Boolean"]["output"];
  rotateAPIKey: Scalars["String"]["output"];
  startTimeSeriesExport: TimeSeriesExport;
  startTimeSeriesExportAggregateKPI: TimeSeriesExport;
  statisticsMutate: Scalars["Boolean"]["output"];
  updateAPIKey: Scalars["Boolean"]["output"];
  updateAPIKeyPermissions: ApiKey;
  updateAPIKeyRestrictions: ApiKey;
  updateKPIDefinition: KpiDefinition;
  updatePermissionLabel: Permission;
  updateRole: Role;
  updateSDInstance: SdInstance;
  updateSDInstanceGroup: SdInstanceGroup;
  updateUser: User;
  updateUserConfig: UserConfig;
};

export type MutationAssignRoleToUserArgs = {
  input: AssignRoleInput;
};

export type MutationCancelTimeSeriesExportArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationCloneRoleArgs = {
  label: Scalars["String"]["input"];
  uid: Scalars["String"]["input"];
};

export type MutationCreateApiKeyArgs = {
  input: ApiKeyInput;
};

export type MutationCreateKpiDefinitionArgs = {
  input: KpiDefinitionInput;
};

export type MutationCreateRoleArgs = {
  input: RoleInput;
};

export type MutationCreateSdInstanceGroupArgs = {
  input: SdInstanceGroupInput;
};

export type MutationCreateSdTypeArgs = {
  input: SdTypeInput;
};

export type MutationDeleteApiKeyArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationDeleteKpiDefinitionArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationDeleteRoleArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationDeleteSdInstanceGroupArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationDeleteSdTypeArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationDisableUserArgs = {
  reason?: InputMaybe<Scalars["String"]["input"]>;
  uid: Scalars["String"]["input"];
};

export type MutationEnableUserArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationRevokeApiKeyArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationRevokeAllSessionsForUserArgs = {
  userUID: Scalars["String"]["input"];
};

export type MutationRevokeSessionArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationRevokeUserSessionsArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationRotateApiKeyArgs = {
  uid: Scalars["String"]["input"];
};

export type MutationStartTimeSeriesExportArgs = {
  input: TimeSeriesReadInput;
};

export type MutationStartTimeSeriesExportAggregateKpiArgs = {
  input: TimeSeriesReadAggregateKpiInput;
};

export type MutationStatisticsMutateArgs = {
  inputData: InputData;
};

export type MutationUpdateApiKeyArgs = {
  input: ApiKeyInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateApiKeyPermissionsArgs = {
  permissionUIDs: Array<Scalars["String"]["input"]>;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateApiKeyRestrictionsArgs = {
  input: ApiKeyRestrictionsInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateKpiDefinitionArgs = {
  input: KpiDefinitionInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdatePermissionLabelArgs = {
  label: Scalars["String"]["input"];
  uid: Scalars["String"]["input"];
};

export type MutationUpdateRoleArgs = {
  input: RoleInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateSdInstanceArgs = {
  input: SdInstanceUpdateInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateSdInstanceGroupArgs = {
  input: SdInstanceGroupInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateUserArgs = {
  input: UserUpdateInput;
  uid: Scalars["String"]["input"];
};

export type MutationUpdateUserConfigArgs = {
  input: UserConfigInput;
};

export type NumericEqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericGeqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericGEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericGtAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericGTAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericLeqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericLEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericLtAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericLTAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericNeqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericNEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    numericReferenceValue: Scalars["Float"]["output"];
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type NumericNotExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "NumericNotExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type OutputData = {
  __typename?: "OutputData";
  data: Scalars["JSON"]["output"];
  deviceId: Scalars["String"]["output"];
  deviceType?: Maybe<Scalars["String"]["output"]>;
  time: Scalars["Date"]["output"];
};

export type ParameterRole = "field" | "meta" | "tag" | "time";

export type Permission = {
  __typename?: "Permission";
  label: Scalars["String"]["output"];
  uid: Scalars["String"]["output"];
};

export type Query = {
  __typename?: "Query";
  apiKey: ApiKey;
  apiKeys: Array<ApiKey>;
  apiKeysByUser: Array<ApiKey>;
  kpiDefinition: KpiDefinition;
  kpiDefinitions: Array<KpiDefinition>;
  kpiDefinitionsBySdInstance: Array<KpiDefinition>;
  kpiDefinitionsBySdType: Array<KpiDefinition>;
  kpiResult: KpiFulfillmentCheckResult;
  kpiResults: Array<KpiFulfillmentCheckResult>;
  kpiResultsByKPI: Array<KpiFulfillmentCheckResult>;
  me: User;
  permissions: Array<Permission>;
  rawDataPoint: RawDataPoint;
  rawDataPointsBySDType: Array<RawDataPoint>;
  role?: Maybe<Role>;
  roles: Array<Role>;
  sdInstance: SdInstance;
  sdInstanceGroup: SdInstanceGroup;
  sdInstanceGroups: Array<SdInstanceGroup>;
  sdInstances: Array<SdInstance>;
  sdInstancesByKpiDefinition: Array<SdInstance>;
  sdInstancesByType: Array<SdInstance>;
  sdType: SdType;
  sdTypes: Array<SdType>;
  sessions: Array<UserSession>;
  sessionsByUser: Array<UserSession>;
  statisticsQuerySensorsWithFields: Array<OutputData>;
  statisticsQuerySimpleSensors: Array<OutputData>;
  timeSeriesDistinctTagValues: TimeSeriesDistinctTagValuesResponse;
  timeSeriesExport: TimeSeriesExport;
  timeSeriesRead: TimeSeriesReadResponse;
  timeSeriesReadAggregateKPI: TimeSeriesReadResponse;
  user: User;
  userConfig: UserConfig;
  userRole?: Maybe<Role>;
  users: Array<User>;
};

export type QueryApiKeyArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryApiKeysByUserArgs = {
  userUID: Scalars["String"]["input"];
};

export type QueryKpiDefinitionArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryKpiDefinitionsBySdInstanceArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryKpiDefinitionsBySdTypeArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryKpiResultArgs = {
  request: KpiFulfillmentCheckResultRequest;
};

export type QueryKpiResultsByKpiArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryRawDataPointArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryRawDataPointsBySdTypeArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryRoleArgs = {
  uid?: InputMaybe<Scalars["String"]["input"]>;
};

export type QuerySdInstanceArgs = {
  uid: Scalars["String"]["input"];
};

export type QuerySdInstanceGroupArgs = {
  uid: Scalars["String"]["input"];
};

export type QuerySdInstancesByKpiDefinitionArgs = {
  uid: Scalars["String"]["input"];
};

export type QuerySdInstancesByTypeArgs = {
  uid: Scalars["String"]["input"];
};

export type QuerySdTypeArgs = {
  uid: Scalars["String"]["input"];
};

export type QuerySessionsByUserArgs = {
  userUID: Scalars["String"]["input"];
};

export type QueryStatisticsQuerySensorsWithFieldsArgs = {
  request?: InputMaybe<StatisticsInput>;
  sensors: SensorsWithFields;
};

export type QueryStatisticsQuerySimpleSensorsArgs = {
  request?: InputMaybe<StatisticsInput>;
  sensors: SimpleSensors;
};

export type QueryTimeSeriesDistinctTagValuesArgs = {
  request: TimeSeriesDistinctTagValuesInput;
};

export type QueryTimeSeriesExportArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryTimeSeriesReadArgs = {
  request: TimeSeriesReadInput;
};

export type QueryTimeSeriesReadAggregateKpiArgs = {
  request: TimeSeriesReadAggregateKpiInput;
};

export type QueryUserArgs = {
  uid: Scalars["String"]["input"];
};

export type QueryUserRoleArgs = {
  uid: Scalars["String"]["input"];
};

export type RawDataPoint = {
  __typename?: "RawDataPoint";
  eventTime: Scalars["String"]["output"];
  payload: Scalars["JSON"]["output"];
  sdInstanceUID: Scalars["String"]["output"];
  sdTypeUID: Scalars["String"]["output"];
};

export type RawDataPointArrivedFilter = {
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
};

export type Role = {
  __typename?: "Role";
  label: Scalars["String"]["output"];
  permissions?: Maybe<Array<Permission>>;
  system: Scalars["Boolean"]["output"];
  uid: Scalars["String"]["output"];
};

export type RoleInput = {
  label: Scalars["String"]["input"];
  permissionUIDs: Array<Scalars["String"]["input"]>;
};

export type SdInstance = {
  __typename?: "SDInstance";
  confirmedByUser: Scalars["Boolean"]["output"];
  label: Scalars["String"]["output"];
  type: SdType;
  uid: Scalars["String"]["output"];
  userIdentifier: Scalars["String"]["output"];
};

export type SdInstanceGroup = {
  __typename?: "SDInstanceGroup";
  label: Scalars["String"]["output"];
  sdInstanceUIDs: Array<Scalars["String"]["output"]>;
  uid: Scalars["String"]["output"];
  userIdentifier: Scalars["String"]["output"];
};

export type SdInstanceGroupInput = {
  label: Scalars["String"]["input"];
  sdInstanceUIDs: Array<Scalars["String"]["input"]>;
  uid: Scalars["String"]["input"];
  userIdentifier: Scalars["String"]["input"];
};

export type SdInstanceMode = "all" | "selected";

export type SdInstanceRegisteredFilter = {
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
};

export type SdInstanceUpdateInput = {
  confirmedByUser?: InputMaybe<Scalars["Boolean"]["input"]>;
  label?: InputMaybe<Scalars["String"]["input"]>;
  userIdentifier?: InputMaybe<Scalars["String"]["input"]>;
};

export type SdParameter = {
  __typename?: "SDParameter";
  denotation: Scalars["String"]["output"];
  label: Scalars["String"]["output"];
  role: SdParameterRole;
  type: SdParameterType;
};

export type SdParameterInput = {
  denotation: Scalars["String"]["input"];
  label: Scalars["String"]["input"];
  role: SdParameterRole;
  type: SdParameterType;
};

export type SdParameterRole = "field" | "tag";

export type SdParameterType = "boolean" | "number" | "string";

export type SdType = {
  __typename?: "SDType";
  label: Scalars["String"]["output"];
  parameters: Array<SdParameter>;
  uid: Scalars["String"]["output"];
};

export type SdTypeInput = {
  label: Scalars["String"]["input"];
  parameters: Array<SdParameterInput>;
  uid: Scalars["String"]["input"];
};

export type SensorField = {
  key: Scalars["String"]["input"];
  values: Array<Scalars["String"]["input"]>;
};

export type SensorsWithFields = {
  sensors: Array<SensorField>;
};

export type SimpleSensors = {
  sensors: Array<Scalars["String"]["input"]>;
};

/** Data used for querying the selected bucket */
export type StatisticsInput = {
  /**
   * Amount of minutes to aggregate by
   * For example if the queried range has 1 hour and aggregateMinutes is set to 10 the aggregation will result in 6 points
   */
  aggregateSeconds?: InputMaybe<Scalars["Int"]["input"]>;
  /** Start of the querying window */
  from?: InputMaybe<Scalars["Date"]["input"]>;
  /** Aggregation operator to use, if needed */
  operation?: InputMaybe<StatisticsOperation>;
  /**
   * Timezone override default UTC.
   * For more details why and how this affects queries see: https://www.influxdata.com/blog/time-zones-in-flux/.
   * In most cases you can ignore this and some edge aggregations can be influenced.
   * If you need a precise result or the aggregation uses high amount of minutes provide the target time zone.
   */
  timezone?: InputMaybe<Scalars["String"]["input"]>;
  /** End of the querying window */
  to?: InputMaybe<Scalars["Date"]["input"]>;
};

export type StatisticsOperation =
  | "count"
  | "first"
  | "integral"
  | "last"
  | "max"
  | "mean"
  | "median"
  | "min"
  | "mode"
  | "none"
  | "quantile"
  | "reduce"
  | "skew"
  | "spread"
  | "stddev"
  | "sum"
  | "timeweightedavg";

export type StringEqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "StringEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
    stringReferenceValue: Scalars["String"]["output"];
  };

export type StringExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "StringExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type StringNeqAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "StringNEQAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
    stringReferenceValue: Scalars["String"]["output"];
  };

export type StringNotExistsAtomKpiNode = AtomKpiNode &
  KpiNode & {
    __typename?: "StringNotExistsAtomKPINode";
    comparedRecordOffset?: Maybe<Scalars["Int"]["output"]>;
    comparedSDParameterSpecification?: Maybe<Scalars["String"]["output"]>;
    id: Scalars["ID"]["output"];
    nodeType: KpiNodeType;
    parentNodeID?: Maybe<Scalars["ID"]["output"]>;
    referenceMode: KpiReferenceMode;
    sdParameterSpecification: Scalars["String"]["output"];
  };

export type Subscription = {
  __typename?: "Subscription";
  onKPIFulfillmentChecked: Array<KpiFulfillmentCheckResult>;
  onRawDataPointArrived: Array<RawDataPoint>;
  onSDInstanceRegistered: SdInstance;
  onTimeSeriesExportUpdated: TimeSeriesExport;
};

export type SubscriptionOnKpiFulfillmentCheckedArgs = {
  filter?: InputMaybe<KpiFulfillmentCheckedFilter>;
};

export type SubscriptionOnRawDataPointArrivedArgs = {
  filter?: InputMaybe<RawDataPointArrivedFilter>;
};

export type SubscriptionOnSdInstanceRegisteredArgs = {
  filter?: InputMaybe<SdInstanceRegisteredFilter>;
};

export type SubscriptionOnTimeSeriesExportUpdatedArgs = {
  filter: TimeSeriesExportFilter;
};

export type TimeSeriesCursor = {
  __typename?: "TimeSeriesCursor";
  kpiDefinitionUID?: Maybe<Scalars["String"]["output"]>;
  sdInstanceUID: Scalars["String"]["output"];
  time: Scalars["Date"]["output"];
};

export type TimeSeriesCursorInput = {
  kpiDefinitionUID?: InputMaybe<Scalars["String"]["input"]>;
  sdInstanceUID: Scalars["String"]["input"];
  time: Scalars["Date"]["input"];
};

export type TimeSeriesDataPoint = {
  __typename?: "TimeSeriesDataPoint";
  data: Scalars["JSON"]["output"];
  tags: Scalars["JSON"]["output"];
  time: Scalars["Date"]["output"];
};

export type TimeSeriesDistinctTagValuesInput = {
  filters?: InputMaybe<FilterNodeInput>;
  from?: InputMaybe<Scalars["Date"]["input"]>;
  kpiDefinitionUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUID?: InputMaybe<Scalars["String"]["input"]>;
  tag: Scalars["String"]["input"];
  to?: InputMaybe<Scalars["Date"]["input"]>;
  type: TimeSeriesType;
};

export type TimeSeriesDistinctTagValuesResponse = {
  __typename?: "TimeSeriesDistinctTagValuesResponse";
  error?: Maybe<Scalars["String"]["output"]>;
  values: Array<Scalars["String"]["output"]>;
};

export type TimeSeriesExport = {
  __typename?: "TimeSeriesExport";
  createdAt: Scalars["Date"]["output"];
  downloadUrl?: Maybe<Scalars["String"]["output"]>;
  error?: Maybe<Scalars["String"]["output"]>;
  expiresAt?: Maybe<Scalars["Date"]["output"]>;
  status: ExportStatus;
  uid: Scalars["String"]["output"];
};

export type TimeSeriesExportFilter = {
  uids?: InputMaybe<Array<Scalars["String"]["input"]>>;
};

export type TimeSeriesParameter = {
  __typename?: "TimeSeriesParameter";
  denotation: Scalars["String"]["output"];
  label: Scalars["String"]["output"];
  role: ParameterRole;
};

export type TimeSeriesReadAggregateKpiInput = {
  aggregateSeconds: Scalars["Int"]["input"];
  batch?: InputMaybe<Scalars["Int"]["input"]>;
  cursor?: InputMaybe<TimeSeriesCursorInput>;
  filters?: InputMaybe<FilterNodeInput>;
  from?: InputMaybe<Scalars["Date"]["input"]>;
  kpiDefinitionUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUID?: InputMaybe<Scalars["String"]["input"]>;
  to?: InputMaybe<Scalars["Date"]["input"]>;
};

export type TimeSeriesReadInput = {
  batch?: InputMaybe<Scalars["Int"]["input"]>;
  cursor?: InputMaybe<TimeSeriesCursorInput>;
  filters?: InputMaybe<FilterNodeInput>;
  from?: InputMaybe<Scalars["Date"]["input"]>;
  kpiDefinitionUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  sdInstanceUIDs?: InputMaybe<Array<Scalars["String"]["input"]>>;
  sdTypeUID?: InputMaybe<Scalars["String"]["input"]>;
  sortDesc?: InputMaybe<Scalars["Boolean"]["input"]>;
  to?: InputMaybe<Scalars["Date"]["input"]>;
  type: TimeSeriesType;
};

export type TimeSeriesReadResponse = {
  __typename?: "TimeSeriesReadResponse";
  base: Scalars["JSON"]["output"];
  data?: Maybe<Array<TimeSeriesDataPoint>>;
  error?: Maybe<Scalars["String"]["output"]>;
  hasMoreBatches: Scalars["Boolean"]["output"];
  hasMoreData: Scalars["Boolean"]["output"];
  nextCursor?: Maybe<TimeSeriesCursor>;
  parameters?: Maybe<Array<TimeSeriesParameter>>;
};

export type TimeSeriesType = "kpi" | "raw";

export type User = {
  __typename?: "User";
  disabled: Scalars["Boolean"]["output"];
  disabledAt?: Maybe<Scalars["Date"]["output"]>;
  disabledReason?: Maybe<Scalars["String"]["output"]>;
  email: Scalars["String"]["output"];
  lastLoginAt?: Maybe<Scalars["Date"]["output"]>;
  name?: Maybe<Scalars["String"]["output"]>;
  oauth2Provider?: Maybe<Scalars["String"]["output"]>;
  profileImageURL?: Maybe<Scalars["String"]["output"]>;
  role?: Maybe<Role>;
  uid: Scalars["String"]["output"];
  username: Scalars["String"]["output"];
};

export type UserConfig = {
  __typename?: "UserConfig";
  config: Scalars["JSON"]["output"];
  userUID: Scalars["String"]["output"];
};

export type UserConfigInput = {
  config: Scalars["JSON"]["input"];
};

export type UserSession = {
  __typename?: "UserSession";
  createdAt: Scalars["Date"]["output"];
  expiresAt: Scalars["Date"]["output"];
  ipAddress: Scalars["String"]["output"];
  revoked: Scalars["Boolean"]["output"];
  uid: Scalars["String"]["output"];
  updatedAt: Scalars["Date"]["output"];
  userAgent: Scalars["String"]["output"];
  userUID: Scalars["String"]["output"];
};

export type UserUpdateInput = {
  name?: InputMaybe<Scalars["String"]["input"]>;
  profileImageURL?: InputMaybe<Scalars["String"]["input"]>;
  username?: InputMaybe<Scalars["String"]["input"]>;
};

export type ApiKeysQueryVariables = Exact<{ [key: string]: never }>;

export type ApiKeysQuery = {
  __typename?: "Query";
  apiKeys: Array<{
    __typename?: "APIKey";
    uid: string;
    label: string;
    revoked: boolean;
    expiresAt?: any | null;
    permissions: Array<string>;
    ipRestrictions: Array<string>;
  }>;
};

export type CreateApiKeyMutationVariables = Exact<{
  input: ApiKeyInput;
}>;

export type CreateApiKeyMutation = {
  __typename?: "Mutation";
  createAPIKey: string;
};

export type UpdateApiKeyMutationVariables = Exact<{
  uid: Scalars["String"]["input"];
  input: ApiKeyInput;
}>;

export type UpdateApiKeyMutation = {
  __typename?: "Mutation";
  updateAPIKey: boolean;
};

export type DeleteApiKeyMutationVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type DeleteApiKeyMutation = {
  __typename?: "Mutation";
  deleteAPIKey: boolean;
};

export type UserConfigQueryVariables = Exact<{ [key: string]: never }>;

export type UserConfigQuery = {
  __typename?: "Query";
  userConfig: { __typename?: "UserConfig"; userUID: string; config: any };
};

export type UpdateUserConfigMutationVariables = Exact<{
  input: UserConfigInput;
}>;

export type UpdateUserConfigMutation = {
  __typename?: "Mutation";
  updateUserConfig: { __typename?: "UserConfig"; userUID: string; config: any };
};

export type RoleQueryVariables = Exact<{ [key: string]: never }>;

export type RoleQuery = {
  __typename?: "Query";
  role?: {
    __typename?: "Role";
    permissions?: Array<{
      __typename?: "Permission";
      uid: string;
      label: string;
    }> | null;
  } | null;
};

export type KpiDefinitionsQueryVariables = Exact<{ [key: string]: never }>;

export type KpiDefinitionsQuery = {
  __typename?: "Query";
  kpiDefinitions: Array<{
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
    sdTypeUID: string;
    sdInstanceMode: SdInstanceMode;
  }>;
};

export type KpiDefinitionsBySdTypeQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type KpiDefinitionsBySdTypeQuery = {
  __typename?: "Query";
  kpiDefinitionsBySdType: Array<{
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
    sdTypeUID: string;
    sdInstanceMode: SdInstanceMode;
  }>;
};

export type CreateKpiDefinitionMutationVariables = Exact<{
  input: KpiDefinitionInput;
}>;

export type CreateKpiDefinitionMutation = {
  __typename?: "Mutation";
  createKPIDefinition: {
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
  };
};

export type UpdateKpiDefinitionMutationVariables = Exact<{
  uid: Scalars["String"]["input"];
  input: KpiDefinitionInput;
}>;

export type UpdateKpiDefinitionMutation = {
  __typename?: "Mutation";
  updateKPIDefinition: {
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
  };
};

export type KpiResultQueryVariables = Exact<{
  request: KpiFulfillmentCheckResultRequest;
}>;

export type KpiResultQuery = {
  __typename?: "Query";
  kpiResult: {
    __typename?: "KPIFulfillmentCheckResult";
    eventTime: string;
    fulfilled: boolean;
  };
};

export type OnKpiFulfillmentCheckedSubscriptionVariables = Exact<{
  filter?: InputMaybe<KpiFulfillmentCheckedFilter>;
}>;

export type OnKpiFulfillmentCheckedSubscription = {
  __typename?: "Subscription";
  onKPIFulfillmentChecked: Array<{
    __typename?: "KPIFulfillmentCheckResult";
    sdTypeUID: string;
    sdInstanceUID: string;
    kpiDefinitionUID: string;
    eventTime: string;
    fulfilled: boolean;
  }>;
};

export type KpiDefinitionsBySdInstanceQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type KpiDefinitionsBySdInstanceQuery = {
  __typename?: "Query";
  kpiDefinitionsBySdInstance: Array<{
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
    sdTypeUID: string;
    sdInstanceMode: SdInstanceMode;
  }>;
};

export type DeleteKpiDefinitionMutationVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type DeleteKpiDefinitionMutation = {
  __typename?: "Mutation";
  deleteKPIDefinition: boolean;
};

export type KpiDefinitionQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type KpiDefinitionQuery = {
  __typename?: "Query";
  kpiDefinition: {
    __typename?: "KPIDefinition";
    uid?: string | null;
    label: string;
    sdTypeUID: string;
    userIdentifier: string;
    sdInstanceMode: SdInstanceMode;
    selectedSDInstanceUIDs: Array<string>;
    nodes: Array<
      | {
          __typename?: "BooleanEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          booleanReferenceValue: boolean;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "BooleanExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "BooleanNEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          booleanReferenceValue: boolean;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "BooleanNotExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "LogicalOperationKPINode";
          type: LogicalOperationType;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericGEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericGTAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericLEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericLTAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericNEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          numericReferenceValue: number;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "NumericNotExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "StringEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          stringReferenceValue: string;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "StringExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "StringNEQAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          stringReferenceValue: string;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
      | {
          __typename?: "StringNotExistsAtomKPINode";
          sdParameterSpecification: string;
          referenceMode: KpiReferenceMode;
          comparedSDParameterSpecification?: string | null;
          comparedRecordOffset?: number | null;
          id: string;
          parentNodeID?: string | null;
          nodeType: KpiNodeType;
        }
    >;
  };
};

export type SdInstancesQueryVariables = Exact<{ [key: string]: never }>;

export type SdInstancesQuery = {
  __typename?: "Query";
  sdInstances: Array<{ __typename?: "SDInstance"; uid: string; label: string }>;
};

export type SdInstancesByTypeQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type SdInstancesByTypeQuery = {
  __typename?: "Query";
  sdInstancesByType: Array<{
    __typename?: "SDInstance";
    uid: string;
    label: string;
    type: { __typename?: "SDType"; uid: string };
  }>;
};

export type SdInstancesByKpiDefinitionQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type SdInstancesByKpiDefinitionQuery = {
  __typename?: "Query";
  sdInstancesByKpiDefinition: Array<{
    __typename?: "SDInstance";
    uid: string;
    label: string;
  }>;
};

export type RawDataPointQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type RawDataPointQuery = {
  __typename?: "Query";
  rawDataPoint: {
    __typename?: "RawDataPoint";
    sdInstanceUID: string;
    payload: any;
    eventTime: string;
  };
};

export type OnRawDataPointArrivedSubscriptionVariables = Exact<{
  filter?: InputMaybe<RawDataPointArrivedFilter>;
}>;

export type OnRawDataPointArrivedSubscription = {
  __typename?: "Subscription";
  onRawDataPointArrived: Array<{
    __typename?: "RawDataPoint";
    sdTypeUID: string;
    sdInstanceUID: string;
    payload: any;
    eventTime: string;
  }>;
};

export type SdInstanceQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type SdInstanceQuery = {
  __typename?: "Query";
  sdInstance: {
    __typename?: "SDInstance";
    uid: string;
    label: string;
    confirmedByUser: boolean;
    userIdentifier: string;
    type: {
      __typename?: "SDType";
      uid: string;
      label: string;
      parameters: Array<{
        __typename?: "SDParameter";
        label: string;
        denotation: string;
        type: SdParameterType;
        role: SdParameterRole;
      }>;
    };
  };
};

export type SdTypesQueryVariables = Exact<{ [key: string]: never }>;

export type SdTypesQuery = {
  __typename?: "Query";
  sdTypes: Array<{
    __typename?: "SDType";
    uid: string;
    label: string;
    parameters: Array<{
      __typename?: "SDParameter";
      label: string;
      denotation: string;
      type: SdParameterType;
      role: SdParameterRole;
    }>;
  }>;
};

export type SdTypeByUidQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type SdTypeByUidQuery = {
  __typename?: "Query";
  sdType: {
    __typename?: "SDType";
    uid: string;
    label: string;
    parameters: Array<{
      __typename?: "SDParameter";
      label: string;
      denotation: string;
      type: SdParameterType;
      role: SdParameterRole;
    }>;
  };
};

export type TimeSeriesReadQueryVariables = Exact<{
  request: TimeSeriesReadInput;
}>;

export type TimeSeriesReadQuery = {
  __typename?: "Query";
  timeSeriesRead: {
    __typename?: "TimeSeriesReadResponse";
    base: any;
    hasMoreBatches: boolean;
    hasMoreData: boolean;
    error?: string | null;
    data?: Array<{
      __typename?: "TimeSeriesDataPoint";
      time: any;
      tags: any;
      data: any;
    }> | null;
    parameters?: Array<{
      __typename?: "TimeSeriesParameter";
      denotation: string;
      label: string;
      role: ParameterRole;
    }> | null;
    nextCursor?: {
      __typename?: "TimeSeriesCursor";
      time: any;
      sdInstanceUID: string;
      kpiDefinitionUID?: string | null;
    } | null;
  };
};

export type StartTimeSeriesExportMutationVariables = Exact<{
  input: TimeSeriesReadInput;
}>;

export type StartTimeSeriesExportMutation = {
  __typename?: "Mutation";
  startTimeSeriesExport: {
    __typename?: "TimeSeriesExport";
    uid: string;
    status: ExportStatus;
    downloadUrl?: string | null;
    createdAt: any;
    expiresAt?: any | null;
    error?: string | null;
  };
};

export type CancelTimeSeriesExportMutationVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type CancelTimeSeriesExportMutation = {
  __typename?: "Mutation";
  cancelTimeSeriesExport: {
    __typename?: "TimeSeriesExport";
    uid: string;
    status: ExportStatus;
    downloadUrl?: string | null;
    createdAt: any;
    expiresAt?: any | null;
    error?: string | null;
  };
};

export type TimeSeriesExportQueryVariables = Exact<{
  uid: Scalars["String"]["input"];
}>;

export type TimeSeriesExportQuery = {
  __typename?: "Query";
  timeSeriesExport: {
    __typename?: "TimeSeriesExport";
    uid: string;
    status: ExportStatus;
    downloadUrl?: string | null;
    createdAt: any;
    expiresAt?: any | null;
    error?: string | null;
  };
};

export type OnTimeSeriesExportUpdatedSubscriptionVariables = Exact<{
  filter: TimeSeriesExportFilter;
}>;

export type OnTimeSeriesExportUpdatedSubscription = {
  __typename?: "Subscription";
  onTimeSeriesExportUpdated: {
    __typename?: "TimeSeriesExport";
    uid: string;
    status: ExportStatus;
    downloadUrl?: string | null;
    createdAt: any;
    expiresAt?: any | null;
    error?: string | null;
  };
};

export type TimeSeriesReadAggregateKpiQueryVariables = Exact<{
  request: TimeSeriesReadAggregateKpiInput;
}>;

export type TimeSeriesReadAggregateKpiQuery = {
  __typename?: "Query";
  timeSeriesReadAggregateKPI: {
    __typename?: "TimeSeriesReadResponse";
    base: any;
    hasMoreBatches: boolean;
    hasMoreData: boolean;
    error?: string | null;
    parameters?: Array<{
      __typename?: "TimeSeriesParameter";
      denotation: string;
      label: string;
      role: ParameterRole;
    }> | null;
    data?: Array<{
      __typename?: "TimeSeriesDataPoint";
      time: any;
      tags: any;
      data: any;
    }> | null;
    nextCursor?: {
      __typename?: "TimeSeriesCursor";
      time: any;
      sdInstanceUID: string;
      kpiDefinitionUID?: string | null;
    } | null;
  };
};

export const ApiKeysDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "ApiKeys" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "apiKeys" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                { kind: "Field", name: { kind: "Name", value: "revoked" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
                { kind: "Field", name: { kind: "Name", value: "permissions" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "ipRestrictions" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<ApiKeysQuery, ApiKeysQueryVariables>;
export const CreateApiKeyDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateApiKey" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "APIKeyInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createAPIKey" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateApiKeyMutation,
  CreateApiKeyMutationVariables
>;
export const UpdateApiKeyDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateApiKey" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "APIKeyInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateAPIKey" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateApiKeyMutation,
  UpdateApiKeyMutationVariables
>;
export const DeleteApiKeyDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "DeleteApiKey" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "deleteAPIKey" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DeleteApiKeyMutation,
  DeleteApiKeyMutationVariables
>;
export const UserConfigDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "UserConfig" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "userConfig" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "userUID" } },
                { kind: "Field", name: { kind: "Name", value: "config" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<UserConfigQuery, UserConfigQueryVariables>;
export const UpdateUserConfigDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateUserConfig" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "UserConfigInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateUserConfig" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "userUID" } },
                { kind: "Field", name: { kind: "Name", value: "config" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateUserConfigMutation,
  UpdateUserConfigMutationVariables
>;
export const RoleDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "Role" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "role" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "permissions" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "uid" } },
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<RoleQuery, RoleQueryVariables>;
export const KpiDefinitionsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "KpiDefinitions" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "kpiDefinitions" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceMode" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<KpiDefinitionsQuery, KpiDefinitionsQueryVariables>;
export const KpiDefinitionsBySdTypeDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "KpiDefinitionsBySdType" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "kpiDefinitionsBySdType" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceMode" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  KpiDefinitionsBySdTypeQuery,
  KpiDefinitionsBySdTypeQueryVariables
>;
export const CreateKpiDefinitionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateKpiDefinition" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "KPIDefinitionInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createKPIDefinition" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateKpiDefinitionMutation,
  CreateKpiDefinitionMutationVariables
>;
export const UpdateKpiDefinitionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateKpiDefinition" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "KPIDefinitionInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateKPIDefinition" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateKpiDefinitionMutation,
  UpdateKpiDefinitionMutationVariables
>;
export const KpiResultDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "KpiResult" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "request" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "KPIFulfillmentCheckResultRequest" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "kpiResult" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "request" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "request" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "eventTime" } },
                { kind: "Field", name: { kind: "Name", value: "fulfilled" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<KpiResultQuery, KpiResultQueryVariables>;
export const OnKpiFulfillmentCheckedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "OnKPIFulfillmentChecked" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "filter" },
          },
          type: {
            kind: "NamedType",
            name: { kind: "Name", value: "KPIFulfillmentCheckedFilter" },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "onKPIFulfillmentChecked" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "filter" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "filter" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceUID" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "kpiDefinitionUID" },
                },
                { kind: "Field", name: { kind: "Name", value: "eventTime" } },
                { kind: "Field", name: { kind: "Name", value: "fulfilled" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  OnKpiFulfillmentCheckedSubscription,
  OnKpiFulfillmentCheckedSubscriptionVariables
>;
export const KpiDefinitionsBySdInstanceDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "KpiDefinitionsBySdInstance" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "kpiDefinitionsBySdInstance" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceMode" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  KpiDefinitionsBySdInstanceQuery,
  KpiDefinitionsBySdInstanceQueryVariables
>;
export const DeleteKpiDefinitionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "DeleteKpiDefinition" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "deleteKPIDefinition" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DeleteKpiDefinitionMutation,
  DeleteKpiDefinitionMutationVariables
>;
export const KpiDefinitionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "KpiDefinition" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "kpiDefinition" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "userIdentifier" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceMode" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "selectedSDInstanceUIDs" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "nodes" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "parentNodeID" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "nodeType" },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "LogicalOperationKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "type" },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "AtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "referenceMode" },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "comparedSDParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "comparedRecordOffset",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "StringEQAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "stringReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "StringNEQAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "stringReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "StringExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "StringNotExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "BooleanEQAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "booleanReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "BooleanNEQAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "booleanReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "BooleanExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "BooleanNotExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "NumericEQAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "NumericNEQAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "NumericGTAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "NumericGEQAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "NumericLTAtomKPINode" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "NumericLEQAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "numericReferenceValue",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "NumericExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "NumericNotExistsAtomKPINode",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "sdParameterSpecification",
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<KpiDefinitionQuery, KpiDefinitionQueryVariables>;
export const SdInstancesDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdInstances" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdInstances" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<SdInstancesQuery, SdInstancesQueryVariables>;
export const SdInstancesByTypeDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdInstancesByType" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdInstancesByType" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "type" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "uid" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SdInstancesByTypeQuery,
  SdInstancesByTypeQueryVariables
>;
export const SdInstancesByKpiDefinitionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdInstancesByKpiDefinition" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdInstancesByKpiDefinition" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SdInstancesByKpiDefinitionQuery,
  SdInstancesByKpiDefinitionQueryVariables
>;
export const RawDataPointDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "RawDataPoint" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "rawDataPoint" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceUID" },
                },
                { kind: "Field", name: { kind: "Name", value: "payload" } },
                { kind: "Field", name: { kind: "Name", value: "eventTime" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<RawDataPointQuery, RawDataPointQueryVariables>;
export const OnRawDataPointArrivedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "OnRawDataPointArrived" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "filter" },
          },
          type: {
            kind: "NamedType",
            name: { kind: "Name", value: "RawDataPointArrivedFilter" },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "onRawDataPointArrived" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "filter" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "filter" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "sdTypeUID" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "sdInstanceUID" },
                },
                { kind: "Field", name: { kind: "Name", value: "payload" } },
                { kind: "Field", name: { kind: "Name", value: "eventTime" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  OnRawDataPointArrivedSubscription,
  OnRawDataPointArrivedSubscriptionVariables
>;
export const SdInstanceDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdInstance" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdInstance" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "confirmedByUser" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "userIdentifier" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "type" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "uid" } },
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "parameters" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "label" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "denotation" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "type" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "role" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<SdInstanceQuery, SdInstanceQueryVariables>;
export const SdTypesDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdTypes" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdTypes" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "parameters" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "denotation" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "type" } },
                      { kind: "Field", name: { kind: "Name", value: "role" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<SdTypesQuery, SdTypesQueryVariables>;
export const SdTypeByUidDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SdTypeByUID" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sdType" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "label" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "parameters" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "denotation" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "type" } },
                      { kind: "Field", name: { kind: "Name", value: "role" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<SdTypeByUidQuery, SdTypeByUidQueryVariables>;
export const TimeSeriesReadDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "TimeSeriesRead" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "request" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "TimeSeriesReadInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timeSeriesRead" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "request" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "request" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "base" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "data" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "time" } },
                      { kind: "Field", name: { kind: "Name", value: "tags" } },
                      { kind: "Field", name: { kind: "Name", value: "data" } },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "parameters" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "denotation" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                      { kind: "Field", name: { kind: "Name", value: "role" } },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "hasMoreBatches" },
                },
                { kind: "Field", name: { kind: "Name", value: "hasMoreData" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "nextCursor" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "time" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "sdInstanceUID" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "kpiDefinitionUID" },
                      },
                    ],
                  },
                },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<TimeSeriesReadQuery, TimeSeriesReadQueryVariables>;
export const StartTimeSeriesExportDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "StartTimeSeriesExport" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "TimeSeriesReadInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "startTimeSeriesExport" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "downloadUrl" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  StartTimeSeriesExportMutation,
  StartTimeSeriesExportMutationVariables
>;
export const CancelTimeSeriesExportDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CancelTimeSeriesExport" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "cancelTimeSeriesExport" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "downloadUrl" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CancelTimeSeriesExportMutation,
  CancelTimeSeriesExportMutationVariables
>;
export const TimeSeriesExportDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "TimeSeriesExport" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "uid" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timeSeriesExport" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "uid" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "uid" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "downloadUrl" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  TimeSeriesExportQuery,
  TimeSeriesExportQueryVariables
>;
export const OnTimeSeriesExportUpdatedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "OnTimeSeriesExportUpdated" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "filter" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "TimeSeriesExportFilter" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "onTimeSeriesExportUpdated" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "filter" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "filter" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "uid" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "downloadUrl" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  OnTimeSeriesExportUpdatedSubscription,
  OnTimeSeriesExportUpdatedSubscriptionVariables
>;
export const TimeSeriesReadAggregateKpiDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "TimeSeriesReadAggregateKpi" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "request" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "TimeSeriesReadAggregateKPIInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timeSeriesReadAggregateKPI" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "request" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "request" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "parameters" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "denotation" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "label" } },
                      { kind: "Field", name: { kind: "Name", value: "role" } },
                    ],
                  },
                },
                { kind: "Field", name: { kind: "Name", value: "base" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "data" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "time" } },
                      { kind: "Field", name: { kind: "Name", value: "tags" } },
                      { kind: "Field", name: { kind: "Name", value: "data" } },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "hasMoreBatches" },
                },
                { kind: "Field", name: { kind: "Name", value: "hasMoreData" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "nextCursor" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "time" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "sdInstanceUID" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "kpiDefinitionUID" },
                      },
                    ],
                  },
                },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  TimeSeriesReadAggregateKpiQuery,
  TimeSeriesReadAggregateKpiQueryVariables
>;
