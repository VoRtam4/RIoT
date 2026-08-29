/**
 * @file timeSeriesHelpers.go
 * @brief Pomocné funkce pro přípravu požadavků, filtrů a RPC streamů časových řad.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé pomocné vrstvy pro time series dotazy, agregace a autorizaci požadavků.
 *
 * @ingroup riot_backend_core
 */
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
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
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

type resolvedTimeSeriesInput struct {
	sdTypeUID         string
	sdInstanceUIDs    []string
	kpiDefinitionUIDs []string
	kpiDefinitionIDs  []uint32
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
	resolved, err := resolveAndValidate(
		userID,
		input.Type,
		input.SdTypeUID,
		input.SdInstanceUIDs,
		input.KpiDefinitionUIDs,
	)
	if err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		if err := validateRequestedKPIDefinitions(userID, resolved.kpiDefinitionIDs); err != nil {
			return sharedModel.TimeSeriesReadRequest{}, err
		}
	}
	input.SdInstanceUIDs = resolved.sdInstanceUIDs
	return gql2dll.ToDLLTimeSeriesReadRequest(input, resolved.sdTypeUID, resolved.kpiDefinitionUIDs), nil
}

func buildTimeSeriesAggregateReadRequest(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) (sharedModel.TimeSeriesReadRequest, error) {
	resolved, err := resolveAndValidate(
		userID,
		graphQLModel.TimeSeriesTypeKpi,
		input.SdTypeUID,
		input.SdInstanceUIDs,
		input.KpiDefinitionUIDs,
	)
	if err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	if err := validateRequestedKPIDefinitions(userID, resolved.kpiDefinitionIDs); err != nil {
		return sharedModel.TimeSeriesReadRequest{}, err
	}
	input.SdInstanceUIDs = resolved.sdInstanceUIDs
	return gql2dll.ToDLLTimeSeriesReadKPIRequest(input, resolved.sdTypeUID, resolved.kpiDefinitionUIDs), nil
}

func buildTimeSeriesDistinctTagValuesRequest(userID uint32, input graphQLModel.TimeSeriesDistinctTagValuesInput) (sharedModel.TimeSeriesDistinctTagValuesRequest, error) {
	resolved, err := resolveAndValidate(
		userID,
		input.Type,
		input.SdTypeUID,
		input.SdInstanceUIDs,
		input.KpiDefinitionUIDs,
	)
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesRequest{}, err
	}
	if err := validateDistinctTagInput(input, resolved.sdTypeUID); err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesRequest{}, err
	}
	input.SdInstanceUIDs = resolved.sdInstanceUIDs
	return gql2dll.ToDLLTimeSeriesDistinctTagValuesRequest(input, resolved.sdTypeUID, resolved.kpiDefinitionUIDs), nil
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

