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

// Package k8s provides Kubernetes client utilities and tools.
package k8s

import (
	"fmt"
	"os"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Clients bundles the various K8s client types.
type Clients struct {
	ClientSet     kubernetes.Interface
	DynamicClient dynamic.Interface
	Discovery     discovery.DiscoveryInterface
	Config        *rest.Config
	Namespace     string
}

// NewClients creates K8s clients. Tries in-cluster config first,
// falls back to kubeconfig (KUBECONFIG env or ~/.kube/config).
func NewClients() (*Clients, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath())
		if err != nil {
			return nil, fmt.Errorf("failed to create K8s config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	namespace := resolveNamespace()

	return &Clients{
		ClientSet:     clientset,
		DynamicClient: dynClient,
		Discovery:     discoveryClient,
		Config:        config,
		Namespace:     namespace,
	}, nil
}

func kubeconfigPath() string {
	if v := os.Getenv("KUBECONFIG"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home + "/.kube/config"
}

func resolveNamespace() string {
	if v := os.Getenv("KUBE_AGENTS_NAMESPACE"); v != "" {
		return v
	}
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		return string(data)
	}
	return "default"
}
