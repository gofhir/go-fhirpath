package existence

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"go-fhirpath/fhirpath/handlers"
	"regexp"
	"strings"
)

type LogicOperator string

const (
	AND LogicOperator = "and"
	OR  LogicOperator = "or"
)

type ExistType int

const (
	ExistNone ExistType = iota
	ExistCondition
	ExistTypeCheck
)

type ExistQuery struct {
	Conditions map[string]string
	Logic      LogicOperator
	TypeCheck  string
	Mode       ExistType
}

// parseExistClause Parses any valid exists(...) query:
// plain exists(), conditions, or $this is type checks
func parseExistClause(query string) (ExistQuery, error) {
	re := regexp.MustCompile(`exists\((.*?)\)`)
	matches := re.FindStringSubmatch(query)
	if len(matches) != 2 {
		return ExistQuery{}, fmt.Errorf("invalid exists clause format")
	}

	conditionStr := strings.TrimSpace(matches[1])

	// Handle plain exists()
	if conditionStr == "" {
		return ExistQuery{Mode: ExistNone}, nil
	}

	// Handle $this is TYPE
	if strings.Contains(conditionStr, "$thisis") || strings.Contains(conditionStr, "$this is") {
		reType := regexp.MustCompile(`\$this\s*is\s*(\w+)`)
		m := reType.FindStringSubmatch(conditionStr)
		if len(m) == 2 {
			return ExistQuery{
				TypeCheck: m[1],
				Mode:      ExistTypeCheck,
			}, nil
		}
		return ExistQuery{}, fmt.Errorf("invalid type check expression")
	}

	// Handle normal AND / OR conditions
	var op = AND
	if strings.Contains(conditionStr, "or") {
		op = OR
		conditionStr = strings.ReplaceAll(conditionStr, "or", "|")
	} else if strings.Contains(conditionStr, "and") {
		op = AND
		conditionStr = strings.ReplaceAll(conditionStr, "and", "|")
	}

	conditions := make(map[string]string)
	reKeyValue := regexp.MustCompile(`(\w+)='([^']+)'`)
	for _, part := range strings.Split(conditionStr, "|") {
		part = strings.TrimSpace(part)
		match := reKeyValue.FindStringSubmatch(part)
		if len(match) != 3 {
			return ExistQuery{}, fmt.Errorf("invalid key-value pair inside exists: %s", part)
		}
		conditions[match[1]] = match[2]
	}

	return ExistQuery{
		Conditions: conditions,
		Logic:      op,
		Mode:       ExistCondition,
	}, nil
}

func handleExistClause(jsonData string, lastKey string) string {

	query, err := parseExistClause(lastKey)
	if err != nil {
		return "[false]"
	}

	// Case 1: exists()
	if query.Mode == ExistNone {
		return "[" + fmt.Sprintf("%t", jsonData != "[]") + "]" // Return true if the array is non-empty
	}

	// Case 2: exists($this is ResourceType)
	if query.Mode == ExistTypeCheck {
		match := false
		_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
			ref, _ := jsonparser.GetString(data, "reference")
			if strings.HasPrefix(ref, query.TypeCheck+"/") {
				match = true
			}
		})

		return "[" + fmt.Sprintf("%t", match) + "]" // Return true if the type check matches
	}

	// Case 3: exists(condition with and/or)
	if query.Mode == ExistCondition {
		matched := false
		_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return
			}

			if query.Logic == AND {
				allMatch := true
				for key, val := range query.Conditions {
					if result[key] != val {
						allMatch = false
						break
					}
				}
				if allMatch {
					matched = true
				}
			} else if query.Logic == OR {
				for key, val := range query.Conditions {
					if result[key] == val {
						matched = true
						break
					}
				}
			}
		})

		return "[" + fmt.Sprintf("%t", matched) + "]" // Return true if the key-value pair exists
	}

	return "[false]"
}

func init() {
	handlers.RegisterHandler("exists", "Exist Clause", handleExistClause)
}