func loadParametersFromDB(timeSeriesType graphQLModel.TimeSeriesType, sdTypeUID *string) []graphQLModel.TimeSeriesParameter {
	if sdTypeUID == nil {
		return nil
	}
	normalizedSDTypeUID, normalizeErr := normalizeSDTypeUID(*sdTypeUID)
	if normalizeErr != nil {
		return nil
	}
	db := dbClient.GetRelationalDatabaseClientInstance()
	res := db.LoadSDTypeBasedOnUID(normalizedSDTypeUID)
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
				Denotation: "kpiDefinitionUID",
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
			Time:             resp.NextCursor.Time.Format(time.RFC3339Nano),
			SdInstanceUID:    resp.NextCursor.SDInstanceUID,
			KpiDefinitionUID: resp.NextCursor.KPIDefinitionUID,
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

func validateDistinctTagInput(input graphQLModel.TimeSeriesDistinctTagValuesInput, sdTypeUID string) error {
	if strings.TrimSpace(input.Tag) == "" {
		return fmt.Errorf("tag must be specified")
	}

	db := dbClient.GetRelationalDatabaseClientInstance()
	normalizedSDTypeUID, normalizeErr := normalizeSDTypeUID(sdTypeUID)
	if normalizeErr != nil {
		return normalizeErr
	}
	sdTypeResult := db.LoadSDTypeBasedOnUID(normalizedSDTypeUID)
	if sdTypeResult.IsFailure() {
		return sdTypeResult.GetError()
	}
	sdType := sdTypeResult.GetPayload()
	allowedTags := map[string]bool{
		"sdInstanceUID": true,
	}
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		allowedTags["kpiDefinitionUID"] = true
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

func resolveAndValidate(userID uint32, timeSeriesType graphQLModel.TimeSeriesType, sdTypeUIDInput *string, sdInstanceUIDInputs []string, kpiDefinitionUIDInputs []string) (resolvedTimeSeriesInput, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	var sdTypeID *uint32
	var sdTypeUID string
	if sdTypeUIDInput != nil && strings.TrimSpace(*sdTypeUIDInput) != "" {
		normalizedUID, normalizeErr := normalizeSDTypeUID(*sdTypeUIDInput)
		if normalizeErr != nil {
			return resolvedTimeSeriesInput{}, normalizeErr
		}
		res := db.LoadSDTypeBasedOnUID(normalizedUID)
		if res.IsFailure() {
			return resolvedTimeSeriesInput{}, res.GetError()
		}
		loadedSDType := res.GetPayload()
		loadedID := loadedSDType.ID.GetPayload()
		sdTypeID = &loadedID
		sdTypeUID = loadedSDType.UID
	}
	var kpis []sharedModel.KPIDefinition
	kpiDefinitionUIDs := make([]string, 0, len(kpiDefinitionUIDInputs))
	kpiDefinitionIDs := make([]uint32, 0, len(kpiDefinitionUIDInputs))
	if len(kpiDefinitionUIDInputs) > 0 {
		for _, uid := range kpiDefinitionUIDInputs {
			trimmedUID := strings.TrimSpace(uid)
			if trimmedUID == "" {
				return resolvedTimeSeriesInput{}, fmt.Errorf("kpiDefinitionUID must not be empty")
			}
			res := db.LoadKPIDefinitionByUID(userID, trimmedUID)
			if res.IsFailure() {
				return resolvedTimeSeriesInput{}, res.GetError()
			}
			kpi := res.GetPayload()
			if kpi.ID == nil {
				return resolvedTimeSeriesInput{}, fmt.Errorf("KPI definition loaded by UID has no internal ID: %s", trimmedUID)
			}
			kpis = append(kpis, kpi)
			kpiDefinitionUIDs = append(kpiDefinitionUIDs, trimmedUID)
			kpiDefinitionIDs = append(kpiDefinitionIDs, *kpi.ID)
		}
	}
	if len(kpis) > 0 {
		kpiSDType := kpis[0].SDTypeID
		for _, kpi := range kpis {
			if kpi.SDTypeID != kpiSDType {
				return resolvedTimeSeriesInput{}, fmt.Errorf("all KPI definitions must belong to the same SDType")
			}
		}
		if sdTypeID == nil {
			sdTypeID = &kpiSDType
			res := db.LoadSDType(kpiSDType)
			if res.IsFailure() {
				return resolvedTimeSeriesInput{}, res.GetError()
			}
			sdTypeUID = res.GetPayload().UID
		} else if *sdTypeID != kpiSDType {
			return resolvedTimeSeriesInput{}, fmt.Errorf("sdType does not match KPI definitions")
		}
	}
	if timeSeriesType == graphQLModel.TimeSeriesTypeRaw && sdTypeID == nil {
		return resolvedTimeSeriesInput{}, fmt.Errorf("sdTypeUID must be specified for RAW")
	}
	if timeSeriesType == graphQLModel.TimeSeriesTypeKpi && sdTypeID == nil && len(kpis) == 0 {
		return resolvedTimeSeriesInput{}, fmt.Errorf("either sdTypeUID or kpiDefinitionUIDs must be specified")
	}
	sdInstanceUIDs := make([]string, 0, len(sdInstanceUIDInputs))
	for _, uid := range sdInstanceUIDInputs {
		trimmedUID := strings.TrimSpace(uid)
		if trimmedUID == "" {
			return resolvedTimeSeriesInput{}, fmt.Errorf("sdInstanceUID must not be empty")
		}
		if sdTypeUID != "" {
			normalizedUID, _, normalizeErr := sharedUtils.NormalizeScopedUID(trimmedUID, sdTypeUID, "sdi", "SD instance")
			if normalizeErr != nil {
				return resolvedTimeSeriesInput{}, normalizeErr
			}
			trimmedUID = normalizedUID
		}
		res := db.LoadSDInstanceBasedOnUID(trimmedUID)
		if res.IsFailure() {
			return resolvedTimeSeriesInput{}, res.GetError()
		}
		if res.GetPayload().IsEmpty() {
			return resolvedTimeSeriesInput{}, fmt.Errorf("couldn't find SD instance for UID: %s", trimmedUID)
		}
		instance := res.GetPayload().GetPayload()
		if sdTypeID != nil && instance.SDType.ID.GetPayload() != *sdTypeID {
			return resolvedTimeSeriesInput{}, fmt.Errorf("sdInstance %s does not belong to sdType", trimmedUID)
		}
		sdInstanceUIDs = append(sdInstanceUIDs, instance.UID)
	}
	return resolvedTimeSeriesInput{
		sdTypeUID:         sdTypeUID,
		sdInstanceUIDs:    sdInstanceUIDs,
		kpiDefinitionUIDs: kpiDefinitionUIDs,
		kpiDefinitionIDs:  kpiDefinitionIDs,
	}, nil
}
