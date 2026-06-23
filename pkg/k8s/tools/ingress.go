package tools

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/k8s"
)

// IngressListReq is the input for listing ingresses.
type IngressListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// IngressListRsp is the output.
type IngressListRsp struct {
	Ingresses []IngressSummary `json:"ingresses"`
	Err       string           `json:"error,omitempty"`
}

// IngressSummary is a simplified ingress view.
type IngressSummary struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Hosts     []string `json:"hosts"`
	Class     string   `json:"class"`
	Age       string   `json:"age"`
}

// NewIngressListTool creates a tool for listing ingresses.
func NewIngressListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req IngressListReq) (IngressListRsp, error) {
			return ingressList(ctx, clients, req)
		},
		function.WithName("ingress_list"),
		function.WithDescription("List ingresses with hosts and ingress class."),
	)
}

func ingressList(ctx context.Context, c *k8s.Clients, req IngressListReq) (IngressListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	list, err := c.ClientSet.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return IngressListRsp{Err: err.Error()}, nil
	}
	summaries := make([]IngressSummary, 0, len(list.Items))
	for _, ing := range list.Items {
		hosts := make([]string, 0, len(ing.Spec.Rules))
		for _, rule := range ing.Spec.Rules {
			if rule.Host != "" {
				hosts = append(hosts, rule.Host)
			}
		}
		class := ""
		if ing.Spec.IngressClassName != nil {
			class = *ing.Spec.IngressClassName
		}
		summaries = append(summaries, IngressSummary{
			Name:      ing.Name,
			Namespace: ing.Namespace,
			Hosts:     hosts,
			Class:     class,
			Age:       ageStr(ing.CreationTimestamp.Time),
		})
	}
	return IngressListRsp{Ingresses: summaries}, nil
}

// IngressGetReq is the input for describing an ingress.
type IngressGetReq struct {
	Name      string `json:"name" jsonschema:"description=ingress name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// IngressGetRsp is the output.
type IngressGetRsp struct {
	Ingress *IngressDetail `json:"ingress"`
	Err     string         `json:"error,omitempty"`
}

// IngressDetail contains key ingress information.
type IngressDetail struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Class     string            `json:"class"`
	Hosts     []string          `json:"hosts"`
	TLS       []string          `json:"tls,omitempty"`
	Rules     []IngressRuleInfo `json:"rules"`
	Labels    map[string]string `json:"labels"`
}

// IngressRuleInfo describes a single ingress rule.
type IngressRuleInfo struct {
	Host  string            `json:"host"`
	Paths []IngressPathInfo `json:"paths"`
}

// IngressPathInfo describes a single path within a rule.
type IngressPathInfo struct {
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	ServicePort string `json:"service_port"`
}

// NewIngressGetTool creates a tool for describing an ingress.
func NewIngressGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req IngressGetReq) (IngressGetRsp, error) {
			return ingressGet(ctx, clients, req)
		},
		function.WithName("ingress_get"),
		function.WithDescription("Get detailed information about an ingress including hosts, TLS, and routing rules."),
	)
}

func ingressGet(ctx context.Context, c *k8s.Clients, req IngressGetReq) (IngressGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	ing, err := c.ClientSet.NetworkingV1().Ingresses(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return IngressGetRsp{Err: err.Error()}, nil
	}

	hosts := make([]string, 0)
	tlsHosts := make([]string, 0)
	rules := make([]IngressRuleInfo, 0)

	for _, t := range ing.Spec.TLS {
		tlsHosts = append(tlsHosts, t.Hosts...)
	}

	for _, rule := range ing.Spec.Rules {
		host := rule.Host
		if host != "" {
			hosts = append(hosts, host)
		}
		paths := make([]IngressPathInfo, 0)
		if rule.HTTP != nil {
			for _, p := range rule.HTTP.Paths {
				svcPort := ""
				if p.Backend.Service.Port.Number > 0 {
					svcPort = fmt.Sprintf("%d", p.Backend.Service.Port.Number)
				} else if p.Backend.Service.Port.Name != "" {
					svcPort = p.Backend.Service.Port.Name
				}
				paths = append(paths, IngressPathInfo{
					Path:        p.Path,
					ServiceName: p.Backend.Service.Name,
					ServicePort: svcPort,
				})
			}
		}
		rules = append(rules, IngressRuleInfo{Host: host, Paths: paths})
	}

	class := ""
	if ing.Spec.IngressClassName != nil {
		class = *ing.Spec.IngressClassName
	}

	return IngressGetRsp{Ingress: &IngressDetail{
		Name:      ing.Name,
		Namespace: ing.Namespace,
		Class:     class,
		Hosts:     hosts,
		TLS:       tlsHosts,
		Rules:     rules,
		Labels:    ing.Labels,
	}}, nil
}
