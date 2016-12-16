package config

const (
	NEW_SUBSCRIPTION_SUFFIX = "_new_subscriptions"
	MO_TARIFFICATE          = "_mo_tarifficate"
	REQUESTS_SUFFIX         = "_requests"
	RESPONSES_SUFFIX        = "_responses"
	SMS_REQUEST_SUFFIX      = "_sms_requests"
	SMS_RESPONSE_SUFFIX     = "_sms_responses"
)

type ConsumeQueueConfig struct {
	Name          string `yaml:"name"`
	PrefetchCount int    `yaml:"prefetch_count" default:"600"`
	ThreadsCount  int    `yaml:"threads_count" default:"60"`
}

type OperatorConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `default:"true" yaml:"enabled"`
	Retries struct {
		Enabled     bool `default:"false" yaml:"enabled,omitempty"`
		QueueSize   int  `default:"1200" yaml:"queue_size,omitempty"`
		FromDBCount int  `default:"1200" yaml:"from_db_count,omitempty"`
	} `yaml:"retries"`
}

func (oc OperatorConfig) NewSubscriptionQueueName() string {
	return NewSubscriptionQueueName(oc.Name)
}
func NewSubscriptionQueueName(operatorName string) string {
	return operatorName + NEW_SUBSCRIPTION_SUFFIX
}
func (oc OperatorConfig) GetMOQueueName() string {
	return GetMOQueueName(oc.Name)
}
func GetMOQueueName(operatorName string) string {
	return operatorName + MO_TARIFFICATE
}
func (oc OperatorConfig) GetRequestsQueueName() string {
	return RequestQueue(oc.Name)
}
func (oc OperatorConfig) GetResponsesQueueName() string {
	return ResponsesQueue(oc.Name)
}
func ResponsesQueue(operatorName string) string {
	return operatorName + RESPONSES_SUFFIX
}
func RequestQueue(operatorName string) string {
	return operatorName + REQUESTS_SUFFIX
}
func (oc OperatorConfig) GetSMSRequestsQueueName() string {
	return SMSRequestQueue(oc.Name)
}
func SMSRequestQueue(operatorName string) string {
	return operatorName + SMS_REQUEST_SUFFIX
}
func (oc OperatorConfig) GetSMSResponsesQueueName() string {
	return SMSResponsesQueue(oc.Name)
}
func SMSResponsesQueue(operatorName string) string {
	return operatorName + SMS_RESPONSE_SUFFIX
}
