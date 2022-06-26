package main

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/plugins"
	yaml "gopkg.in/yaml.v2"
)

type plugin struct {
	definition plugins.PluginDefinition
}

// YAMLConfig ...
type YAMLConfig struct {
	PluginName                  string            `yaml:"PluginName"`
	Version                     string            `yaml:"Version"`
	JSScripts                   []string          `yaml:"JSScripts"`
	DefaultEnvironmentVariables map[string]string `yaml:"DefaultEnvironmentVariables"`
}

func newPlugin(definition plugins.PluginDefinition) context.ContextedInstance {
	return &plugin{
		definition: definition,
	}
}

func (plugin *plugin) Go(current context.Context) {

	config := &YAMLConfig{}
	err := yaml.Unmarshal([]byte(plugin.definition.GetBody()), &config)
	if err != nil {
		current.Log(1, err.Error())
		return
	}

	scripts := []*script{}

	for _, resource := range config.JSScripts {
		body, err := plugin.definition.GetResource(resource)
		if err != nil {
			current.Log(1, err.Error())
			return
		}
		scripts = append(scripts, newScript(resource, string(*body)))
		current.Log(101, fmt.Sprintf("%v loaded", resource))
	}

	eventLoop := newEventLoop(goja.New(), scripts)
	eventLoop.addAPI(apiConsole)
	eventLoop.addAPI(apiScheduler)

	_, err = current.NewContextFor(eventLoop, config.PluginName, "eventLoop")
	if err != nil {
		current.Log(105, "NewContextFor", "skipping")
	} else {
		current.Log(101, "loop started")
	loop:
		for {
			select {
			case <-time.After(1 * time.Second):
				if plugin.definition.Outdated() {
					current.Cancel()
					break
				}
				break
			case _, opened := <-current.Opened():
				if !opened {
					break loop
				}
			}
		}
	}

}
