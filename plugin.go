package main

import (
	"time"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/plugins"
)

type plugin struct {
	definition plugins.PluginDefinition
}

func newPlugin(definition plugins.PluginDefinition) context.ContextedInstance {
	return &plugin{
		definition: definition,
	}
}

func (plugin *plugin) Go(current context.Context) {
loop:
	for {
		select {
		case <-time.After(1 * time.Second):
			if plugin.definition.Outdated() {
				break loop
			}
			break
		case <-current.OnDone():
			break loop
		}

	}
}

func (plugin *plugin) Dispose(current context.Context) {}
