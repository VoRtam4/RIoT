package sharedConstants

const (
	BuiltInFanoutExchangeName                             = "amq.fanout"
	MessageProcessingUnitConnectionNotificationsQueueName = "message-processing-unit-connection-notifications"
	KPIFulfillmentCheckRequestsQueueName                  = "kpi-fulfillment-check-requests"
	KPIFulfillmentCheckResultsQueueName                   = "kpi-fulfillment-check-results"
	SDTypeRegistrationRequestsQueueName                   = "sd-type-registration-requests"
	SDInstanceRegistrationRequestsQueueName               = "sd-instance-registration-requests"
	SetOfSDInstancesUpdatesQueueName                      = "set-of-sd-instances-updates"
	SetOfSDTypesUpdatesQueueName                          = "set-of-sd-types-updates"
	TimeSeriesRawDataQueueName                            = "time-series-raw-data"
	TimeSeriesKPIResultQueueName                          = "time-series-kpi-results"
	TimeSeriesReadRequestQueueName                        = "time-series-read-request"
	TimeSeriesReadResponseQueueName                       = "time-series-read-response"
	KPIReprocessRequestQueueName                          = "kpi-reprocess-requests"
	TimeSeriesDeleteRequestQueueName                      = "time-series-delete-requests"
)
