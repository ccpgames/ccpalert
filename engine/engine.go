package engine

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/stvp/pager"
)

//a map of all registered rules
var rules map[ruleID]alert

type (
	//AlertRule represents a single rule for checking a single metric
	//including the condition and action to take when the rule is triggered
	alert struct {
		name      string
		condition Rule
		message   AlertMessage
		text      string
	}

	//AlertMessage is an interface defining a single function,
	//Alert, which is intended to notify a recipient of the alert
	AlertMessage interface {
		Send(text string) error
	}

	//Rule is a function type for checking whether a rule is met
	//if Rule returns true, an alert is triggered
	Rule func(float64) bool

	//EmailAlert represents an alert sent via SMTP
	EmailAlert struct {
		Recipient   string
		Username    string
		Password    string
		EmailServer string
		Port        int
	}

	//PagerDutyAlert represents an alert sent via the PagerDuty API
	PagerDutyAlert struct {
		APIKey string
	}

	ruleID struct {
		ruleName, series, field string
	}
)

//Alert sends an alet via PagerDuty
func (pagerDuty PagerDutyAlert) Alert(text string) error {
	pager.ServiceKey = pagerDuty.APIKey
	_, err := pager.Trigger(text)
	return err
}

//Alert sends a alert via SMTP
func (email EmailAlert) Alert(text string) error {
	fmt.Println(text)

	auth := smtp.PlainAuth(
		"",
		email.Username,
		email.Password,
		email.EmailServer,
	)

	err := smtp.SendMail(
		email.EmailServer+":"+strconv.Itoa(email.Port),
		auth,
		email.Username,
		[]string{email.Recipient},
		[]byte(text),
	)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

//AddAlert creates a new AlertRule and registers it
func AddAlert(ruleName string, series string, field string, text string, condition Rule, message AlertMessage) {
	if rules == nil {
		rules = make(map[ruleID]alert)
	}

	rule := new(alert)
	rule.name = ruleName
	rule.condition = condition
	rule.text = text
	rule.message = message
	rules[ruleID{ruleName, series, field}] = *rule

	fmt.Println(rules)
}

//Check a datapoint against a rule
func Check(ruleName string, series string, field string, value float64) error {
	id := new(ruleID)
	id.ruleName = ruleName
	id.series = series
	rule := rules[ruleID{ruleName, series, field}]

	if rule.name == "" {
		return errors.New("rule not found")
	}

	if rule.condition(value) {
		err := rule.message.Send(rule.text)
		if err != nil {
			return err
		}
	}

	return nil
}
