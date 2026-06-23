package tools

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/k8s"
)

// HPAListReq is the input for listing HPAs.
type HPAListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// HPAListRsp is the output.
type HPAListRsp struct {
	HPAs []HPASummary `json:"hpas"`
	Err  string       `json:"error,omitempty"`
}

// HPASummary is a simplified HPA view.
type HPASummary struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	TargetKind  string `json:"target_kind"`
	TargetName  string `json:"target_name"`
	MinReplicas int32  `json:"min_replicas"`
	MaxReplicas int32  `json:"max_replicas"`
	CurrentCPU  string `json:"current_cpu"`
	Age         string `json:"age"`
}

// NewHPAListTool creates a tool for listing HPAs.
func NewHPAListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req HPAListReq) (HPAListRsp, error) {
			return hpaList(ctx, clients, req)
		},
		function.WithName("hpa_list"),
		function.WithDescription("List Horizontal Pod Autoscalers with target, min/max replicas, and current CPU usage."),
	)
}

func hpaList(ctx context.Context, c *k8s.Clients, req HPAListReq) (HPAListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	list, err := c.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return HPAListRsp{Err: err.Error()}, nil
	}
	summaries := make([]HPASummary, 0, len(list.Items))
	for _, h := range list.Items {
		targetKind := h.Spec.ScaleTargetRef.Kind
		targetName := h.Spec.ScaleTargetRef.Name
		currentCPU := ""
		for _, m := range h.Status.CurrentMetrics {
			if m.Resource != nil && m.Resource.Name == "cpu" {
				if m.Resource.Current.AverageUtilization != nil {
					currentCPU = fmt.Sprintf("%d%%", *m.Resource.Current.AverageUtilization)
				}
			}
		}
		var minReplicas, maxReplicas int32
		if h.Spec.MinReplicas != nil {
			minReplicas = *h.Spec.MinReplicas
		}
		maxReplicas = h.Spec.MaxReplicas

		summaries = append(summaries, HPASummary{
			Name:        h.Name,
			Namespace:   h.Namespace,
			TargetKind:  targetKind,
			TargetName:  targetName,
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			CurrentCPU:  currentCPU,
			Age:         ageStr(h.CreationTimestamp.Time),
		})
	}
	return HPAListRsp{HPAs: summaries}, nil
}

// HPAGetReq is the input for describing an HPA.
type HPAGetReq struct {
	Name      string `json:"name" jsonschema:"description=hpa name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// HPAGetRsp is the output.
type HPAGetRsp struct {
	HPA *HPADetail `json:"hpa"`
	Err string     `json:"error,omitempty"`
}

// HPADetail contains key HPA information.
type HPADetail struct {
	Name        string   `json:"name"`
	Namespace   string   `json:"namespace"`
	TargetKind  string   `json:"target_kind"`
	TargetName  string   `json:"target_name"`
	MinReplicas int32    `json:"min_replicas"`
	MaxReplicas int32    `json:"max_replicas"`
	Metrics     []string `json:"metrics"`
	Current     []string `json:"current"`
	Conditions  []string `json:"conditions"`
}

// NewHPAGetTool creates a tool for describing an HPA.
func NewHPAGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req HPAGetReq) (HPAGetRsp, error) {
			return hpaGet(ctx, clients, req)
		},
		function.WithName("hpa_get"),
		function.WithDescription("Get detailed information about an HPA including metrics, current load, and conditions."),
	)
}

func hpaGet(ctx context.Context, c *k8s.Clients, req HPAGetReq) (HPAGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	h, err := c.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return HPAGetRsp{Err: err.Error()}, nil
	}

	var minReplicas int32
	if h.Spec.MinReplicas != nil {
		minReplicas = *h.Spec.MinReplicas
	}

	metrics := make([]string, 0, len(h.Spec.Metrics))
	for _, m := range h.Spec.Metrics {
		if m.Resource != nil {
			var target string
			if m.Resource.Target.AverageUtilization != nil {
				target = fmt.Sprintf("%d%%", *m.Resource.Target.AverageUtilization)
			} else if !m.Resource.Target.AverageValue.IsZero() {
				target = m.Resource.Target.AverageValue.String()
			}
			metrics = append(metrics, fmt.Sprintf("%s: %s", m.Resource.Name, target))
		}
	}

	current := make([]string, 0, len(h.Status.CurrentMetrics))
	for _, m := range h.Status.CurrentMetrics {
		if m.Resource != nil {
			var val string
			if m.Resource.Current.AverageUtilization != nil {
				val = fmt.Sprintf("%d%%", *m.Resource.Current.AverageUtilization)
			} else if !m.Resource.Current.AverageValue.IsZero() {
				val = m.Resource.Current.AverageValue.String()
			}
			current = append(current, fmt.Sprintf("%s: %s", m.Resource.Name, val))
		}
	}

	conditions := make([]string, 0, len(h.Status.Conditions))
	for _, cond := range h.Status.Conditions {
		conditions = append(conditions, fmt.Sprintf("%s: %s (%s)",
			cond.Type, cond.Status, cond.Reason))
	}

	return HPAGetRsp{HPA: &HPADetail{
		Name:        h.Name,
		Namespace:   h.Namespace,
		TargetKind:  h.Spec.ScaleTargetRef.Kind,
		TargetName:  h.Spec.ScaleTargetRef.Name,
		MinReplicas: minReplicas,
		MaxReplicas: h.Spec.MaxReplicas,
		Metrics:     metrics,
		Current:     current,
		Conditions:  conditions,
	}}, nil
}
