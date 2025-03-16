package other

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"go-fhirpath/fhirpath/handlers"
)

// HandleFirstElement Handles `first()` function extraction
func handleFirstElement(jsonData string) string {
	value, _, _, err := jsonparser.Get([]byte(jsonData), "[0]")
	if err != nil {
		return `[]`
	}

	return formatArray(value)
}

// Format JSON value as an array
func formatArray(value []byte) string {
	var result []string

	result = []string{string(value)}

	encoded, _ := json.Marshal(result)
	return string(encoded)
}

func init() {
	handlers.RegisterHandler("first", "First Element", func(result string, key string) string {
		return handleFirstElement(result)
	})
}
