package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/plugins"
)

var (
	bindAddrFlag                 *string
	pluginsPathFlag              *string
	sleepBetweenPluginUpdatesSec *int
	exitCode                     int
	exitException                string
)

// LogExitError ...
func logExitError(err error) {
	exitCode = 1
	exitException = err.Error()
}

func main() {
	bindAddrFlag = flag.String("bindAddr", "127.0.0.1:8080", "bind address")
	pluginsPathFlag = flag.String("pluginsPath", "plugins", "path to plugins")
	sleepBetweenPluginUpdatesSec = flag.Int("sleepBetweenPluginUpdatesSec", 3, "pause in seconds between plugins updates")

	flag.Parse()

	ctrlC := make(chan os.Signal, 1)

	rootContext := context.NewRootContext(context.NewConsoleLogDebugger())

	var apiServer = NewAPIServer(*bindAddrFlag, func(err error) {
		rootContext.Log(2, err.Error())
		logExitError(err)
		rootContext.Terminate()
	})

	pluginsProvider := plugins.NewPluginsFromFilesProvider(*pluginsPathFlag, "*.yaml")
	pluginsManager := plugins.NewPluginsManager(pluginsProvider, 3, newPlugin)

	apiServerContext := rootContext.NewContextFor(apiServer, *bindAddrFlag, "apiServer")
	apiServerContext.NewContextFor(pluginsManager, "pluginsManager", "pluginsManager")

	{ // handle ctrl+c for gracefully shutdown using context
		signal.Notify(ctrlC, os.Interrupt)
		go func() {
			<-ctrlC
			rootContext.Log(2, "CTRL+C signal")
			rootContext.Terminate()
		}()
	}

	rootContext.Wait()

	rootContext.Log(0, fmt.Sprintf("exitCode=%v %v", exitCode, exitException))
	os.Exit(exitCode)
}
