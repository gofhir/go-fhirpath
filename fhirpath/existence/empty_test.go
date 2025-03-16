package existence

import (
	"go-fhirpath/fhirpath/handlers"
	"testing"
)

func Test_handleEmptyClause(t *testing.T) {
	type args struct {
		jsonData string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "Empty array",
			args: args{
				jsonData: "[]",
			},
			expected: "[true]",
		},
		{
			name: "Non-empty array",
			args: args{
				jsonData: "[1]",
			},
			expected: "[false]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleEmptyClause(tt.args.jsonData); got != tt.expected {
				t.Errorf("handleEmptyClause() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_EmptyHandlerIsRegistered(t *testing.T) {

	if ok := handlers.IsHandlerRegistered("empty"); !ok {
		t.Error("empty handler is not registered")
	}
}
