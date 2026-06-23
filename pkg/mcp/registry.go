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

// Package mcp provides MCP (Model Context Protocol) tool registry.
package mcp

import (
	"fmt"
	"log/slog"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/mcp"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// Registry manages MCP tool server connections.
type Registry struct {
	toolSets map[string]tool.ToolSet
}

// NewRegistry creates an empty MCP registry.
func NewRegistry() *Registry {
	return &Registry{toolSets: make(map[string]tool.ToolSet)}
}

// Discover connects to configured MCP servers.
func (r *Registry) Discover(servers []config.MCPServer) error {
	for _, srv := range servers {
		ts, err := createToolSet(srv)
		if err != nil {
			slog.Warn("MCP connect failed", "server", srv.Name, "error", err)
			continue
		}
		r.toolSets[srv.Name] = ts
		slog.Info("MCP connected", "server", srv.Name, "transport", srv.Transport)
	}
	return nil
}

// ToolSets returns all discovered MCP tool sets.
func (r *Registry) ToolSets() []tool.ToolSet {
	sets := make([]tool.ToolSet, 0, len(r.toolSets))
	for _, ts := range r.toolSets {
		sets = append(sets, ts)
	}
	return sets
}

// Close disconnects all MCP sessions.
func (r *Registry) Close() {
	for name, ts := range r.toolSets {
		if err := ts.Close(); err != nil {
			slog.Warn("MCP close error", "server", name, "error", err)
		}
	}
}

func createToolSet(srv config.MCPServer) (tool.ToolSet, error) {
	switch srv.Transport {
	case "stdio":
		if srv.Command == "" {
			return nil, fmt.Errorf("stdio transport requires command")
		}
		return mcp.NewMCPToolSet(mcp.ConnectionConfig{
			Transport: srv.Transport,
			Command:   srv.Command,
			Args:      srv.Args,
		}, mcp.WithName(srv.Name)), nil
	case "sse", "streamable":
		if srv.ServerURL == "" {
			return nil, fmt.Errorf("%s transport requires server_url", srv.Transport)
		}
		return mcp.NewMCPToolSet(mcp.ConnectionConfig{
			Transport: srv.Transport,
			ServerURL: srv.ServerURL,
		}, mcp.WithName(srv.Name)), nil
	}
	return nil, fmt.Errorf("unsupported transport: %s", srv.Transport)
}
