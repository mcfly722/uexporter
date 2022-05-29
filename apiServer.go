package main

import (
	originalContext "context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

type apiServer struct {
	logger  *logger.Logger
	router  *mux.Router
	http    *http.Server
	working sync.WaitGroup
	done    chan bool
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

	var apiServer = &apiServer{
		logger: log,
		router: router,
		http:   httpServer,
		done:   done,
	}

	go func() {
		log.LogEvent(logger.EventTypeInfo, "webServer", fmt.Sprintf("starting on %v", bindAddr))

		apiServer.working.Add(1)

		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.LogEvent(logger.EventTypeInfo, "webServer", err.Error())
			} else {
				log.LogEvent(logger.EventTypeException, "webServer", err.Error())
				exitCode = 1
			}
		}

		go func() {
			done <- true
		}()

		apiServer.working.Done()

	}()

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

	go func() {
		apiServer.http.Shutdown(originalContext.Background())
	}()

	// wait till would be totally disposed
	apiServer.working.Wait()
}
