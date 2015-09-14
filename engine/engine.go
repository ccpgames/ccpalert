package engine

import (
	"log"
	"net/smtp"
	"strconv"

	"github.com/ccpgames/ccpalert/config"
	"github.com/stvp/pager"
)

var (
	//A table (map of maps) of rules. The first map is indexed by the metric key,
	//the second table by an alert name. This allows each metric to have multiple
	//alert rules associated with it.
	rules = make(map[string]map[string]Rule)
)

type (
	//Rule represents a single rule for checking a single metric
	//including the condition and action to take when the rule is triggered
	Rule struct {
		name      string
		condition AlertCondition
		text      string
	}

	//AlertMessage is an interface defining a single function,
	//Alert, which is intended to notify a recipient of the alert
	AlertMessage interface {
		Send(text string) error
	}

	//AlertCondition is a function type for checking whether a rule is met
	//if Rule returns true, an alert is triggered
	AlertCondition func(float64) bool
)

//Send sends an alert
func (alert Rule) Send() error {

	if config.PagerDuty {
		pager.ServiceKey = config.PagerDutyAPIKey
		_, err := pager.Trigger(alert.text)
		return err
	}

	if config.Email {
		auth := smtp.PlainAuth(
			"",
			config.EmailConfig.Username,
			config.EmailConfig.Password,
			config.EmailConfig.EmailServer,
		)

		err := smtp.SendMail(
			config.EmailConfig.EmailServer+":"+strconv.Itoa(config.EmailConfig.Port),
			auth,
			config.EmailConfig.Username,
			[]string{config.EmailConfig.Recipient},
			[]byte(alert.text),
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

//AddAlert creates a new AlertRule and registers it
func AddAlert(ruleName string, key string, text string, condition AlertCondition) {
	rule := new(Rule)
	rule.name = ruleName
	rule.condition = condition
	rule.text = text

	if rules[key] == nil {
		rules[key] = make(map[string]Rule)
	}

	rules[key][ruleName] = *rule
}

//Check a datapoint against a rule
func Check(key string, value float64) (bool, error) {

	relatedRules := rules[key]
	ruleTriggered := false

	for _, rule := range relatedRules {
		ruleTriggered = rule.condition(value)
		if ruleTriggered {
			err := rule.Send()
			if err != nil {
				return ruleTriggered, nil
			}
		}
	}

	return ruleTriggered, nil
}
