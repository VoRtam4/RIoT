package domainLogicLayer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type timeSeriesRPCSession struct {
	client        rabbitmq.Client
	channel       *amqp091.Channel
	messages      <-chan amqp091.Delivery
	correlationID string
	replyTo       string
}

func newTimeSeriesRPCSession(correlationID string) (*timeSeriesRPCSession, error) {
	client := rabbitmq.NewClient()
	channel := client.GetChannel()
	replyQueue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		client.Dispose()
		return nil, err
	}
	messages, err := channel.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		client.Dispose()
		return nil, err
	}
	return &timeSeriesRPCSession{
		client:        client,
		channel:       channel,
		messages:      messages,
		correlationID: correlationID,
		replyTo:       replyQueue.Name,
	}, nil
}

func newGeneratedTimeSeriesRPCSession() (*timeSeriesRPCSession, error) {
	return newTimeSeriesRPCSession(uuid.New().String())
}

func (s *timeSeriesRPCSession) Close() {
	if s == nil || s.client == nil {
		return
	}
	s.client.Dispose()
}

func (s *timeSeriesRPCSession) publish(queueName string, payload []byte) error {
	return s.channel.PublishWithContext(context.Background(), "", queueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		Body:          payload,
		CorrelationId: s.correlationID,
		ReplyTo:       s.replyTo,
	})
}

func (s *timeSeriesRPCSession) PublishReadRequest(req sharedModel.TimeSeriesReadRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return s.publish(sharedConstants.TimeSeriesReadRequestQueueName, payload)
}

func (s *timeSeriesRPCSession) PublishDistinctTagValuesRequest(req sharedModel.TimeSeriesDistinctTagValuesRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return s.publish(sharedConstants.TimeSeriesDistinctTagValuesRequestQueueName, payload)
}

func buildTimeSeriesReadRequest(userID uint32, input graphQLModel.TimeSeriesReadInput) (sharedModel.TimeSeriesReadRequest, error) {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		if err := validateRequestedKPIDefinitions(userID, input.KpiDefinitionIDs); err != nil {
			return sharedModel.TimeSeriesReadRequest{}, err
		}
	}
	return gql2dll.ToDLLTimeSeriesReadRequest(input, sdTypeUID, instanceUIDs), nil
}

func buildTimeSeriesAggregateReadRequest(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) (sharedModel.TimeSeriesReadRequest, error) {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, graphQLModel.TimeSeriesTypeKpi, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	if err := validateRequestedKPIDefinitions(userID, input.KpiDefinitionIDs); err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	return gql2dll.ToDLLTimeSeriesReadKPIRequest(input, sdTypeUID, instanceUIDs), nil
}

func buildTimeSeriesDistinctTagValuesRequest(userID uint32, input graphQLModel.TimeSeriesDistinctTagValuesInput) (sharedModel.TimeSeriesDistinctTagValuesRequest, error) {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesRequest{}, err
	}
	if err := validateDistinctTagInput(userID, input); err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesRequest{}, err
	}
	return gql2dll.ToDLLTimeSeriesDistinctTagValuesRequest(input, sdTypeUID, instanceUIDs), nil
}

func validateRequestedKPIDefinitions(userID uint32, kpiDefinitionIDs []uint32) error {
	db := dbClient.GetRelationalDatabaseClientInstance()
	for _, id := range kpiDefinitionIDs {
		res := db.LoadKPIDefinition(userID, id)
		if res.IsFailure() {
			return res.GetError()
		}
	}
	return nil
}

