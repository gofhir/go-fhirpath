package existence

import (
	"testing"
)

func Test_parseExistClause(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name     string
		args     args
		expected ExistQuery
		wantErr  bool
	}{
		{
			name: "Valid exists() clause",
			args: args{
				query: "exists($this is Practitioner)",
			},
			expected: ExistQuery{
				TypeCheck: "Practitioner",
				Mode:      ExistTypeCheck,
			},
			wantErr: false,
		},
		{
			name: "Invalid exists() clause",
			args: args{
				query: "exists($this",
			},
			expected: ExistQuery{},
			wantErr:  true,
		},
		{
			name: "Invalid exists() clause",
			args: args{
				query: "exists($this)",
			},
			expected: ExistQuery{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExistClause(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExistClause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Mode != tt.expected.Mode {
				t.Errorf("parseExistClause() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_handleExistClause(t *testing.T) {
	type args struct {
		jsonData string
		fhirPath string
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
				fhirPath: "exists()",
			},
			expected: "[false]",
		},
		{
			name: "Non-empty array",
			args: args{
				jsonData: "[1]",
				fhirPath: "exists()",
			},
			expected: "[true]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleExistClause(tt.args.jsonData, tt.args.fhirPath); got != tt.expected {
				t.Errorf("handleExistClause() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
