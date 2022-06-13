package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

// API ...
type API interface {
	Init(current context.Context, jsRuntime *JSRuntime)
}

type callback struct {
	function *goja.Callable
	args     []goja.Value
}

// JSRuntime ...
type JSRuntime struct {
	name      string
	logger    *logger.Logger
	VM        *goja.Runtime
	callbacks chan *callback
	context   context.Context
}

// NewJSRuntime ...
func NewJSRuntime(name string) *JSRuntime {
	runtime := &JSRuntime{
		name:      name,
		logger:    logger.NewLogger(5),
		VM:        goja.New(),
		callbacks: make(chan *callback),
	}

	runtime.VM.SetFieldNameMapper(goja.UncapFieldNameMapper())

	return runtime
}

// SetLogger ...
func (jsRuntime *JSRuntime) SetLogger(logger *logger.Logger) {
	jsRuntime.logger = logger
}

// Go ...
func (jsRuntime *JSRuntime) Go(current context.Context) {
	jsRuntime.context = current
	jsRuntime.logger.LogEvent(logger.EventTypeInfo, "jsRuntime", fmt.Sprintf("%v started", jsRuntime.name))

loop:
	for {
		select {
		case callback := <-jsRuntime.callbacks:
			_, err := (*callback.function)(nil, callback.args...)
			if err != nil {
				jsRuntime.CallHandlerException(err)
			}
		case <-current.OnDone():
			break loop
		}
	}
}

// AddAPI ...
func (jsRuntime *JSRuntime) AddAPI(api API) {
	api.Init(jsRuntime.context, jsRuntime)
}

// Dispose ...
func (jsRuntime *JSRuntime) Dispose() {
	jsRuntime.logger.LogEvent(logger.EventTypeInfo, "jsRuntime", fmt.Sprintf("%v disposed", jsRuntime.name))
}

// CallHandlerException ...
func (jsRuntime *JSRuntime) CallHandlerException(err error) {
	jsRuntime.logger.LogEvent(logger.EventTypeInfo, "jsRuntime", fmt.Sprintf("issued handler exception: %v", err))
}
