package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dop251/goja"
	"github.com/gorilla/mux"
	"github.com/mcfly722/goPackages/logger"
	"github.com/mcfly722/goPackages/plugins"
)

// Plugin ...
type Plugin struct {
	*plugins.Plugin
	logger *logger.Logger
	router *mux.Router
}

// OnLoad ...
func (plugin *Plugin) OnLoad() {
	plugin.logger = logger.NewLogger(30)

	file, err := os.Open(plugin.RelativeName)
	if err != nil {
		plugin.logger.LogEvent(logger.EventTypeException, plugin.RelativeName, "could not open file for reading")
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			plugin.logger.LogEvent(logger.EventTypeException, plugin.RelativeName, "could close file after reading")
		}
	}()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		plugin.logger.LogEvent(logger.EventTypeException, plugin.RelativeName, "could not read file")
		return
	}

	vm := goja.New()
	_, err = vm.RunString(string(body))
	if err != nil {
		plugin.logger.LogEvent(logger.EventTypeException, plugin.RelativeName, err.Error())
		return
	}

	log.Println(fmt.Sprintf("loaded plugin: %v", plugin.RelativeName))
}

// OnUpdate ...
func (plugin *Plugin) OnUpdate() {
	log.Println(fmt.Sprintf("updated plugin: %v", plugin.RelativeName))
}

// OnUnload ...
func (plugin *Plugin) OnUnload() {
	log.Println(fmt.Sprintf("uloaded plugin: %v", plugin.RelativeName))
}
