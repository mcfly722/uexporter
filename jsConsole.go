package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

type jsConsole struct {
	name   string
	logger *logger.Logger
}

// NewJSConsole ...
func NewJSConsole(name string, logger *logger.Logger) API {
	return &jsConsole{
		name:   name,
		logger: logger,
	}
}

// Init ...
func (jsConsole *jsConsole) Init(current context.Context, jsRuntime *JSRuntime) {
	console := &Console{name: jsConsole.name, logger: jsConsole.logger}
	jsRuntime.VM.Set("Console", console)
}

// Console ...
type Console struct {
	name   string
	logger *logger.Logger
}

// Log ...
func (console *Console) Log(objects ...goja.Value) {
	for _, obj := range objects {
		console.logger.LogEvent(logger.EventTypeInfo, console.name, fmt.Sprintf("%v", obj))
	}
}
