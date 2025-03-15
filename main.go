package main

import (
	"encoding/json"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/buger/jsonparser"
	parser "go-fhirpath/parser/grammar"
	"regexp"
	"strings"
)

// Custom listener for extracting a FHIRPath value
type ExtractValueListener struct {
	*parser.BasefhirpathListener
	ResourceJSON []byte
	Result       string
	DataType     jsonparser.ValueType
}

func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	fhirPath := strings.TrimPrefix(ctx.GetText(), "Patient.")
	fmt.Printf("FHIRPath: %s\n", fhirPath)

	// Try extracting value directly
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(fhirPath, ".")...)
	if err == nil {
		l.Result = formatValue(value, dataType)
		l.DataType = dataType
		return
	}

	// Handle nested paths
	pathParts := strings.Split(fhirPath, ".")
	if len(pathParts) < 2 {
		l.Result = `"field not exist"`
		l.DataType = jsonparser.String
		return
	}

	parentPath := strings.Join(pathParts[:len(pathParts)-1], ".")
	lastKey := pathParts[len(pathParts)-1]
	fmt.Printf("Parent Path: %s, Last Key: %s\n", parentPath, lastKey)
	fmt.Printf("L.Result: %s\n", l.Result)

	// Extract value from the computed parent path
	value, dataType, _, err = jsonparser.Get([]byte(l.Result), lastKey)
	fmt.Printf("Result Value: %s\n", formatValue(value, dataType))
	fmt.Printf("Last Key: %s\n", lastKey)
	if err == nil {
		l.Result = formatValue(value, dataType)
		l.DataType = dataType
		return
	}

	// Handle `where(...)` condition
	if strings.Contains(lastKey, "where") {
		l.Result = handleWhereClause(l.Result, lastKey)
		l.DataType = jsonparser.String
		fmt.Printf("Where Clause: %s\n", l.Result)
		return
	}

	// Handle `first()`
	if strings.Contains(lastKey, "first") {
		l.Result = handleFirstElement(l.Result)
		return
	}

	// Handle nested object and array extraction
	l.Result = extractFromObjectsAndArrays(l.Result, lastKey)
}

// Format JSON value as an array
func formatArray(value []byte, dataType jsonparser.ValueType) string {
	var result []string

	if dataType == jsonparser.String {
		result = []string{string(value)}
	} else {
		result = []string{string(value)}
	}

	encoded, _ := json.Marshal(result)
	return string(encoded)
}

// Format JSON value appropriately
func formatValue(value []byte, dataType jsonparser.ValueType) string {
	if dataType == jsonparser.String {
		return `["` + string(value) + `"]`
	}
	return string(value)
}

// Extracts values when handling `where(...)` clauses
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

	return strings.Join(extractedValues, ", ")
}

// Handles `first()` function extraction
func handleFirstElement(jsonData string) string {
	value, dataType, _, err := jsonparser.Get([]byte(jsonData), "[0]")
	if err != nil {
		return `[]`
	}

	return formatArray(value, dataType)
}

