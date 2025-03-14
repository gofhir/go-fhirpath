package main

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/buger/jsonparser"
	parser "go-fhirpath/parser/grammar"
	"strings"
)

// Custom listener for extracting a FHIRPath value
type ExtractValueListener struct {
	*parser.BasefhirpathListener
	ResourceJSON []byte
	Result       string
}

/*
func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	fhirPath := ctx.GetText()                           // Example: "Patient.name.given"
	fhirPath = strings.TrimPrefix(fhirPath, "Patient.") // Remove "Patient" prefix
	//jsonPath := strings.ReplaceAll(fhirPath, ".", "/")  // Convert FHIRPath to JSON path

	// Extract the value directly
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(fhirPath, ".")...)

	fmt.Printf("Value: %s, Type: %s, Error: %v\n", value, dataType, err)

	// Case 1: Direct match (primitive values or full object)
	if err == nil {
		if dataType == jsonparser.String {
			l.Result = `"` + string(value) + `"` // Return primitive values properly formatted
		} else {
			l.Result = string(value) // Return JSON arrays/objects as-is
		}
		return
	}

	fmt.Printf("l.result: %s\n", l.Result)

	// Case 2: Handling nested arrays (e.g., name[0].given)
	parentPathParts := strings.Split(fhirPath, ".")

	fmt.Printf("parentPathParts: %v\n", parentPathParts)
	if len(parentPathParts) < 2 {
		l.Result = `"field not exist"`
		return
	}
	parentPath := strings.Join(parentPathParts[:len(parentPathParts)-1], ".")

	fmt.Printf("parentPath: %s\n", parentPath)
	lastKey := parentPathParts[len(parentPathParts)-1]

	fmt.Printf("lastKey: %s\n", lastKey)

	parentValue, parentType, _, parentErr := jsonparser.Get(l.ResourceJSON, strings.Split(parentPath, ".")...)

	fmt.Printf("parentValue: %s, parentType: %s, parentErr: %v\n", parentValue, parentType, parentErr)
	if parentErr == nil && parentType == jsonparser.Array {
		var extractedValues []string
		_, _ = jsonparser.ArrayEach(parentValue, func(data []byte, dataType jsonparser.ValueType, _ int, _ error) {

			fmt.Printf("data: %s\n", string(data))
			subValue, subType, _, subErr := jsonparser.Get(data, lastKey)

			fmt.Printf("subValue: %s, subType: %s, Error: %v\n", subValue, subType, subErr)
			if subErr == nil {
				if subType == jsonparser.Array {
					// Extract inner array elements and preserve array structure
					_, _ = jsonparser.ArrayEach(subValue, func(innerValue []byte, innerType jsonparser.ValueType, _ int, _ error) {
						if innerType == jsonparser.String {
							extractedValues = append(extractedValues, `"`+string(innerValue)+`"`)
						}
					})
				} else if subType == jsonparser.String {
					extractedValues = append(extractedValues, `"`+string(subValue)+`"`)

					fmt.Printf("extractedValues: %v\n", extractedValues)
				}
			}
		})

		// Automatically determine if result should be an array or a single value
		if len(extractedValues) > 1 {
			l.Result = "[" + strings.Join(extractedValues, ",") + "]"
		} else if len(extractedValues) == 1 {
			if parentType == jsonparser.Array {
				l.Result = "[" + extractedValues[0] + "]" // Preserve array structure
			} else {
				l.Result = extractedValues[0] // Single value case
			}
		}
		return
	}

	// Case 3: Field does not exist
	l.Result = `"field not exist"`
}

*/

