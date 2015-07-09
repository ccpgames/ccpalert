package api

import (
	"errors"
	"strconv"
	"strings"
)

//ParseRule parses an alert rule in the form of [variable] [operator] [threshold]
//and returns a Go function representing that rule
func ParseRule(rule string) (string, func(value float64) bool, error) {
	//split rule into field, operator, threshold and check for correctness
	ruleComponents := strings.Split(rule, " ")

	if len(ruleComponents) > 3 {
		return "", nil, errors.New("malformed rule")
	}

	field := ruleComponents[0]
	operator := ruleComponents[1]
	threshold, err := strconv.ParseFloat(ruleComponents[2], 64)

	if err != nil {
		return "", nil, errors.New("unable to parse threshold")
	}
	var alertRule func(value float64) bool

	//create a function to check the alert rule against values
	if operator == "=" {
		alertRule = func(value float64) bool {
			if value == threshold {
				return true
			}
			return false
		}
	} else if operator == ">" {
		alertRule = func(value float64) bool {
			if value > threshold {
				return true
			}
			return false
		}
	} else if operator == "<" {
		alertRule = func(value float64) bool {
			if value < threshold {
				return true
			}
			return false
		}
	}

	return field, alertRule, nil
}
