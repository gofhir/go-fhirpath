package handlers

import (
	"testing"
)

func Test_IsHandlerRegistered(t *testing.T) {
	Clear() // reset state

	// Pre-register "empty" handler before running the tests
	RegisterHandler("empty", "Empty Clause", func(result string, _ string) string {
		return result
	})

	type args struct {
		handlerName string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Registered handler",
			args: args{
				handlerName: "empty",
			},
			expected: true,
		},
		{
			name: "Unregistered handler",
			args: args{
				handlerName: "nonexistent",
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHandlerRegistered(tt.args.handlerName); got != tt.expected {
				t.Errorf("IsHandlerRegistered() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_RegisterHandler(t *testing.T) {
	Clear() // reset state

	// Pre-register "empty" handler before running the tests
	RegisterHandler("empty", "Empty Clause", func(result string, _ string) string {
		return result
	})

	type args struct {
		handlerName string
		description string
		handler     HandlerFunc
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Register new handler",
			args: args{
				handlerName: "newHandler",
				description: "New Handler",
				handler: func(result string, _ string) string {
					return result
				},
			},
			expected: true,
		},
		{
			name: "Register existing handler",
			args: args{
				handlerName: "empty",
				description: "Empty Clause",
				handler: func(result string, _ string) string {
					return result
				},
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterHandler(tt.args.handlerName, tt.args.description, tt.args.handler)
			if got := IsHandlerRegistered(tt.args.handlerName); got != tt.expected {
				t.Errorf("RegisterHandler() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
