package tools

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	"github.com/codefuture-io/kube-agents/pkg/k8s"
)

// ConfigMapListReq is the input for listing ConfigMaps.
type ConfigMapListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
	Label     string `json:"label" jsonschema:"description=label selector,omitempty"`
}

// ConfigMapListRsp is the output.
type ConfigMapListRsp struct {
	ConfigMaps []ConfigMapSummary `json:"configmaps"`
	Err        string             `json:"error,omitempty"`
}

// ConfigMapSummary is a simplified ConfigMap view.
type ConfigMapSummary struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Keys      int    `json:"keys"`
	Age       string `json:"age"`
}

// NewConfigMapListTool creates a tool for listing ConfigMaps.
func NewConfigMapListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req ConfigMapListReq) (ConfigMapListRsp, error) {
			return configMapList(ctx, clients, req)
		},
		function.WithName("configmap_list"),
		function.WithDescription("List ConfigMaps with key count."),
	)
}

func configMapList(ctx context.Context, c *k8s.Clients, req ConfigMapListReq) (ConfigMapListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	opts := metav1.ListOptions{}
	if req.Label != "" {
		opts.LabelSelector = req.Label
	}
	list, err := c.ClientSet.CoreV1().ConfigMaps(ns).List(ctx, opts)
	if err != nil {
		return ConfigMapListRsp{Err: err.Error()}, nil
	}
	summaries := make([]ConfigMapSummary, 0, len(list.Items))
	for _, cm := range list.Items {
		summaries = append(summaries, ConfigMapSummary{
			Name:      cm.Name,
			Namespace: cm.Namespace,
			Keys:      len(cm.Data) + len(cm.BinaryData),
			Age:       ageStr(cm.CreationTimestamp.Time),
		})
	}
	return ConfigMapListRsp{ConfigMaps: summaries}, nil
}

// ConfigMapGetReq is the input for describing a ConfigMap.
type ConfigMapGetReq struct {
	Name      string `json:"name" jsonschema:"description=configmap name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// ConfigMapGetRsp is the output.
type ConfigMapGetRsp struct {
	ConfigMap *ConfigMapDetail `json:"configmap"`
	Err       string           `json:"error,omitempty"`
}

// ConfigMapDetail contains ConfigMap information.
type ConfigMapDetail struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Keys       []string          `json:"keys"`
	Data       map[string]string `json:"data,omitempty"`
	BinaryKeys []string          `json:"binary_keys,omitempty"`
	Labels     map[string]string `json:"labels"`
}

// NewConfigMapGetTool creates a tool for describing a ConfigMap.
func NewConfigMapGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req ConfigMapGetReq) (ConfigMapGetRsp, error) {
			return configMapGet(ctx, clients, req)
		},
		function.WithName("configmap_get"),
		function.WithDescription("Get detailed information about a ConfigMap. Data values over 256 characters are truncated."),
	)
}

func configMapGet(ctx context.Context, c *k8s.Clients, req ConfigMapGetReq) (ConfigMapGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	cm, err := c.ClientSet.CoreV1().ConfigMaps(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return ConfigMapGetRsp{Err: err.Error()}, nil
	}

	keys := make([]string, 0, len(cm.Data)+len(cm.BinaryData))
	truncated := make(map[string]string, len(cm.Data))
	binaryKeys := make([]string, 0, len(cm.BinaryData))

	for k, v := range cm.Data {
		keys = append(keys, k)
		if len(v) > 256 {
			truncated[k] = v[:256] + "...(truncated)"
		} else {
			truncated[k] = v
		}
	}
	for k := range cm.BinaryData {
		keys = append(keys, k)
		binaryKeys = append(binaryKeys, k)
	}

	return ConfigMapGetRsp{ConfigMap: &ConfigMapDetail{
		Name:       cm.Name,
		Namespace:  cm.Namespace,
		Keys:       keys,
		Data:       truncated,
		BinaryKeys: binaryKeys,
		Labels:     cm.Labels,
	}}, nil
}

// SecretListReq is the input for listing Secrets.
type SecretListReq struct {
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
	Label     string `json:"label" jsonschema:"description=label selector,omitempty"`
}

// SecretListRsp is the output.
type SecretListRsp struct {
	Secrets []SecretSummary `json:"secrets"`
	Err     string          `json:"error,omitempty"`
}

// SecretSummary is a simplified Secret view.
type SecretSummary struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Keys      int    `json:"keys"`
	Age       string `json:"age"`
}

// NewSecretListTool creates a tool for listing Secrets.
func NewSecretListTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req SecretListReq) (SecretListRsp, error) {
			return secretList(ctx, clients, req)
		},
		function.WithName("secret_list"),
		function.WithDescription("List Secrets with type and key count. Secret values are never exposed."),
	)
}

func secretList(ctx context.Context, c *k8s.Clients, req SecretListReq) (SecretListRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	opts := metav1.ListOptions{}
	if req.Label != "" {
		opts.LabelSelector = req.Label
	}
	list, err := c.ClientSet.CoreV1().Secrets(ns).List(ctx, opts)
	if err != nil {
		return SecretListRsp{Err: err.Error()}, nil
	}
	summaries := make([]SecretSummary, 0, len(list.Items))
	for _, s := range list.Items {
		summaries = append(summaries, SecretSummary{
			Name:      s.Name,
			Namespace: s.Namespace,
			Type:      string(s.Type),
			Keys:      len(s.Data),
			Age:       ageStr(s.CreationTimestamp.Time),
		})
	}
	return SecretListRsp{Secrets: summaries}, nil
}

// SecretGetReq is the input for describing a Secret.
type SecretGetReq struct {
	Name      string `json:"name" jsonschema:"description=secret name,required"`
	Namespace string `json:"namespace" jsonschema:"description=namespace,omitempty"`
}

// SecretGetRsp is the output.
type SecretGetRsp struct {
	Secret *SecretDetail `json:"secret"`
	Err    string        `json:"error,omitempty"`
}

// SecretDetail contains Secret metadata (never exposes values).
type SecretDetail struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`
	Keys      []string          `json:"keys"`
	Labels    map[string]string `json:"labels"`
}

// NewSecretGetTool creates a tool for describing a Secret.
func NewSecretGetTool(clients *k8s.Clients) tool.Tool {
	return function.NewFunctionTool(
		func(ctx context.Context, req SecretGetReq) (SecretGetRsp, error) {
			return secretGet(ctx, clients, req)
		},
		function.WithName("secret_get"),
		function.WithDescription("Get metadata about a Secret including type and key names. Secret values are never exposed."),
	)
}

func secretGet(ctx context.Context, c *k8s.Clients, req SecretGetReq) (SecretGetRsp, error) {
	ns := req.Namespace
	if ns == "" {
		ns = c.Namespace
	}
	s, err := c.ClientSet.CoreV1().Secrets(ns).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return SecretGetRsp{Err: err.Error()}, nil
	}

	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}

	return SecretGetRsp{Secret: &SecretDetail{
		Name:      s.Name,
		Namespace: s.Namespace,
		Type:      string(s.Type),
		Keys:      keys,
		Labels:    s.Labels,
	}}, nil
}
