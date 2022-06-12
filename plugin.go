package main

import (
	"fmt"
	"log"

	"github.com/mcfly722/goPackages/logger"
	yaml "gopkg.in/yaml.v2"
)

// Plugin ...
type Plugin struct {
	Log *logger.Logger
}

// NewPlugin ...
func NewPlugin() *Plugin {
	log := logger.NewLogger(100)
	log.SetOutputToConsole(true)
	return &Plugin{
		Log: log,
	}
}

// YAMLConfig ...
type YAMLConfig struct {
	PluginName                  string            `yaml:"PluginName"`
	Version                     string            `yaml:"Version"`
	JSScripts                   []string          `yaml:"JSScripts"`
	DefaultEnvironmentVariables map[string]string `yaml:"DefaultEnvironmentVariables"`
}

// OnLoad ...
func (plugin *Plugin) OnLoad(relativeName string, body string) {
	config := &YAMLConfig{}

	err := yaml.Unmarshal([]byte(body), &config)
	if err != nil {
		plugin.Log.LogEvent(logger.EventTypeException, relativeName, err.Error())
		return
	}

	plugin.Log.LogEvent(logger.EventTypeInfo, relativeName, fmt.Sprintf("PluginName                 : %v", config.PluginName))
	plugin.Log.LogEvent(logger.EventTypeInfo, relativeName, fmt.Sprintf("Version                    : %v", config.Version))
	plugin.Log.LogEvent(logger.EventTypeInfo, relativeName, fmt.Sprintf("JSScripts                  : %v", config.JSScripts))
	plugin.Log.LogEvent(logger.EventTypeInfo, relativeName, fmt.Sprintf("DefaultEnvironmentVariables: %v", config.DefaultEnvironmentVariables))
	plugin.Log.LogEvent(logger.EventTypeInfo, relativeName, "loaded")
}

// OnUpdate ...
func (plugin *Plugin) OnUpdate(relativeName string, body string) {
	log.Println(fmt.Sprintf("%v updated", relativeName))
}

// OnDispose ...
func (plugin *Plugin) OnDispose(relativeName string) {
	log.Println(fmt.Sprintf("%v uloaded", relativeName))
}

// UpdateRequired ...
func (plugin *Plugin) UpdateRequired() bool {
	return false
}

/*
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

*/