func loadParametersFromDB(timeSeriesType graphQLModel.TimeSeriesType, sdTypeID *uint32) []graphQLModel.TimeSeriesParameter {
	if sdTypeID == nil {
		return nil
	}
	res := dbClient.GetRelationalDatabaseClientInstance().LoadSDType(*sdTypeID)
	if res.IsFailure() {
		return nil
	}
	sdType := res.GetPayload()
	parametersCount := 5
	if timeSeriesType == graphQLModel.TimeSeriesTypeRaw {
		parametersCount = len(sdType.Parameters) + 4
	}
	params := make([]graphQLModel.TimeSeriesParameter, 0, parametersCount)
	params = append(params,
		graphQLModel.TimeSeriesParameter{
			Denotation: "type",
			Label:      "Type",
			Role:       graphQLModel.ParameterRoleMeta,
		},
		graphQLModel.TimeSeriesParameter{
			Denotation: "sdTypeUID",
			Label:      "SDType",
			Role:       graphQLModel.ParameterRoleMeta,
		},
		graphQLModel.TimeSeriesParameter{
			Denotation: "time",
			Label:      "Time",
			Role:       graphQLModel.ParameterRoleTime,
		},
		graphQLModel.TimeSeriesParameter{
			Denotation: "sdInstanceUID",
			Label:      "SDInstance",
			Role:       graphQLModel.ParameterRoleMeta,
		},
	)
	if timeSeriesType == graphQLModel.TimeSeriesTypeKpi {
		params = append(params,
			graphQLModel.TimeSeriesParameter{
				Denotation: "kpiDefinitionID",
				Label:      "KpiDefinition",
				Role:       graphQLModel.ParameterRoleMeta,
			},
			graphQLModel.TimeSeriesParameter{
				Denotation: "fulfilled",
				Label:      "Fulfilled",
				Role:       graphQLModel.ParameterRoleField,
			},
		)
	}
	for _, parameter := range sdType.Parameters {
		if timeSeriesType == graphQLModel.TimeSeriesTypeKpi && parameter.Role == dllModel.SDParameterRoleField {
			continue
		}
		params = append(params, graphQLModel.TimeSeriesParameter{
			Denotation: parameter.Denotation,
			Label:      parameter.Label,
			Role:       graphQLModel.ParameterRole(parameter.Role),
		})
	}
	roleOrder := map[string]int{
		string(graphQLModel.ParameterRoleTime):  0,
		string(graphQLModel.ParameterRoleMeta):  1,
		string(graphQLModel.ParameterRoleTag):   2,
		string(graphQLModel.ParameterRoleField): 3,
	}
	sort.SliceStable(params, func(i, j int) bool {
		r1 := strings.ToLower(string(params[i].Role))
		r2 := strings.ToLower(string(params[j].Role))
		o1 := roleOrder[r1]
		o2 := roleOrder[r2]
		if _, ok := roleOrder[r1]; !ok {
			o1 = 99
		}
		if _, ok := roleOrder[r2]; !ok {
			o2 = 99
		}
		if o1 != o2 {
			return o1 < o2
		}
		return params[i].Denotation < params[j].Denotation
	})
	return params
}

func mapToGraphQLResponse(resp sharedModel.TimeSeriesReadResponse, params []graphQLModel.TimeSeriesParameter) graphQLModel.TimeSeriesReadResponse {
	data := make([]graphQLModel.TimeSeriesDataPoint, 0, len(resp.Data))
	for _, dataPoint := range resp.Data {
		dataJSON, _ := json.Marshal(dataPoint.Data)
		tagsJSON, _ := json.Marshal(dataPoint.Tags)
		data = append(data, graphQLModel.TimeSeriesDataPoint{
			Time: dataPoint.Time.Format(time.RFC3339Nano),
			Tags: string(tagsJSON),
			Data: string(dataJSON),
		})
	}
	baseString, _ := json.Marshal(resp.Base)
	var cursor *graphQLModel.TimeSeriesCursor
	if resp.NextCursor != nil {
		cursor = &graphQLModel.TimeSeriesCursor{
			Time:            resp.NextCursor.Time.Format(time.RFC3339Nano),
			SdInstanceUID:   resp.NextCursor.SDInstanceUID,
			KpiDefinitionID: resp.NextCursor.KPIDefinitionID,
		}
	}
	return graphQLModel.TimeSeriesReadResponse{
		Parameters:     params,
		Base:           string(baseString),
		Data:           data,
		HasMoreBatches: resp.HasMoreBatches,
		HasMoreData:    resp.HasMoreData,
		NextCursor:     cursor,
		Error:          toErrorPtr(resp.Error),
	}
}

