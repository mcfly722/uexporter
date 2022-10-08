package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/jsEngine"
)

var (
	bindAddrFlag           *string
	pluginFlag             *string
	passwordSHA256hashFlag *string
	userNameFlag           *string

	exitCode      int
	exitException string
)

// LogExitError ...
func logExitError(err error) {
	exitCode = 1
	exitException = err.Error()
}

func main() {
	var passwordHash string
	var pluginContent string

	bindAddrFlag = flag.String("bindAddr", ":8080", "bind address")
	pluginFlag = flag.String("plugin", "plugin.js", "JavaScript plugin")
	passwordSHA256hashFlag = flag.String("passwordSHA256hash", "", "password sha256 hash")
	userNameFlag = flag.String("username", "uexporter", "user name")

	flag.Parse()

	{ // set passwordHash
		if os.Getenv("UEXPORTER_PASSWORDSHA256HASH") != "" {
			passwordHash = os.Getenv("UEXPORTER_PASSWORDSHA256HASH")
		} else {
			if *passwordSHA256hashFlag != "" {
				passwordHash = *passwordSHA256hashFlag
			} else {
				log.Fatal("You have to specify -passwordSHA256hash value, or UEXPORTER_PASSWORDSHA256HASH environment variable")
			}
		}
	}

	{ // read plugin content
		content, err := ioutil.ReadFile(*pluginFlag)
		if err != nil {
			log.Fatal(err)
		}
		pluginContent = string(content)

		//fmt.Println(pluginContent)
	}

	ctrlC := make(chan os.Signal, 1)

	rootContext := context.NewRootContext(context.NewConsoleLogDebugger(100, true))

	var httpServer = newHTTPServer(*bindAddrFlag, func(err error) {
		rootContext.Log(2, err.Error())
		logExitError(err)
		rootContext.Cancel()
	}, *userNameFlag, passwordHash)

	_, err := rootContext.NewContextFor(httpServer, *bindAddrFlag, "apiServer")
	if err == nil {

		{ // starting JavaScript plugin EventLoop
			scripts := []jsEngine.Script{jsEngine.NewScript(*pluginFlag, pluginContent)}
			eventLoop := jsEngine.NewEventLoop(goja.New(), scripts)

			eventLoop.Import(jsEngine.Console{})
			eventLoop.Import(jsEngine.Scheduler{})
			eventLoop.Import(jsEngine.Exec{})
			eventLoop.Import(UExporter{httpServer: httpServer})

			_, err := rootContext.NewContextFor(eventLoop, *pluginFlag, "eventLoop")
			if err != nil {
				log.Fatal(err)
			}
		}

		{ // handle ctrl+c for gracefully shutdown using context
			signal.Notify(ctrlC, os.Interrupt)
			go func() {
				<-ctrlC
				rootContext.Log(2, "CTRL+C signal")
				rootContext.Cancel()
			}()
		}

		rootContext.Wait()
	}

	rootContext.Log(0, fmt.Sprintf("exitCode=%v %v", exitCode, exitException))

	os.Exit(exitCode)
}
