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

package tools

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/k8s"
)

// NamespaceListReq is the input for listing namespaces.
type NamespaceListReq struct{}

// NamespaceListRsp is the output.
type NamespaceListRsp struct {
	Namespaces []NsSummary `json:"namespaces"`
	Err        string      `json:"error,omitempty"`
}

// NsSummary is a simplified namespace view.
type NsSummary struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Age    string `json:"age"`
}

// NewNamespaceListTool creates a tool for listing namespaces.
func NewNamespaceListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req NamespaceListReq) (NamespaceListRsp, error) {
			return nsList(ctx, clients)
		},
		function.WithName("namespace_list"),
		function.WithDescription("List all namespaces in the cluster."),
	)
}

func nsList(ctx context.Context, c *k8s.Clients) (NamespaceListRsp, error) {
	list, err := c.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return NamespaceListRsp{Err: err.Error()}, nil
	}
	summaries := make([]NsSummary, 0, len(list.Items))
	for _, ns := range list.Items {
		summaries = append(summaries, NsSummary{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			Age:    ageStr(ns.CreationTimestamp.Time),
		})
	}
	return NamespaceListRsp{Namespaces: summaries}, nil
}

// NamespaceGetReq is the input for describing a namespace.
type NamespaceGetReq struct {
	Name string `json:"name" jsonschema:"description=namespace name,required"`
}

// NamespaceGetRsp is the output.
type NamespaceGetRsp struct {
	Namespace *NsDetail `json:"namespace"`
	Err       string    `json:"error,omitempty"`
}

// NsDetail contains namespace details.
type NsDetail struct {
	Name   string            `json:"name"`
	Status string            `json:"status"`
	Labels map[string]string `json:"labels"`
}

// NewNamespaceGetTool creates a tool for describing a namespace.
func NewNamespaceGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req NamespaceGetReq) (NamespaceGetRsp, error) {
			return nsGet(ctx, clients, req)
		},
		function.WithName("namespace_get"),
		function.WithDescription("Get details of a namespace including labels."),
	)
}

func nsGet(ctx context.Context, c *k8s.Clients, req NamespaceGetReq) (NamespaceGetRsp, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return NamespaceGetRsp{Err: err.Error()}, nil
	}
	return NamespaceGetRsp{Namespace: &NsDetail{
		Name:   ns.Name,
		Status: string(ns.Status.Phase),
		Labels: ns.Labels,
	}}, nil
}

// SetNamespaceReq sets the current session namespace.
type SetNamespaceReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace to set as current,required"`
}

// SetNamespaceRsp returns confirmation.
type SetNamespaceRsp struct {
	Message string `json:"message"`
}

// NewSetNamespaceTool creates a tool for changing the active namespace.
func NewSetNamespaceTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req SetNamespaceReq) (SetNamespaceRsp, error) {
			prev := clients.Namespace
			clients.Namespace = req.Namespace
			return SetNamespaceRsp{
				Message: fmt.Sprintf("namespace changed from %s to %s", prev, req.Namespace),
			}, nil
		},
		function.WithName("namespace_set"),
		function.WithDescription("Set the current namespace for subsequent operations."),
	)
}
