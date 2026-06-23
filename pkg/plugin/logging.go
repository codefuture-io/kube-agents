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

package plugin

import (
	"context"
	"log/slog"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/plugin"
)

// LoggingPlugin logs agent lifecycle events.
type LoggingPlugin struct{}

// NewLoggingPlugin creates a logging plugin.
func NewLoggingPlugin() *LoggingPlugin { return &LoggingPlugin{} }

// Name returns the plugin identifier.
func (p *LoggingPlugin) Name() string { return "logging" }

// Register hooks into the plugin registry.
func (p *LoggingPlugin) Register(r *plugin.Registry) {
	r.BeforeAgent(p.beforeAgent)
	r.AfterAgent(p.afterAgent)
}

func (p *LoggingPlugin) beforeAgent(
	ctx context.Context,
	args *agent.BeforeAgentArgs,
) (*agent.BeforeAgentResult, error) {
	slog.Info("agent start", "session", args.Invocation.Session.ID, "agent", args.Invocation.AgentName)
	return &agent.BeforeAgentResult{Context: ctx}, nil
}

func (p *LoggingPlugin) afterAgent(
	ctx context.Context,
	args *agent.AfterAgentArgs,
) (*agent.AfterAgentResult, error) {
	slog.Info("agent done", "session", args.Invocation.Session.ID, "agent", args.Invocation.AgentName)
	return &agent.AfterAgentResult{Context: ctx}, nil
}

var _ plugin.Plugin = (*LoggingPlugin)(nil)
