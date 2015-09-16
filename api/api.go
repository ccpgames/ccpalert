package api

import (
	"encoding/json"
	"net/http"

	"github.com/ccpgames/ccpalert/ccpalertql"
	"github.com/ccpgames/ccpalert/engine"
)

type (
	//CCPAlertAPI represents an instance of the API
	CCPAlertAPI struct {
		Config *Config
	}

	//Config represents the configuration for the API
	Config struct {
		Engine *engine.AlertEngine
	}

	ruleRequest struct {
		RawAlertStatement string
	}

	checkRequest struct {
		Key   string
		Value float64
	}
)

//NewAPI returns a new isntance of CCPAlertAPI
func NewAPI(c *Config) *CCPAlertAPI {
	return &CCPAlertAPI{Config: c}
}

func (api *CCPAlertAPI) addRule(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var query ruleRequest
	err := decoder.Decode(&query)

	if err != nil {
		http.Error(w, "invalid rule", 500)
		return
	}

	rule, err := ccpalertql.ParseAlertStatement(query.RawAlertStatement)
	if err != nil {
		http.Error(w, "malformed rule", 500)
	} else {
		api.Config.Engine.AddRule(rule)
	}

}

//
// func check(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var request checkRequest
// 	err := decoder.Decode(&request)
//
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	if err != nil {
// 		http.Error(w, "invalid request", 500)
// 	}
//
// 	err = engine.Check(request.Key, request.Value)
//
// 	if err != nil {
// 		http.Error(w, "rule not found", 500)
// 	}
// }
//

//ServeAPI serves the ccpalert API on port 8080
func (api *CCPAlertAPI) ServeAPI() {
	server := http.NewServeMux()
	server.HandleFunc("/addRule", api.addRule)
	//server.HandleFunc("/check", check)
	http.ListenAndServe(":8080", server)
}
