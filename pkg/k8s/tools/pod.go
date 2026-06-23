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
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/k8s"
)

// PodListReq is the input for listing pods.
type PodListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace to list pods from,omitempty"`
	Label     string `json:"label" jsonschema:"description=label selector (e.g. app=nginx),omitempty"`
}

// PodListRsp is the output for the pod list.
type PodListRsp struct {
	Pods []PodSummary `json:"pods"`
	Err  string       `json:"error,omitempty"`
}

// PodSummary is a simplified pod view.
type PodSummary struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Node      string `json:"node"`
	IP        string `json:"ip"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
}

// NewPodListTool creates a tool for listing pods.
func NewPodListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req PodListReq) (PodListRsp, error) {
			return podList(ctx, clients, req)
		},
		function.WithName("pod_list"),
		function.WithDescription("List pods in a namespace. Returns name, status, node, IP, restarts."),
	)
}

func podList(ctx context.Context, c *k8s.Clients, req PodListReq) (PodListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	opts := metav1.ListOptions{}
	if req.Label != "" {
		opts.LabelSelector = req.Label
	}
	list, err := c.ClientSet.CoreV1().Pods(ns).List(ctx, opts)
	if err != nil {
		return PodListRsp{Err: err.Error()}, nil
	}
	summaries := make([]PodSummary, 0, len(list.Items))
	for _, p := range list.Items {
		summaries = append(summaries, PodSummary{
			Name:      p.Name,
			Namespace: p.Namespace,
			Status:    string(p.Status.Phase),
			Node:      p.Spec.NodeName,
			IP:        p.Status.PodIP,
			Restarts:  restartCount(p),
			Age:       ageStr(p.CreationTimestamp.Time),
		})
	}
	return PodListRsp{Pods: summaries}, nil
}

// PodGetReq is the input for describing a pod.
type PodGetReq struct {
	Name      string `json:"name" jsonschema:"description=pod name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// PodGetRsp is the output for a pod detail.
type PodGetRsp struct {
	Pod *PodDetail `json:"pod"`
	Err string     `json:"error,omitempty"`
}

// PodDetail contains key pod information.
type PodDetail struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Status     string            `json:"status"`
	Node       string            `json:"node"`
	IP         string            `json:"ip"`
	Labels     map[string]string `json:"labels"`
	Containers []ContainerInfo   `json:"containers"`
	Conditions []string          `json:"conditions"`
	Events     []string          `json:"events,omitempty"`
}

// ContainerInfo summarizes a container.
type ContainerInfo struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Ready   bool   `json:"ready"`
	Restart int32  `json:"restarts"`
}

// NewPodGetTool creates a tool for describing a pod.
func NewPodGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req PodGetReq) (PodGetRsp, error) {
			return podGet(ctx, clients, req)
		},
		function.WithName("pod_get"),
		function.WithDescription("Get detailed information about a pod including containers, conditions, and recent events."),
	)
}

func podGet(ctx context.Context, c *k8s.Clients, req PodGetReq) (PodGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	p, err := c.ClientSet.CoreV1().Pods(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return PodGetRsp{Err: err.Error()}, nil
	}

	containers := make([]ContainerInfo, 0, len(p.Status.ContainerStatuses))
	for _, cs := range p.Status.ContainerStatuses {
		containers = append(containers, ContainerInfo{
			Name:    cs.Name,
			Image:   cs.Image,
			Ready:   cs.Ready,
			Restart: cs.RestartCount,
		})
	}

	conditions := make([]string, 0, len(p.Status.Conditions))
	for _, cond := range p.Status.Conditions {
		if cond.Status == corev1.ConditionTrue {
			conditions = append(conditions, fmt.Sprintf("%s: %s", cond.Type, cond.Reason))
		}
	}

	events, _ := getPodEvents(ctx, c, ns, p.Name)

	return PodGetRsp{Pod: &PodDetail{
		Name:       p.Name,
		Namespace:  p.Namespace,
		Status:     string(p.Status.Phase),
		Node:       p.Spec.NodeName,
		IP:         p.Status.PodIP,
		Labels:     p.Labels,
		Containers: containers,
		Conditions: conditions,
		Events:     events,
	}}, nil
}

// PodLogsReq is the input for fetching pod logs.
type PodLogsReq struct {
	Name      string `json:"name" jsonschema:"description=pod name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
	Container string `json:"container" jsonschema:"description=container name,omitempty"`
	Tail      int64  `json:"tail" jsonschema:"description=number of lines from end (default 50),omitempty"`
}

// PodLogsRsp holds log output.
type PodLogsRsp struct {
	Logs string `json:"logs"`
	Err  string `json:"error,omitempty"`
}

// NewPodLogsTool creates a tool for fetching pod logs.
func NewPodLogsTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req PodLogsReq) (PodLogsRsp, error) {
			return podLogs(ctx, clients, req)
		},
		function.WithName("pod_logs"),
		function.WithDescription("Fetch recent logs from a pod container. Defaults to last 50 lines."),
	)
}

func podLogs(ctx context.Context, c *k8s.Clients, req PodLogsReq) (PodLogsRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	tail := req.Tail
	if tail == 0 {
		tail = 50
	}
	opts := &corev1.PodLogOptions{TailLines: &tail}
	if req.Container != "" {
		opts.Container = req.Container
	}
	logs, err := c.ClientSet.CoreV1().Pods(ns).GetLogs(req.Name, opts).Do(ctx).Raw()
	if err != nil {
		return PodLogsRsp{Err: err.Error()}, nil
	}
	return PodLogsRsp{Logs: string(logs)}, nil
}

// PodDeleteReq is the input for deleting a pod.
type PodDeleteReq struct {
	Name      string `json:"name" jsonschema:"description=pod name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// PodDeleteRsp is the output after deletion.
type PodDeleteRsp struct {
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

// NewPodDeleteTool creates a tool for deleting a pod.
func NewPodDeleteTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req PodDeleteReq) (PodDeleteRsp, error) {
			return podDelete(ctx, clients, req)
		},
		function.WithName("pod_delete"),
		function.WithDescription("Delete a pod by name. This is a destructive operation."),
	)
}

func podDelete(ctx context.Context, c *k8s.Clients, req PodDeleteReq) (PodDeleteRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	err := c.ClientSet.CoreV1().Pods(ns).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return PodDeleteRsp{Err: err.Error()}, nil
	}
	return PodDeleteRsp{Message: fmt.Sprintf("pod %s/%s deleted", ns, req.Name)}, nil
}

// --- helpers ---

func restartCount(p corev1.Pod) int32 {
	var total int32
	for _, s := range p.Status.ContainerStatuses {
		total += s.RestartCount
	}
	return total
}

func getPodEvents(ctx context.Context, c *k8s.Clients, ns, name string) ([]string, error) {
	events, err := c.ClientSet.CoreV1().Events(ns).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(events.Items))
	for _, e := range events.Items {
		lines = append(lines, fmt.Sprintf("[%s] %s: %s",
			e.Type, e.Reason, e.Message))
	}
	return lines, nil
}

func ageStr(t interface{ String() string }) string {
	s := t.String()
	if idx := strings.IndexByte(s, '.'); idx > 0 {
		s = s[:idx] + "Z"
	}
	return s
}
