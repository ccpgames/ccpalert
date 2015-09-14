package ccpalertql

import "testing"

//TestScan tests scan
func TestParse(t *testing.T) {
	_, result, _ := ParseAlertStatement("ALERT foo IF bar > 10 TEXT \"HELLO WOLRD\"")

	ruleCheckResult := result.Condition(11)

	if ruleCheckResult != true {
		t.Error("Expected rule to be triggered")
	}
}
