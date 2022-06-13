package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
	yaml "gopkg.in/yaml.v2"
)

// Plugin ...
type Plugin struct {
	log         *logger.Logger
	jsScripts   map[string]time.Time
	jsRuntime   *JSRuntime
	rootContext context.RootContext
}

// NewPlugin ...
func NewPlugin() *Plugin {
	log := logger.NewLogger(100)
	log.SetOutputToConsole(true)

	return &Plugin{
		log:       log,
		jsScripts: make(map[string]time.Time),
	}
}

// YAMLConfig ...
type YAMLConfig struct {
	PluginName                  string            `yaml:"PluginName"`
	Version                     string            `yaml:"Version"`
	JSScripts                   []string          `yaml:"JSScripts"`
	DefaultEnvironmentVariables map[string]string `yaml:"DefaultEnvironmentVariables"`
}

func (plugin *Plugin) load(pluginsFullPath string, relativeName string, body string) error {

	{ // start jsRuntime anyway
		plugin.rootContext = context.NewRootContext()
		plugin.jsRuntime = NewJSRuntime(relativeName)
		plugin.jsRuntime.SetLogger(plugin.log)
		plugin.rootContext.NewContextFor(plugin.jsRuntime)
	}

	plugin.jsRuntime.AddAPI(NewJSConsole(relativeName, plugin.log))
	plugin.jsRuntime.AddAPI(NewJSScheduler(relativeName, plugin.log))

	pluginRootPath := filepath.Dir(filepath.Join(pluginsFullPath, relativeName))

	config := &YAMLConfig{}

	err := yaml.Unmarshal([]byte(body), &config)
	if err != nil {
		plugin.log.LogEvent(logger.EventTypeException, relativeName, err.Error())
		return err
	}

	plugin.jsScripts = make(map[string]time.Time)

	for _, jsScript := range config.JSScripts {
		fullFilePath := filepath.Join(pluginRootPath, jsScript)
		file, err := os.Stat(fullFilePath)
		if err != nil {
			plugin.log.LogEvent(logger.EventTypeException, relativeName, err.Error())
			return err
		}
		plugin.jsScripts[jsScript] = file.ModTime()

		data, err := os.ReadFile(fullFilePath)
		if err != nil {
			plugin.log.LogEvent(logger.EventTypeException, relativeName, err.Error())
			return err
		}

		_, err = plugin.jsRuntime.VM.RunString(string(data))
		if err != nil {
			plugin.log.LogEvent(logger.EventTypeException, relativeName, err.Error())
			return err
		}

	}

	return nil
}

// OnLoad ...
func (plugin *Plugin) OnLoad(pluginsFullPath string, relativeName string, body string) {
	plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "loading")
	if plugin.load(pluginsFullPath, relativeName, body) == nil {
		plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "loaded")
	}
}

// OnUpdate ...
func (plugin *Plugin) OnUpdate(pluginsFullPath string, relativeName string, body string) {
	plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "updating")
	plugin.rootContext.Terminate()
	if plugin.load(pluginsFullPath, relativeName, body) == nil {
		plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "updated")
	}
}

// OnDispose ...
func (plugin *Plugin) OnDispose(pluginsFullPath string, relativeName string) {
	plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "disposing")
	plugin.rootContext.Terminate()
	plugin.log.LogEvent(logger.EventTypeInfo, relativeName, "disposed")
}

// UpdateRequired ...
func (plugin *Plugin) UpdateRequired(pluginsFullPath string, relativeName string) bool {

	pluginRootPath := filepath.Dir(filepath.Join(pluginsFullPath, relativeName))

	{ // check all jsScripts for changes, if one of scripts has been changed, we need plugin update
		for jsScript, lastModTime := range plugin.jsScripts {
			file, err := os.Stat(filepath.Join(pluginRootPath, jsScript))
			if err != nil {
				return true
			}
			if file.ModTime() != lastModTime {
				return true
			}
		}
	}
	return false
}