/*
func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	fhirPath := ctx.GetText()                           // Example: "Patient.name.given"
	fhirPath = strings.TrimPrefix(fhirPath, "Patient.") // Remove "Patient" prefix
	jsonPath := strings.ReplaceAll(fhirPath, ".", "/")  // Convert FHIRPath to JSON path

	fmt.Printf("jsonPath: %s\n", jsonPath)

	// Extract the value directly
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(jsonPath, "/")...)

	// Case 1: Direct match
	if err == nil {
		if dataType == jsonparser.String {
			l.Result = string(value) // Return primitive values as strings
		} else if dataType == jsonparser.Array || dataType == jsonparser.Object {
			l.Result = string(value) // Return JSON arrays/objects as a string
		}
		return
	}

	// Case 2: Handling nested arrays (e.g., name[0].given)
	parentPath := strings.Join(strings.Split(jsonPath, "/")[:len(strings.Split(jsonPath, "/"))-1], "/")

	fmt.Printf("Parent path is: %s\n", parentPath)
	lastKey := strings.Split(jsonPath, "/")[len(strings.Split(jsonPath, "/"))-1]

	fmt.Printf("Last key is: %s\n", lastKey)

	parentValue, parentType, _, parentErr := jsonparser.Get(l.ResourceJSON, strings.Split(parentPath, "/")...)

	fmt.Printf("Parent type is: %s\n", parentType)
	if parentErr == nil && parentType == jsonparser.Array {
		// Iterate through array and extract values
		var extractedValues []string
		_, _ = jsonparser.ArrayEach(parentValue, func(data []byte, dataType jsonparser.ValueType, _ int, _ error) {

			subValue, subType, _, subErr := jsonparser.Get(data, lastKey)
			if subErr == nil {
				if subType == jsonparser.Array {
					// Extract elements from inner array to avoid double nesting
					_, _ = jsonparser.ArrayEach(subValue, func(innerValue []byte, innerType jsonparser.ValueType, _ int, _ error) {
						if innerType == jsonparser.String {
							extractedValues = append(extractedValues, string(innerValue))
						}
					})
				} else if subType == jsonparser.String {
					extractedValues = append(extractedValues, string(subValue))
				} else if subType == jsonparser.Object {
					subValue, _, _, subErr := jsonparser.Get(data, lastKey)
					if subErr == nil {
						extractedValues = append(extractedValues, string(subValue))
					}
				}
			}
		})

		// Convert extracted values to JSON array
		if len(extractedValues) > 0 {
			l.Result = "[" + strings.Join(extractedValues, ",") + "]"
			return
		}
	}

	// Case 3: Field does not exist
	l.Result = `"field not exist"`
}

*/

/*
func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	fhirPath := ctx.GetText() // Example: "Patient.identifier.type.coding.system"

	fmt.Printf("fhirPath: %s\n", fhirPath)
	fhirPath = strings.TrimPrefix(fhirPath, "Patient.") // Remove "Patient" prefix
	jsonPath := strings.ReplaceAll(fhirPath, ".", "/")  // Convert FHIRPath to JSON path

	// Attempt direct extraction first
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(jsonPath, "/")...)

	if err == nil {
		l.Result = formatJSON(value, dataType)
		return
	}

	// Handle nested fields inside arrays (e.g., identifier[].type.coding)
	pathParts := strings.Split(jsonPath, "/")
	if len(pathParts) < 2 {
		l.Result = `"field not exist"`
		return
	}

	parentPath := strings.Join(pathParts[:len(pathParts)-1], "/") // Parent path (e.g., "identifier")
	lastKey := pathParts[len(pathParts)-1]                        // Last key (e.g., "type.coding")

	// Extract parent node
	parentValue, parentType, _, parentErr := jsonparser.Get(l.ResourceJSON, strings.Split(parentPath, "/")...)

	if parentErr == nil {
		var extractedValues []string

		// ✅ **If Parent is an array (`identifier[]`), extract each object's nested fields (`type.coding[]`).**
		if parentType == jsonparser.Array {
			_, _ = jsonparser.ArrayEach(parentValue, func(data []byte, dataType jsonparser.ValueType, _ int, _ error) {
				subValue, subType, _, subErr := jsonparser.Get(data, strings.Split(lastKey, "/")...)
				if subErr == nil {
					extractedValues = append(extractedValues, formatJSON(subValue, subType))
				}
			})
		} else {
			// ✅ **If Parent is an object, directly extract the value.**
			subValue, subType, _, subErr := jsonparser.Get(parentValue, strings.Split(lastKey, "/")...)
			if subErr == nil {
				l.Result = formatJSON(subValue, subType)
				return
			}
		}

		// ✅ **If multiple values exist, return them as an array.**
		if len(extractedValues) > 1 {
			l.Result = "[" + strings.Join(extractedValues, ",") + "]"
		} else if len(extractedValues) == 1 {
			l.Result = extractedValues[0] // Preserve object/array structure
		} else {
			l.Result = `"field not exist"`
		}
		return
	}

	// ❌ If the field is still not found, return `"field not exist"`
	l.Result = `"field not exist"`
}

*/

