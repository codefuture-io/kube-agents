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

// EventListReq is the input for listing events.
type EventListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace filter,omitempty"`
	Type      string `json:"type" jsonschema:"description=event type (Normal or Warning),omitempty"`
	Limit     int64  `json:"limit" jsonschema:"description=max events to return (default 20),omitempty"`
}

// EventListRsp is the output.
type EventListRsp struct {
	Events []EventSummary `json:"events"`
	Err    string         `json:"error,omitempty"`
}

// EventSummary is a simplified event view.
type EventSummary struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Object  string `json:"object"`
	Message string `json:"message"`
	Age     string `json:"age"`
}

// NewEventListTool creates a tool for listing cluster events.
func NewEventListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req EventListReq) (EventListRsp, error) {
			return eventList(ctx, clients, req)
		},
		function.WithName("event_list"),
		function.WithDescription("List recent cluster events, optionally filtered by namespace and type (Normal/Warning)."),
	)
}

func eventList(ctx context.Context, c *k8s.Clients, req EventListReq) (EventListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}

	opts := metav1.ListOptions{}
	if limit := req.Limit; limit > 0 {
		opts.Limit = limit
	} else {
		opts.Limit = 20
	}
	if req.Type != "" {
		opts.FieldSelector = fmt.Sprintf("type=%s", req.Type)
	}

	list, err := c.ClientSet.CoreV1().Events(ns).List(ctx, opts)
	if err != nil {
		return EventListRsp{Err: err.Error()}, nil
	}
	summaries := make([]EventSummary, 0, len(list.Items))
	for _, e := range list.Items {
		summaries = append(summaries, EventSummary{
			Type:    e.Type,
			Reason:  e.Reason,
			Object:  fmt.Sprintf("%s/%s", e.InvolvedObject.Kind, e.InvolvedObject.Name),
			Message: e.Message,
			Age:     ageStr(e.LastTimestamp.Time),
		})
	}
	return EventListRsp{Events: summaries}, nil
}
