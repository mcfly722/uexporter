package main

import (
	originalContext "context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/context"
)

// OnErrorHandler ...
type OnErrorHandler func(err error)

type httpServer struct {
	router             *mux.Router
	httpServer         *http.Server
	error              chan error
	onErrorHandler     OnErrorHandler
	userName           string
	passwordSHA256Hash string
}

func (httpServer *httpServer) isAuthenticated(username string, password string) bool {
	if username != httpServer.userName {
		return false
	}
	h := sha256.New()
	h.Write([]byte(password))

	hash := hex.EncodeToString(h.Sum(nil))

	return hash == httpServer.passwordSHA256Hash
}

// newHTTPServer ...
func newHTTPServer(bindAddr string, onErrorHandler OnErrorHandler, userName string, passwordSHA256Hash string) *httpServer {

	router := mux.NewRouter()

	server := &httpServer{
		router:             router,
		httpServer:         &http.Server{Addr: bindAddr, Handler: router},
		error:              make(chan error),
		onErrorHandler:     onErrorHandler,
		userName:           userName,
		passwordSHA256Hash: strings.ToLower(passwordSHA256Hash),
	}

	return server
}

// Go ...
func (httpServer *httpServer) Go(current context.Context) {

	go func() {
		if err := httpServer.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				httpServer.error <- err
			}
		}
	}()

loop:
	for {
		select {
		case err := <-httpServer.error:
			httpServer.onErrorHandler(err)
			break
		case _, opened := <-current.Opened():
			if !opened {
				break loop
			}
		}
	}
	httpServer.httpServer.Shutdown(originalContext.Background())
}
