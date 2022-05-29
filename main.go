package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/gorilla/mux"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
	"github.com/mcfly722/goPackages/plugins"
)

// Server ...
type webServer struct {
	logger         *logger.Logger
	addr           string
	router         *mux.Router
	pluginsManager *plugins.Manager
}

var (
	bindAddrFlag                 *string
	pluginsPathFlag              *string
	sleepBetweenPluginUpdatesSec *int
	exitCode                     int
)

func main() {
	bindAddrFlag = flag.String("bindAddr", "127.0.0.1:8080", "bind address")
	pluginsPathFlag = flag.String("pluginsPath", "plugins", "path to plugins")
	sleepBetweenPluginUpdatesSec = flag.Int("sleepBetweenPluginUpdatesSec", 3, "pause in seconds between plugins updates")

	flag.Parse()

	var log = logger.NewLogger(100)
	log.SetOutputToConsole(true)

	var apiServer = newAPIServer(*bindAddrFlag, log)

	ctx := context.NewContextFor(apiServer)

	// handle ctrl+c for gracefully shutdown using context
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.LogEvent(logger.EventTypeInfo, "webServer", "CTRL+C signal")
		ctx.OnDone() <- true
	}()

	/*
		pluginsConstructor := func(plugin *plugins.Plugin) plugins.IPlugin {
			return &Plugin{
				Plugin: plugin,
				router: router,
			}
		}

		pluginsManager, err := plugins.NewPluginsManager(ctx, *pluginsPathFlag, "*.js", *sleepBetweenPluginUpdatesSec, pluginsConstructor)
		if err != nil {
			log.Fatal(err)
		}
	*/

	ctx.Wait()

	log.LogEvent(logger.EventTypeInfo, "webServer", "finished")

	os.Exit(exitCode)
}
