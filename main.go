package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/plugins"
)

// APIServer ...
type APIServer struct {
	router         *mux.Router
	pluginsManager *plugins.Manager
	scheduler      *Scheduler
}

// NewAPIServer ...
func NewAPIServer(pluginsManager *plugins.Manager) *APIServer {
	return &APIServer{
		router:         mux.NewRouter(),
		pluginsManager: pluginsManager,
	}
}

// Start ...
func (s *APIServer) Start(bindAddr string) error {
	s.router.HandleFunc("/", s.handlePluginsManager())

	log.Println(fmt.Sprintf("starting server on %v", bindAddr))
	return http.ListenAndServe(bindAddr, s.router)
}

func (s *APIServer) handlePluginsManager() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		out := fmt.Sprintf("<html>%v</html>", s.pluginsManager.ToHTML())
		io.WriteString(w, out)
	}
}

var (
	bindAddrFlag                 *string
	pluginsPathFlag              *string
	sleepBetweenPluginUpdatesSec *int
)

func main() {

	bindAddrFlag = flag.String("bindAddr", "127.0.0.1:8080", "bind address")
	pluginsPathFlag = flag.String("pluginsPath", "plugins", "path to plugins")
	sleepBetweenPluginUpdatesSec = flag.Int("sleepBetweenPluginUpdatesSec", 3, "pause in seconds between plugins updates")

	pluginsManager, err := plugins.NewPluginsManager(*pluginsPathFlag, *sleepBetweenPluginUpdatesSec, pluginsConstructor)
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(pluginsManager)
	if err := server.Start(*bindAddrFlag); err != nil {
		log.Fatal(err)
	}

}
