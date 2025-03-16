package fhirpath

import (
	"encoding/json"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/gofhir/go-fhirpath/fhirpath/parser/grammar"
)

func Evaluate(resourceJSON string, fhirPathExpr string) []byte {
	resourceBytes := []byte(resourceJSON)

	input := antlr.NewInputStream(fhirPathExpr)
	lexer := parser.NewfhirpathLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewfhirpathParser(tokenStream)
	p.BuildParseTrees = true

	listener := &ExtractValueListener{ResourceJSON: resourceBytes}

	antlr.ParseTreeWalkerDefault.Walk(listener, p.EntireExpression())

	// Parse JSON without escaping characters
	r, err := parseJSONString(listener.Result)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte("[]")
	}

	return r
}

// Parse a JSON string while preserving its structure without escaping characters
func parseJSONString(input string) ([]byte, error) {
	var jsonData json.RawMessage // Allows storing raw JSON without modification

	// Unmarshal the JSON string into a raw JSON structure
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		return []byte("[]"), err
	}

	// Convert back to a JSON-formatted string without escaping characters
	parsedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return []byte("[]"), err
	}

	return parsedJSON, nil
}
