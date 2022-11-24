package main

import (
	originalContext "context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

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
	path               string
	hostName           string

	content map[string]string
	ready   sync.Mutex
}

func (httpServer *httpServer) isAuthenticated(username string, password string) bool {

	if username != httpServer.userName {
		return false
	}

	hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))                      // original one hash
	hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v\n", password)))) // this hash with additional \n symbol, when you using "echo ... | sha256sum" this symbol added automatically

	return hash1 == httpServer.passwordSHA256Hash || hash2 == httpServer.passwordSHA256Hash
}

// newHTTPServer ...
func newHTTPServer(bindAddr string, path string, onErrorHandler OnErrorHandler, userName string, passwordSHA256Hash string, hostName string) *httpServer {

	router := mux.NewRouter()

	server := &httpServer{
		router:             router,
		httpServer:         &http.Server{Addr: bindAddr, Handler: router},
		error:              make(chan error),
		onErrorHandler:     onErrorHandler,
		userName:           userName,
		passwordSHA256Hash: strings.ToLower(passwordSHA256Hash),
		content:            make(map[string]string),
		path:               path,
		hostName:           hostName,
	}

	return server
}

// Go ...
func (httpServer *httpServer) Go(current context.Context) {

	httpServer.router.HandleFunc(httpServer.path, httpServer.getHTTPHandler())

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

func (httpServer *httpServer) getHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if (httpServer.passwordSHA256Hash != ""){ // -skipAuth
			username, password, ok := r.BasicAuth()

			if !ok {
				w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%v: Give username and password"`, httpServer.hostName))
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "No basic auth present"}`))
				return
			}

			if !httpServer.isAuthenticated(username, password) {
				w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%v: Give username and password"`, httpServer.hostName))
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf(`{"message": "Invalid username or password", server:"%v"}`, httpServer.hostName)))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")

		{ // write content from content map, sort it by key name
			httpServer.ready.Lock()
			keys := make([]string, 0, len(httpServer.content))
			for k := range httpServer.content {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			w.Write([]byte(fmt.Sprintf("# host: %v\n# time: %v\n\n", httpServer.hostName, time.Now())))

			for _, key := range keys {
				w.Write([]byte(fmt.Sprintf("# %v\n%v\n\n", key, httpServer.content[key])))
			}
			httpServer.ready.Unlock()
		}

		return
	}

}

func (httpServer *httpServer) publish(name, body string) {
	httpServer.ready.Lock()
	httpServer.content[name] = body
	httpServer.ready.Unlock()
}
