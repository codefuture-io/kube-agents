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
	"time"

	"trpc.group/trpc-go/trpc-agent-go/plugin"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// AuditPlugin logs all tool invocations for audit trails.
type AuditPlugin struct{}

// NewAuditPlugin creates an audit plugin.
func NewAuditPlugin() *AuditPlugin { return &AuditPlugin{} }

// Name returns the plugin identifier.
func (p *AuditPlugin) Name() string { return "audit" }

// Register hooks into the plugin registry.
func (p *AuditPlugin) Register(r *plugin.Registry) {
	r.AfterTool(p.afterTool)
}

func (p *AuditPlugin) afterTool(
	ctx context.Context,
	args *tool.AfterToolArgs,
) (*tool.AfterToolResult, error) {
	slog.Info("tool invocation",
		"tool", args.ToolName,
		"callID", args.ToolCallID,
		"time", time.Now().Format(time.RFC3339),
		"hasResult", args.Result != nil,
		"error", args.Error)
	return &tool.AfterToolResult{Context: ctx}, nil
}

var _ plugin.Plugin = (*AuditPlugin)(nil)
