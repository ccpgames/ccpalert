package engine

import "testing"

func TestBasicAlertRules(t *testing.T) {
	alertCondition1 := func(value float64) bool {
		if value > 10 {
			return true
		}
		return false
	}

	AddAlert("alert1", "metric1", "alert triggered", alertCondition1)

	result1, _ := Check("metric1", 11)
	if !result1 {
		t.Error("Alert rule should have been triggered")
	}

	result2, _ := Check("metric1", 3)
	if result2 {
		t.Error("Alert rule should not have been triggered")
	}

}
