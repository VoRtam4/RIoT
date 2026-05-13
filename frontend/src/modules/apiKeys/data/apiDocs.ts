/**
 * @file apiDocs.ts
 * @brief Datový popis API dokumentace zobrazované ve frontendové části správy API klíčů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
export type ApiTechnology = "graphql" | "rest" | "websocket";

export type ApiVariantDoc = {
  technology: ApiTechnology;
  label: string;
  summary: string;
  authExample?: string;
  requestExample: string;
  responseExample?: string;
  notes?: string[];
};

export type ApiActionDoc = {
  id: string;
  title: string;
  summary: string;
  variants: ApiVariantDoc[];
};

export type ApiFeatureDoc = {
  id: string;
  title: string;
  summary: string;
  actions: ApiActionDoc[];
};

const gql = "GraphQL";
const rest = "REST";
const ws = "WebSocket";

export const apiFeatures: ApiFeatureDoc[] = [
  {
    id: "sd-types",
    title: "Device Types",
    summary:
      "Define device schemas, inspect one type or the full catalog, and remove obsolete types.",
    actions: [
      {
        id: "list-types",
        title: "List Types",
        summary: "Load the full catalog of device types including parameter metadata.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use one query to load all types and every parameter needed for editors or validation.",
            requestExample: `query ListDeviceTypes {
  sdTypes {
    id
    uid
    label
    parameters {
      id
      label
      denotation
      type
      role
    }
  }
}`,
            responseExample: `{
  "data": {
    "sdTypes": [
      {
        "id": "4",
        "uid": "thermostat",
        "label": "Thermostat",
        "parameters": [
          {
            "id": "21",
            "label": "Temperature",
            "denotation": "temperature",
            "type": "number",
            "role": "field"
          },
          {
            "id": "22",
            "label": "Location",
            "denotation": "location",
            "type": "string",
            "role": "tag"
          }
        ]
      }
    ]
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "The REST collection endpoint returns the same catalog in plain JSON.",
            requestExample: `GET /rest/sd-types
Accept: application/json
Authorization: Bearer <token>`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Request the catalog through the message envelope used by the interactive WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "sd-types-list",
  "action": "get_sd_types"
}`,
            responseExample: `{
  "type": "response",
  "id": "sd-types-list",
  "success": true,
  "payload": [
    {
      "id": 4,
      "uid": "thermostat",
      "label": "Thermostat",
      "parameters": [
        {
          "id": 21,
          "label": "Temperature",
          "denotation": "temperature",
          "type": "number",
          "role": "field"
        }
      ]
    }
  ]
}`,
          },
        ],
      },
      {
        id: "type-detail",
        title: "Get Type Detail",
        summary: "Fetch one specific type including its parameter schema.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Resolve one type by ID for editors, previews, or device-type-aware forms.",
            requestExample: `query DeviceTypeDetail {
  sdType(id: 4) {
    id
    uid
    label
    parameters {
      id
      label
      denotation
      type
      role
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read one type by path parameter.",
            requestExample: `GET /rest/sd-types/4
Accept: application/json`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Pass the target ID in the WebSocket payload.",
            requestExample: `{
  "type": "request",
  "id": "sd-type-detail",
  "action": "get_sd_type",
  "payload": {
    "id": 4
  }
}`,
          },
        ],
      },
      {
        id: "create-type",
        title: "Create Type",
        summary: "Create a new device type with a full parameter definition.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "GraphQL is the richest option when building type editors because the input stays strongly typed.",
            requestExample: `mutation CreateDeviceType {
  createSDType(
    input: {
      uid: "power-meter"
      label: "Power Meter"
      parameters: [
        {
          label: "Voltage"
          denotation: "voltage"
          type: number
          role: field
        }
        {
          label: "Current"
          denotation: "current"
          type: number
          role: field
        }
        {
          label: "Floor"
          denotation: "floor"
          type: string
          role: tag
        }
      ]
    }
  ) {
    id
    uid
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Submit the full type schema as JSON.",
            requestExample: `POST /rest/sd-types
Content-Type: application/json

{
  "uid": "power-meter",
  "label": "Power Meter",
  "parameters": [
    {
      "label": "Voltage",
      "denotation": "voltage",
      "type": "number",
      "role": "field"
    },
    {
      "label": "Current",
      "denotation": "current",
      "type": "number",
      "role": "field"
    },
    {
      "label": "Floor",
      "denotation": "floor",
      "type": "string",
      "role": "tag"
    }
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Wrap the same JSON payload in a WebSocket request message.",
            requestExample: `{
  "type": "request",
  "id": "sd-type-create",
  "action": "create_sd_type",
  "payload": {
    "uid": "power-meter",
    "label": "Power Meter",
    "parameters": [
      {
        "label": "Voltage",
        "denotation": "voltage",
        "type": "number",
        "role": "field"
      },
      {
        "label": "Floor",
        "denotation": "floor",
        "type": "string",
        "role": "tag"
      }
    ]
  }
}`,
          },
        ],
      },
      {
        id: "delete-type",
        title: "Delete Type",
        summary: "Remove a type that should no longer be available.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete by ID and receive a boolean success flag.",
            requestExample: `mutation DeleteDeviceType {
  deleteSDType(id: 12)
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the resource directly.",
            requestExample: `DELETE /rest/sd-types/12`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target ID in the delete action payload.",
            requestExample: `{
  "type": "request",
  "id": "sd-type-delete",
  "action": "delete_sd_type",
  "payload": {
    "id": 12
  }
}`,
          },
        ],
      },
    ],
  },
  {
    id: "sd-instances",
    title: "Devices",
    summary:
      "Browse devices, filter them by type or KPI, inspect one device, update user-managed fields, and subscribe to registrations.",
    actions: [
      {
        id: "list-devices",
        title: "List Devices",
        summary: "Load every registered device instance.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Query the full device list including type metadata for overview pages.",
            requestExample: `query ListDevices {
  sdInstances {
    id
    uid
    label
    confirmedByUser
    userIdentifier
    type {
      id
      uid
      label
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the collection endpoint for the unfiltered list.",
            requestExample: `GET /rest/sd-instances`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "The WebSocket action mirrors the REST collection endpoint.",
            requestExample: `{
  "type": "request",
  "id": "devices-list",
  "action": "get_sd_instances"
}`,
          },
        ],
      },
      {
        id: "device-detail",
        title: "Get Device Detail",
        summary: "Load one device instance and its associated type.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Resolve one device in a single query for detail pages.",
            requestExample: `query DeviceDetail {
  sdInstance(id: 17) {
    id
    uid
    label
    confirmedByUser
    userIdentifier
    type {
      id
      uid
      label
      parameters {
        id
        denotation
        type
        role
      }
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "The current REST router exposes device detail via `/rest/sd-instance`.",
            requestExample: `GET /rest/sd-instance?id=17`,
            notes: [
              "This endpoint uses a query parameter instead of a path parameter in the current backend.",
            ],
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Pass the device ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "device-detail",
  "action": "get_sd_instance",
  "payload": {
    "id": 17
  }
}`,
          },
        ],
      },
      {
        id: "filter-devices",
        title: "Filter Devices",
        summary: "Filter by device type or by KPI definition.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "GraphQL exposes dedicated queries for both filter dimensions.",
            requestExample: `query FilterDevices {
  byType: sdInstancesByType(id: 4) {
    id
    label
    uid
  }
  byKpi: sdInstancesByKpiDefinition(id: 9) {
    id
    label
    uid
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "REST exposes separate endpoints for type-based and KPI-based filtering.",
            requestExample: `GET /rest/sd-instances/type/4
GET /rest/sd-instances/kpi/9`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Use one action per filter mode and send the selected ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "devices-by-type",
  "action": "get_sd_instances_by_type",
  "payload": { "id": 4 }
}

{
  "type": "request",
  "id": "devices-by-kpi",
  "action": "get_sd_instances_by_kpi",
  "payload": { "id": 9 }
}`,
          },
        ],
      },
      {
        id: "update-device",
        title: "Update Device",
        summary:
          "Edit user-controlled fields such as the display label or manual confirmation flag.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use the update mutation to rename the device, set a user identifier, or confirm it.",
            requestExample: `mutation UpdateDevice {
  updateSDInstance(
    id: 17
    input: {
      label: "Boiler Sensor A-17"
      userIdentifier: "asset-plant-a-17"
      confirmedByUser: true
    }
  ) {
    id
    label
    userIdentifier
    confirmedByUser
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PATCH the mutable fields of one device.",
            requestExample: `PATCH /rest/sd-instances/17
Content-Type: application/json

{
  "label": "Boiler Sensor A-17",
  "userIdentifier": "asset-plant-a-17",
  "confirmedByUser": true
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the device ID and the update input in one message.",
            requestExample: `{
  "type": "request",
  "id": "device-update",
  "action": "update_sd_instance",
  "payload": {
    "id": 17,
    "input": {
      "label": "Boiler Sensor A-17",
      "userIdentifier": "asset-plant-a-17",
      "confirmedByUser": true
    }
  }
}`,
          },
        ],
      },
      {
        id: "device-registration-stream",
        title: "Subscribe To Device Registrations",
        summary:
          "Listen for newly registered devices to update dashboards or onboarding flows in real time.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Subscribe directly to registration events with an optional filter.",
            requestExample: `subscription DeviceRegistrationStream {
  onSDInstanceRegistered(filter: { sdTypeID: 4 }) {
    id
    uid
    label
    confirmedByUser
  }
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket API uses `type: subscribe` and an event topic.",
            requestExample: `{
  "type": "subscribe",
  "id": "device-registration-stream",
  "topic": "sd_instance_registered",
  "payload": {
    "sdTypeID": 4
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "For one-way server push, use the REST SSE endpoint under the same REST interface family.",
            authExample: `curl -N \\
  -H "X-API-Key: riot_live_abc123_secret_value" \\
  -H "Accept: text/event-stream" \\
  http://localhost:9090/rest/subscribe/sd_instance_registered`,
            requestExample: `GET /rest/subscribe/sd_instance_registered
Accept: text/event-stream

Body:
{
  "sdTypeID": 4
}`,
          },
        ],
      },
    ],
  },
  {
    id: "raw-data-points",
    title: "Raw Data Points",
    summary:
      "Read the latest raw payloads either for a whole device type or for one specific device instance.",
    actions: [
      {
        id: "list-raw-by-sdtype",
        title: "List Latest Raw Data By Device Type",
        summary:
          "Load the latest raw data points for all devices belonging to one device type.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use the SD type ID to fetch the latest raw data points across all matching device instances.",
            requestExample: `query RawDataByType {
  rawDataPointsBySDType(id: 4) {
    sdTypeID
    sdInstanceID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "rawDataPointsBySDType": [
      {
        "sdTypeID": "4",
        "sdInstanceID": "17",
        "payload": {
          "temperature": 82.4,
          "isOnline": true,
          "location": "boiler-room-a"
        },
        "eventTime": "2026-04-14T12:30:00Z"
      },
      {
        "sdTypeID": "4",
        "sdInstanceID": "18",
        "payload": {
          "temperature": 79.1,
          "isOnline": true,
          "location": "boiler-room-b"
        },
        "eventTime": "2026-04-14T12:29:42Z"
      }
    ]
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Read the latest raw data set for one SD type from the REST path-based endpoint.",
            requestExample: `GET /rest/raw/4
Accept: application/json
X-API-Key: riot_live_abc123_secret_value`,
            responseExample: `[
  {
    "sdTypeID": 4,
    "sdInstanceID": 17,
    "payload": {
      "temperature": 82.4,
      "isOnline": true,
      "location": "boiler-room-a"
    },
    "eventTime": "2026-04-14T12:30:00Z"
  },
  {
    "sdTypeID": 4,
    "sdInstanceID": 18,
    "payload": {
      "temperature": 79.1,
      "isOnline": true,
      "location": "boiler-room-b"
    },
    "eventTime": "2026-04-14T12:29:42Z"
  }
]`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Request the latest raw data set by SD type through the WebSocket request/response channel.",
            requestExample: `{
  "type": "request",
  "id": "raw-by-type",
  "action": "get_raw_by_sdtype",
  "payload": {
    "id": 4
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "raw-by-type",
  "success": true,
  "payload": [
    {
      "sdTypeID": 4,
      "sdInstanceID": 17,
      "payload": {
        "temperature": 82.4,
        "isOnline": true
      },
      "eventTime": "2026-04-14T12:30:00Z"
    }
  ]
}`,
          },
        ],
      },
      {
        id: "latest-raw-by-instance",
        title: "Get Latest Raw Data For One Device",
        summary:
          "Load the most recent raw data point for a single device instance.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Query the latest raw data point for one device instance by its ID.",
            requestExample: `query LatestRawDataPoint {
  rawDataPoint(id: 17) {
    sdTypeID
    sdInstanceID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "rawDataPoint": {
      "sdTypeID": "4",
      "sdInstanceID": "17",
      "payload": {
        "temperature": 82.4,
        "isOnline": true,
        "location": "boiler-room-a"
      },
      "eventTime": "2026-04-14T12:30:00Z"
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Use the single raw-data REST endpoint with the device instance ID.",
            requestExample: `GET /rest/raw?id=17
Accept: application/json
X-API-Key: riot_live_abc123_secret_value`,
            responseExample: `{
  "sdTypeID": 4,
  "sdInstanceID": 17,
  "payload": {
    "temperature": 82.4,
    "isOnline": true,
    "location": "boiler-room-a"
  },
  "eventTime": "2026-04-14T12:30:00Z"
}`,
            notes: [
              "The current REST router exposes this endpoint as `/rest/raw` and expects the instance ID through request parsing logic.",
            ],
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Request the latest raw point for one device instance through the low-level WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "raw-by-instance",
  "action": "get_raw_by_instance",
  "payload": {
    "id": 17
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "raw-by-instance",
  "success": true,
  "payload": {
    "sdTypeID": 4,
    "sdInstanceID": 17,
    "payload": {
      "temperature": 82.4,
      "isOnline": true
    },
    "eventTime": "2026-04-14T12:30:00Z"
  }
}`,
          },
        ],
      },
      {
        id: "raw-data-stream",
        title: "Subscribe To Raw Data Points",
        summary:
          "Receive incoming raw data points in real time for a selected device type or device instance.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Subscribe to incoming raw data and optionally filter by device type or device instance.",
            requestExample: `subscription RawDataPointStream {
  onRawDataPointArrived(filter: { sdTypeID: 4, sdInstanceID: 17 }) {
    sdTypeID
    sdInstanceID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "onRawDataPointArrived": [
      {
        "sdTypeID": "4",
        "sdInstanceID": "17",
        "payload": {
          "temperature": 83.2,
          "isOnline": true
        },
        "eventTime": "2026-04-14T12:31:05Z"
      }
    ]
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Use the REST SSE stream when the client only needs server-to-client event delivery.",
            authExample: `curl -N \\
  -H "X-API-Key: riot_live_abc123_secret_value" \\
  -H "Accept: text/event-stream" \\
  http://localhost:9090/rest/subscribe/raw_data_point`,
            requestExample: `GET /rest/subscribe/raw_data_point
Accept: text/event-stream

Body:
{
  "sdTypeID": 4,
  "sdInstanceID": 17
}`,
            responseExample: `event: raw_data_point
data: [{"sdTypeID":4,"sdInstanceID":17,"payload":{"temperature":83.2,"isOnline":true},"eventTime":"2026-04-14T12:31:05Z"}]`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Use the low-level WebSocket subscribe envelope with the raw-data event topic.",
            requestExample: `{
  "type": "subscribe",
  "id": "raw-data-stream",
  "topic": "raw_data_point",
  "payload": {
    "sdTypeID": 4,
    "sdInstanceID": 17
  }
}`,
            responseExample: `{
  "type": "event",
  "id": "raw-data-stream",
  "topic": "raw_data_point",
  "payload": [
    {
      "sdTypeID": 4,
      "sdInstanceID": 17,
      "payload": {
        "temperature": 83.2,
        "isOnline": true
      },
      "eventTime": "2026-04-14T12:31:05Z"
    }
  ]
}`,
          },
        ],
      },
    ],
  },
  {
    id: "kpi-definitions",
    title: "KPI Definitions",
    summary:
      "Create, inspect, filter, update, and delete KPI definitions including logical node trees and instance targeting.",
    actions: [
      {
        id: "list-kpis",
        title: "List KPI Definitions",
        summary: "Load the full KPI definition catalog.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Query all KPI definitions with the selected instance mode and node tree.",
            requestExample: `query ListKpis {
  kpiDefinitions {
    id
    label
    sdTypeID
    sdTypeUID
    userIdentifier
    sdInstanceMode
    selectedSDInstanceIDs
    nodes {
      id
      parentNodeID
      nodeType
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read the KPI catalog from the REST collection endpoint.",
            requestExample: `GET /rest/kpi-definitions`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "The WebSocket action returns the same catalog payload.",
            requestExample: `{
  "type": "request",
  "id": "kpi-list",
  "action": "get_kpi_definitions"
}`,
          },
        ],
      },
      {
        id: "kpi-detail",
        title: "Get KPI Detail",
        summary: "Load one KPI definition including all nodes.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Fetch one KPI for detail pages or rule editors.",
            requestExample: `query KpiDetail {
  kpiDefinition(id: 9) {
    id
    label
    sdTypeID
    sdTypeUID
    sdInstanceMode
    selectedSDInstanceIDs
    nodes {
      id
      parentNodeID
      nodeType
      ... on LogicalOperationKPINode {
        type
      }
      ... on NumericGTAtomKPINode {
        sdParameterID
        sdParameterSpecification
        numericReferenceValue
      }
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read one KPI definition by path parameter.",
            requestExample: `GET /rest/kpi-definitions/9`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target KPI ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "kpi-detail",
  "action": "get_kpi_definition",
  "payload": {
    "id": 9
  }
}`,
          },
        ],
      },
      {
        id: "filter-kpis",
        title: "Filter KPI Definitions",
        summary: "Filter KPI definitions by device type or device instance.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use dedicated filter queries for editors bound to a device type or a selected device.",
            requestExample: `query FilterKpis {
  byType: kpiDefinitionsBySdType(id: 4) {
    id
    label
  }
  byDevice: kpiDefinitionsBySdInstance(id: 17) {
    id
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "REST exposes separate endpoints for both filter dimensions.",
            requestExample: `GET /rest/kpi-definitions/type/4
GET /rest/kpi-definitions/instance/17`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use one request action per filter type.",
            requestExample: `{
  "type": "request",
  "id": "kpi-by-type",
  "action": "get_kpi_definitions_by_type",
  "payload": { "id": 4 }
}

{
  "type": "request",
  "id": "kpi-by-device",
  "action": "get_kpi_definitions_by_sd_instance",
  "payload": { "id": 17 }
}`,
          },
        ],
      },
      {
        id: "create-kpi",
        title: "Create KPI Definition",
        summary:
          "Create a KPI definition including targeting mode and a rule tree composed of nodes.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "This example creates a KPI for one device type and targets selected instances.",
            requestExample: `mutation CreateKpi {
  createKPIDefinition(
    input: {
      label: "Overheat Alarm"
      sdTypeID: 4
      sdTypeUID: "thermostat"
      userIdentifier: "user-42"
      sdInstanceMode: selected
      selectedSDInstanceIDs: [17, 18]
      nodes: [
        {
          id: 1
          parentNodeID: null
          nodeType: logicalOperation
          logicalOperationType: and
        }
        {
          id: 2
          parentNodeID: 1
          nodeType: numericGTAtom
          sdParameterID: 21
          sdParameterSpecification: "temperature"
          numericReferenceValue: 85
        }
        {
          id: 3
          parentNodeID: 1
          nodeType: booleanEQAtom
          sdParameterID: 25
          sdParameterSpecification: "isOnline"
          booleanReferenceValue: true
        }
      ]
    }
  ) {
    id
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Send the same rule tree in JSON to the create endpoint.",
            requestExample: `POST /rest/kpi-definitions
Content-Type: application/json

{
  "label": "Overheat Alarm",
  "sdTypeID": 4,
  "sdTypeUID": "thermostat",
  "userIdentifier": "user-42",
  "sdInstanceMode": "selected",
  "selectedSDInstanceIDs": [17, 18],
  "nodes": [
    {
      "id": 1,
      "parentNodeID": null,
      "nodeType": "logicalOperation",
      "logicalOperationType": "and"
    },
    {
      "id": 2,
      "parentNodeID": 1,
      "nodeType": "numericGTAtom",
      "sdParameterID": 21,
      "sdParameterSpecification": "temperature",
      "numericReferenceValue": 85
    }
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Wrap the same create payload in the WebSocket envelope.",
            requestExample: `{
  "type": "request",
  "id": "kpi-create",
  "action": "create_kpi_definition",
  "payload": {
    "label": "Overheat Alarm",
    "sdTypeID": 4,
    "sdTypeUID": "thermostat",
    "userIdentifier": "user-42",
    "sdInstanceMode": "selected",
    "selectedSDInstanceIDs": [17, 18],
    "nodes": [
      {
        "id": 1,
        "parentNodeID": null,
        "nodeType": "logicalOperation",
        "logicalOperationType": "and"
      }
    ]
  }
}`,
          },
        ],
      },
      {
        id: "update-kpi",
        title: "Update KPI Definition",
        summary:
          "Replace the full KPI definition when editing label, targets, or node tree structure.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "The update mutation uses the same rich input shape as creation.",
            requestExample: `mutation UpdateKpi {
  updateKPIDefinition(
    id: 9
    input: {
      label: "Overheat Alarm v2"
      sdTypeID: 4
      sdTypeUID: "thermostat"
      userIdentifier: "user-42"
      sdInstanceMode: all
      selectedSDInstanceIDs: []
      nodes: [
        {
          id: 1
          parentNodeID: null
          nodeType: logicalOperation
          logicalOperationType: and
        },
        {
          id: 2
          parentNodeID: 1
          nodeType: numericGTAtom
          sdParameterID: 21
          sdParameterSpecification: "temperature"
          numericReferenceValue: 90
        }
      ]
    }
  ) {
    id
    label
    sdInstanceMode
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PUT the complete definition to keep server and editor state aligned.",
            requestExample: `PUT /rest/kpi-definitions/9
Content-Type: application/json

{
  "label": "Overheat Alarm v2",
  "sdTypeID": 4,
  "sdTypeUID": "thermostat",
  "userIdentifier": "user-42",
  "sdInstanceMode": "all",
  "selectedSDInstanceIDs": [],
  "nodes": [
    {
      "id": 1,
      "parentNodeID": null,
      "nodeType": "logicalOperation",
      "logicalOperationType": "and"
    },
    {
      "id": 2,
      "parentNodeID": 1,
      "nodeType": "numericGTAtom",
      "sdParameterID": 21,
      "sdParameterSpecification": "temperature",
      "numericReferenceValue": 90
    }
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target ID plus the full replacement input.",
            requestExample: `{
  "type": "request",
  "id": "kpi-update",
  "action": "update_kpi_definition",
  "payload": {
    "id": 9,
    "input": {
      "label": "Overheat Alarm v2",
      "sdTypeID": 4,
      "sdTypeUID": "thermostat",
      "userIdentifier": "user-42",
      "sdInstanceMode": "all",
      "selectedSDInstanceIDs": [],
      "nodes": []
    }
  }
}`,
          },
        ],
      },
      {
        id: "delete-kpi",
        title: "Delete KPI Definition",
        summary: "Delete an obsolete KPI definition by ID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete through a simple mutation and check the boolean result.",
            requestExample: `mutation DeleteKpi {
  deleteKPIDefinition(id: 9)
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the KPI resource directly.",
            requestExample: `DELETE /rest/kpi-definitions/9`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the delete action with the KPI ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "kpi-delete",
  "action": "delete_kpi_definition",
  "payload": {
    "id": 9
  }
}`,
          },
        ],
      },
    ],
  },
  {
    id: "kpi-results",
    title: "KPI Results",
    summary:
      "Inspect evaluation results, retrieve a specific result, and subscribe to fresh KPI outcomes in real time.",
    actions: [
      {
        id: "list-kpi-results",
        title: "List KPI Results",
        summary: "Load the latest KPI evaluation results across the system.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Fetch all results or include only the fields needed by the current view.",
            requestExample: `query ListKpiResults {
  kpiResults {
    id
    kpiDefinitionID
    sdInstanceID
    fulfilled
    payload
    eventTime
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "The REST endpoint returns the result collection in JSON.",
            requestExample: `GET /rest/kpi-results`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the `get_kpi_results` request action.",
            requestExample: `{
  "type": "request",
  "id": "kpi-results-list",
  "action": "get_kpi_results"
}`,
          },
        ],
      },
      {
        id: "results-by-kpi",
        title: "List Results For One KPI",
        summary: "Filter the result history by KPI definition.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Query a KPI-scoped result history.",
            requestExample: `query ResultsByKpi {
  kpiResultsByKPI(id: 9) {
    id
    sdInstanceID
    fulfilled
    eventTime
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the KPI ID as a path parameter.",
            requestExample: `GET /rest/kpi-results/9`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "The action mirrors the REST path-based filter.",
            requestExample: `{
  "type": "request",
  "id": "kpi-results-by-kpi",
  "action": "get_kpi_results_by_kpi",
  "payload": {
    "id": 9
  }
}`,
          },
        ],
      },
      {
        id: "single-kpi-result",
        title: "Get One Result",
        summary:
          "Resolve one evaluation result from a structured request object, for example from a selected timestamp.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use a rich request object when the result must be resolved by more than just a numeric ID.",
            requestExample: `query OneKpiResult {
  kpiResult(
    request: {
      kpiDefinitionID: 9
      sdInstanceID: 17
      eventTime: "2026-04-14T12:30:00Z"
    }
  ) {
    id
    fulfilled
    payload
    eventTime
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "POST the structured request payload to resolve one exact result.",
            requestExample: `POST /rest/kpi-result
Content-Type: application/json

{
  "kpiDefinitionID": 9,
  "sdInstanceID": 17,
  "eventTime": "2026-04-14T12:30:00Z"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the same request object via WebSocket.",
            requestExample: `{
  "type": "request",
  "id": "kpi-result-one",
  "action": "get_kpi_result",
  "payload": {
    "kpiDefinitionID": 9,
    "sdInstanceID": 17,
    "eventTime": "2026-04-14T12:30:00Z"
  }
}`,
          },
        ],
      },
      {
        id: "kpi-result-stream",
        title: "Subscribe To KPI Result Events",
        summary: "Receive new KPI evaluations as they happen.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Subscribe directly from Apollo or another GraphQL client.",
            requestExample: `subscription KpiResultStream {
  onKPIFulfillmentChecked(filter: { kpiDefinitionID: 9, sdInstanceID: 17 }) {
    id
    fulfilled
    eventTime
    payload
  }
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Low-level WebSocket subscriptions use a topic and a filter payload.",
            requestExample: `{
  "type": "subscribe",
  "id": "kpi-result-stream",
  "topic": "kpi_fulfillment_checked",
  "payload": {
    "kpiDefinitionID": 9,
    "sdInstanceID": 17
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "For one-way streaming over HTTP, use the SSE endpoint that belongs to the REST interface.",
            authExample: `curl -N \\
  -H "X-API-Key: riot_live_abc123_secret_value" \\
  -H "Accept: text/event-stream" \\
  http://localhost:9090/rest/subscribe/kpi_fulfillment_checked`,
            requestExample: `GET /rest/subscribe/kpi_fulfillment_checked
Accept: text/event-stream

Body:
{
  "kpiDefinitionID": 9,
  "sdInstanceID": 17
}`,
          },
        ],
      },
    ],
  },
  {
    id: "device-groups",
    title: "Device Groups",
    summary:
      "Manage logical groups of devices for organization, targeting, or dashboard composition.",
    actions: [
      {
        id: "list-groups",
        title: "List Groups",
        summary: "Load all device groups.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Load the group collection with its member IDs.",
            requestExample: `query ListGroups {
  sdInstanceGroups {
    id
    label
    sdInstanceIDs
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read the group collection from the REST endpoint.",
            requestExample: `GET /rest/sd-instance-groups`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the corresponding collection request action.",
            requestExample: `{
  "type": "request",
  "id": "groups-list",
  "action": "get_sd_instance_groups"
}`,
          },
        ],
      },
      {
        id: "group-detail",
        title: "Get Group Detail",
        summary: "Fetch one group and inspect its membership.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Read one group by ID.",
            requestExample: `query GroupDetail {
  sdInstanceGroup(id: 3) {
    id
    label
    sdInstanceIDs
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use a path parameter to fetch one group.",
            requestExample: `GET /rest/sd-instance-groups/3`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the group ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "group-detail",
  "action": "get_sd_instance_group",
  "payload": { "id": 3 }
}`,
          },
        ],
      },
      {
        id: "create-group",
        title: "Create Group",
        summary: "Create a new group and attach one or more devices.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Create a group with its initial device membership.",
            requestExample: `mutation CreateGroup {
  createSDInstanceGroup(
    input: {
      label: "Boiler Room"
      sdInstanceIDs: [17, 18, 19]
    }
  ) {
    id
    label
    sdInstanceIDs
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Send the group definition in JSON.",
            requestExample: `POST /rest/sd-instance-groups
Content-Type: application/json

{
  "label": "Boiler Room",
  "sdInstanceIDs": [17, 18, 19]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Wrap the same create payload in a request message.",
            requestExample: `{
  "type": "request",
  "id": "group-create",
  "action": "create_sd_instance_group",
  "payload": {
    "label": "Boiler Room",
    "sdInstanceIDs": [17, 18, 19]
  }
}`,
          },
        ],
      },
      {
        id: "update-group",
        title: "Update Group",
        summary: "Change the group label or replace its membership.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Replace the group state in one mutation.",
            requestExample: `mutation UpdateGroup {
  updateSDInstanceGroup(
    id: 3
    input: {
      label: "Boiler Room - North"
      sdInstanceIDs: [17, 18, 22]
    }
  ) {
    id
    label
    sdInstanceIDs
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PUT the full group payload.",
            requestExample: `PUT /rest/sd-instance-groups/3
Content-Type: application/json

{
  "label": "Boiler Room - North",
  "sdInstanceIDs": [17, 18, 22]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Provide the group ID together with the replacement input.",
            requestExample: `{
  "type": "request",
  "id": "group-update",
  "action": "update_sd_instance_group",
  "payload": {
    "id": 3,
    "input": {
      "label": "Boiler Room - North",
      "sdInstanceIDs": [17, 18, 22]
    }
  }
}`,
          },
        ],
      },
      {
        id: "delete-group",
        title: "Delete Group",
        summary: "Delete a group that is no longer needed.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete by group ID and check the boolean result.",
            requestExample: `mutation DeleteGroup {
  deleteSDInstanceGroup(id: 3)
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the group resource directly.",
            requestExample: `DELETE /rest/sd-instance-groups/3`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Issue the delete action with the group ID.",
            requestExample: `{
  "type": "request",
  "id": "group-delete",
  "action": "delete_sd_instance_group",
  "payload": {
    "id": 3
  }
}`,
          },
        ],
      },
    ],
  },
  {
    id: "time-series",
    title: "History And Export",
    summary:
      "Read raw or KPI history, resolve distinct tag values, paginate through batches, use aggregate KPI reads, and trigger exports for offline analysis.",
    actions: [
      {
        id: "read-history",
        title: "Read History",
        summary:
          "Read raw or KPI time-series data with filters, sorting, limits, and cursors.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "This query reads KPI history for one device and carries the cursor needed for the next request.",
            requestExample: `query ReadHistory {
  timeSeriesRead(
    request: {
      type: kpi
      sdTypeID: 4
      sdInstanceIDs: [17]
      kpiDefinitionIDs: [9]
      sortDesc: true
      limit: 50
      cursor: {
        time: "2026-04-14T12:00:00Z"
        sdInstanceUID: "boiler-sensor-a-17"
      }
    }
  ) {
    parameters {
      denotation
      label
      role
    }
    base
    data {
      time
      tags
      data
    }
    hasMoreBatches
    hasMoreData
    nextCursor {
      time
      sdInstanceUID
      kpiDefinitionID
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "The REST handler expects the request object in the body. POST is the practical method for this shape.",
            requestExample: `POST /rest/time-series
Content-Type: application/json

{
  "type": "kpi",
  "sdTypeID": 4,
  "sdInstanceIDs": [17],
  "kpiDefinitionIDs": [9],
  "sortDesc": true,
  "limit": 50,
  "cursor": {
    "time": "2026-04-14T12:00:00Z",
    "sdInstanceUID": "boiler-sensor-a-17"
  }
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket API streams responses in one or more batches using the same message ID.",
            requestExample: `{
  "type": "request",
  "id": "history-read",
  "action": "time-series",
  "payload": {
    "type": "kpi",
    "sdTypeID": 4,
    "sdInstanceIDs": [17],
    "kpiDefinitionIDs": [9],
    "sortDesc": true,
    "limit": 50
  }
}`,
          },
        ],
      },
      {
        id: "distinct-tag-values",
        title: "Get Distinct Tag Values",
        summary:
          "Resolve the available values for one tag with optional filters over raw or KPI history.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use this query to populate filter dropdowns from the current history slice.",
            requestExample: `query DistinctTagValues {
  timeSeriesDistinctTagValues(
    request: {
      type: raw
      sdTypeID: 4
      sdInstanceIDs: [17, 18]
      from: "2026-04-01T00:00:00Z"
      to: "2026-04-14T23:59:59Z"
      tag: "location"
    }
  ) {
    values
    error
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Submit the same request object to the dedicated distinct-tag-values endpoint.",
            requestExample: `POST /rest/time-series/distinct-tag-values
Content-Type: application/json

{
  "type": "raw",
  "sdTypeID": 4,
  "sdInstanceIDs": [17, 18],
  "from": "2026-04-01T00:00:00Z",
  "to": "2026-04-14T23:59:59Z",
  "tag": "location"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket action returns the same distinct-values payload.",
            requestExample: `{
  "type": "request",
  "id": "history-distinct-tag-values",
  "action": "time-series-distinct-tag-values",
  "payload": {
    "type": "raw",
    "sdTypeID": 4,
    "sdInstanceIDs": [17, 18],
    "from": "2026-04-01T00:00:00Z",
    "to": "2026-04-14T23:59:59Z",
    "tag": "location"
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-distinct-tag-values",
  "success": true,
  "payload": {
    "values": ["floor-1", "floor-2"],
    "error": ""
  }
}`,
          },
        ],
      },
      {
        id: "read-aggregate-kpi",
        title: "Read Aggregate KPI History",
        summary:
          "Request aggregated KPI history for charting and analytical overviews.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "GraphQL already exposes the aggregate KPI reader as a first-class query.",
            requestExample: `query ReadAggregateKpi {
  timeSeriesReadAggregateKPI(
    request: {
      sdTypeID: 4
      kpiDefinitionIDs: [9, 11]
      sdInstanceIDs: [17, 18]
      from: "2026-04-01T00:00:00Z"
      to: "2026-04-14T23:59:59Z"
      aggregateSeconds: 3600
      batch: 200
      limit: 200
    }
  ) {
    parameters {
      denotation
      label
      role
    }
    base
    data {
      time
      tags
      data
    }
    hasMoreBatches
    hasMoreData
    nextCursor {
      time
      sdInstanceUID
      kpiDefinitionID
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "A dedicated aggregate endpoint can be exposed through the existing aggregate handler.",
            requestExample: `POST /rest/time-series/aggregate-kpi
Content-Type: application/json

{
  "sdTypeID": 4,
  "kpiDefinitionIDs": [9, 11],
  "sdInstanceIDs": [17, 18],
  "from": "2026-04-01T00:00:00Z",
  "to": "2026-04-14T23:59:59Z",
  "aggregateSeconds": 3600,
  "batch": 200,
  "limit": 200
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The aggregate handler streams batches using the same message ID.",
            requestExample: `{
  "type": "request",
  "id": "history-aggregate",
  "action": "time-series-aggregate-kpi",
  "payload": {
    "sdTypeID": 4,
    "kpiDefinitionIDs": [9, 11],
    "sdInstanceIDs": [17, 18],
    "from": "2026-04-01T00:00:00Z",
    "to": "2026-04-14T23:59:59Z",
    "aggregateSeconds": 3600,
    "batch": 200,
    "limit": 200
  }
}`,
          },
        ],
      },
      {
        id: "export-history",
        title: "Start History Export",
        summary:
          "Start an asynchronous export job and receive the initial export snapshot.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "GraphQL starts the job and returns the current TimeSeriesExport state immediately.",
            requestExample: `mutation StartExport {
  startTimeSeriesExport(
    input: {
      type: raw
      sdTypeID: 4
      sdInstanceIDs: [17, 18]
      sortDesc: true
      limit: 1000
    }
  ) {
    id
    status
    downloadUrl
    createdAt
    expiresAt
    error
  }
}`,
            responseExample: `{
  "data": {
    "startTimeSeriesExport": {
      "id": 41,
      "status": "pending",
      "downloadUrl": null,
      "createdAt": "2026-05-02T23:53:23.626667429Z",
      "expiresAt": null,
      "error": null
    }
  }
}`,
            notes: [
              "Use the returned id for status checks, cancellation, or subscriptions.",
              "The CSV file is still downloaded through REST once the job reaches done.",
            ],
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "POST the history request and receive the same export snapshot shape as GraphQL.",
            requestExample: `POST /rest/time-series/export
Content-Type: application/json

{
  "type": "raw",
  "sdTypeID": 4,
  "sdInstanceIDs": [17, 18],
  "sortDesc": true,
  "limit": 1000
}`,
            responseExample: `{
  "id": 41,
  "status": "pending",
  "downloadUrl": null,
  "createdAt": "2026-05-02T23:53:23.626667429Z",
  "expiresAt": null,
  "error": null
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket request returns the same TimeSeriesExport payload in the response.",
            requestExample: `{
  "type": "request",
  "id": "history-export",
  "action": "time-series-export",
  "payload": {
    "type": "raw",
    "sdTypeID": 4,
    "sdInstanceIDs": [17, 18],
    "limit": 1000
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export",
  "success": true,
  "payload": {
    "id": 41,
    "status": "pending",
    "downloadUrl": null,
    "createdAt": "2026-05-02T23:53:23.626667429Z",
    "expiresAt": null,
    "error": null
  }
}`,
          },
        ],
      },
      {
        id: "export-status",
        title: "Get Export Status",
        summary:
          "Load the current state of one export job without waiting for the file download endpoint.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use the export ID to fetch a single snapshot that can be pending, processing, done, failed, cancelled, or expired.",
            requestExample: `query ExportStatus {
  timeSeriesExport(id: 41) {
    id
    status
    downloadUrl
    createdAt
    expiresAt
    error
  }
}`,
            responseExample: `{
  "data": {
    "timeSeriesExport": {
      "id": 41,
      "status": "done",
      "downloadUrl": "/rest/time-series/export/41",
      "createdAt": "2026-05-02T23:53:23.626667429Z",
      "expiresAt": "2026-05-03T00:23:31.688347473Z",
      "error": null
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "The status endpoint returns the same export shape as the GraphQL query and WebSocket status request.",
            requestExample: `GET /rest/time-series/export/41/status`,
            responseExample: `{
  "id": 41,
  "status": "done",
  "downloadUrl": "/rest/time-series/export/41",
  "createdAt": "2026-05-02T23:53:23.626667429Z",
  "expiresAt": "2026-05-03T00:23:31.688347473Z",
  "error": null
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level request/response WebSocket API exposes an explicit status action as well.",
            requestExample: `{
  "type": "request",
  "id": "history-export-status",
  "action": "time-series-export-status",
  "payload": {
    "id": 41
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-status",
  "success": true,
  "payload": {
    "id": 41,
    "status": "done",
    "downloadUrl": "/rest/time-series/export/41",
    "createdAt": "2026-05-02T23:53:23.626667429Z",
    "expiresAt": "2026-05-03T00:23:31.688347473Z",
    "error": null
  }
}`,
          },
        ],
      },
      {
        id: "cancel-export",
        title: "Cancel Export",
        summary:
          "Request cancellation of one export job and receive its latest snapshot.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "GraphQL cancellation returns the terminal or already-completed export snapshot.",
            requestExample: `mutation CancelExport {
  cancelTimeSeriesExport(id: 41) {
    id
    status
    downloadUrl
    createdAt
    expiresAt
    error
  }
}`,
            responseExample: `{
  "data": {
    "cancelTimeSeriesExport": {
      "id": 41,
      "status": "cancelled",
      "downloadUrl": null,
      "createdAt": "2026-05-02T23:53:23.626667429Z",
      "expiresAt": "2026-05-03T00:23:31.688347473Z",
      "error": null
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Use DELETE on the export resource to request cancellation.",
            requestExample: `DELETE /rest/time-series/export/41`,
            responseExample: `{
  "id": 41,
  "status": "cancelled",
  "downloadUrl": null,
  "createdAt": "2026-05-02T23:53:23.626667429Z",
  "expiresAt": "2026-05-03T00:23:31.688347473Z",
  "error": null
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket API uses a dedicated cancel action with the export id in the payload.",
            requestExample: `{
  "type": "request",
  "id": "history-export-cancel",
  "action": "time-series-export-cancel",
  "payload": {
    "id": 41
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-cancel",
  "success": true,
  "payload": {
    "id": 41,
    "status": "cancelled",
    "downloadUrl": null,
    "createdAt": "2026-05-02T23:53:23.626667429Z",
    "expiresAt": "2026-05-03T00:23:31.688347473Z",
    "error": null
  }
}`,
          },
        ],
      },
      {
        id: "subscribe-export",
        title: "Subscribe To Export Updates",
        summary:
          "Receive push updates when the export moves through pending, processing, done, failed, cancelled, or expired.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Apollo-style GraphQL subscriptions stream the same TimeSeriesExport payload filtered by export IDs.",
            requestExample: `subscription OnExportUpdated {
  onTimeSeriesExportUpdated(filter: { ids: [41] }) {
    id
    status
    downloadUrl
    createdAt
    expiresAt
    error
  }
}`,
            responseExample: `{
  "data": {
    "onTimeSeriesExportUpdated": {
      "id": 41,
      "status": "done",
      "downloadUrl": "/rest/time-series/export/41",
      "createdAt": "2026-05-02T23:53:23.626667429Z",
      "expiresAt": "2026-05-03T00:23:31.688347473Z",
      "error": null
    }
  }
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level WebSocket API subscribes to the shared time_series_export_updated topic and pushes event messages.",
            requestExample: `{
  "type": "subscribe",
  "id": "history-export-events",
  "topic": "time_series_export_updated",
  "payload": {
    "ids": [41]
  }
}`,
            responseExample: `{
  "type": "event",
  "id": "history-export-events",
  "topic": "time_series_export_updated",
  "payload": {
    "id": 41,
    "status": "done",
    "downloadUrl": "/rest/time-series/export/41",
    "createdAt": "2026-05-02T23:53:23.626667429Z",
    "expiresAt": "2026-05-03T00:23:31.688347473Z",
    "error": null
  }
}`,
          },
        ],
      },
      {
        id: "download-export",
        title: "Download Export File",
        summary:
          "Download the generated CSV file through REST once the export status reaches done.",
        variants: [
          {
            technology: "rest",
            label: rest,
            summary:
              "The download endpoint serves the file. If called too early it waits for job completion, so the recommended flow is start -> status/subscription -> download.",
            requestExample: `GET /rest/time-series/export/41`,
            responseExample: `HTTP/1.1 200 OK
Content-Type: text/csv
Content-Disposition: attachment; filename="time-series-export-41.csv"

time,sdInstanceUID,temperature
2026-05-02T23:53:24Z,boiler-sensor-a-17,71.2
2026-05-02T23:53:25Z,boiler-sensor-a-17,71.3`,
          },
        ],
      },
      {
        id: "export-aggregate-kpi",
        title: "Start Aggregate KPI Export",
        summary:
          "Export aggregate KPI history for reporting, charts, or external analysis.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "The aggregate export uses the same asynchronous TimeSeriesExport lifecycle as raw history export.",
            requestExample: `mutation StartAggregateExport {
  startTimeSeriesExportAggregateKPI(
    input: {
      sdTypeID: 4
      kpiDefinitionIDs: [9]
      sdInstanceIDs: [17, 18]
      from: "2026-04-01T00:00:00Z"
      to: "2026-04-14T23:59:59Z"
      aggregateSeconds: 3600
    }
  ) {
    id
    status
    downloadUrl
    createdAt
    expiresAt
    error
  }
}`,
            responseExample: `{
  "data": {
    "startTimeSeriesExportAggregateKPI": {
      "id": 56,
      "status": "pending",
      "downloadUrl": null,
      "createdAt": "2026-05-02T23:58:11.054112902Z",
      "expiresAt": null,
      "error": null
    }
  }
}`,
            notes: [
              "Use the same status, cancellation, subscription, and REST download operations as normal history export.",
            ],
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "POST the aggregate export request to the dedicated REST endpoint and receive the same export snapshot payload.",
            requestExample: `POST /rest/time-series/export/aggregate-kpi
Content-Type: application/json

{
  "sdTypeID": 4,
  "kpiDefinitionIDs": [9],
  "sdInstanceIDs": [17, 18],
  "from": "2026-04-01T00:00:00Z",
  "to": "2026-04-14T23:59:59Z",
  "aggregateSeconds": 3600
}`,
            responseExample: `{
  "id": 56,
  "status": "pending",
  "downloadUrl": null,
  "createdAt": "2026-05-02T23:58:11.054112902Z",
  "expiresAt": null,
  "error": null
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "The low-level aggregate export action also returns the common TimeSeriesExport payload.",
            requestExample: `{
  "type": "request",
  "id": "history-export-aggregate",
  "action": "time-series-export-aggregate-kpi",
  "payload": {
    "sdTypeID": 4,
    "kpiDefinitionIDs": [9],
    "sdInstanceIDs": [17, 18],
    "aggregateSeconds": 3600
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-aggregate",
  "success": true,
  "payload": {
    "id": 56,
    "status": "pending",
    "downloadUrl": null,
    "createdAt": "2026-05-02T23:58:11.054112902Z",
    "expiresAt": null,
    "error": null
  }
}`,
          },
        ],
      },
    ],
  },
  {
    id: "api-keys",
    title: "API Keys",
    summary:
      "Create keys, inspect them, update restrictions or revocation state, and remove obsolete credentials.",
    actions: [
      {
        id: "list-api-keys",
        title: "List API Keys",
        summary: "Load every key visible to the current user.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Load the list with permissions and restrictions for management screens.",
            requestExample: `query ListApiKeys {
  apiKeys {
    id
    label
    expiresAt
    revoked
    rateLimit
    lastUsedAt
    permissions
    ipRestrictions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read the full key collection from the REST endpoint.",
            requestExample: `GET /rest/api-keys`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Request the key list over the interactive WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "api-keys-list",
  "action": "get_api_keys"
}`,
          },
        ],
      },
      {
        id: "api-key-detail",
        title: "Get API Key Detail",
        summary: "Inspect one key before editing it.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Read the editable state of a single key.",
            requestExample: `query ApiKeyDetail {
  apiKey(id: 5) {
    id
    label
    expiresAt
    revoked
    rateLimit
    lastUsedAt
    permissions
    ipRestrictions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read a single key from the REST resource path.",
            requestExample: `GET /rest/api-keys/5`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the detail action with the target ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "api-key-detail",
  "action": "get_api_key",
  "payload": {
    "id": 5
  }
}`,
          },
        ],
      },
      {
        id: "create-api-key",
        title: "Create API Key",
        summary:
          "Create a key with permissions, expiration, rate limit, and IP restrictions.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "The GraphQL mutation returns the newly generated secret string exactly once.",
            requestExample: `mutation CreateApiKey {
  createAPIKey(
    input: {
      label: "Grafana Import"
      expiresAt: "2026-12-31"
      rateLimit: 1000
      permissions: [
        "read time_series"
        "read kpi_results"
        "subscribe raw_data"
      ]
      ipRestrictions: [
        "10.0.0.0/24"
        "192.168.1.10"
      ]
    }
  )
}`,
            responseExample: `{
  "data": {
    "createAPIKey": "riot_live_abc123_secret_value"
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "The REST endpoint returns the created key secret or creation result payload.",
            requestExample: `POST /rest/api-keys
Content-Type: application/json

{
  "label": "Grafana Import",
  "expiresAt": "2026-12-31",
  "rateLimit": 1000,
  "permissions": [
    "read time_series",
    "read kpi_results",
    "subscribe raw_data"
  ],
  "ipRestrictions": [
    "10.0.0.0/24",
    "192.168.1.10"
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Create the key through the WebSocket action and keep the returned secret immediately.",
            requestExample: `{
  "type": "request",
  "id": "api-key-create",
  "action": "create_api_key",
  "payload": {
    "label": "Grafana Import",
    "expiresAt": "2026-12-31",
    "rateLimit": 1000,
    "permissions": [
      "read time_series",
      "read kpi_results",
      "subscribe raw_data"
    ],
    "ipRestrictions": [
      "10.0.0.0/24",
      "192.168.1.10"
    ]
  }
}`,
          },
        ],
      },
      {
        id: "update-api-key",
        title: "Update API Key",
        summary:
          "Change label, permissions, expiration, rate limit, IP restrictions, or revocation state.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Update the full editable state of an existing key.",
            requestExample: `mutation UpdateApiKey {
  updateAPIKey(
    id: 5
    input: {
      label: "Grafana Import - Read Only"
      expiresAt: "2027-01-31"
      revoked: false
      rateLimit: 500
      permissions: [
        "read time_series"
        "read kpi_results"
      ]
      ipRestrictions: [
        "10.0.0.0/24"
      ]
    }
  )
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PUT the replacement state for the API key.",
            requestExample: `PUT /rest/api-keys/5
Content-Type: application/json

{
  "label": "Grafana Import - Read Only",
  "expiresAt": "2027-01-31",
  "revoked": false,
  "rateLimit": 500,
  "permissions": [
    "read time_series",
    "read kpi_results"
  ],
  "ipRestrictions": [
    "10.0.0.0/24"
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target key ID plus the replacement input.",
            requestExample: `{
  "type": "request",
  "id": "api-key-update",
  "action": "update_api_key",
  "payload": {
    "id": 5,
    "input": {
      "label": "Grafana Import - Read Only",
      "expiresAt": "2027-01-31",
      "revoked": false,
      "rateLimit": 500,
      "permissions": [
        "read time_series",
        "read kpi_results"
      ],
      "ipRestrictions": [
        "10.0.0.0/24"
      ]
    }
  }
}`,
          },
        ],
      },
      {
        id: "delete-api-key",
        title: "Delete API Key",
        summary: "Delete a key that should no longer be used.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete by ID and verify the boolean result.",
            requestExample: `mutation DeleteApiKey {
  deleteAPIKey(id: 5)
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the key resource directly.",
            requestExample: `DELETE /rest/api-keys/5`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Issue the delete action with the key ID.",
            requestExample: `{
  "type": "request",
  "id": "api-key-delete",
  "action": "delete_api_key",
  "payload": {
    "id": 5
  }
}`,
          },
        ],
      },
    ],
  },
  {
    id: "user-config",
    title: "User Preferences",
    summary:
      "Load, store, and clear user-specific UI preferences such as favorite KPI definitions and favorite devices.",
    actions: [
      {
        id: "get-user-config",
        title: "Get User Preferences",
        summary: "Load the current user's saved UI preferences.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "This is the same data used by the dashboard and favorites logic in the frontend.",
            requestExample: `query UserPreferences {
  userConfig {
    favoriteKpis
    favoriteSdInstances
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the REST resource for a simple JSON response.",
            requestExample: `GET /rest/user-config`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Load the config via the WebSocket request action.",
            requestExample: `{
  "type": "request",
  "id": "user-config-get",
  "action": "get_user_config"
}`,
          },
        ],
      },
      {
        id: "update-user-config",
        title: "Update User Preferences",
        summary: "Save favorites or other UI-specific user settings.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Update both favorites lists in one mutation.",
            requestExample: `mutation UpdateUserPreferences {
  updateUserConfig(
    input: {
      favoriteKpis: [9, 11, 14]
      favoriteSdInstances: [17, 18]
    }
  ) {
    favoriteKpis
    favoriteSdInstances
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "POST the replacement preference payload.",
            requestExample: `POST /rest/user-config
Content-Type: application/json

{
  "favoriteKpis": [9, 11, 14],
  "favoriteSdInstances": [17, 18]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the same payload inside the WebSocket request envelope.",
            requestExample: `{
  "type": "request",
  "id": "user-config-update",
  "action": "update_user_config",
  "payload": {
    "favoriteKpis": [9, 11, 14],
    "favoriteSdInstances": [17, 18]
  }
}`,
          },
        ],
      },
      {
        id: "delete-user-config",
        title: "Delete User Preferences",
        summary: "Remove the stored preferences and fall back to defaults.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete the stored preference record.",
            requestExample: `mutation DeleteUserPreferences {
  deleteUserConfig
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the user-config resource.",
            requestExample: `DELETE /rest/user-config`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Trigger the delete request action without additional payload.",
            requestExample: `{
  "type": "request",
  "id": "user-config-delete",
  "action": "delete_user_config"
}`,
          },
        ],
      },
    ],
  },
  {
    id: "roles",
    title: "Roles And Permissions",
    summary:
      "Inspect role definitions, resolve the effective role of the current or another user, and assign a role when permitted.",
    actions: [
      {
        id: "list-roles",
        title: "List Roles",
        summary: "Load every available role and its permissions.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "The role list is useful for admin forms and permission explainers.",
            requestExample: `query Roles {
  roles {
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the REST collection endpoint for all roles.",
            requestExample: `GET /rest/user-roles`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Fetch the role catalog via the WebSocket action.",
            requestExample: `{
  "type": "request",
  "id": "roles-list",
  "action": "get_user_roles"
}`,
          },
        ],
      },
      {
        id: "get-current-role",
        title: "Get Current User Role",
        summary: "Resolve the role assigned to the current authenticated user.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Use `role` to fetch the current user's effective role.",
            requestExample: `query CurrentRole {
  role {
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "The current-user role endpoint is separate from role lookup by user ID.",
            requestExample: `GET /rest/user-roles/user`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Trigger the `get_role` action without a payload.",
            requestExample: `{
  "type": "request",
  "id": "role-current",
  "action": "get_role"
}`,
          },
        ],
      },
      {
        id: "get-user-role",
        title: "Get Another User Role",
        summary: "Resolve the assigned role for a specific user ID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Query another user's role by ID.",
            requestExample: `query UserRole {
  userRole(id: 4) {
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the user ID in the REST path.",
            requestExample: `GET /rest/user-roles/user/4`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the user ID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "role-user",
  "action": "get_user_role",
  "payload": {
    "id": 4
  }
}`,
          },
        ],
      },
      {
        id: "assign-role",
        title: "Assign Role",
        summary: "Assign a new role to a user through the admin flow.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Use a dedicated input object with user ID and role label.",
            requestExample: `mutation AssignRole {
  assignRoleToUser(
    input: {
      userID: 4
      roleLabel: "Admin"
    }
  )
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Submit the assignment payload to the update endpoint.",
            requestExample: `PUT /rest/user-roles/user
Content-Type: application/json

{
  "userID": 4,
  "roleLabel": "Admin"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the same assignment payload inside the request envelope.",
            requestExample: `{
  "type": "request",
  "id": "role-assign",
  "action": "update_user_role",
  "payload": {
    "userID": 4,
    "roleLabel": "Admin"
  }
}`,
          },
        ],
      },
    ],
  },
];

export const technologyOrder: ApiTechnology[] = [
  "graphql",
  "rest",
  "websocket",
];
