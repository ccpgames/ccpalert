package engine

import (
	"testing"
	"time"
)

type testAlert struct {
	result chan bool
}

//alert is a dummy alert which writes a bool to a channel and
//does not use the text string
func (dummy testAlert) Send(text string) error {
	dummy.result <- true
	return nil
}

func TestBasicAlertRules(t *testing.T) {
	resultChan := make(chan bool, 1)

	alertCondition1 := func(value float64) bool {
		if value > 10 {
			return true
		}
		return false
	}

	alert1 := new(testAlert)
	alert1.result = resultChan

	AddAlert("testRule1", "metric1", "", alertCondition1, alert1)

	Check("testRule1", "metric1", 11)
	if !checkResult(resultChan) {
		t.Error("Alert rule should have been triggered")
	}

	Check("testRule1", "metric1", 3)
	if checkResult(resultChan) {
		t.Error("Alert rule should not have been triggered")
	}

}

func checkResult(resultChan chan bool) bool {
	select {
	case res := <-resultChan:
		return res
	case <-time.After(time.Second * 1):
		return false
	}
}
