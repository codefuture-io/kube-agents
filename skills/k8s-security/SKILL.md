---
name: k8s-security
description: Security auditing for Kubernetes: RBAC review, pod security contexts, network policy validation, and secret management.
---

# K8s Security

## Workflow

1. Use `resource_list` to list RBAC roles and bindings
2. Use `pod_get` to check pod security contexts
3. Use `service_list` to identify externally exposed services
4. Use `namespace_get` to review namespace configurations

## Check Items

- RBAC: Identify overly permissive ClusterRoleBindings
- Pods: Flag privileged containers, hostPID/hostNetwork usage
- Services: Identify LoadBalancer services with broad exposure
- Secrets: Verify secrets are not exposed in pod environment variables
