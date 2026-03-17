package domainLogicLayer

import (
	"sync"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
)

var (
	dllRabbitMQClient rabbitmq.Client
	once              sync.Once
)

func getDLLRabbitMQClient() rabbitmq.Client {
	once.Do(func() {
		dllRabbitMQClient = rabbitmq.NewClient()
	})
	return dllRabbitMQClient
}
