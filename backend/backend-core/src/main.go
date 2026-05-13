/**
 * @file main.go
 * @brief Vstupní bod modulu Backend Core platformy RIoT a spuštění jeho interních komunikačních smyček.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ vstupního bodu a inicializace serverové části.
 * - Vojtěch Hubáček: doplnění runISCLoop a napojení zpracování příchozích raw datových bodů přes ProcessIncomingRawDataPoints.
 *
 * @defgroup riot_backend_core Backend Core
 * @ingroup riot
 * @see ../README.md
 */
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/isc"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func waitForDependencies() {
	rawRabbitMQURL := sharedUtils.GetEnvironmentVariableValue("RABBITMQ_URL").GetPayloadOrDefault("amqp://guest:guest@rabbitmq:5672")
	parsedRabbitMQURL, err := url.Parse(rawRabbitMQURL)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the RabbitMQ URL: %s", rawRabbitMQURL))
	rabbitMQ := sharedUtils.NewPairOf(parsedRabbitMQURL.Hostname(), parsedRabbitMQURL.Port())
	rawPostgresURL := sharedUtils.GetEnvironmentVariableValue("POSTGRES_URL").GetPayloadOrDefault("postgres://admin:password@postgres:5432/postgres-db")
	parsedPostgresURL, err := url.Parse(rawPostgresURL)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the Postgres URL: %s", rawPostgresURL))
	postgreSQL := sharedUtils.NewPairOf(parsedPostgresURL.Hostname(), parsedPostgresURL.Port())
	sharedUtils.TerminateOnError(sharedUtils.WaitForDSs(time.Minute, rabbitMQ, postgreSQL), "Some dependencies of this application are inaccessible")
}

func kickstartISC() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	isc.SetupRabbitMQInfrastructureForISC(rabbitMQClient)
	isc.EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient)
	isc.EnqueueMessageRepresentingCurrentSDInstanceConfiguration(rabbitMQClient)
	go runISCLoop(isc.ProcessIncomingMessageProcessingUnitConnectionNotifications)
	go runISCLoop(isc.ProcessIncomingSDTypeRegistrationRequests)
	go runISCLoop(isc.ProcessIncomingSDInstanceRegistrationRequests)
	go runISCLoop(isc.ProcessIncomingRawDataPoints)
	go runISCLoop(isc.ProcessIncomingKPIFulfillmentCheckResults)
}

func runISCLoop(worker func()) {
	for {
		worker()
		time.Sleep(time.Second)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.Println("Waiting for dependencies...")
	waitForDependencies()
	log.Println("Dependencies ready...")
	err := dbClient.GetRelationalDatabaseClientInstance().PerformOnStartupOperations(auth.RolePermissions, auth.FormatPermissionLabel)
	sharedUtils.TerminateOnError(err, "Unable to perform on-startup database operations")
	//sharedUtils.StartLoggingProfilingInformationPeriodically(time.Minute)
	kickstartISC()
	api.StartServer()
}
