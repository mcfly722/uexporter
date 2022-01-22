package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// APIServer ...
type APIServer struct {
	router    *mux.Router
	scheduler *Scheduler
}

// NewAPIServer ...
func NewAPIServer(pluginsPath string) *APIServer {
	return &APIServer{
		router:    mux.NewRouter(),
		scheduler: NewScheduler(pluginsPath),
	}
}

// Start ...
func (s *APIServer) Start(bindAddr string) error {
	s.router.HandleFunc("/", s.handleTasks())

	fmt.Println(fmt.Sprintf("starting server on %v", bindAddr))
	return http.ListenAndServe(bindAddr, s.router)
}

func (s *APIServer) handleTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		out := fmt.Sprintf("<html>%v</html>", s.scheduler.ToHTML())
		io.WriteString(w, out)
	}
}

var (
	bindAddrFlag    *string
	pluginsPathFlag *string
)

func main() {

	bindAddrFlag = flag.String("bindAddr", "127.0.0.1:8080", "bind address")
	pluginsPathFlag = flag.String("pluginsPath", "plugins", "path to plugins")

	server := NewAPIServer(*pluginsPathFlag)
	if err := server.Start(*bindAddrFlag); err != nil {
		log.Fatal(err)
	}

}
