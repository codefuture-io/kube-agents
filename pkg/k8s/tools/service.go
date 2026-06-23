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

// ServiceListReq is the input for listing services.
type ServiceListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// ServiceListRsp is the output for the service list.
type ServiceListRsp struct {
	Services []SvcSummary `json:"services"`
	Err      string       `json:"error,omitempty"`
}

// SvcSummary is a simplified service view.
type SvcSummary struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Type      string   `json:"type"`
	ClusterIP string   `json:"cluster_ip"`
	Ports     []string `json:"ports"`
	Age       string   `json:"age"`
}

// NewServiceListTool creates a tool for listing services.
func NewServiceListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req ServiceListReq) (ServiceListRsp, error) {
			return svcList(ctx, clients, req)
		},
		function.WithName("service_list"),
		function.WithDescription("List services with type, cluster IP, and ports."),
	)
}

func svcList(ctx context.Context, c *k8s.Clients, req ServiceListReq) (ServiceListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	list, err := c.ClientSet.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return ServiceListRsp{Err: err.Error()}, nil
	}
	summaries := make([]SvcSummary, 0, len(list.Items))
	for _, s := range list.Items {
		ports := make([]string, 0, len(s.Spec.Ports))
		for _, p := range s.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		summaries = append(summaries, SvcSummary{
			Name:      s.Name,
			Namespace: s.Namespace,
			Type:      string(s.Spec.Type),
			ClusterIP: s.Spec.ClusterIP,
			Ports:     ports,
			Age:       ageStr(s.CreationTimestamp.Time),
		})
	}
	return ServiceListRsp{Services: summaries}, nil
}

// ServiceGetReq is the input for describing a service.
type ServiceGetReq struct {
	Name      string `json:"name" jsonschema:"description=service name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// ServiceGetRsp is the output for a service detail.
type ServiceGetRsp struct {
	Service *SvcDetail `json:"service"`
	Err     string     `json:"error,omitempty"`
}

// SvcDetail contains key service information.
type SvcDetail struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Type       string            `json:"type"`
	ClusterIP  string            `json:"cluster_ip"`
	ExternalIP string            `json:"external_ip,omitempty"`
	Ports      []string          `json:"ports"`
	Selector   map[string]string `json:"selector"`
}

// NewServiceGetTool creates a tool for describing a service.
func NewServiceGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req ServiceGetReq) (ServiceGetRsp, error) {
			return svcGet(ctx, clients, req)
		},
		function.WithName("service_get"),
		function.WithDescription("Get details of a service including type, ports, and selector."),
	)
}

func svcGet(ctx context.Context, c *k8s.Clients, req ServiceGetReq) (ServiceGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	s, err := c.ClientSet.CoreV1().Services(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return ServiceGetRsp{Err: err.Error()}, nil
	}

	ports := make([]string, 0, len(s.Spec.Ports))
	for _, p := range s.Spec.Ports {
		ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
	}

	extIP := ""
	if len(s.Status.LoadBalancer.Ingress) > 0 {
		extIP = s.Status.LoadBalancer.Ingress[0].IP
	}

	return ServiceGetRsp{Service: &SvcDetail{
		Name:       s.Name,
		Namespace:  s.Namespace,
		Type:       string(s.Spec.Type),
		ClusterIP:  s.Spec.ClusterIP,
		ExternalIP: extIP,
		Ports:      ports,
		Selector:   s.Spec.Selector,
	}}, nil
}
