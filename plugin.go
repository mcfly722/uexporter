package main

import (
	"fmt"
	"log"

	"github.com/mcfly722/goPackages/plugins"
)

// Plugin ...
type Plugin struct {
	*plugins.Plugin
}

// Constructor ...
func pluginsConstructor(plugin *plugins.Plugin) plugins.IPlugin {
	return &Plugin{Plugin: plugin}
}

// OnLoad ...
func (plugin *Plugin) OnLoad() {
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
