package main

import (
	"flag"

	"github.com/ccpgames/ccpalert/api"
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
	apiConfig := api.Config{Engine: engineInstance}
	apiInstance := api.NewAPI(&apiConfig)

	go apiInstance.ServeAPI()
	go dbscheduler.Schedule()

}
