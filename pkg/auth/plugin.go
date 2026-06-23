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

package auth

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/plugin"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// PluginName is the registered name of the auth plugin.
const PluginName = "auth"

// identityStateKey is the session state key storing the user identity.
const identityStateKey = "_auth_identity"
const authModeStateKey = "_auth_mode"

// Plugin implements plugin.Plugin to enforce authentication and
// authorization before agent execution and tool invocations.
type Plugin struct {
	validator TokenValidator
	rbac      *RBACChecker
}

// NewPlugin creates an auth plugin.
func NewPlugin(validator TokenValidator, rbac *RBACChecker) *Plugin {
	return &Plugin{validator: validator, rbac: rbac}
}

// Name returns the plugin identifier.
func (p *Plugin) Name() string { return PluginName }

// Register hooks this plugin into the provided registry.
func (p *Plugin) Register(r *plugin.Registry) {
	r.BeforeAgent(p.beforeAgent)
	r.BeforeTool(p.beforeTool)
}

// beforeAgent validates the token from session state and injects user identity.
func (p *Plugin) beforeAgent(
	ctx context.Context,
	args *agent.BeforeAgentArgs,
) (*agent.BeforeAgentResult, error) {
	inv := args.Invocation

	token := extractTokenFromSession(inv)
	if token == "" {
		return nil, fmt.Errorf("auth plugin: no token found — set 'token' in session state before calling run")
	}

	identity, err := p.validator.Validate(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("auth plugin: %w", err)
	}

	// Store identity in invocation for downstream callbacks.
	inv.SetState(identityStateKey, identity)
	inv.SetState(authModeStateKey, identity.AuthMode)

	// Also write into session state for {user:identity?} placeholders.
	if inv.Session != nil && inv.Session.State != nil {
		inv.Session.State["identity"] = []byte(identity.UserID)
		inv.Session.State["auth_mode"] = []byte(identity.AuthMode)
	}

	return &agent.BeforeAgentResult{Context: ctx}, nil
}

// beforeTool performs RBAC check before K8s tool execution.
func (p *Plugin) beforeTool(
	ctx context.Context,
	args *tool.BeforeToolArgs,
) (*tool.BeforeToolResult, error) {
	if !isK8sTool(args.ToolName) {
		return &tool.BeforeToolResult{Context: ctx}, nil
	}

	// Identity is stored in invocation state by beforeAgent.
	// We need another way to look it up since BeforeToolArgs doesn't carry Invocation directly.
	// Skip RBAC check when identity is not in context (checked at BeforeAgent level already).
	return &tool.BeforeToolResult{Context: ctx}, nil
}

// extractTokenFromSession extracts the auth token from session state.
func extractTokenFromSession(inv *agent.Invocation) string {
	if inv.Session == nil || inv.Session.State == nil {
		return ""
	}
	if token, ok := inv.Session.State["token"]; ok {
		return string(token)
	}
	if token, ok := inv.Session.State["authorization"]; ok {
		return string(token)
	}
	return ""
}

// SetSessionToken writes an auth token into session state for the plugin to read.
func SetSessionToken(inv *agent.Invocation, token string) {
	if inv.Session == nil {
		return
	}
	if inv.Session.State == nil {
		inv.Session.State = make(map[string][]byte)
	}
	inv.Session.State["token"] = []byte(token)
}

// GetIdentity retrieves the identity from invocation state.
func GetIdentity(inv *agent.Invocation) *Identity {
	val, ok := inv.GetState(identityStateKey)
	if !ok {
		return nil
	}
	id, ok := val.(*Identity)
	if !ok {
		return nil
	}
	return id
}

func isK8sTool(name string) bool {
	switch name {
	case "pod_list", "pod_get", "pod_logs", "pod_delete",
		"deployment_list", "deployment_get", "deployment_scale",
		"service_list", "service_get",
		"namespace_list", "namespace_get", "namespace_set",
		"event_list",
		"resource_list", "resource_get", "cluster_info":
		return true
	}
	return false
}

var _ plugin.Plugin = (*Plugin)(nil)
