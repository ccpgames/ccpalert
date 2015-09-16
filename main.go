package main

import (
	"flag"

	"github.com/ccpgames/ccpalert/api"
	"github.com/ccpgames/ccpalert/config"
)

func main() {
	configFilePath := flag.String("config", "", "configuration file")
	flag.Parse()

	configFile, err := config.ReadConfig(*configFilePath)

	if err != nil {
		panic("Unable to read config")
	}

	config.ParseConfig(configFile)
	api.ServeAPI()
}
