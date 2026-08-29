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
const canonicalSDInstanceUIDNote =
  "Read, filter, group and update APIs expect canonical SD instance UIDs in the `sdt:<type>.sdi:<instance>` form.";
const sdInstanceRegistrationUIDNote =
  "When an ingest or registration message carries `sdTypeUID`, `sdInstanceUID` may be the local suffix; the system stores and returns the canonical `sdt:<type>.sdi:<instance>` form.";
const kpiUIDSuffixNote =
  "For create/update inputs, `uid` is the KPI UID suffix. The system stores and returns the canonical `sdt:<type>.kpi:<kpi>` KPI UID.";

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
        summary:
          "Load the full catalog of device types including parameter metadata.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use one query to load all types and every parameter needed for editors or validation.",
            requestExample: `query ListDeviceTypes {
  sdTypes {
    uid
    label
    parameters {
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
        "uid": "sdt:thermostat",
        "label": "Thermostat",
        "parameters": [
          {
            "label": "Temperature",
            "denotation": "temperature",
            "type": "number",
            "role": "field"
          },
          {
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
      "uid": "sdt:thermostat",
      "label": "Thermostat",
      "parameters": [
        {
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
              "Resolve one type by UID for editors, previews, or device-type-aware forms.",
            requestExample: `query DeviceTypeDetail {
  sdType(uid: "sdt:thermostat") {
    uid
    label
    parameters {
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
            requestExample: `GET /rest/sd-types/sdt:thermostat
Accept: application/json`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Pass the target UID in the WebSocket payload.",
            requestExample: `{
  "type": "request",
  "id": "sd-type-detail",
  "action": "get_sd_type",
  "payload": {
    "uid": "sdt:thermostat"
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
            summary: "Delete by UID and receive a boolean success flag.",
            requestExample: `mutation DeleteDeviceType {
  deleteSDType(uid: "sdt:thermostat")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the resource directly.",
            requestExample: `DELETE /rest/sd-types/sdt:thermostat`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target UID in the delete action payload.",
            requestExample: `{
  "type": "request",
  "id": "sd-type-delete",
  "action": "delete_sd_type",
  "payload": {
    "uid": "sdt:thermostat"
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
            summary:
              "Query the full device list including type metadata for overview pages.",
            requestExample: `query ListDevices {
  sdInstances {
    uid
    label
    confirmedByUser
    userIdentifier
    type {
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
            summary:
              "The WebSocket action mirrors the REST collection endpoint.",
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
  sdInstance(uid: "sdt:thermostat.sdi:device-a-17") {
    uid
    label
    confirmedByUser
    userIdentifier
    type {
      uid
      label
      parameters {
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
            requestExample: `GET /rest/sd-instance?uid=sdt:thermostat.sdi:device-a-17`,
            notes: [
              canonicalSDInstanceUIDNote,
              "This endpoint uses a query parameter instead of a path parameter in the current backend.",
            ],
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Pass the device UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "device-detail",
  "action": "get_sd_instance",
  "payload": {
    "uid": "sdt:thermostat.sdi:device-a-17"
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
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
  byType: sdInstancesByType(uid: "sdt:thermostat") {
    label
    uid
  }
  byKpi: sdInstancesByKpiDefinition(uid: "sdt:thermostat.kpi:overheat_alarm") {
    label
    uid
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "REST exposes separate endpoints for type-based and KPI-based filtering.",
            requestExample: `GET /rest/sd-instances/type/sdt:thermostat
GET /rest/sd-instances/kpi/sdt:thermostat.kpi:overheat_alarm`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Use one action per filter mode and send the selected UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "devices-by-type",
  "action": "get_sd_instances_by_type",
  "payload": { "uid": "sdt:thermostat" }
}

{
  "type": "request",
  "id": "devices-by-kpi",
  "action": "get_sd_instances_by_kpi",
  "payload": { "uid": "sdt:thermostat.kpi:overheat_alarm" }
}`,
            notes: [canonicalSDInstanceUIDNote],
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
    uid: "sdt:thermostat.sdi:device-a-17"
    input: {
      label: "Boiler Sensor A-17"
      userIdentifier: "asset-plant-a-17"
      confirmedByUser: true
    }
  ) {
    uid
    label
    userIdentifier
    confirmedByUser
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
          },
          {
            technology: "rest",
            label: rest,
            summary: "PATCH the mutable fields of one device.",
            requestExample: `PATCH /rest/sd-instances/sdt:thermostat.sdi:device-a-17
Content-Type: application/json

{
  "label": "Boiler Sensor A-17",
  "userIdentifier": "asset-plant-a-17",
  "confirmedByUser": true
}`,
            notes: [canonicalSDInstanceUIDNote],
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the device UID and the update input in one message.",
            requestExample: `{
  "type": "request",
  "id": "device-update",
  "action": "update_sd_instance",
  "payload": {
    "uid": "sdt:thermostat.sdi:device-a-17",
    "input": {
      "label": "Boiler Sensor A-17",
      "userIdentifier": "asset-plant-a-17",
      "confirmedByUser": true
    }
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
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
            summary:
              "Subscribe directly to registration events with an optional filter.",
            requestExample: `subscription DeviceRegistrationStream {
  onSDInstanceRegistered(filter: { sdTypeUID: "sdt:thermostat" }) {
    uid
    label
    confirmedByUser
  }
}`,
            notes: [sdInstanceRegistrationUIDNote],
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
    "sdTypeUID": "sdt:thermostat"
  }
}`,
            notes: [sdInstanceRegistrationUIDNote],
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
  "sdTypeUID": "sdt:thermostat"
}`,
            notes: [sdInstanceRegistrationUIDNote],
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
              "Use the SD type UID to fetch the latest raw data points across all matching device instances.",
            requestExample: `query RawDataByType {
  rawDataPointsBySDType(uid: "sdt:thermostat") {
    sdTypeUID
    sdInstanceUID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "rawDataPointsBySDType": [
      {
        "sdTypeUID": "sdt:thermostat",
        "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
        "payload": {
          "temperature": 82.4,
          "isOnline": true,
          "location": "sdt:thermostat.sdi:boiler-room-a"
        },
        "eventTime": "2026-04-14T12:30:00Z"
      },
      {
        "sdTypeUID": "sdt:thermostat",
        "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-b",
        "payload": {
          "temperature": 79.1,
          "isOnline": true,
          "location": "sdt:thermostat.sdi:boiler-room-b"
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
              "Read the latest raw data set for one SD type UID from the REST path-based endpoint.",
            requestExample: `GET /rest/raw/type/sdt:thermostat
Accept: application/json
X-API-Key: riot_live_abc123_secret_value`,
            responseExample: `[
  {
    "sdTypeUID": "sdt:thermostat",
    "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
    "payload": {
      "temperature": 82.4,
      "isOnline": true,
      "location": "sdt:thermostat.sdi:boiler-room-a"
    },
    "eventTime": "2026-04-14T12:30:00Z"
  },
  {
    "sdTypeUID": "sdt:thermostat",
    "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-b",
    "payload": {
      "temperature": 79.1,
      "isOnline": true,
      "location": "sdt:thermostat.sdi:boiler-room-b"
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
    "uid": "sdt:thermostat"
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "raw-by-type",
  "success": true,
  "payload": [
    {
      "sdTypeUID": "sdt:thermostat",
      "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
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
              "Query the latest raw data point for one device instance by its UID.",
            requestExample: `query LatestRawDataPoint {
  rawDataPoint(uid: "sdt:thermostat.sdi:boiler-room-a") {
    sdTypeUID
    sdInstanceUID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "rawDataPoint": {
      "sdTypeUID": "sdt:thermostat",
      "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
      "payload": {
        "temperature": 82.4,
        "isOnline": true,
        "location": "sdt:thermostat.sdi:boiler-room-a"
      },
      "eventTime": "2026-04-14T12:30:00Z"
    }
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Use the single raw-data REST endpoint with the device instance UID.",
            requestExample: `GET /rest/raw?uid=sdt:thermostat.sdi:boiler-room-a
Accept: application/json
X-API-Key: riot_live_abc123_secret_value`,
            responseExample: `{
  "sdTypeUID": "sdt:thermostat",
  "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
  "payload": {
    "temperature": 82.4,
    "isOnline": true,
    "location": "sdt:thermostat.sdi:boiler-room-a"
  },
  "eventTime": "2026-04-14T12:30:00Z"
}`,
            notes: [
              canonicalSDInstanceUIDNote,
              "The current REST router exposes this endpoint as `/rest/raw` and expects the instance UID in the `uid` query parameter.",
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
    "uid": "sdt:thermostat.sdi:boiler-room-a"
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "raw-by-instance",
  "success": true,
  "payload": {
    "sdTypeUID": "sdt:thermostat",
    "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
    "payload": {
      "temperature": 82.4,
      "isOnline": true
    },
    "eventTime": "2026-04-14T12:30:00Z"
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
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
  onRawDataPointArrived(filter: { sdTypeUIDs: ["sdt:thermostat"], sdInstanceUIDs: ["sdt:thermostat.sdi:boiler-room-a"] }) {
    sdTypeUID
    sdInstanceUID
    payload
    eventTime
  }
}`,
            responseExample: `{
  "data": {
    "onRawDataPointArrived": [
      {
        "sdTypeUID": "sdt:thermostat",
        "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
        "payload": {
          "temperature": 83.2,
          "isOnline": true
        },
        "eventTime": "2026-04-14T12:31:05Z"
      }
    ]
  }
}`,
            notes: [canonicalSDInstanceUIDNote],
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
  "sdTypeUIDs": ["sdt:thermostat"],
  "sdInstanceUIDs": ["sdt:thermostat.sdi:boiler-room-a"]
}`,
            responseExample: `event: raw_data_point
data: [{"sdTypeUID":"sdt:thermostat","sdInstanceUID":"sdt:thermostat.sdi:boiler-room-a","payload":{"temperature":83.2,"isOnline":true},"eventTime":"2026-04-14T12:31:05Z"}]`,
            notes: [canonicalSDInstanceUIDNote],
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
    "sdTypeUIDs": ["sdt:thermostat"],
    "sdInstanceUIDs": ["sdt:thermostat.sdi:boiler-room-a"]
  }
}`,
            responseExample: `{
  "type": "event",
  "id": "raw-data-stream",
  "topic": "raw_data_point",
  "payload": [
    {
      "sdTypeUID": "sdt:thermostat",
      "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a",
      "payload": {
        "temperature": 83.2,
        "isOnline": true
      },
      "eventTime": "2026-04-14T12:31:05Z"
    }
  ]
}`,
            notes: [canonicalSDInstanceUIDNote],
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
    uid
    label
    sdTypeUID
    userIdentifier
    sdInstanceMode
    selectedSDInstanceUIDs
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
  kpiDefinition(uid: "sdt:thermostat.kpi:overheat_alarm") {
    uid
    label
    sdTypeUID
    sdInstanceMode
    selectedSDInstanceUIDs
    nodes {
      id
      parentNodeID
      nodeType
      ... on LogicalOperationKPINode {
        type
      }
      ... on NumericGTAtomKPINode {
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
            requestExample: `GET /rest/kpi-definitions/sdt:thermostat.kpi:overheat_alarm`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target KPI UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "kpi-detail",
  "action": "get_kpi_definition",
  "payload": {
    "uid": "sdt:thermostat.kpi:overheat_alarm"
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
  byType: kpiDefinitionsBySdType(uid: "sdt:thermostat") {
    uid
    label
  }
  byDevice: kpiDefinitionsBySdInstance(uid: "sdt:thermostat.sdi:device-a-17") {
    uid
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "REST exposes separate endpoints for both filter dimensions.",
            requestExample: `GET /rest/kpi-definitions/type/sdt:thermostat
GET /rest/kpi-definitions/instance/sdt:thermostat.sdi:device-a-17`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use one request action per filter type.",
            requestExample: `{
  "type": "request",
  "id": "kpi-by-type",
  "action": "get_kpi_definitions_by_type",
    "payload": { "uid": "sdt:thermostat" }
}

{
  "type": "request",
  "id": "kpi-by-device",
  "action": "get_kpi_definitions_by_sd_instance",
  "payload": { "uid": "sdt:thermostat.sdi:device-a-17" }
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
              "This example creates a KPI for one device type and targets selected instances. The input UID is a suffix.",
            requestExample: `mutation CreateKpi {
  createKPIDefinition(
    input: {
      uid: "overheat_alarm"
      label: "Overheat Alarm"
      sdTypeUID: "sdt:thermostat"
      userIdentifier: "user-42"
      sdInstanceMode: selected
      selectedSDInstanceUIDs: ["sdt:thermostat.sdi:thermostat-a-17", "sdt:thermostat.sdi:thermostat-a-18"]
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
          sdParameterSpecification: "temperature"
          numericReferenceValue: 85
        }
        {
          id: 3
          parentNodeID: 1
          nodeType: booleanEQAtom
          sdParameterSpecification: "isOnline"
          booleanReferenceValue: true
        }
      ]
    }
  ) {
    uid
    label
  }
}`,
            responseExample: `{
  "data": {
    "createKPIDefinition": {
      "uid": "sdt:thermostat.kpi:overheat_alarm",
      "label": "Overheat Alarm"
    }
  }
}`,
            notes: [kpiUIDSuffixNote, canonicalSDInstanceUIDNote],
          },
          {
            technology: "rest",
            label: rest,
            summary: "Send the same rule tree in JSON to the create endpoint.",
            requestExample: `POST /rest/kpi-definitions
Content-Type: application/json

{
  "uid": "overheat_alarm",
  "label": "Overheat Alarm",
  "sdTypeUID": "sdt:thermostat",
  "userIdentifier": "user-42",
  "sdInstanceMode": "selected",
  "selectedSDInstanceUIDs": ["sdt:thermostat.sdi:thermostat-a-17", "sdt:thermostat.sdi:thermostat-a-18"],
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
      "sdParameterSpecification": "temperature",
      "numericReferenceValue": 85
    }
  ]
}`,
            responseExample: `{
  "uid": "sdt:thermostat.kpi:overheat_alarm",
  "label": "Overheat Alarm"
}`,
            notes: [kpiUIDSuffixNote, canonicalSDInstanceUIDNote],
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
    "uid": "overheat_alarm",
    "label": "Overheat Alarm",
    "sdTypeUID": "sdt:thermostat",
    "userIdentifier": "user-42",
    "sdInstanceMode": "selected",
    "selectedSDInstanceUIDs": ["sdt:thermostat.sdi:thermostat-a-17", "sdt:thermostat.sdi:thermostat-a-18"],
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
            responseExample: `{
  "type": "response",
  "id": "kpi-create",
  "success": true,
  "payload": {
    "uid": "sdt:thermostat.kpi:overheat_alarm",
    "label": "Overheat Alarm"
  }
}`,
            notes: [kpiUIDSuffixNote, canonicalSDInstanceUIDNote],
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
              "The update mutation uses the canonical KPI UID in the target argument and the UID suffix in the replacement input.",
            requestExample: `mutation UpdateKpi {
  updateKPIDefinition(
    uid: "sdt:thermostat.kpi:overheat_alarm"
    input: {
      uid: "overheat_alarm"
      label: "Overheat Alarm v2"
      sdTypeUID: "sdt:thermostat"
      userIdentifier: "user-42"
      sdInstanceMode: all
      selectedSDInstanceUIDs: []
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
          sdParameterSpecification: "temperature"
          numericReferenceValue: 90
        }
      ]
    }
  ) {
    uid
    label
    sdInstanceMode
  }
}`,
            notes: [kpiUIDSuffixNote],
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "PUT the complete definition to keep server and editor state aligned.",
            requestExample: `PUT /rest/kpi-definitions/sdt:thermostat.kpi:overheat_alarm
Content-Type: application/json

{
  "uid": "overheat_alarm",
  "label": "Overheat Alarm v2",
  "sdTypeUID": "sdt:thermostat",
  "userIdentifier": "user-42",
  "sdInstanceMode": "all",
  "selectedSDInstanceUIDs": [],
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
      "sdParameterSpecification": "temperature",
      "numericReferenceValue": 90
    }
  ]
}`,
            responseExample: `{
  "uid": "sdt:thermostat.kpi:overheat_alarm",
  "label": "Overheat Alarm v2",
  "sdInstanceMode": "all"
}`,
            notes: [kpiUIDSuffixNote],
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target UID plus the full replacement input.",
            requestExample: `{
  "type": "request",
  "id": "kpi-update",
  "action": "update_kpi_definition",
  "payload": {
    "uid": "sdt:thermostat.kpi:overheat_alarm",
    "input": {
      "uid": "overheat_alarm",
      "label": "Overheat Alarm v2",
      "sdTypeUID": "sdt:thermostat",
      "userIdentifier": "user-42",
      "sdInstanceMode": "all",
      "selectedSDInstanceUIDs": [],
      "nodes": []
    }
  }
}`,
            notes: [kpiUIDSuffixNote],
          },
        ],
      },
      {
        id: "delete-kpi",
        title: "Delete KPI Definition",
        summary: "Delete an obsolete KPI definition by UID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Delete through a simple mutation and check the boolean result.",
            requestExample: `mutation DeleteKpi {
  deleteKPIDefinition(uid: "sdt:thermostat.kpi:overheat_alarm")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the KPI resource directly.",
            requestExample: `DELETE /rest/kpi-definitions/sdt:thermostat.kpi:overheat_alarm`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the delete action with the KPI UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "kpi-delete",
  "action": "delete_kpi_definition",
  "payload": {
    "uid": "sdt:thermostat.kpi:overheat_alarm"
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
            summary:
              "Fetch all results or include only the fields needed by the current view.",
            requestExample: `query ListKpiResults {
  kpiResults {
    kpiDefinitionUID
    sdTypeUID
    sdInstanceUID
    fulfilled
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
  kpiResultsByKPI(uid: "sdt:thermostat.kpi:room_temperature") {
    kpiDefinitionUID
    sdInstanceUID
    fulfilled
    eventTime
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the KPI UID as a path parameter.",
            requestExample: `GET /rest/kpi-results/sdt:thermostat.kpi:room_temperature`,
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
	    "uid": "sdt:thermostat.kpi:room_temperature"
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
              "Use a UID request object when resolving a specific KPI result.",
            requestExample: `query OneKpiResult {
  kpiResult(
    request: {
      kpiDefinitionUID: "sdt:thermostat.kpi:room_temperature"
      sdInstanceUID: "sdt:thermostat.sdi:boiler-room-a"
    }
  ) {
    kpiDefinitionUID
    sdInstanceUID
    fulfilled
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
  "kpiDefinitionUID": "sdt:thermostat.kpi:room_temperature",
  "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a"
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
    "kpiDefinitionUID": "sdt:thermostat.kpi:room_temperature",
    "sdInstanceUID": "sdt:thermostat.sdi:boiler-room-a"
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
            summary:
              "Subscribe directly from Apollo or another GraphQL client.",
            requestExample: `subscription KpiResultStream {
  onKPIFulfillmentChecked(filter: { kpiDefinitionUIDs: ["sdt:thermostat.kpi:room_temperature"], sdInstanceUIDs: ["sdt:thermostat.sdi:boiler-room-a"] }) {
    kpiDefinitionUID
    sdInstanceUID
    fulfilled
    eventTime
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
    "kpiDefinitionUIDs": ["sdt:thermostat.kpi:room_temperature"],
    "sdInstanceUIDs": ["sdt:thermostat.sdi:boiler-room-a"]
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
  "kpiDefinitionUIDs": ["sdt:thermostat.kpi:room_temperature"],
  "sdInstanceUIDs": ["sdt:thermostat.sdi:boiler-room-a"]
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
            summary: "Load the group collection with its member UIDs.",
            requestExample: `query ListGroups {
  sdInstanceGroups {
    uid
    label
    sdInstanceUIDs
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
            summary: "Read one group by UID.",
            requestExample: `query GroupDetail {
  sdInstanceGroup(uid: "grp:boiler_room") {
    uid
    label
    sdInstanceUIDs
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use a path parameter to fetch one group.",
            requestExample: `GET /rest/sd-instance-groups/grp:boiler_room`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the group UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "group-detail",
  "action": "get_sd_instance_group",
  "payload": { "uid": "grp:boiler_room" }
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
	      uid: "grp:boiler_room"
	      label: "Boiler Room"
	      userIdentifier: "BR"
	      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-c"]
	    }
	  ) {
	    uid
	    label
	    sdInstanceUIDs
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
  "uid": "grp:boiler_room",
  "label": "Boiler Room",
  "userIdentifier": "BR",
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-c"]
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
    "uid": "grp:boiler_room",
    "label": "Boiler Room",
    "userIdentifier": "BR",
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-c"]
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
	    uid: "grp:boiler_room"
	    input: {
	      uid: "grp:boiler_room"
	      label: "Boiler Room - North"
	      userIdentifier: "BR-N"
	      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-north"]
	    }
	  ) {
	    uid
	    label
	    sdInstanceUIDs
	  }
	}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PUT the full group payload.",
            requestExample: `PUT /rest/sd-instance-groups/grp:boiler_room
Content-Type: application/json

{
  "uid": "grp:boiler_room",
  "label": "Boiler Room - North",
  "userIdentifier": "BR-N",
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-north"]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Provide the group UID together with the replacement input.",
            requestExample: `{
  "type": "request",
  "id": "group-update",
  "action": "update_sd_instance_group",
  "payload": {
    "uid": "grp:boiler_room",
    "input": {
      "uid": "grp:boiler_room",
      "label": "Boiler Room - North",
      "userIdentifier": "BR-N",
      "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-a", "sdt:boiler.sdi:boiler-b", "sdt:boiler.sdi:boiler-north"]
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
            summary: "Delete by group UID and check the boolean result.",
            requestExample: `mutation DeleteGroup {
  deleteSDInstanceGroup(uid: "grp:boiler_room")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the group resource directly.",
            requestExample: `DELETE /rest/sd-instance-groups/grp:boiler_room`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Issue the delete action with the group UID.",
            requestExample: `{
  "type": "request",
  "id": "group-delete",
  "action": "delete_sd_instance_group",
  "payload": {
    "uid": "grp:boiler_room"
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
      sdTypeUID: "sdt:boiler"
      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-sensor-a-17"]
      kpiDefinitionUIDs: ["sdt:boiler.kpi:overheat"]
      sortDesc: true
      limit: 50
      cursor: {
        time: "2026-04-14T12:00:00Z"
        sdInstanceUID: "sdt:boiler.sdi:boiler-sensor-a-17"
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
      kpiDefinitionUID
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
  "sdTypeUID": "sdt:boiler",
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17"],
  "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat"],
  "sortDesc": true,
  "limit": 50,
  "cursor": {
    "time": "2026-04-14T12:00:00Z",
    "sdInstanceUID": "sdt:boiler.sdi:boiler-sensor-a-17"
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
    "sdTypeUID": "sdt:boiler",
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17"],
    "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat"],
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
      sdTypeUID: "sdt:boiler"
      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"]
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
  "sdTypeUID": "sdt:boiler",
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
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
    "sdTypeUID": "sdt:boiler",
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
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
      sdTypeUID: "sdt:boiler"
      kpiDefinitionUIDs: ["sdt:boiler.kpi:overheat", "sdt:boiler.kpi:pressure-drop"]
      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"]
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
      kpiDefinitionUID
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
  "sdTypeUID": "sdt:boiler",
  "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat", "sdt:boiler.kpi:pressure-drop"],
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
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
    "sdTypeUID": "sdt:boiler",
    "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat", "sdt:boiler.kpi:pressure-drop"],
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
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
      sdTypeUID: "sdt:boiler"
      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"]
      sortDesc: true
      limit: 1000
    }
  ) {
    uid
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
      "uid": "exp:J4K8M2Q9R7T6W3P1",
      "status": "pending",
      "downloadUrl": null,
      "createdAt": "2026-05-02T23:53:23.626667429Z",
      "expiresAt": null,
      "error": null
    }
  }
}`,
            notes: [
              "Use the returned uid for status checks, cancellation, or subscriptions.",
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
  "sdTypeUID": "sdt:boiler",
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
  "sortDesc": true,
  "limit": 1000
}`,
            responseExample: `{
  "uid": "exp:J4K8M2Q9R7T6W3P1",
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
    "sdTypeUID": "sdt:boiler",
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
    "limit": 1000
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export",
  "success": true,
  "payload": {
    "uid": "exp:J4K8M2Q9R7T6W3P1",
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
              "Use the export UID to fetch a single snapshot that can be pending, processing, done, failed, cancelled, or expired.",
            requestExample: `query ExportStatus {
  timeSeriesExport(uid: "exp:J4K8M2Q9R7T6W3P1") {
    uid
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
      "uid": "exp:J4K8M2Q9R7T6W3P1",
      "status": "done",
      "downloadUrl": "/rest/time-series/export/exp:J4K8M2Q9R7T6W3P1",
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
            requestExample: `GET /rest/time-series/export/exp:J4K8M2Q9R7T6W3P1/status`,
            responseExample: `{
  "uid": "exp:J4K8M2Q9R7T6W3P1",
  "status": "done",
  "downloadUrl": "/rest/time-series/export/exp:J4K8M2Q9R7T6W3P1",
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
    "uid": "exp:J4K8M2Q9R7T6W3P1"
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-status",
  "success": true,
  "payload": {
    "uid": "exp:J4K8M2Q9R7T6W3P1",
    "status": "done",
    "downloadUrl": "/rest/time-series/export/exp:J4K8M2Q9R7T6W3P1",
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
  cancelTimeSeriesExport(uid: "exp:J4K8M2Q9R7T6W3P1") {
    uid
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
      "uid": "exp:J4K8M2Q9R7T6W3P1",
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
            requestExample: `DELETE /rest/time-series/export/exp:J4K8M2Q9R7T6W3P1`,
            responseExample: `{
  "uid": "exp:J4K8M2Q9R7T6W3P1",
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
              "The low-level WebSocket API uses a dedicated cancel action with the export uid in the payload.",
            requestExample: `{
  "type": "request",
  "id": "history-export-cancel",
  "action": "time-series-export-cancel",
  "payload": {
    "uid": "exp:J4K8M2Q9R7T6W3P1"
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-cancel",
  "success": true,
  "payload": {
    "uid": "exp:J4K8M2Q9R7T6W3P1",
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
              "Apollo-style GraphQL subscriptions stream the same TimeSeriesExport payload filtered by export UIDs.",
            requestExample: `subscription OnExportUpdated {
  onTimeSeriesExportUpdated(filter: { uids: ["exp:J4K8M2Q9R7T6W3P1"] }) {
    uid
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
      "uid": "exp:J4K8M2Q9R7T6W3P1",
      "status": "done",
      "downloadUrl": "/rest/time-series/export/exp:J4K8M2Q9R7T6W3P1",
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
    "uids": ["exp:J4K8M2Q9R7T6W3P1"]
  }
}`,
            responseExample: `{
  "type": "event",
  "id": "history-export-events",
  "topic": "time_series_export_updated",
  "payload": {
    "uid": "exp:J4K8M2Q9R7T6W3P1",
    "status": "done",
    "downloadUrl": "/rest/time-series/export/exp:J4K8M2Q9R7T6W3P1",
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
            requestExample: `GET /rest/time-series/export/exp:J4K8M2Q9R7T6W3P1`,
            responseExample: `HTTP/1.1 200 OK
Content-Type: text/csv
Content-Disposition: attachment; filename="time-series-export-41.csv"

time,sdInstanceUID,temperature
2026-05-02T23:53:24Z,sdt:boiler.sdi:boiler-sensor-a-17,71.2
2026-05-02T23:53:25Z,sdt:boiler.sdi:boiler-sensor-a-17,71.3`,
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
      sdTypeUID: "sdt:boiler"
      kpiDefinitionUIDs: ["sdt:boiler.kpi:overheat"]
      sdInstanceUIDs: ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"]
      from: "2026-04-01T00:00:00Z"
      to: "2026-04-14T23:59:59Z"
      aggregateSeconds: 3600
    }
  ) {
    uid
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
      "uid": "exp:T8P1K6R3M9Q2W5D4",
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
  "sdTypeUID": "sdt:boiler",
  "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat"],
  "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
  "from": "2026-04-01T00:00:00Z",
  "to": "2026-04-14T23:59:59Z",
  "aggregateSeconds": 3600
}`,
            responseExample: `{
  "uid": "exp:T8P1K6R3M9Q2W5D4",
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
    "sdTypeUID": "sdt:boiler",
    "kpiDefinitionUIDs": ["sdt:boiler.kpi:overheat"],
    "sdInstanceUIDs": ["sdt:boiler.sdi:boiler-sensor-a-17", "sdt:boiler.sdi:boiler-sensor-a-18"],
    "aggregateSeconds": 3600
  }
}`,
            responseExample: `{
  "type": "response",
  "id": "history-export-aggregate",
  "success": true,
  "payload": {
    "uid": "exp:T8P1K6R3M9Q2W5D4",
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
    id: "users",
    title: "Users And Sessions",
    summary:
      "Inspect users, update account metadata, disable or re-enable accounts, and revoke active sessions.",
    actions: [
      {
        id: "list-users",
        title: "List Users",
        summary: "Load all user accounts visible to the current administrator.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Read users with role and disabled state in one query.",
            requestExample: `query Users {
  users {
    uid
    email
    name
    disabled
    disabledAt
    disabledReason
    role {
      uid
      label
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Read the user collection from the REST endpoint.",
            requestExample: `GET /rest/users`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Request the user list over the interactive WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "users-list",
  "action": "get_users"
}`,
          },
        ],
      },
      {
        id: "get-user",
        title: "Get User",
        summary: "Inspect one account by public user UID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Resolve one user by UID.",
            requestExample: `query UserDetail {
  user(uid: "usr:H8P2R5T9K1M4Q7D6") {
    uid
    email
    name
    disabled
    role {
      uid
      label
    }
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the user UID in the REST path.",
            requestExample: `GET /rest/users/usr:H8P2R5T9K1M4Q7D6`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target user UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "user-detail",
  "action": "get_user",
  "payload": {
    "uid": "usr:H8P2R5T9K1M4Q7D6"
  }
}`,
          },
        ],
      },
      {
        id: "update-user",
        title: "Update User",
        summary: "Update editable user profile fields.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Patch the user through the dedicated update mutation.",
            requestExample: `mutation UpdateUser {
  updateUser(
    uid: "usr:H8P2R5T9K1M4Q7D6"
    input: {
      name: "Ops Console User"
      profileImageURL: "https://example.test/avatar.png"
    }
  ) {
    uid
    name
    profileImageURL
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PATCH the editable fields on the user resource.",
            requestExample: `PATCH /rest/users/usr:H8P2R5T9K1M4Q7D6
Content-Type: application/json

{
  "name": "Ops Console User",
  "profileImageURL": "https://example.test/avatar.png"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the user UID and update input in the payload.",
            requestExample: `{
  "type": "request",
  "id": "user-update",
  "action": "update_user",
  "payload": {
    "uid": "usr:H8P2R5T9K1M4Q7D6",
    "input": {
      "name": "Ops Console User",
      "profileImageURL": "https://example.test/avatar.png"
    }
  }
}`,
          },
        ],
      },
      {
        id: "disable-user",
        title: "Disable User",
        summary: "Block future authentication for a user and store the reason.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Disable the account by UID.",
            requestExample: `mutation DisableUser {
  disableUser(
    uid: "usr:H8P2R5T9K1M4Q7D6"
    reason: "Access no longer required"
  ) {
    uid
    disabled
    disabledReason
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "POST the disable command to the user resource.",
            requestExample: `POST /rest/users/usr:H8P2R5T9K1M4Q7D6/disable
Content-Type: application/json

{
  "reason": "Access no longer required"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the disable command through the WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "user-disable",
  "action": "disable_user",
  "payload": {
    "uid": "usr:H8P2R5T9K1M4Q7D6",
    "reason": "Access no longer required"
  }
}`,
          },
        ],
      },
      {
        id: "enable-user",
        title: "Enable User",
        summary: "Allow a previously disabled account to authenticate again.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Enable the account by UID.",
            requestExample: `mutation EnableUser {
  enableUser(uid: "usr:H8P2R5T9K1M4Q7D6") {
    uid
    disabled
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "POST the enable command to the user resource.",
            requestExample: `POST /rest/users/usr:H8P2R5T9K1M4Q7D6/enable`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the enable command through the WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "user-enable",
  "action": "enable_user",
  "payload": {
    "uid": "usr:H8P2R5T9K1M4Q7D6"
  }
}`,
          },
        ],
      },
      {
        id: "list-sessions",
        title: "List Sessions",
        summary:
          "Inspect active sessions for the current user or a selected user.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Use `sessions` for own sessions and `sessionsByUser` for admin lookup.",
            requestExample: `query Sessions {
  sessions {
    uid
    userUID
    expiresAt
    revoked
    createdAt
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "Read own sessions or path-based sessions for a selected user.",
            requestExample: `GET /rest/sessions

GET /rest/users/usr:H8P2R5T9K1M4Q7D6/sessions`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Use the matching WebSocket session actions.",
            requestExample: `{
  "type": "request",
  "id": "sessions-list",
  "action": "get_sessions"
}`,
          },
        ],
      },
      {
        id: "revoke-session",
        title: "Revoke Session",
        summary: "Revoke one session or all sessions for a selected user.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary:
              "Revoke by session UID or revoke all sessions for a user UID.",
            requestExample: `mutation RevokeSession {
  revokeSession(uid: "ses:M8Q2W5D9R1T6P4K7")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete one session or revoke all sessions for a user.",
            requestExample: `DELETE /rest/sessions/ses:M8Q2W5D9R1T6P4K7

POST /rest/users/usr:H8P2R5T9K1M4Q7D6/sessions/revoke`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the session UID in the revoke action payload.",
            requestExample: `{
  "type": "request",
  "id": "session-revoke",
  "action": "revoke_session",
  "payload": {
    "uid": "ses:M8Q2W5D9R1T6P4K7"
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
            summary:
              "Load the list with permissions and restrictions for management screens.",
            requestExample: `query ListApiKeys {
  apiKeys {
    uid
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
  apiKey(uid: "ak:F7M2Q9R4T8K1P6D3") {
    uid
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
            requestExample: `GET /rest/api-keys/ak:F7M2Q9R4T8K1P6D3`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Use the detail action with the target UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "api-key-detail",
  "action": "get_api_key",
  "payload": {
    "uid": "ak:F7M2Q9R4T8K1P6D3"
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
        "time_series.read"
        "kpi_results.read"
        "raw_data.subscribe"
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
            summary:
              "The REST endpoint returns the created key secret or creation result payload.",
            requestExample: `POST /rest/api-keys
Content-Type: application/json

{
  "label": "Grafana Import",
  "expiresAt": "2026-12-31",
  "rateLimit": 1000,
  "permissions": [
    "time_series.read",
    "kpi_results.read",
    "raw_data.subscribe"
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
            summary:
              "Create the key through the WebSocket action and keep the returned secret immediately.",
            requestExample: `{
  "type": "request",
  "id": "api-key-create",
  "action": "create_api_key",
  "payload": {
    "label": "Grafana Import",
    "expiresAt": "2026-12-31",
    "rateLimit": 1000,
    "permissions": [
      "time_series.read",
      "kpi_results.read",
      "raw_data.subscribe"
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
    uid: "ak:F7M2Q9R4T8K1P6D3"
    input: {
      label: "Grafana Import - Read Only"
      expiresAt: "2027-01-31"
      revoked: false
      rateLimit: 500
      permissions: [
        "time_series.read"
        "kpi_results.read"
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
            requestExample: `PUT /rest/api-keys/ak:F7M2Q9R4T8K1P6D3
Content-Type: application/json

{
  "label": "Grafana Import - Read Only",
  "expiresAt": "2027-01-31",
  "revoked": false,
  "rateLimit": 500,
  "permissions": [
    "time_series.read",
    "kpi_results.read"
  ],
  "ipRestrictions": [
    "10.0.0.0/24"
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the target key UID plus the replacement input.",
            requestExample: `{
  "type": "request",
  "id": "api-key-update",
  "action": "update_api_key",
  "payload": {
    "uid": "ak:F7M2Q9R4T8K1P6D3",
    "input": {
      "label": "Grafana Import - Read Only",
      "expiresAt": "2027-01-31",
      "revoked": false,
      "rateLimit": 500,
      "permissions": [
        "time_series.read",
        "kpi_results.read"
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
            summary: "Delete by UID and verify the boolean result.",
            requestExample: `mutation DeleteApiKey {
  deleteAPIKey(uid: "ak:F7M2Q9R4T8K1P6D3")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the key resource directly.",
            requestExample: `DELETE /rest/api-keys/ak:F7M2Q9R4T8K1P6D3`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Issue the delete action with the key UID.",
            requestExample: `{
  "type": "request",
  "id": "api-key-delete",
  "action": "delete_api_key",
  "payload": {
    "uid": "ak:F7M2Q9R4T8K1P6D3"
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
            summary:
              "This is the same data used by the dashboard and favorites logic in the frontend.",
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
            summary:
              "Use the same payload inside the WebSocket request envelope.",
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
            summary:
              "Trigger the delete request action without additional payload.",
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
            summary:
              "The role list is useful for admin forms and permission explainers.",
            requestExample: `query Roles {
  roles {
    uid
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the REST collection endpoint for all roles.",
            requestExample: `GET /rest/roles`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Fetch the role catalog via the WebSocket action.",
            requestExample: `{
  "type": "request",
  "id": "roles-list",
  "action": "get_roles"
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
    uid
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary:
              "The current-user role endpoint is separate from role lookup by user UID.",
            requestExample: `GET /rest/roles/current`,
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
        summary: "Resolve the assigned role for a specific user UID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Query another user's role by UID.",
            requestExample: `query UserRole {
  userRole(uid: "usr:H8P2R5T9K1M4Q7D6") {
    uid
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the user UID in the REST path.",
            requestExample: `GET /rest/users/usr:H8P2R5T9K1M4Q7D6/role`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the user UID in the payload.",
            requestExample: `{
  "type": "request",
  "id": "role-user",
  "action": "get_user_role",
  "payload": {
    "uid": "usr:H8P2R5T9K1M4Q7D6"
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
            summary: "Use a dedicated input object with user UID and role UID.",
            requestExample: `mutation AssignRole {
  assignRoleToUser(
    input: {
      userUID: "usr:H8P2R5T9K1M4Q7D6"
      roleUID: "role:P9K2M6Q1T8R4D7W3"
    }
  )
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Submit the assignment payload to the update endpoint.",
            requestExample: `PUT /rest/users/usr:H8P2R5T9K1M4Q7D6/role
Content-Type: application/json

{
  "roleUID": "role:P9K2M6Q1T8R4D7W3"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary:
              "Send the same assignment payload inside the request envelope.",
            requestExample: `{
  "type": "request",
  "id": "role-assign",
  "action": "update_user_role",
  "payload": {
    "userUID": "usr:H8P2R5T9K1M4Q7D6",
    "roleUID": "role:P9K2M6Q1T8R4D7W3"
  }
}`,
          },
        ],
      },
      {
        id: "list-permissions",
        title: "List Permissions",
        summary: "Load the available permission catalog.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Read every permission UID and its editable label.",
            requestExample: `query Permissions {
  permissions {
    uid
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Use the REST permission catalog endpoint.",
            requestExample: `GET /rest/permissions`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Fetch permissions through the WebSocket API.",
            requestExample: `{
  "type": "request",
  "id": "permissions-list",
  "action": "get_permissions"
}`,
          },
        ],
      },
      {
        id: "create-role",
        title: "Create Role",
        summary:
          "Create a custom role from a label and selected permission UIDs.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Create a non-system role.",
            requestExample: `mutation CreateRole {
  createRole(
    input: {
      label: "Read Only Integrations"
      permissionUIDs: [
        "time_series.read"
        "kpi_results.read"
      ]
    }
  ) {
    uid
    label
    system
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "POST the new role definition to the role collection.",
            requestExample: `POST /rest/roles
Content-Type: application/json

{
  "label": "Read Only Integrations",
  "permissionUIDs": [
    "time_series.read",
    "kpi_results.read"
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Create the role through the WebSocket action.",
            requestExample: `{
  "type": "request",
  "id": "role-create",
  "action": "create_role",
  "payload": {
    "label": "Read Only Integrations",
    "permissionUIDs": [
      "time_series.read",
      "kpi_results.read"
    ]
  }
}`,
          },
        ],
      },
      {
        id: "update-role",
        title: "Update Role",
        summary:
          "Replace the editable label and permission set for a custom role.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Update a role by UID.",
            requestExample: `mutation UpdateRole {
  updateRole(
    uid: "role:P9K2M6Q1T8R4D7W3"
    input: {
      label: "Operations Read Only"
      permissionUIDs: [
        "sd_types.read"
        "sd_instances.read"
        "time_series.read"
      ]
    }
  ) {
    uid
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PUT the replacement role definition.",
            requestExample: `PUT /rest/roles/role:P9K2M6Q1T8R4D7W3
Content-Type: application/json

{
  "label": "Operations Read Only",
  "permissionUIDs": [
    "sd_types.read",
    "sd_instances.read",
    "time_series.read"
  ]
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the role UID and replacement input.",
            requestExample: `{
  "type": "request",
  "id": "role-update",
  "action": "update_role",
  "payload": {
    "uid": "role:P9K2M6Q1T8R4D7W3",
    "label": "Operations Read Only",
    "permissionUIDs": [
      "sd_types.read",
      "sd_instances.read",
      "time_series.read"
    ]
  }
}`,
          },
        ],
      },
      {
        id: "delete-role",
        title: "Delete Role",
        summary: "Delete a custom role that is not assigned to users.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Delete a role by UID.",
            requestExample: `mutation DeleteRole {
  deleteRole(uid: "role:P9K2M6Q1T8R4D7W3")
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "Delete the role resource directly.",
            requestExample: `DELETE /rest/roles/role:P9K2M6Q1T8R4D7W3`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the role UID in the delete action payload.",
            requestExample: `{
  "type": "request",
  "id": "role-delete",
  "action": "delete_role",
  "payload": {
    "uid": "role:P9K2M6Q1T8R4D7W3"
  }
}`,
          },
        ],
      },
      {
        id: "clone-role",
        title: "Clone Role",
        summary:
          "Create a custom role by copying an existing role's permissions.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Clone by source role UID and target label.",
            requestExample: `mutation CloneRole {
  cloneRole(
    uid: "role:user"
    label: "User Copy"
  ) {
    uid
    label
    permissions
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "POST a clone command under the source role.",
            requestExample: `POST /rest/roles/role:user/clone
Content-Type: application/json

{
  "label": "User Copy"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the source UID and new label.",
            requestExample: `{
  "type": "request",
  "id": "role-clone",
  "action": "clone_role",
  "payload": {
    "uid": "role:user",
    "label": "User Copy"
  }
}`,
          },
        ],
      },
      {
        id: "update-permission-label",
        title: "Update Permission Label",
        summary:
          "Change the display label of a permission without changing its UID.",
        variants: [
          {
            technology: "graphql",
            label: gql,
            summary: "Update the label by permission UID.",
            requestExample: `mutation UpdatePermissionLabel {
  updatePermissionLabel(
    uid: "time_series.read"
    label: "Read Time Series"
  ) {
    uid
    label
  }
}`,
          },
          {
            technology: "rest",
            label: rest,
            summary: "PATCH the permission label resource.",
            requestExample: `PATCH /rest/permissions/time_series.read/label
Content-Type: application/json

{
  "label": "Read Time Series"
}`,
          },
          {
            technology: "websocket",
            label: ws,
            summary: "Send the permission UID and replacement label.",
            requestExample: `{
  "type": "request",
  "id": "permission-label-update",
  "action": "update_permission_label",
  "payload": {
    "uid": "time_series.read",
    "label": "Read Time Series"
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