func validateDistinctTagInput(userID uint32, input graphQLModel.TimeSeriesDistinctTagValuesInput) error {
	if strings.TrimSpace(input.Tag) == "" {
		return fmt.Errorf("tag must be specified")
	}

	db := dbClient.GetRelationalDatabaseClientInstance()
	var sdType dllModel.SDType
	var loaded bool

	if input.SdTypeID != nil {
		res := db.LoadSDType(*input.SdTypeID)
		if res.IsFailure() {
			return res.GetError()
		}
		sdType = res.GetPayload()
		loaded = true
	} else if input.Type == graphQLModel.TimeSeriesTypeKpi && len(input.KpiDefinitionIDs) > 0 {
		res := db.LoadKPIDefinition(userID, input.KpiDefinitionIDs[0])
		if res.IsFailure() {
			return res.GetError()
		}
		kpiDefinition := res.GetPayload()
		sdTypeResult := db.LoadSDType(kpiDefinition.SDTypeID)
		if sdTypeResult.IsFailure() {
			return sdTypeResult.GetError()
		}
		sdType = sdTypeResult.GetPayload()
		loaded = true
	}
	if !loaded {
		return fmt.Errorf("unable to resolve sdType for tag validation")
	}
	allowedTags := map[string]bool{
		"sdInstanceUID": true,
	}
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		allowedTags["kpiDefinitionID"] = true
	}
	for _, parameter := range sdType.Parameters {
		if parameter.Role == dllModel.SDParameterRoleTag {
			allowedTags[parameter.Denotation] = true
		}
	}
	if !allowedTags[input.Tag] {
		return fmt.Errorf("tag %q is not allowed for sdType %s", input.Tag, sdType.UID)
	}
	return nil
}

func toErrorPtr(err string) *string {
	if err == "" {
		return nil
	}
	return &err
}

func resolveAndValidate(userID uint32, timeSeriesType graphQLModel.TimeSeriesType, sdTypeIDInput *uint32, sdInstanceIDs []uint32, kpiDefinitionIDs []uint32) (string, []string, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	var sdTypeID *uint32
	var sdTypeUID string
	if sdTypeIDInput != nil {
		res := db.LoadSDType(*sdTypeIDInput)
		if res.IsFailure() {
			return "", nil, res.GetError()
		}
		sdTypeID = sdTypeIDInput
		sdTypeUID = res.GetPayload().UID
	}
	var kpis []sharedModel.KPIDefinition
	if len(kpiDefinitionIDs) > 0 {
		for _, id := range kpiDefinitionIDs {
			res := db.LoadKPIDefinition(userID, id)
			if res.IsFailure() {
				return "", nil, res.GetError()
			}
			kpis = append(kpis, res.GetPayload())
		}
		kpiSDType := kpis[0].SDTypeID
		for _, kpi := range kpis {
			if kpi.SDTypeID != kpiSDType {
				return "", nil, fmt.Errorf("all KPI definitions must belong to the same SDType")
			}
		}
		if sdTypeID == nil {
			sdTypeID = &kpiSDType
			res := db.LoadSDType(kpiSDType)
			if res.IsFailure() {
				return "", nil, res.GetError()
			}
			sdTypeUID = res.GetPayload().UID
		} else if *sdTypeID != kpiSDType {
			return "", nil, fmt.Errorf("sdType does not match KPI definitions")
		}
	}
	if timeSeriesType == graphQLModel.TimeSeriesTypeRaw && sdTypeID == nil {
		return "", nil, fmt.Errorf("sdTypeID must be specified for RAW")
	}
	if timeSeriesType == graphQLModel.TimeSeriesTypeKpi && sdTypeID == nil && len(kpis) == 0 {
		return "", nil, fmt.Errorf("either sdTypeID or kpiDefinitionIDs must be specified")
	}
	instanceUIDs := make([]string, 0, len(sdInstanceIDs))
	for _, id := range sdInstanceIDs {
		res := db.LoadSDInstance(id)
		if res.IsFailure() {
			return "", nil, res.GetError()
		}
		instance := res.GetPayload()
		if sdTypeID != nil && instance.SDType.ID.GetPayload() != *sdTypeID {
			return "", nil, fmt.Errorf("sdInstance %d does not belong to sdType", id)
		}
		instanceUIDs = append(instanceUIDs, instance.UID)
	}
	return sdTypeUID, instanceUIDs, nil
}
