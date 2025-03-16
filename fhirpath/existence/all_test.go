package existence

import (
	"testing"
)

func Test_parseAllClause(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name     string
		args     args
		expected string
		wantErr  bool
	}{
		{
			name: "Valid all() clause",
			args: args{
				query: "all($this is Practitioner)",
			},
			expected: "Practitioner",
			wantErr:  false,
		},
		{
			name: "Invalid all() clause",
			args: args{
				query: "all($this is)",
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "Invalid all() clause",
			args: args{
				query: "all($this is Practitioner",
			},
			expected: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAllClause(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAllClause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parseAllClause() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_handleAllClause(t *testing.T) {
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
			name: "All elements match",
			args: args{
				jsonData: `[
					{"reference": "Practitioner/1"},
					{"reference": "Practitioner/2"},
					{"reference": "Practitioner/3"}
				]`,
				fhirPath: "all($this is Practitioner)",
			},
			expected: "[true]",
		},
		{
			name: "Not all elements match",
			args: args{
				jsonData: `[
					{"reference": "Practitioner/1"},
					{"reference": "Practitioner/2"},
					{"reference": "Patient/3"}
				]`,
				fhirPath: "all($this is Practitioner)",
			},
			expected: "[false]",
		},
		{
			name: "Empty collection",
			args: args{
				jsonData: "[]",
				fhirPath: "all($this is Practitioner)",
			},
			expected: "[true]",
		},
		{
			name: "Invalid all() clause",
			args: args{
				jsonData: "[]",
				fhirPath: "all($this is)",
			},
			expected: "[false]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleAllClause(tt.args.jsonData, tt.args.fhirPath); got != tt.expected {
				t.Errorf("handleAllClause() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
