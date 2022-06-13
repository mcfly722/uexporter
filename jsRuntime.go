package main

import (
	"fmt"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

// JSRuntime ...
type JSRuntime struct {
	name   string
	logger *logger.Logger
}

// NewJSRuntime ...
func NewJSRuntime(name string) *JSRuntime {
	engine := &JSRuntime{
		name:   name,
		logger: logger.NewLogger(5),
	}

	return engine
}

// SetLogger ...
func (jsRuntime *JSRuntime) SetLogger(logger *logger.Logger) {
	jsRuntime.logger = logger
}

// Go ...
func (jsRuntime *JSRuntime) Go(current context.Context) {

	jsRuntime.logger.LogEvent(logger.EventTypeInfo, "jsRuntime", fmt.Sprintf("%v started", jsRuntime.name))

loop:
	for {
		select {
		case <-current.OnDone():
			break loop
		}
	}
}

// Dispose ...
func (jsRuntime *JSRuntime) Dispose() {
	jsRuntime.logger.LogEvent(logger.EventTypeInfo, "jsRuntime", fmt.Sprintf("%v disposed", jsRuntime.name))
}
