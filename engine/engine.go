package engine

import (
	"log"
	"net/smtp"
	"strconv"

	"github.com/stvp/pager"
)

type (
	//AlertEngine represents the central component of CCP Alert which
	//stores alerting rules, checks rules and triggers alerts
	AlertEngine struct {
		Config *Config
		Rules  map[string]map[string]Rule
	}

	//Config represents the configuration for the AlertEngine
	Config struct {
		EmailRecipient  string
		EmailUsername   string
		EmailPassword   string
		EmailServer     string
		EmailPort       int
		PagerDutyAPIKey string
	}

	//Rule represents a single rule for checking a single metric
	//including the condition and action to take when the rule is triggered
	Rule struct {
		Name      string
		MetricKey string
		Condition AlertCondition
		Text      string
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

//NewAlertEngine returns a new instance of AlertEngine
func NewAlertEngine(c *Config) *AlertEngine {
	s := &AlertEngine{Config: c}
	return s
}

//Send sends an alert
func (engine *AlertEngine) Send(triggeredRule Rule) error {
	if len(engine.Config.PagerDutyAPIKey) > 0 {
		pager.ServiceKey = engine.Config.PagerDutyAPIKey
		_, err := pager.Trigger(triggeredRule.Text)
		return err
	}

	if len(engine.Config.EmailServer) > 0 {
		auth := smtp.PlainAuth(
			"",
			engine.Config.EmailUsername,
			engine.Config.EmailPassword,
			engine.Config.EmailServer,
		)

		err := smtp.SendMail(
			engine.Config.EmailServer+":"+strconv.Itoa(engine.Config.EmailPort),
			auth,
			engine.Config.EmailUsername,
			[]string{engine.Config.EmailRecipient},
			[]byte(triggeredRule.Text),
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

//CreateRule creates a new AlertRule and registers it
func (engine *AlertEngine) CreateRule(ruleName string, key string, text string, condition AlertCondition) {
	rule := new(Rule)
	rule.Name = ruleName
	rule.Condition = condition
	rule.Text = text

	engine.AddRule(*rule)
}

//AddRule adds a new rule
func (engine *AlertEngine) AddRule(newRule Rule) {
	if engine.Rules[newRule.MetricKey] == nil {
		engine.Rules[newRule.MetricKey] = make(map[string]Rule)
	}

	engine.Rules[newRule.MetricKey][newRule.Name] = newRule
}

//Check a datapoint against a rule
func (engine *AlertEngine) Check(key string, value float64) (bool, error) {
	relatedRules := engine.Rules[key]
	ruleTriggered := false

	for _, rule := range relatedRules {
		ruleTriggered = rule.Condition(value)
		if ruleTriggered {
			err := engine.Send(rule)
			if err != nil {
				return ruleTriggered, nil
			}
		}
	}

	return ruleTriggered, nil
}