func (l *ExtractValueListener) ExitInvocationExpression(ctx *parser.InvocationExpressionContext) {
	fhirPath := ctx.GetText()
	fmt.Printf("fhirPath: %s\n", fhirPath)
	fhirPath = strings.TrimPrefix(fhirPath, "Patient.")

	// Extract the value directly
	value, dataType, _, err := jsonparser.Get(l.ResourceJSON, strings.Split(fhirPath, ".")...)

	if err == nil {
		if dataType == jsonparser.String {
			l.Result = `"` + string(value) + `"`
		} else {
			l.Result = string(value)
		}
		return
	}

	pathParts := strings.Split(fhirPath, ".")
	if len(pathParts) < 2 {
		// TODO: Response maybe an empty array??
		l.Result = `"field not exist"`
		return
	}

	parentPath := strings.Join(pathParts[:len(pathParts)-1], ".")
	fmt.Printf("Parent path: %s\n", parentPath)
	lastKey := pathParts[len(pathParts)-1]
	fmt.Printf("Last key: %s\n", lastKey)

	var extractedValues []string

	_, _ = jsonparser.ArrayEach([]byte(l.Result), func(data []byte, dataType jsonparser.ValueType, offset int, err error) {

		if dataType == jsonparser.Object {
			subValue, _, _, subErr := jsonparser.Get(data, strings.Split(lastKey, ".")...)
			if subErr == nil {
				extractedValues = append(extractedValues, string(subValue))
			}
		}
	})

	fmt.Printf("Extracted values: %s\n", strings.Join(extractedValues, ", "))

	if len(extractedValues) > 0 {
		l.Result = "[" + strings.Join(extractedValues, ",") + "]" // Return as an array
	}

}

// ✅ **Helper function to properly format JSON output**
func formatJSON(value []byte, dataType jsonparser.ValueType) string {
	if dataType == jsonparser.String {
		return `"` + string(value) + `"`
	}
	return string(value) // Return raw JSON objects or arrays
}

// Function to extract a value from a FHIR resource using FHIRPath
func extractFHIRPathValue(resourceJSON string, fhirPathExpr string) string {

	resourceBytes := []byte(resourceJSON)

	// Create an ANTLR input stream
	input := antlr.NewInputStream(fhirPathExpr)

	// Create a lexer
	lexer := parser.NewfhirpathLexer(input)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create a parser
	p := parser.NewfhirpathParser(tokenStream)
	p.BuildParseTrees = true

	// Attach a custom listener
	listener := &ExtractValueListener{ResourceJSON: resourceBytes}

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
      "use": "official",
      "given": [
        "Peter",
        "James"
      ],
      "family": "Chalmers"
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
	fhirPathExpr := "Patient.name.given"

	// Extract value using FHIRPath
	result := extractFHIRPathValue(patientJSON, fhirPathExpr)

	// Print the result
	fmt.Println("Extracted Value:", result)
}