// Extracts values from objects and arrays, ensuring a JSON array is always returned
func extractFromObjectsAndArrays(jsonData string, lastKey string) string {
	var extractedValues []interface{}

	fmt.Printf("Extracting from: %s\n", jsonData)
	fmt.Printf("Last Key: %s\n", lastKey)

	_, _ = jsonparser.ArrayEach([]byte(jsonData), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || dataType != jsonparser.Object {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		subValue, subType, _, subErr := jsonparser.Get(data, strings.Split(lastKey, ".")...)
		if subErr == nil {
			if subType == jsonparser.Array {
				_, _ = jsonparser.ArrayEach(subValue, func(innerValue []byte, innerType jsonparser.ValueType, _ int, _ error) {
					if innerType == jsonparser.String {
						extractedValues = append(extractedValues, string(innerValue))
					}
				})
			} else if subType == jsonparser.Object {
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

	// Convert extracted values into a valid JSON array
	encoded, _ := json.Marshal(extractedValues)
	return string(encoded)
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

// Extracts a value from a FHIR resource using FHIRPath
func extractFHIRPathValue(resourceJSON string, fhirPathExpr string) string {
	resourceBytes := []byte(resourceJSON)

	// Create ANTLR input stream
	input := antlr.NewInputStream(fhirPathExpr)
	lexer := parser.NewfhirpathLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewfhirpathParser(tokenStream)
	p.BuildParseTrees = true

	// Attach custom listener
	listener := &ExtractValueListener{ResourceJSON: resourceBytes}

	// Attach error listeners
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.RemoveErrorListeners()
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	antlr.ParseTreeWalkerDefault.Walk(listener, p.EntireExpression())

	return listener.Result
}

func main() {
	// Example FHIR Patient resource (JSON)
	patientJSON := `{
  "resourceType": "Patient",
  "id": "example",
  "address": [
    {
      "use": "home",
      "city": "PleasantVille",
      "type": "both",
      "state": "Vic",
      "line": [
        "534 Erewhon St"
      ],
      "postalCode": "3999",
      "period": {
        "start": "1974-12-25"
      },
      "district": "Rainbow",
      "text": "534 Erewhon St PeasantVille, Rainbow, Vic  3999"
    }
  ],
  "managingOrganization": {
    "reference": "Organization/1"
  },
  "name": [
    {
      "use": "usual",
      "given": [
        "Peter",
        "James"
      ],
      "family": "Chalmers",
      "period": {
		"start": "2002",
		"end": "2004"
	  }
    },
    {
      "use": "usual",
      "given": [
        "Jim"
      ]
    },
    {
      "use": "maiden",
      "given": [
        "Peter",
        "James"
      ],
      "family": "Windsor",
      "period": {
        "end": "2002"
      }
    }
  ],
  "birthDate": "1974-12-25",
  "deceased": {
    "boolean": false
  },
  "active": true,
  "identifier": [
    {
      "use": "usual",
      "type": {
        "coding": [
          {
            "code": "MR",
            "system": "http://hl7.org/fhir/v2/0203"
          }
        ]
      },
      "value": "12345",
      "period": {
        "start": "2001-05-06"
      },
      "system": "urn:oid:1.2.36.146.595.217.0.1",
      "assigner": {
        "display": "Acme Healthcare"
      }
    }
  ],
  "telecom": [
    {
      "use": "home"
    },
    {
      "use": "work",
      "rank": 1,
      "value": "(03) 5555 6473",
      "system": "phone"
    },
    {
      "use": "mobile",
      "rank": 2,
      "value": "(03) 3410 5613",
      "system": "phone"
    },
    {
      "use": "old",
      "value": "(03) 5555 8834",
      "period": {
        "end": "2014"
      },
      "system": "phone"
    }
  ],
  "gender": "male",
  "contact": [
    {
      "name": {
        "given": [
          "Bénédicte"
        ],
        "family": "du Marché",
        "_family": {
          "extension": [
            {
              "url": "http://hl7.org/fhir/StructureDefinition/humanname-own-prefix",
              "valueString": "VV"
            }
          ]
        }
      },
      "gender": "female",
      "period": {
        "start": "2012"
      },
      "address": {
        "use": "home",
        "city": "PleasantVille",
        "line": [
          "534 Erewhon St"
        ],
        "type": "both",
        "state": "Vic",
        "period": {
          "start": "1974-12-25"
        },
        "district": "Rainbow",
        "postalCode": "3999"
      },
      "telecom": [
        {
          "value": "+33 (237) 998327",
          "system": "phone"
        }
      ],
      "relationship": [
        {
          "coding": [
            {
              "code": "N",
              "system": "http://hl7.org/fhir/v2/0131"
            }
          ]
        }
      ]
    }
  ]
}`

	// Example FHIRPath expression
	fhirPathExpr := "Patient.name.where(use='usual').given"

	// Extract value using FHIRPath
	result := extractFHIRPathValue(patientJSON, fhirPathExpr)

	// Print the result
	fmt.Println("Extracted Value:", result)
}
