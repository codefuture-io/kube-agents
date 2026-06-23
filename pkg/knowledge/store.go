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

// Package knowledge provides RAG knowledge retrieval for kube-agents.
package knowledge

import (
	"context"
	"log/slog"

	"trpc.group/trpc-go/trpc-agent-go/knowledge"
	"trpc.group/trpc-go/trpc-agent-go/knowledge/source/file"
	"trpc.group/trpc-go/trpc-agent-go/knowledge/source/url"
	knowledgetool "trpc.group/trpc-go/trpc-agent-go/knowledge/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// Store manages the RAG knowledge base.
type Store struct {
	kb     knowledge.Knowledge
	search tool.Tool
}

// NewStore creates a knowledge store from configured sources.
func NewStore(cfg config.KnowledgeConfig) (*Store, error) {
	var filePaths []string
	var urls []string
	for _, s := range cfg.Sources {
		switch s.Type {
		case "file":
			filePaths = append(filePaths, s.Path)
		case "url":
			urls = append(urls, s.Path)
		}
	}

	if len(filePaths) == 0 && len(urls) == 0 {
		slog.Info("knowledge skipped: no sources")
		return &Store{}, nil
	}

	kb := knowledge.New()
	ctx := context.Background()

	if len(filePaths) > 0 {
		kb.AddSource(ctx, file.New(filePaths))
	}
	if len(urls) > 0 {
		kb.AddSource(ctx, url.New(urls))
	}

	go func() {
		if err := kb.Load(ctx); err != nil {
			slog.Warn("knowledge load error", "error", err)
		}
	}()

	search := knowledgetool.NewKnowledgeSearchTool(kb,
		knowledgetool.WithToolName("knowledge_search"),
		knowledgetool.WithToolDescription("Search K8s documentation and reference materials."),
	)

	slog.Info("knowledge initialized", "fileSources", len(filePaths), "urlSources", len(urls))
	return &Store{kb: kb, search: search}, nil
}

// SearchTool returns the knowledge search tool for agent registration.
func (s *Store) SearchTool() tool.Tool {
	return s.search
}
