package ccpalertql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ccpgames/ccpalert/db"
	"github.com/ccpgames/ccpalert/engine"
)

type (
	//Parser represents an instance of the CCPAlertQL parser
	Parser struct {
		Scheduler *db.Scheduler
		Engine    *engine.AlertEngine
	}

	//Result represents any outcome of parsing which should be returned to the user
	Result struct {
		OK         bool
		err        error
		ResultList []string
	}
)

//NewParser returns a new instance of the CCPAlertQL parser
func NewParser(engine *engine.AlertEngine, scheduler *db.Scheduler) *Parser {
	return &Parser{Scheduler: scheduler, Engine: engine}
}

//Parse identifies the query and calls the apppropriate parser function
func (p *Parser) Parse(query string) error {
	if len(query) == 0 {
		err := fmt.Errorf("Unable to parse query")
		return err
	}

	switch strings.Fields(query)[0] {
	case "ALERT":
		p.ParseAlertStatement(query)
	}

	return nil
}

//ParseScheduleStatement parses a schedule query and schedules the contained InfluxDB query
//A schedule statement takes the form of:
//SCHEDULE INFLUXDB <influxdb query>
//To give examples:
//SCHEDULE INFLUXDB "SELECT last(value) from myseries"
func (p *Parser) ParseScheduleStatement(scheduleStatment string) error {
	scanner := NewScanner(scheduleStatment)
	tokens := scanner.scan()

	if tokens[0].tokenType != ALERT {
		err := fmt.Errorf("found %q, expected SCHEDULE", tokens[0].literal)
		return err
	}

	if tokens[1].tokenType != IDENTIFIER {
		err := fmt.Errorf("found %q, expected IDENTIFIER", tokens[0].literal)
		return err
	}
	key := tokens[1].literal

	if tokens[2].tokenType != INFLUXDB {
		err := fmt.Errorf("found %q, expected INFLUXDB", tokens[0].literal)
		return err
	}

	if tokens[3].tokenType != STRING {
		err := fmt.Errorf("found %q, expected INFLUXDB", tokens[0].literal)
		return err
	}

	query := tokens[3].literal

	if len(tokens) > 4 {
		err := fmt.Errorf("trailing characters %q", tokens[8].literal)
		return err
	}

	//Check the the encapsulted InfluxDB query is valid
	if _, err := p.Scheduler.ExecuteQuery(query); err != nil {
		return err
	}

	p.Scheduler.AddQuery(key, query)
	return nil
}

//ParseAlertStatement takes a raw alert statement query and parses it to a Rule struct
//An alert statement stakes the form:
//ALERT <alert name> IF <metric name> <operator> <threshold value> TEXT <description of alert>
//To give examples:
//ALERT cpuOnFireAlert IF superImportantServer.cpuUsage > 100 TEXT "Critical production server is heavily loaded"
//ALERT noplayers IF tq.currentPlayers == 0 TEXT "something has gone badly wrong"
func (p *Parser) ParseAlertStatement(alertStatement string) (engine.Rule, error) {
	scanner := NewScanner(alertStatement)
	tokens := scanner.scan()
	newRule := new(engine.Rule)

	if tokens[0].tokenType != ALERT {
		err := fmt.Errorf("found %q, expected ALERT", tokens[0].literal)
		return engine.Rule{}, err
	}

	if tokens[1].tokenType == IDENTIFIER {
		newRule.Name = tokens[1].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[1].literal)
		return engine.Rule{}, err
	}

	if tokens[2].tokenType != IF {
		err := fmt.Errorf("found %q, expected IF", tokens[2].literal)
		return engine.Rule{}, err
	}

	if tokens[3].tokenType == IDENTIFIER {
		newRule.MetricKey = tokens[3].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[3].literal)
		return engine.Rule{}, err
	}

	if tokens[4].tokenType != OP {
		err := fmt.Errorf("found %q, expected <,> or ==", tokens[4].literal)
		return engine.Rule{}, err
	}

	if tokens[5].tokenType != VALUE {
		err := fmt.Errorf("found %q, expected value", tokens[5].literal)
		return engine.Rule{}, err
	}

	threshold, err := strconv.ParseFloat(tokens[5].literal, 64)

	if err != nil {
		return engine.Rule{}, err
	}

	condition, err := NewCondition(tokens[4].literal, threshold)

	if err == nil {
		newRule.Condition = condition
	} else {
		return engine.Rule{}, err
	}

	if tokens[6].tokenType != TEXT {
		err := fmt.Errorf("found %q, expected TEXT", tokens[6].literal)
		return engine.Rule{}, err
	}

	if tokens[7].tokenType != STRING {
		err := fmt.Errorf("found %q, expected string", tokens[7].literal)
		return engine.Rule{}, err
	}

	if len(tokens) > 8 {
		err := fmt.Errorf("trailing characters %q", tokens[8].literal)
		return engine.Rule{}, err
	}

	return *newRule, nil
}
