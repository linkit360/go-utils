package structs

import (
	"time"
)

type AccessCampaignNotify struct {
	Msisdn       string    `json:"msisdn,omitempty"`
	CampaignHash string    `json:"campaign_hash,omitempty"`
	Tid          string    `json:"tid,omitempty"`
	IP           string    `json:"ip,omitempty"`
	OperatorCode int64     `json:"operator_code,omitempty"`
	CountryCode  int64     `json:"country_code,omitempty"`
	ServiceCode  string    `json:"service_code,omitempty"`
	CampaignId   string    `json:"campaign_id,omitempty"`
	ContentCode  string    `json:"content_code,omitempty"`
	Supported    bool      `json:"supported,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Referer      string    `json:"referer,omitempty"`
	UrlPath      string    `json:"url_path,omitempty"`
	Method       string    `json:"method,omitempty"`
	Headers      string    `json:"headers,omitempty"`
	Error        string    `json:"err,omitempty"`
	SentAt       time.Time `json:"sent_at,omitempty"`
}

type ContentSentProperties struct {
	ContentId      string    `json:"id_content,omitempty"`
	SentAt         time.Time `json:"sent_at,omitempty"`
	Msisdn         string    `json:"msisdn,omitempty"`
	Tid            string    `json:"tid,omitempty"`
	UniqueUrl      string    `json:"unique_url,omitempty"`
	ContentPath    string    `json:"content_path,omitempty"`
	ContentName    string    `json:"content_name,omitempty"`
	CapmaignHash   string    `json:"capmaign_hash,omitempty"`
	CampaignId     string    `json:"campaign_id,omitempty"`
	ServiceCode    string    `json:"service_code,omitempty"`
	SubscriptionId int64     `json:"subscription_id,omitempty"`
	CountryCode    int64     `json:"country_code,omitempty"`
	OperatorCode   int64     `json:"operator_code,omitempty"`
	Publisher      string    `json:"publisher,omitempty"`
	Error          string    `json:"error,omitempty"`
}

func (t ContentSentProperties) Key() string {
	return t.Msisdn + "-" + t.CampaignId
}

type EventNotifyContentSent struct {
	EventName string                `json:"event_name,omitempty"`
	EventData ContentSentProperties `json:"event_data,omitempty"`
}
