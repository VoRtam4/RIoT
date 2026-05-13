/**
 * @file misc.go
 * @brief Sdílené pomocné prvky doménové vrstvy Backend Core.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace sdíleného základu doménové vrstvy.
 *
 * @ingroup riot_backend_core
 */
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
