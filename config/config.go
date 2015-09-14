package config

//PagerDutyAPIKey stores the api key for pagerduty
var PagerDutyAPIKey string

//PagerDuty indicates if alerts should be sent via PagerDuty
var PagerDuty bool

//Email indicates if alerts should be sent via email
var Email bool

//EmailConfig proviedes the configuration for sending an alert via email
var EmailConfig emailConfig

type (
	emailConfig struct {
		Recipient   string
		Username    string
		Password    string
		EmailServer string
		Port        int
	}
)
