package input

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ccpgames/ccpalert/config"
	"github.com/influxdb/influxdb/client"
)

//ExecuteQuery executes a query against an InfluxDB database
func ExecuteQuery(query, database string, influxConfig config.InfluxDBConfigStruct) (float64, error) {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", influxConfig.Host, influxConfig.Port))
	if err != nil {
		log.Fatal(err)
	}
	con, err := client.NewClient(client.Config{
		URL:      *host,
		Username: influxConfig.Username,
		Password: influxConfig.Password,
	})

	if err != nil {
		return 0, err
	}

	q := client.Query{
		Command:  query,
		Database: database,
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
