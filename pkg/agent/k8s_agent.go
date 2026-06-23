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

// Package agent provides agent creation and configuration for kube-agents.
package agent

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// Options configures the K8s agent.
type Options struct {
	ModelName  string
	APIKey     string
	BaseURL    string
	Tools      []tool.Tool
	ToolSets   []tool.ToolSet
	Generation model.GenerationConfig
}

// NewK8sAgent creates the primary Kubernetes assistant LLMAgent.
func NewK8sAgent(opts Options) (agent.Agent, error) {
	var modelOpts []openai.Option
	modelOpts = append(modelOpts, openai.WithAPIKey(opts.APIKey))
	if opts.BaseURL != "" {
		modelOpts = append(modelOpts, openai.WithBaseURL(opts.BaseURL))
	}
	modelOpts = append(modelOpts, openai.WithVariant(openai.VariantDeepSeek))

	modelInstance := openai.New(opts.ModelName, modelOpts...)

	genConfig := opts.Generation
	if genConfig.Stream == false && genConfig.Temperature == nil && genConfig.MaxTokens == nil {
		genConfig = model.GenerationConfig{
			Stream:      true,
			Temperature: ptrOf(0.7),
		}
	}

	llmAgent := llmagent.New("k8s-assistant",
		llmagent.WithModel(modelInstance),
		llmagent.WithDescription("Kubernetes cluster operations assistant"),
		llmagent.WithInstruction(
			"You are a Kubernetes operations assistant. "+
				"You help users manage and diagnose K8s resources. "+
				"Use the available tools to list, describe, create, and troubleshoot resources. "+
				"Always verify the cluster state before making changes. "+
				"Current user identity: {user:identity?}. "+
				"Current namespace: {temp:namespace?}.",
		),
		llmagent.WithTools(opts.Tools),
		llmagent.WithToolSets(opts.ToolSets),
		llmagent.WithGenerationConfig(genConfig),
	)

	return llmAgent, nil
}

// CalculatorTool returns the basic arithmetic calculator tool.
func CalculatorTool() tool.Tool {
	return function.NewFunctionTool(
		calculator,
		function.WithName("calculator"),
		function.WithDescription(
			"Performs basic arithmetic: add, sub, mul, div. "+
				"Parameters: a, b are numbers, op is one of add/sub/mul/div",
		),
	)
}

// CalculatorReq is the input for the calculator tool.
type CalculatorReq struct {
	A  float64 `json:"a" jsonschema:"description=first operand,required"`
	B  float64 `json:"b" jsonschema:"description=second operand,required"`
	Op string  `json:"op" jsonschema:"description=operation,enum=add,enum=sub,enum=mul,enum=div,required"`
}

// CalculatorRsp is the output of the calculator tool.
type CalculatorRsp struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func calculator(_ context.Context, req CalculatorReq) (CalculatorRsp, error) {
	var result float64
	switch req.Op {
	case "add":
		result = req.A + req.B
	case "sub":
		result = req.A - req.B
	case "mul":
		result = req.A * req.B
	case "div":
		if req.B == 0 {
			return CalculatorRsp{Error: "division by zero"}, nil
		}
		result = req.A / req.B
	default:
		return CalculatorRsp{Error: "unknown operation: " + req.Op}, nil
	}
	return CalculatorRsp{Result: result}, nil
}

func ptrOf[T any](v T) *T {
	return &v
}

// MustNewLLMAgent creates a simple LLMAgent from pre-built pieces.
// This is a convenience wrapper for CLI/interactive use.
func MustNewLLMAgent(m model.Model, tools []tool.Tool, genCfg model.GenerationConfig) agent.Agent {
	return llmagent.New("assistant",
		llmagent.WithModel(m),
		llmagent.WithDescription("General-purpose assistant"),
		llmagent.WithInstruction("You are a helpful assistant."),
		llmagent.WithTools(tools),
		llmagent.WithGenerationConfig(genCfg),
	)
}
