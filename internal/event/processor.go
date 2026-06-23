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

// Package event provides streaming event processing utilities.
package event

import (
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/event"
)

// StreamProcessor buffers reasoning content from DeepSeek and formats
// streaming output for display. DeepSeek emits reasoning content in
// many tiny chunks; buffering them into a single block produces clean output.
type StreamProcessor struct {
	reasoningBuf     strings.Builder
	reasoningFlushed bool
}

// NewStreamProcessor creates a new StreamProcessor.
func NewStreamProcessor() *StreamProcessor {
	return &StreamProcessor{}
}

// Process handles a single streaming event chunk.
// It returns the text that should be printed immediately, or an empty string
// if the content should be buffered.
func (p *StreamProcessor) Process(ev *event.Event) string {
	if ev.Object != "chat.completion.chunk" {
		return ""
	}

	delta := ev.Response.Choices[0].Delta

	// Buffer reasoning content — emitted in tiny chunks by DeepSeek.
	if delta.ReasoningContent != "" {
		p.reasoningBuf.WriteString(delta.ReasoningContent)
	}

	if delta.Content != "" {
		var out string
		if !p.reasoningFlushed {
			p.reasoningFlushed = true
			if p.reasoningBuf.Len() > 0 {
				out = fmt.Sprintf("\n[Reasoning]\n%s\n[/Reasoning]\n\n",
					p.reasoningBuf.String())
			}
		}
		out += delta.Content
		return out
	}

	return ""
}

// Finalize returns any remaining buffered content (e.g., reasoning-only
// responses) that should be printed after the event stream ends.
func (p *StreamProcessor) Finalize() string {
	if !p.reasoningFlushed && p.reasoningBuf.Len() > 0 {
		return fmt.Sprintf("\n[Reasoning]\n%s\n[/Reasoning]\n",
			p.reasoningBuf.String())
	}
	return ""
}
