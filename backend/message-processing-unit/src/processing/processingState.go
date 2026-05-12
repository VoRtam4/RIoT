package processing

import (
	"log"
	"sync"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/google/uuid"
)

var (
	kpiDefinitionsBySDTypeDenotationMap      = make(map[string][]sharedModel.KPIDefinitionMPU)
	kpiDefinitionsBySDTypeDenotationMapMutex sync.RWMutex
	unitUUID                                 string
	limit                                    = 500
	sdTypeDefinitions                        = make(map[string]sharedModel.SDTypeDefinitionCache)
	sdTypeDefinitionsMutex                   sync.RWMutex
	activeReprocessJobs                      sync.Map
	jobSyncMap                               sync.Map
	initialConfigReady                       = make(chan struct{})
	initialConfigReadyOnce                   sync.Once
)

func InitializeProcessing() {
	unitUUID = uuid.New().String()
	log.Printf("[MPU] Processing unit initialized | unitUUID=%s", unitUUID)
}

func DenotationMapUpdates(rabbitMQClient rabbitmq.Client) {
	log.Printf("[MPU] Waiting for KPI configuration updates | unit=%s", unitUUID)
	err := rabbitmq.ConsumeJSONMessagesFromFanoutExchange[sharedModel.KPIConfigurationUpdateISCMessage](rabbitMQClient, sharedConstants.BuiltInFanoutExchangeName, func(messagePayload sharedModel.KPIConfigurationUpdateISCMessage) error {
		log.Printf("[MPU] KPI configuration update received | unit=%s entries=%d", unitUUID, len(messagePayload.KpiConfiguration))
		kpiDefinitionsBySDTypeDenotationMapMutex.Lock()
		kpiDefinitionsBySDTypeDenotationMap = messagePayload.KpiConfiguration
		kpiDefinitionsBySDTypeDenotationMapMutex.Unlock()
		log.Printf("[MPU] KPI configuration map updated successfully | unit=%s", unitUUID)
		initialConfigReadyOnce.Do(func() {
			close(initialConfigReady)
			log.Printf("[MPU][BOOTSTRAP] Initial KPI configuration received | unit=%s", unitUUID)
		})
		signalConfigReady(messagePayload.JobID)
		return nil
	})
	if err != nil {
		log.Printf("[MPU] Consumption of messages from the '%s' fanout exchange has failed: %s\n", sharedConstants.BuiltInFanoutExchangeName, err.Error())
	}
}

func UpdateSDType(messages []sharedModel.SDTypeUpdateISCMessage) {
	log.Printf("[MPU][SDTYPE] Received SDType update | count=%d", len(messages))
	newDefinitions := make(map[string]sharedModel.SDTypeDefinitionCache)
	for _, msg := range messages {
		paramMap := make(map[string]sharedModel.SDParameter, len(msg.Parameters))
		for _, p := range msg.Parameters {
			paramMap[p.Denotation] = sharedModel.SDParameter{
				Denotation: p.Denotation,
				Type:       sharedModel.SDParameterType(p.Type),
				Label:      p.Label,
				Role:       sharedModel.SDParameterRole(p.Role),
			}
		}
		newDefinitions[msg.SDTypeUID] = sharedModel.SDTypeDefinitionCache{
			Label:      msg.SDTypeUID,
			Parameters: paramMap,
		}
	}
	sdTypeDefinitionsMutex.Lock()
	sdTypeDefinitions = newDefinitions
	sdTypeDefinitionsMutex.Unlock()
	log.Printf("[MPU][SDTYPE] SDType cache updated | types=%d", len(newDefinitions))
}
