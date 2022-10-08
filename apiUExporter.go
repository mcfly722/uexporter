package main

import (
	"net/http"
	"sync"

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

	content *content
}

type content struct {
	body  string
	ready sync.Mutex
}

func newContent() *content {
	return &content{
		body: "empty",
	}
}

// Constructor ...
func (uexporter UExporter) Constructor(context context.Context, eventLoop jsEngine.EventLoop, runtime *goja.Runtime) {

	newUExporter := &UExporter{
		context:    context,
		eventLoop:  eventLoop,
		runtime:    runtime,
		content:    newContent(),
		httpServer: uexporter.httpServer,
	}

	newUExporter.httpServer.router.HandleFunc("/", uexporter.getHTTPHandler(newUExporter.content))

	runtime.Set("UExporter", newUExporter)
}

// Publish ...
func (uexporter *UExporter) Publish(body string) {
	uexporter.content.ready.Lock()
	uexporter.content.body = body
	uexporter.content.ready.Unlock()
}

func (uexporter *UExporter) getHTTPHandler(content *content) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()

		if !ok {
			w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "No basic auth present"}`))
			return
		}

		if !uexporter.httpServer.isAuthenticated(username, password) {
			w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "Invalid username or password"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")

		content.ready.Lock()
		w.Write([]byte(content.body))
		content.ready.Unlock()

		return
	}

}
