package other

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gofhir/go-fhirpath/fhirpath/handlers"
	"regexp"
	"strings"
)

// HandleWhereClause Extracts values when handling `where(...)` clauses
func handleWhereClause(jsonData string, lastKey string) string {
	key, val, err := parseWhereClause(lastKey)
	if err != nil {
		return `"field not exist"`
	}

	var extractedValues []string
	_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || dataType != jsonparser.Object {
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if result[key] == val {
			fmt.Println("Matched:", result)
			extractedValues = append(extractedValues, string(data))
		}
	})

	return "[" + strings.Join(extractedValues, ",") + "]"
}

// Parses a `where(...)` condition and extracts key-value pairs
func parseWhereClause(query string) (string, string, error) {
	re := regexp.MustCompile(`where\((\w+)='([^']+)'\)`)
	matches := re.FindStringSubmatch(query)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid where clause format")
	}

	return matches[1], matches[2], nil
}

func init() {
	handlers.RegisterHandler("where", "Where Clause", handleWhereClause)
}
