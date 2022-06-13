package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
	"github.com/mcfly722/goPackages/plugins"
)

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

	rootContext := context.NewRootContext()
	defer rootContext.Terminate()

	var apiServer = newAPIServer(*bindAddrFlag, log)
	apiServerContext := rootContext.NewContextFor(apiServer)

	pluginsConstructor := func() plugins.IPlugin {
		return NewPlugin()
	}

	pluginsManager, err := plugins.NewPluginsManager(*pluginsPathFlag, "*.yaml", 3, pluginsConstructor)
	if err != nil {
		log.LogEvent(logger.EventTypeException, "pluginsManager", err.Error())
	} else {

		pluginsManager.SetLogger(log)

		apiServerContext.NewContextFor(pluginsManager)

		// handle ctrl+c for gracefully shutdown using context
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			log.LogEvent(logger.EventTypeInfo, "main", "CTRL+C signal")
			rootContext.Terminate()
		}()

	}

	rootContext.Wait()

	log.LogEvent(logger.EventTypeInfo, "main", fmt.Sprintf("finished exitCode=%v", exitCode))

	os.Exit(exitCode)
}
