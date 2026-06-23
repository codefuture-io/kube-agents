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

// Package memory provides long-term memory persistence for kube-agents.
package memory

import (
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/memory"
	memoryinmemory "trpc.group/trpc-go/trpc-agent-go/memory/inmemory"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// NewService creates a memory service based on the configuration backend.
func NewService(cfg config.StoreConfig) (memory.Service, error) {
	switch cfg.Backend {
	case "redis":
		return nil, fmt.Errorf("redis memory backend not yet implemented")
	default:
		return memoryinmemory.NewMemoryService(), nil
	}
}
