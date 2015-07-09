package main

import "github.com/ccpgames/ccpalert/api"

func main() {
	api.PagerDutyAPIKey = "b424c48523b043e88138cfa874ac70fe"
	api.ServeAPI()
}
