package config

import "strings"

const (
	NEW_SUBSCRIPTION_SUFFIX = "_new_subscritpions"
	MO_TARIFFICATE          = "_mo_tarifficate"
	REQUESTS_SUFFIX         = "_requests"
	RESPONSES_SUFFIX        = "_responses"
	SMS_REQUEST_SUFFIX      = "_sms_request"
	SMS_RESPONSE_SUFFIX     = "_sms_response"
)

type OperatorQueueConfig struct {
	NewSubscription string `yaml:"-"`
	Requests        string `yaml:"-"`
	Responses       string `yaml:"-"`
	SMSRequest      string `yaml:"-"`
	SMSResponse     string `yaml:"-"`
	MOTarifficate   string `yaml:"-"`
}
type OperatorConfig struct {
	RetriesEnabled           bool `default:"false" yaml:"retries_enabled,omitempty"`
	OperatorRequestQueueSize int  `default:"10" yaml:"operator_request_queue_size,omitempty"`
	GetFromDBRetryCount      int  `yaml:"get_from_db_retry_count,omitempty"`
}

func GetOperatorsQueue(enabledOperators map[string]OperatorConfig) map[string]OperatorQueueConfig {
	opConfig := make(map[string]OperatorQueueConfig, len(enabledOperators))
	for operatorName, _ := range enabledOperators {
		name := strings.ToLower(operatorName)
		opConfig[name] = OperatorQueueConfig{
			NewSubscription: name + NEW_SUBSCRIPTION_SUFFIX,
			Requests:        name + REQUESTS_SUFFIX,
			Responses:       name + RESPONSES_SUFFIX,
			SMSRequest:      name + SMS_REQUEST_SUFFIX,
			SMSResponse:     name + SMS_RESPONSE_SUFFIX,
			MOTarifficate:   name + MO_TARIFFICATE,
		}
	}
	return opConfig
}

func GetNewSubscriptionQueueName(operatorName string) string {
	return operatorName + NEW_SUBSCRIPTION_SUFFIX
}
