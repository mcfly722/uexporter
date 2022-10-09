package main

import (
	"github.com/mcfly722/goPackages/context"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/jsEngine"
)

// UExporter ...
type UExporter struct {
	httpServer *httpServer

	context   context.Context
	eventLoop jsEngine.EventLoop
	runtime   *goja.Runtime

	name string
}

// Constructor ...
func (uexporter UExporter) Constructor(context context.Context, eventLoop jsEngine.EventLoop, runtime *goja.Runtime) {

	newUExporter := &UExporter{
		context:    context,
		eventLoop:  eventLoop,
		runtime:    runtime,
		httpServer: uexporter.httpServer,
		name:       uexporter.name,
	}

	runtime.Set("UExporter", newUExporter)
}

// Publish ...
func (uexporter *UExporter) Publish(body string) {
	uexporter.httpServer.publish(uexporter.name, body)
}
