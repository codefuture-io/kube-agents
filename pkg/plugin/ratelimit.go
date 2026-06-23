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
	"fmt"
	"sync"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/plugin"
)

// RateLimitPlugin provides simple per-session request rate limiting.
type RateLimitPlugin struct {
	mu       sync.Mutex
	counters map[string]int
}

// NewRateLimitPlugin creates a rate limiting plugin.
func NewRateLimitPlugin() *RateLimitPlugin {
	return &RateLimitPlugin{counters: make(map[string]int)}
}

// Name returns the plugin identifier.
func (p *RateLimitPlugin) Name() string { return "ratelimit" }

// Register hooks into the plugin registry.
func (p *RateLimitPlugin) Register(r *plugin.Registry) {
	r.BeforeAgent(p.beforeAgent)
}

func (p *RateLimitPlugin) beforeAgent(
	ctx context.Context,
	args *agent.BeforeAgentArgs,
) (*agent.BeforeAgentResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	sessionID := args.Invocation.Session.ID
	p.counters[sessionID]++
	if p.counters[sessionID] > 100 {
		return nil, fmt.Errorf("rate limit exceeded (100 requests per session)")
	}
	return &agent.BeforeAgentResult{Context: ctx}, nil
}

var _ plugin.Plugin = (*RateLimitPlugin)(nil)
