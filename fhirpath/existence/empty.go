package existence

import (
	"fmt"
	"github.com/gofhir/go-fhirpath/fhirpath/handlers"
)

func handleEmptyClause(jsonData string) string {
	return "[" + fmt.Sprintf("%t", jsonData == "[]") + "]" // Return true if the array is empty
}

func init() {
	handlers.RegisterHandler("empty", "Empty Clause", func(result string, _ string) string {
		return handleEmptyClause(result)
	})
}
