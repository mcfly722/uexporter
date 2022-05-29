package main

import (
	originalContext "context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

type apiServer struct {
	logger *logger.Logger
	router *mux.Router
	done   chan bool
	http   *http.Server
}

func (apiServer *apiServer) handlePluginsManager() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := fmt.Sprintf("<html>%v</html>", "hello world")
		io.WriteString(w, out)
	}
}

func newAPIServer(bindAddr string, log *logger.Logger) *apiServer {

	router := mux.NewRouter()

	httpServer := &http.Server{Addr: bindAddr, Handler: router}

	done := make(chan bool)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.LogEvent(logger.EventTypeInfo, "webServer", err.Error())
			} else {
				log.LogEvent(logger.EventTypeException, "webServer", err.Error())
			}
		}
		done <- true
	}()

	log.LogEvent(logger.EventTypeInfo, "webServer", fmt.Sprintf("starting on %v", bindAddr))

	var apiServer = &apiServer{
		logger: log,
		router: router,
		done:   done,
		http:   httpServer,
	}

	router.HandleFunc("/", apiServer.handlePluginsManager())

	return apiServer
}

// Go ...
func (apiServer *apiServer) Go(current context.Context) {
loop:
	for {
		select {
		case <-apiServer.done:
			break loop
		case <-current.OnDone():
			break loop
		}
	}
}

// Dispose ...
func (apiServer *apiServer) Dispose() {
	apiServer.http.Shutdown(originalContext.Background())
}
