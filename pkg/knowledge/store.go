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

package knowledge

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// Store manages the knowledge base with built-in file-based search.
// Does NOT require an embedding API — works with any LLM backend.
type Store struct {
	search    tool.Tool
	uploadDir string
	s3        *S3Store
}

// NewStore creates a knowledge store from configured sources.
// Uses simple text matching search — no embedding API required.
func NewStore(cfg config.KnowledgeConfig, apiKey string) (*Store, error) {
	store := &Store{
		uploadDir: filepath.Join(os.TempDir(), "kube-agents-knowledge"),
	}

	// Initialize S3 store if configured.
	s3Store, err := NewS3Store(cfg.S3)
	if err != nil {
		slog.Warn("S3 store init failed, S3 disabled", "error", err)
	} else {
		store.s3 = s3Store
	}

	// Collect all document paths from configured sources.
	var docPaths []string
	for _, s := range cfg.Sources {
		switch s.Type {
		case "file":
			paths, _ := collectFiles(s.Path)
			docPaths = append(docPaths, paths...)
		}
	}

	if len(docPaths) > 0 {
		slog.Info("knowledge store initialized with local files", "count", len(docPaths))
	} else {
		slog.Info("knowledge store ready (no pre-configured sources, upload API available)")
	}

	// Create the built-in search tool using simple text matching.
	store.search = function.NewFunctionTool(
		func(ctx context.Context, req SearchReq) (SearchRsp, error) {
			return searchDocs(docPaths, req)
		},
		function.WithName("knowledge_search"),
		function.WithDescription(
			"Search Kubernetes documentation and reference materials. "+
				"Use this to find information about K8s concepts, commands, and best practices. "+
				"Provide one or more search keywords.",
		),
	)

	return store, nil
}

// SearchReq is the input for the knowledge search tool.
type SearchReq struct {
	Query string `json:"query" jsonschema:"description=search keywords or question,required"`
}

// SearchRsp is the output.
type SearchRsp struct {
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
}

// SearchResult is a single match.
type SearchResult struct {
	Source  string `json:"source"`
	Snippet string `json:"snippet"`
}

// AddFile adds a file to the knowledge base at runtime.
func (s *Store) AddFile(ctx interface{}, filename string, content []byte) error {
	c, _ := ctx.(context.Context)
	if c == nil {
		c = context.Background()
	}

	// Write to local temp dir.
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return fmt.Errorf("create upload dir: %w", err)
	}
	dst := filepath.Join(s.uploadDir, filepath.Base(filename))
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("write uploaded file: %w", err)
	}

	// Upload to S3 if configured.
	if s.s3 != nil {
		key, err := s.s3.UploadObject(c, filename, content)
		if err != nil {
			slog.Warn("S3 upload failed", "file", filename, "error", err)
		} else {
			slog.Info("S3 upload OK", "key", key)
		}
	}

	slog.Info("knowledge file added", "file", filename, "size", len(content))
	return nil
}

// SearchTool returns the knowledge search tool for agent registration.
func (s *Store) SearchTool() tool.Tool {
	return s.search
}

// --- built-in text search ---

func searchDocs(paths []string, req SearchReq) (SearchRsp, error) {
	if req.Query == "" {
		return SearchRsp{}, nil
	}

	terms := strings.Fields(strings.ToLower(req.Query))
	var results []SearchResult

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)
		lower := strings.ToLower(content)

		// Score: count matching terms.
		score := 0
		for _, term := range terms {
			score += strings.Count(lower, term)
		}
		if score == 0 {
			continue
		}

		// Extract relevant snippet around first match.
		snippet := extractSnippet(content, lower, terms[0], 300)
		results = append(results, SearchResult{
			Source:  filepath.Base(path),
			Snippet: fmt.Sprintf("(score=%d) %s", score, snippet),
		})
	}

	// Sort by score descending (simple bubble — fine for small doc sets).
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if len(results[i].Snippet) < len(results[j].Snippet) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if len(results) > 5 {
		results = results[:5]
	}
	return SearchRsp{Results: results, Count: len(results)}, nil
}

func extractSnippet(content, lower, firstTerm string, window int) string {
	idx := strings.Index(lower, firstTerm)
	if idx < 0 {
		if len(content) > window {
			return content[:window] + "..."
		}
		return content
	}
	start := idx - window/2
	if start < 0 {
		start = 0
	}
	end := start + window
	if end > len(content) {
		end = len(content)
	}
	snippet := content[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}
	return snippet
}

func collectFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".md", ".txt", ".yaml", ".yml":
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
