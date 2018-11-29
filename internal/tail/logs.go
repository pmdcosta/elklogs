package tail

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pmdcosta/elklogs/internal/domain"
)

// default timestamp field name
const timestampField = "@timestamp"

// regexp for parsing out format fields
var formatRegexp = regexp.MustCompile("%[A-Za-z0-9@_.-]+")

// GetFields gets the format fields from a string
func GetFields(format string) []string {
	return formatRegexp.FindAllString(format, -1)
}

func processEntry(e *domain.LogEntry, showTime bool, format string, fields []string) (string, error) {
	// unmarshal the log entry
	var entry map[string]interface{}
	err := json.Unmarshal(*e.Message, &entry)
	if err != nil {
		return "", err
	}

	// if no fields were provided, print the log as is
	if len(fields) == 0 {
		return fmt.Sprintf("%s", *e.Message), nil
	}

	// build the log entry based on the provided output format
	result := format
	for _, f := range fields {
		value, err := evaluateExpression(entry, f[1:])
		if err != nil {
			continue
		}
		result = strings.Replace(result, f, strings.Trim(value, "\n"), -1)
	}

	if showTime {
		value, err := evaluateExpression(entry, timestampField)
		if err != nil {
			return result, nil
		}
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return result, nil
		}
		result = fmt.Sprintf("%s: %s", t.Format("2006-01-02T15:04:05"), result)
	}

	return result, nil
}

// EvaluateExpression Expression evaluation function. It uses map as a model and evaluates expression given as the parameter using dot syntax:
// "foo" evaluates to model[foo]
// "foo.bar" evaluates to model[foo][bar]
// If a key given in the expression does not exist in the model, function will return empty string and an error.
func evaluateExpression(model interface{}, fieldExpression string) (string, error) {
	if fieldExpression == "" {
		return fmt.Sprintf("%v", model), nil
	}
	parts := strings.SplitN(fieldExpression, ".", 2)
	expression := parts[0]
	var nextModel interface{}
	modelMap, ok := model.(map[string]interface{})
	if ok {
		value := modelMap[expression]
		if value != nil {
			nextModel = value
		} else {
			return "", fmt.Errorf("failed to evaluate expression %s on given model (model map does not contain that key?)", fieldExpression)
		}
	} else {
		return "", fmt.Errorf("model on which %s is to be evaluated is not a map", fieldExpression)
	}
	nextExpression := ""
	if len(parts) > 1 {
		nextExpression = parts[1]
	}
	return evaluateExpression(nextModel, nextExpression)
}

func printLogs(entries []string, reverse bool) {
	if !reverse {
		for i := len(entries) - 1; i >= 0; i-- {
			fmt.Println(entries[i])
		}
	} else {
		for i := 0; i < len(entries); i++ {
			fmt.Println(entries[i])
		}
	}
}
