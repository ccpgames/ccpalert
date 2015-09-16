package db

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ccpgames/ccpalert/engine"
	"github.com/influxdb/influxdb/client"
)

type (
	//DBScheduler represents an instance of the query scheduler
	DBScheduler struct {
		InfluxDBConfig *InfluxDBConfig
		Engine         engine.AlertEngine
		Stop           chan struct{}
		Queries        map[string]string
	}

	//InfluxDBConfig provides the config required to pull metrics from InfluxDB
	InfluxDBConfig struct {
		InfluxDBHost     string
		InfluxDBPort     int
		InfluxDBUsername string
		InfluxDBPassword string
		InfluxDBDB       string
	}
)

//NewScheduler returns a new instance of DBScheduler
func NewScheduler(c *InfluxDBConfig, engine engine.AlertEngine) *DBScheduler {
	return &DBScheduler{InfluxDBConfig: c, Engine: engine}
}

//AddQuery adds a query to the scheduler
func (db *DBScheduler) AddQuery(metricKey string, query string) {
	db.Queries[metricKey] = query
}

//Scedule periodically executes predefined InfluxDB queries
func (db *DBScheduler) Schedule() {
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			for key, query := range db.Queries {
				go db.scheduledCheck(key, query)
			}
		case <-db.Stop:
			ticker.Stop()
			return
		}
	}
}

func (db *DBScheduler) scheduledCheck(key, query string) {
	value, err := db.executeQuery(query)
	if err != nil {
		db.Engine.Check(key, value)
	}
}

func (db *DBScheduler) executeQuery(query string) (float64, error) {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", db.InfluxDBConfig.InfluxDBHost, db.InfluxDBConfig.InfluxDBPort))
	if err != nil {
		log.Fatal(err)
	}
	con, err := client.NewClient(client.Config{
		URL:      *host,
		Username: db.InfluxDBConfig.InfluxDBUsername,
		Password: db.InfluxDBConfig.InfluxDBPassword,
	})

	if err != nil {
		return 0, err
	}

	q := client.Query{
		Command:  query,
		Database: db.InfluxDBConfig.InfluxDBDB,
	}

	response, err := con.Query(q)
	if err != nil {
		//This somewhat unpleasant looking line goes through several arrays nested structs
		//to get to the actual value.
		value, ok := response.Results[0].Series[0].Values[0][1].(float64)
		if ok {
			return value, err
		}
		return 0, fmt.Errorf("Unable to parse value from InfluxDB query, ensure that query returns a single value and that the series contains data")
	}
	return 0, err
}
