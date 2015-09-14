package ccpalertql

import (
	"fmt"
	"strconv"

	"github.com/ccpgames/ccpalert/engine"
)

//ParseAlertStatement takes a raw alert statement query and parses it to a Rule struct
func ParseAlertStatement(alertStatement string) (string, engine.Rule, error) {

	scanner := NewScanner(alertStatement)
	tokens := scanner.scan()
	newRule := new(engine.Rule)
	var key string

	if tokens[0].tokenType != ALERT {
		err := fmt.Errorf("found %q, expected ALERT", tokens[0].literal)
		return "", engine.Rule{}, err
	}

	if tokens[1].tokenType == IDENTIFIER {
		newRule.Name = tokens[1].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[1].literal)
		return "", engine.Rule{}, err
	}

	if tokens[2].tokenType != IF {
		err := fmt.Errorf("found %q, expected IF", tokens[2].literal)
		return "", engine.Rule{}, err
	}

	if tokens[3].tokenType == IDENTIFIER {
		key = tokens[3].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[3].literal)
		return "", engine.Rule{}, err
	}

	if tokens[4].tokenType != OP {
		err := fmt.Errorf("found %q, expected <,> or ==", tokens[4].literal)
		return "", engine.Rule{}, err
	}

	if tokens[5].tokenType != VALUE {
		err := fmt.Errorf("found %q, expected value", tokens[5].literal)
		return "", engine.Rule{}, err
	}

	threshold, err := strconv.ParseFloat(tokens[5].literal, 64)

	if err != nil {
		return "", engine.Rule{}, err
	}

	condition, err := NewCondition(tokens[4].literal, threshold)

	if err == nil {
		newRule.Condition = condition
	} else {
		return "", engine.Rule{}, err
	}

	if tokens[6].tokenType != TEXT {
		err := fmt.Errorf("found %q, expected TEXT", tokens[6].literal)
		return "", engine.Rule{}, err
	}

	if tokens[7].tokenType != STRING {
		err := fmt.Errorf("found %q, expected string", tokens[7].literal)
		return "", engine.Rule{}, err
	}

	return key, *newRule, nil
}
