package existence

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gofhir/go-fhirpath/fhirpath/handlers"
	"regexp"
	"strings"
)

func parseAllClause(query string) (string, error) {
	// Match: all($this is Practitioner)
	re := regexp.MustCompile(`all\(\$this\s*is\s*(\w+)\)`)
	matches := re.FindStringSubmatch(query)

	if len(matches) != 2 {
		return "", fmt.Errorf("invalid all() type-check format")
	}

	return matches[1], nil // returns "Practitioner"
}

func handleAllClause(jsonData string, fhirPath string) string {
	typeCheck, err := parseAllClause(fhirPath)
	if err != nil {
		return "[false]"
	}

	count := 0
	allMatch := true

	_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || dataType != jsonparser.Object {
			return
		}

		count++
		ref, _ := jsonparser.GetString(data, "reference")
		if !strings.HasPrefix(ref, typeCheck+"/") {
			allMatch = false
		}
	})

	// Spec rule: empty collection = true
	if count == 0 {
		return "[true]"
	}

	return "[" + fmt.Sprintf("%t", allMatch) + "]" // Return true if all elements match
}

func init() {
	handlers.RegisterHandler("all", "All Clause", handleAllClause)
}
