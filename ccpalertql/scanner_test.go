package ccpalertql

import (
	"fmt"
	"testing"
)

//TestScan tests scan
func TestScan(t *testing.T) {
	scanner := NewScanner("ALERT foo IF bar > 10 TEXT \"HELLO WOLRD\"")
	result := scanner.scan()

	if len(result) != 8 {
		fmt.Println(result)
		t.Error("Was expecting 8 tokens, got", len(result))
	}
}
