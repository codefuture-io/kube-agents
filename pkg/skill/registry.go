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

// Package skill provides a skill registry for kube-agents.
package skill

import (
	"log/slog"

	"trpc.group/trpc-go/trpc-agent-go/skill"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	skilltool "trpc.group/trpc-go/trpc-agent-go/tool/skill"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// Registry manages skill repositories and tools.
type Registry struct {
	repo  *skill.FSRepository
	tools []tool.Tool
}

// NewRegistry creates a skill registry from a directory of SKILL.md files.
func NewRegistry(cfg config.Config) (*Registry, error) {
	dir := cfg.SkillsDir
	if dir == "" {
		dir = "./skills"
	}

	repo, err := skill.NewFSRepository(dir)
	if err != nil {
		return nil, err
	}

	tools := []tool.Tool{
		skilltool.NewLoadTool(repo),
		skilltool.NewRunTool(repo, nil),
		skilltool.NewListDocsTool(repo),
		skilltool.NewSelectDocsTool(repo),
	}

	slog.Info("skills loaded", "dir", dir)

	return &Registry{repo: repo, tools: tools}, nil
}

// Tools returns skill-related tools for agent registration.
func (r *Registry) Tools() []tool.Tool {
	return r.tools
}

// Refresh rescans the skill directory.
func (r *Registry) Refresh() {
	if r.repo != nil {
		r.repo.Refresh()
	}
}
