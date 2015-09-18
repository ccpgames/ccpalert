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
		Engine *engine.AlertEngine
		Parser *ccpalertql.Parser
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
func NewAPI(e *engine.AlertEngine, p *ccpalertql.Parser) *CCPAlertAPI {
	return &CCPAlertAPI{Engine: e, Parser: p}
}

func (api *CCPAlertAPI) query(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var query ruleRequest
	err := decoder.Decode(&query)

	if err != nil {
		http.Error(w, "invalid rule", 500)
		return
	}

	_, err = api.Parser.Parse(query.RawAlertStatement)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//ServeAPI serves the ccpalert API on port 8080
func (api *CCPAlertAPI) ServeAPI() {
	server := http.NewServeMux()
	server.HandleFunc("/query", api.query)
	http.ListenAndServe(":8080", server)
}
