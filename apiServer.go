package main

import (
	originalContext "context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/context"
)

// OnErrorHandler ...
type OnErrorHandler func(err error)

type apiServer struct {
	router         *mux.Router
	httpServer     *http.Server
	error          chan error
	onErrorHandler OnErrorHandler
}

func (apiServer *apiServer) handlePluginsManager() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := fmt.Sprintf("<html>%v</html>", "hello world")
		io.WriteString(w, out)
	}
}

// NewAPIServer ...
func NewAPIServer(bindAddr string, onErrorHandler OnErrorHandler) context.ContextedInstance {

	router := mux.NewRouter()
	httpServer := &http.Server{Addr: bindAddr, Handler: router}

	apiServer := &apiServer{
		router:         router,
		httpServer:     httpServer,
		error:          make(chan error),
		onErrorHandler: onErrorHandler,
	}

	router.HandleFunc("/", apiServer.handlePluginsManager())

	return apiServer
}

// Go ...
func (apiServer *apiServer) Go(current context.Context) {

	go func() {
		if err := apiServer.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				apiServer.error <- err
			}
		}
	}()

loop:
	for {
		select {
		case err := <-apiServer.error:
			apiServer.onErrorHandler(err)
			break
		case _, opened := <-current.Opened():
			if !opened {
				break loop
			}
		}
	}
	apiServer.httpServer.Shutdown(originalContext.Background())
}
