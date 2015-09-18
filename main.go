package main

import (
	"flag"

	"github.com/ccpgames/ccpalert/api"
	"github.com/ccpgames/ccpalert/ccpalertql"
	"github.com/ccpgames/ccpalert/config"
	"github.com/ccpgames/ccpalert/db"
	"github.com/ccpgames/ccpalert/engine"
)

func main() {
	configFilePath := flag.String("config", "", "configuration file")
	flag.Parse()

	configFile, err := config.ReadConfig(*configFilePath)

	if err != nil {
		panic("Unable to read config")
	}

	topLevelConfig := config.ParseConfig(configFile)

	engineInstance := engine.NewAlertEngine(&topLevelConfig.AlertEngineConfig)
	dbscheduler := db.NewScheduler(&topLevelConfig.InfluxDBConfig, *engineInstance)
	parser := ccpalertql.NewParser(engineInstance, dbscheduler)
	apiInstance := api.NewAPI(engineInstance, parser)

	err = parser.ParseAlertStatement("ALERT foobar IF foo < 2 TEXT \"oh gnoes\"")
	if err != nil {
		panic(err)
	}
	dbscheduler.AddQuery("foo", "public", "select last(value) from eve_client_disconnect")

	go dbscheduler.Schedule()
	apiInstance.ServeAPI()

}
