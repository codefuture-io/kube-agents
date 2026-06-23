/*
Copyright 2026 CodeFuture Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package plugin provides a plugin registry for kube-agents.
package plugin

import (
	"trpc.group/trpc-go/trpc-agent-go/plugin"
)

// Registry manages a collection of plugins for the agent system.
type Registry struct {
	plugins []plugin.Plugin
	names   []string
}

// NewRegistry creates an empty plugin registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a plugin to the registry.
func (r *Registry) Register(p plugin.Plugin) {
	r.plugins = append(r.plugins, p)
	r.names = append(r.names, p.Name())
}

// Build returns the ordered list of plugins for runner configuration.
func (r *Registry) Build() []plugin.Plugin {
	return r.plugins
}

// Names returns the names of registered plugins in order.
func (r *Registry) Names() []string {
	return r.names
}

// CreatePlugins creates and registers the default plugins by name.
func CreatePlugins(names []string) *Registry {
	reg := NewRegistry()
	for _, name := range names {
		switch name {
		case "audit":
			reg.Register(NewAuditPlugin())
		case "ratelimit":
			reg.Register(NewRateLimitPlugin())
		case "logging":
			reg.Register(NewLoggingPlugin())
		}
	}
	return reg
}
