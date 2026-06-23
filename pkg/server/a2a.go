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

package server

import (
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	a2aserver "trpc.group/trpc-go/trpc-agent-go/server/a2a"
	"trpc.group/trpc-go/trpc-agent-go/session"

	"github.com/codefuture-io/kube-agents/pkg/config"
)

// StartA2A starts an Agent-to-Agent protocol server.
func StartA2A(cfg config.A2AServerConfig, ag agent.Agent, sessionSvc session.Service) error {
	if !cfg.Enabled {
		return nil
	}

	var opts []a2aserver.Option
	opts = append(opts, a2aserver.WithAgent(ag, true))
	opts = append(opts, a2aserver.WithHost(cfg.Host))
	if sessionSvc != nil {
		opts = append(opts, a2aserver.WithSessionService(sessionSvc))
	}

	srv, err := a2aserver.New(opts...)
	if err != nil {
		return fmt.Errorf("create A2A server: %w", err)
	}

	go func() {
		fmt.Printf("A2A server starting on %s\n", cfg.Host)
		if err := srv.Start(cfg.Host); err != nil {
			fmt.Printf("A2A server error: %v\n", err)
		}
	}()

	return nil
}

// StartA2AWithRunner starts the A2A server using a pre-built runner.
func StartA2AWithRunner(cfg config.A2AServerConfig, r runner.Runner, sessionSvc session.Service) error {
	if !cfg.Enabled {
		return nil
	}

	var opts []a2aserver.Option
	opts = append(opts, a2aserver.WithRunner(r))
	opts = append(opts, a2aserver.WithHost(cfg.Host))
	if sessionSvc != nil {
		opts = append(opts, a2aserver.WithSessionService(sessionSvc))
	}

	srv, err := a2aserver.New(opts...)
	if err != nil {
		return fmt.Errorf("create A2A server: %w", err)
	}

	go func() {
		fmt.Printf("A2A server starting on %s\n", cfg.Host)
		if err := srv.Start(cfg.Host); err != nil {
			fmt.Printf("A2A server error: %v\n", err)
		}
	}()

	return nil
}
