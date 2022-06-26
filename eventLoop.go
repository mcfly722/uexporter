package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
)

type apiConstructor func(context context.Context, eventLoop *eventLoop, runtime *goja.Runtime)

type callback struct {
	function *goja.Callable
	args     []goja.Value
}

type script struct {
	name string
	body string
}

func newScript(name string, body string) *script {
	return &script{
		name: name,
		body: body,
	}
}

type eventLoop struct {
	runtime *goja.Runtime
	apis    []apiConstructor
	scripts []*script
}

func (eventLoop *eventLoop) addAPI(api apiConstructor) {
	eventLoop.apis = append(eventLoop.apis, api)
}

func newEventLoop(runtime *goja.Runtime, scripts []*script) *eventLoop {
	eventLoop := &eventLoop{
		runtime: runtime,
		apis:    []apiConstructor{},
		scripts: scripts,
	}

	return eventLoop
}

// Go ...
func (eventLoop *eventLoop) Go(current context.Context) {

	for _, api := range eventLoop.apis {
		api(current, eventLoop, eventLoop.runtime)
	}

	for _, script := range eventLoop.scripts {
		_, err := eventLoop.runtime.RunString(script.body)
		if err != nil {
			current.Log(1, fmt.Sprintf("%v: %v", script.name, err.Error()))
			return
		}
	}

loop:
	for {
		select {
		case _, opened := <-current.Opened():
			if !opened {
				break loop
			}
		}
	}
}

// Dispose ...
func (eventLoop *eventLoop) Dispose(current context.Context) {}
