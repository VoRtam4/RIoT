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

func StreamTimeSeriesAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, graphQLModel.TimeSeriesTypeKpi, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return err
	}
	db := dbClient.GetRelationalDatabaseClientInstance()
	for _, id := range input.KpiDefinitionIDs {
		res := db.LoadKPIDefinition(userID, id)
		if res.IsFailure() {
			return res.GetError()
		}
	}
	req := gql2dll.ToDLLTimeSeriesReadKPIRequest(input, sdTypeUID, instanceUIDs)
	params := loadParametersFromDB(graphQLModel.TimeSeriesTypeKpi, input.SdTypeID)
	return StreamTimeSeriesBase(userID, req, params, onBatch)
}

func StreamTimeSeries(userID uint32, input graphQLModel.TimeSeriesReadInput, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return err
	}
	db := dbClient.GetRelationalDatabaseClientInstance()
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		for _, id := range input.KpiDefinitionIDs {
			res := db.LoadKPIDefinition(userID, id)
			if res.IsFailure() {
				return res.GetError()
			}
		}
	}
	req := gql2dll.ToDLLTimeSeriesReadRequest(input, sdTypeUID, instanceUIDs)
	params := loadParametersFromDB(input.Type, input.SdTypeID)
	return StreamTimeSeriesBase(userID, req, params, onBatch)
}

func StreamTimeSeriesBase(userID uint32, req sharedModel.TimeSeriesReadRequest, params []graphQLModel.TimeSeriesParameter, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	client := rabbitmq.NewClient()
	defer client.Dispose()
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return err
	}
	correlationID := uuid.New().String()
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), "", sharedConstants.TimeSeriesReadRequestQueueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		Body:          jsonData,
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	})
	if err != nil {
		return err
	}
	return rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](msgs, correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, onBatch(graphQLModel.TimeSeriesReadResponse{
				Error: toErrorPtr(resp.Error),
			})
		}
		gqlResp := mapToGraphQLResponse(resp, params)
		if err := onBatch(gqlResp); err != nil {
			return true, err
		}
		if !resp.HasMoreBatches {
			return true, nil
		}
		return false, nil
	})
}

func ReadTimeSeriesAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) sharedUtils.Result[graphQLModel.TimeSeriesReadResponse] {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, graphQLModel.TimeSeriesTypeKpi, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	db := dbClient.GetRelationalDatabaseClientInstance()
	req := gql2dll.ToDLLTimeSeriesReadKPIRequest(input, sdTypeUID, instanceUIDs)
	batch := 200
	req.Batch = &batch
	for _, id := range input.KpiDefinitionIDs {
		res := db.LoadKPIDefinition(userID, id)
		if res.IsFailure() {
			return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](res.GetError())
		}
	}
	tsResp, err := executeTimeSeriesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	params := loadParametersFromDB(graphQLModel.TimeSeriesTypeKpi, input.SdTypeID)
	gqlResp := mapToGraphQLResponse(tsResp, params)
	return sharedUtils.NewSuccessResult(gqlResp)
}

func ReadTimeSeries(userID uint32, input graphQLModel.TimeSeriesReadInput) sharedUtils.Result[graphQLModel.TimeSeriesReadResponse] {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	db := dbClient.GetRelationalDatabaseClientInstance()
	req := gql2dll.ToDLLTimeSeriesReadRequest(input, sdTypeUID, instanceUIDs)
	batch := 200
	req.Batch = &batch
	if input.Type == graphQLModel.TimeSeriesTypeKpi {
		for _, id := range input.KpiDefinitionIDs {
			res := db.LoadKPIDefinition(userID, id)
			if res.IsFailure() {
				return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](res.GetError())
			}
		}
	}
	tsResp, err := executeTimeSeriesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	params := loadParametersFromDB(input.Type, input.SdTypeID)
	gqlResp := mapToGraphQLResponse(tsResp, params)
	return sharedUtils.NewSuccessResult(gqlResp)
}

func DistinctTimeSeriesTagValues(userID uint32, input graphQLModel.TimeSeriesDistinctTagValuesInput) sharedUtils.Result[graphQLModel.TimeSeriesDistinctTagValuesResponse] {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesDistinctTagValuesResponse](err)
	}
	tagErr := validateDistinctTagInput(userID, input)
	if tagErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesDistinctTagValuesResponse](tagErr)
	}
	req := gql2dll.ToDLLTimeSeriesDistinctTagValuesRequest(input, sdTypeUID, instanceUIDs)
	resp, err := executeTimeSeriesDistinctTagValuesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesDistinctTagValuesResponse](err)
	}
	return sharedUtils.NewSuccessResult(graphQLModel.TimeSeriesDistinctTagValuesResponse{
		Values: resp.Values,
		Error:  toErrorPtr(resp.Error),
	})
}

func executeTimeSeriesQuery(req sharedModel.TimeSeriesReadRequest) (sharedModel.TimeSeriesReadResponse, error) {
	client := rabbitmq.NewClient()
	defer client.Dispose()
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	correlationID := uuid.New().String()
	jsonData, err := json.Marshal(req)
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	err = ch.PublishWithContext(context.Background(), "", sharedConstants.TimeSeriesReadRequestQueueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		Body:          jsonData,
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	})
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	var finalResponse sharedModel.TimeSeriesReadResponse
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](msgs, correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, fmt.Errorf("%s", resp.Error)
		}
		if len(resp.Data) > 0 {
			finalResponse.Data = append(finalResponse.Data, resp.Data...)
		}
		finalResponse.HasMoreData = resp.HasMoreData
		finalResponse.NextCursor = resp.NextCursor
		if !resp.HasMoreBatches {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	return finalResponse, nil
}

