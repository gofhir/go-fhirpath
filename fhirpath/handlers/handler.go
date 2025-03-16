package handlers

import (
	"fmt"
	"sync"
)

type HandlerFunc func(result string, key string) string

type Handler struct {
	Pattern string
	Log     string
	Func    HandlerFunc
}

var (
	registry []Handler
	once     sync.Once
	mu       sync.Mutex
)

func RegisterHandler(pattern string, log string, fn HandlerFunc) {
	mu.Lock()
	defer mu.Unlock()

	// Prevent duplicate registration
	for _, h := range registry {
		if h.Pattern == pattern {
			fmt.Printf("⚠️ Handler for pattern '%s' already registered. Skipping...\n", pattern)
			return
		}
	}

	fmt.Printf("RegisterHandler pattern=%s log=%s\n", pattern, log)

	registry = append(registry, Handler{
		Pattern: pattern,
		Log:     log,
		Func:    fn,
	})
}

func GetHandlers() []Handler {
	mu.Lock()
	defer mu.Unlock()
	return registry
}

func IsHandlerRegistered(pattern string) bool {
	mu.Lock()
	defer mu.Unlock()

	var isRegistered bool

	for _, h := range registry {
		if h.Pattern == pattern {
			isRegistered = true
			break
		}
	}

	return isRegistered
}

// Clear resets the handler registry (use only in tests)
func Clear() {
	mu.Lock()
	defer mu.Unlock()
	
	registry = []Handler{}
}
