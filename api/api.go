package api

type (
	ruleRequest struct {
		Name      string
		Recipient string
		Key       string
		Text      string
		RawRule   string
	}

	checkRequest struct {
		Key   string
		Value float64
	}
)

// func addRule(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var request ruleRequest
// 	err := decoder.Decode(&request)
//
// 	if err != nil {
// 		http.Error(w, "invalid rule", 500)
// 	}
//
// 	rule, err := ParseRule(request.RawRule)
// 	if err != nil {
// 		http.Error(w, "malformed rule", 500)
// 	}
//
// 	//engine.AddAlert(request.Name, request.Key, request.Text, rule)
// }
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
// //ServeAPI serves the ccpalert API on port 8080
// func ServeAPI() {
// 	server := http.NewServeMux()
// 	server.HandleFunc("/addRule", addRule)
// 	server.HandleFunc("/check", check)
// 	http.ListenAndServe(":8080", server)
// }
