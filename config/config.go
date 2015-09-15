package config

import "github.com/spf13/viper"

type (
	emailConfig struct {
		Recipient   string
		Username    string
		Password    string
		EmailServer string
		Port        int
	}
)

//PagerDutyAPIKey stores the api key for pagerduty
var PagerDutyAPIKey string

//PagerDuty indicates if alerts should be sent via PagerDuty
var PagerDuty bool

//Email indicates if alerts should be sent via email
var Email bool

//EmailConfig proviedes the configuration for sending an alert via email
var EmailConfig emailConfig

//ParseConfig reads a json or yaml encoded config file
func ParseConfig(configPath string) {
	viper.SetConfigName(configPath)
	PagerDutyAPIKey = viper.GetString("PagerDutyAPIKey")
	EmailConfig = *new(emailConfig)
	EmailConfig.EmailServer = viper.GetString("email.server")
	EmailConfig.Username = viper.GetString("email.username")
	EmailConfig.Password = viper.GetString("email.password")
	EmailConfig.Port = viper.GetInt("email.port")
	EmailConfig.Recipient = viper.GetString("email.recipient")

}
