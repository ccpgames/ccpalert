package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ccpgames/ccpalert/engine"
)

type (
	ruleRequest struct {
		Name      string
		Recipient string
		Series    string
		Text      string
		RawRule   string
	}

	checkRequest struct {
		Name   string
		Series string
		Field  string
		Value  float64
	}
)

var (
	//PagerDutyAPIKey is used to route messages to pager duty
	PagerDutyAPIKey string
)

type testAlert struct {
	result chan bool
}

//alert is a dummy alert which writes a bool to a channel and
//does not use the text string
func (dummy testAlert) Send(text string) error {
	fmt.Println("ALERT HAS BEEN TRIGGERED.........")
	return nil
}

func addRule(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request ruleRequest
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(w, "invalid rule", 500)
	}

	field, rule, err := ParseRule(request.RawRule)
	if err != nil {
		http.Error(w, "malformed rule", 500)
	}

	if request.Recipient == "" {
		alertMessage := new(engine.PagerDutyAlert)
		alertMessage.APIKey = PagerDutyAPIKey
	}

	alert1 := new(testAlert)

	engine.AddAlert(request.Name, request.Series, field, request.Text, rule, alert1)
}

func check(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request checkRequest
	err := decoder.Decode(&request)

	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		http.Error(w, "invalid request", 500)
	}

	err = engine.Check(request.Name, request.Series, request.Field, request.Value)

	if err != nil {
		http.Error(w, "rule not found", 500)
	}
}

//ServeAPI serves the ccpalert API on port 8080
func ServeAPI() {
	server := http.NewServeMux()
	server.HandleFunc("/addRule", addRule)
	server.HandleFunc("/check", check)
	http.ListenAndServe(":8080", server)
}