func executeTimeSeriesDistinctTagValuesQuery(req sharedModel.TimeSeriesDistinctTagValuesRequest) (sharedModel.TimeSeriesDistinctTagValuesResponse, error) {
	client := rabbitmq.NewClient()
	defer client.Dispose()
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	correlationID := uuid.New().String()
	jsonData, err := json.Marshal(req)
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	err = ch.PublishWithContext(context.Background(), "", sharedConstants.TimeSeriesDistinctTagValuesRequestQueueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		Body:          jsonData,
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	})
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	var finalResponse sharedModel.TimeSeriesDistinctTagValuesResponse
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesDistinctTagValuesResponse](msgs, correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesDistinctTagValuesResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, fmt.Errorf("%s", resp.Error)
		}
		finalResponse = resp
		return true, nil
	})
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	return finalResponse, nil
}

func loadParametersFromDB(Type graphQLModel.TimeSeriesType, SdTypeID *uint32) []graphQLModel.TimeSeriesParameter {
	if SdTypeID == nil {
		return nil
	}
	var sdType dllModel.SDType
	parametersCount := 5
	res := dbClient.GetRelationalDatabaseClientInstance().LoadSDType(*SdTypeID)
	if res.IsFailure() {
		return nil
	}
	sdType = res.GetPayload()
	if Type == graphQLModel.TimeSeriesTypeRaw {
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
	if Type == graphQLModel.TimeSeriesTypeKpi {
		params = append(params, graphQLModel.TimeSeriesParameter{
			Denotation: "kpiDefinitionID",
			Label:      "KpiDefinition",
			Role:       graphQLModel.ParameterRoleMeta,
		})
		params = append(params, graphQLModel.TimeSeriesParameter{
			Denotation: "fulfilled",
			Label:      "Fulfilled",
			Role:       graphQLModel.ParameterRoleField,
		})
	}
	for _, p := range sdType.Parameters {
		if Type == graphQLModel.TimeSeriesTypeKpi && p.Role == dllModel.SDParameterRoleField {
			continue
		}
		params = append(params, graphQLModel.TimeSeriesParameter{
			Denotation: p.Denotation,
			Label:      p.Label,
			Role:       graphQLModel.ParameterRole(p.Role),
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
	for _, dp := range resp.Data {
		dataJSON, _ := json.Marshal(dp.Data)
		tagsJSON, _ := json.Marshal(dp.Tags)
		data = append(data, graphQLModel.TimeSeriesDataPoint{
			Time: dp.Time.Format(time.RFC3339Nano),
			Tags: string(tagsJSON),
			Data: string(dataJSON),
		})
	}
	baseStr, _ := json.Marshal(resp.Base)
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
		Base:           string(baseStr),
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
		kpi := res.GetPayload()
		sdRes := db.LoadSDType(kpi.SDTypeID)
		if sdRes.IsFailure() {
			return sdRes.GetError()
		}
		sdType = sdRes.GetPayload()
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

func resolveAndValidate(userID uint32, Type graphQLModel.TimeSeriesType, SdTypeID *uint32, SdInstanceIDs []uint32, KpiDefinitionIDs []uint32) (string, []string, error) {
	db := dbClient.GetRelationalDatabaseClientInstance()
	var sdTypeID *uint32
	var sdTypeUID string
	if SdTypeID != nil {
		res := db.LoadSDType(*SdTypeID)
		if res.IsFailure() {
			return "", nil, res.GetError()
		}
		sdTypeID = SdTypeID
		sdTypeUID = res.GetPayload().UID
	}
	var kpis []sharedModel.KPIDefinition
	if len(KpiDefinitionIDs) > 0 {
		for _, id := range KpiDefinitionIDs {
			res := db.LoadKPIDefinition(userID, id)
			if res.IsFailure() {
				return "", nil, res.GetError()
			}
			kpis = append(kpis, res.GetPayload())
		}
		kpiSDType := kpis[0].SDTypeID
		for _, k := range kpis {
			if k.SDTypeID != kpiSDType {
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
	if Type == graphQLModel.TimeSeriesTypeRaw {
		if sdTypeID == nil {
			return "", nil, fmt.Errorf("sdTypeID must be specified for RAW")
		}
	}
	if Type == graphQLModel.TimeSeriesTypeKpi {
		if sdTypeID == nil && len(kpis) == 0 {
			return "", nil, fmt.Errorf("either sdTypeID or kpiDefinitionIDs must be specified")
		}
	}
	var instanceUIDs []string
	if len(SdInstanceIDs) > 0 {
		instanceUIDs = make([]string, 0, len(SdInstanceIDs))
		for _, id := range SdInstanceIDs {
			res := db.LoadSDInstance(id)
			if res.IsFailure() {
				return "", nil, res.GetError()
			}
			inst := res.GetPayload()
			if sdTypeID != nil && inst.SDType.ID.GetPayload() != *sdTypeID {
				return "", nil, fmt.Errorf("sdInstance %d does not belong to sdType", id)
			}
			instanceUIDs = append(instanceUIDs, inst.UID)
		}
	}
	return sdTypeUID, instanceUIDs, nil
}
