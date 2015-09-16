package config

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/spf13/viper"
)

type (
	//InfluxDBConfigStruct is a struct which describes the config neccessary
	//to pull metrics from InfluxDB
	InfluxDBConfigStruct struct {
		Host     string
		Port     string
		Username string
		Password string
		DB       string
	}

	//EmailConfigStruct is a struct describing the configuration
	//required to send an alert via email
	EmailConfigStruct struct {
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

//EmailConfig details the configuration for sending an alert via email
var EmailConfig EmailConfigStruct

//influxDBConfig describes the configuration needed to communicate with InflxuDB
var InfluxDBConfig InfluxDBConfigStruct

//ReadConfig takes a file path as a string and returns a string representing
//the contents of that file
func ReadConfig(configFile string) ([]byte, error) {
	//viper accepts config file without extension, so remove extension
	if configFile == "" {
		panic("No config file provided")
	}

	f, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	return f, err
}

//ParseConfig parses a YAML config  file
func ParseConfig(rawConfig []byte) {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(rawConfig))

	PagerDutyAPIKey = viper.GetString("PagerDutyAPIKey")

	EmailConfig = EmailConfigStruct{
		EmailServer: viper.GetString("email.server"),
		Username:    viper.GetString("email.username"),
		Password:    viper.GetString("email.password"),
		Port:        viper.GetInt("email.port"),
		Recipient:   viper.GetString("email.recipient"),
	}

	InfluxDBConfig = InfluxDBConfigStruct{
		Host:     viper.GetString("influx.host"),
		Port:     viper.GetString("influx.port"),
		Username: viper.GetString("influx.username"),
		Password: viper.GetString("influx.password"),
		DB:       viper.GetString("influx.db"),
	}

	if (len(InfluxDBConfig.Host)) == 0 {
		panic("InfluxDB host undefined")
	}

	if (len(InfluxDBConfig.Port)) == 0 {
		panic("InfluxDB port undefined")
	}

	if (len(InfluxDBConfig.Username)) == 0 {
		panic("InfluxDB username undefined")
	}

	if (len(InfluxDBConfig.Password)) == 0 {
		panic("InfluxDB password undefined")
	}

	if (len(InfluxDBConfig.DB)) == 0 {
		panic("InfluxDB db undefined")
	}

}
