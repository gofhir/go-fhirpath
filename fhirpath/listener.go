package fhirpath

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gofhir/go-fhirpath/fhirpath/handlers"
	parser "github.com/gofhir/go-fhirpath/fhirpath/parser/grammar"
	"strings"

	_ "github.com/gofhir/go-fhirpath/fhirpath/existence"
	_ "github.com/gofhir/go-fhirpath/fhirpath/other"
)

// ExtractValueListener Custom listener for extracting a FHIRPath value
type ExtractValueListener struct {
	*parser.BasefhirpathListener
	ResourceJSON []byte
	Result       string
}

func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	
	resourceType, _, _, err := jsonparser.Get(l.ResourceJSON, "resourceType")
	if err != nil {
		l.Result = `[]`
		return
	}

	fhirPath := strings.TrimPrefix(ctx.GetText(), string(resourceType)+".")

	fmt.Printf("FHIRPath: %s\n", fhirPath)

	// Try extracting value directly
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(fhirPath, ".")...)
	if err == nil {
		l.Result = formatValue(value, dataType)
		return
	}

	// Handle nested paths
	pathParts := strings.Split(fhirPath, ".")
	if len(pathParts) < 2 {
		l.Result = `"field not exist"`
		return
	}

	lastKey := pathParts[len(pathParts)-1]

	// Extract value from the computed parent path
	value, dataType, _, err = jsonparser.Get([]byte(l.Result), lastKey)
	if err == nil {
		l.Result = formatValue(value, dataType)
		return
	}

	for _, h := range handlers.GetHandlers() {
		if strings.Contains(lastKey, h.Pattern) {
			if h.Log != "" {
				fmt.Printf("%s: %s\n", h.Log, l.Result)
			}
			l.Result = h.Func(l.Result, lastKey)
			return
		}
	}

	// Handle nested object and array extraction
	l.Result = extractFromObjectsAndArrays(l.Result, lastKey)
}

// Format JSON value appropriately
func formatValue(value []byte, dataType jsonparser.ValueType) string {
	if dataType == jsonparser.String {
		return `["` + string(value) + `"]`
	}
	return string(value)
}

// Extracts values from objects and arrays, ensuring a JSON array is always returned
func extractFromObjectsAndArrays(jsonData string, lastKey string) string {
	var extractedValues []interface{}

	_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || dataType != jsonparser.Object {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Printf("%s: %s\n", lastKey, string(data))

		subValue, subType, _, subErr := jsonparser.Get(data, strings.Split(lastKey, ".")...)
		fmt.Printf("SubValue: %s\n", subValue)
		if subErr == nil {
			if subType == jsonparser.Array {
				fmt.Println("SubValue is an Array")
				_, _ = jsonparser.ArrayEach(subValue, func(innerValue []byte, innerType jsonparser.ValueType, _ int, _ error) {
					if innerType == jsonparser.String {
						extractedValues = append(extractedValues, string(innerValue))
					} else {
						var parsedValue interface{}

						if err := json.Unmarshal(innerValue, &parsedValue); err == nil {
							extractedValues = append(extractedValues, parsedValue) // Store parsed JSON if it's valid
						} else {
							extractedValues = append(extractedValues, string(innerValue)) // Otherwise, store as a raw string

						}
					}
				})
			} else if subType == jsonparser.Object {
				fmt.Println("SubValue is an object")
				var parsedValue interface{}

				if err := json.Unmarshal(subValue, &parsedValue); err == nil {
					extractedValues = append(extractedValues, parsedValue) // Store parsed JSON if it's valid
				} else {
					extractedValues = append(extractedValues, string(subValue)) // Otherwise, store as a raw string
				}
			} else {
				extractedValues = append(extractedValues, string(subValue))
			}
		}
	})

	fmt.Printf("Extracted Values: %v\n", extractedValues)

	if len(extractedValues) == 0 {
		return `[]`
	}

	// Convert extracted values into a valid JSON array
	encoded, _ := json.Marshal(extractedValues)
	return string(encoded)
}
