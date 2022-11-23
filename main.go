package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/jsEngine"
)

var (
	bindAddrFlag           *string
	pluginsFlag            *string
	passwordSHA256hashFlag *string
	userNameFlag           *string
	skipAuth               *bool

	exitCode      int
	exitException string
)

// LogExitError ...
func logExitError(err error) {
	exitCode = 1
	exitException = err.Error()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	var passwordHash string

	bindAddrFlag = flag.String("bindAddr", ":9100", "bind address")
	pluginsFlag = flag.String("plugins", "plugins/topMemory.js,plugins/topCPU.js,plugins/uptime.js", "JavaScript plugins. Use ',' do delimit files")
	passwordSHA256hashFlag = flag.String("passwordSHA256hash", "", "password sha256 hash")
	userNameFlag = flag.String("username", "uexporter", "user name")
	skipAuth = flag.Bool("skipAuth", false, "skip authentication")

	flag.Parse()

	{ // set passwordHash
		if os.Getenv("UEXPORTER_PASSWORDSHA256HASH") != "" {
			passwordHash = os.Getenv("UEXPORTER_PASSWORDSHA256HASH")
		} else {
			if *passwordSHA256hashFlag != "" {
				passwordHash = *passwordSHA256hashFlag
			} else {
				if (*skipAuth) {
					passwordHash = ""
					} else {
						log.Fatal("You have to specify -passwordSHA256hash value, or UEXPORTER_PASSWORDSHA256HASH environment variable, or use -skipAuth to skip authentication")
					}
			}
		}
	}

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	scripts := map[string]jsEngine.Script{}
	{ // read all scripts

		pluginsFileNames := strings.Split(*pluginsFlag, ",")

		for _, pluginName := range pluginsFileNames {

			content, err := ioutil.ReadFile(pluginName)
			if err != nil {
				log.Fatal(err)
			}

			scripts[pluginName] = jsEngine.NewScript(pluginName, string(content))
		}
	}

	ctrlC := make(chan os.Signal, 1)

	rootContext := context.NewRootContext(context.NewConsoleLogDebugger(100, true))

	var httpServer = newHTTPServer(*bindAddrFlag, func(err error) {
		rootContext.Log(2, err.Error())
		logExitError(err)
		rootContext.Cancel()
	}, *userNameFlag, passwordHash, hostName)

	_, err = rootContext.NewContextFor(httpServer, *bindAddrFlag, "apiServer")
	if err == nil {

		for name, script := range scripts {
			eventLoop := jsEngine.NewEventLoop(goja.New(), []jsEngine.Script{script})

			eventLoop.Import(jsEngine.Console{})
			eventLoop.Import(jsEngine.Scheduler{})
			eventLoop.Import(jsEngine.Exec{})
			eventLoop.Import(OperatingSystem{})
			eventLoop.Import(IOUtil{})
			eventLoop.Import(UExporter{httpServer: httpServer, name: name})

			_, err := rootContext.NewContextFor(eventLoop, name, "eventLoop")
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
