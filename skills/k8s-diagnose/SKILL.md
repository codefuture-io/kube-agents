---
name: k8s-diagnose
description: Diagnose common Kubernetes issues: pod crashes, resource bottlenecks, scheduling failures, and networking problems.
---

# K8s Diagnose

## Workflow

1. Use `pod_list` or `pod_get` to identify failing pods
2. Use `pod_logs` to fetch recent logs from crashed containers
3. Use `event_list` filtered by Warning type to find relevant events
4. Use `resource_get` to check node conditions and resource capacity
5. Summarize findings and suggest remediation steps

## Common Patterns

- Crashing pod → check pod_logs for application errors, verify image pull
- Pending pod → check event_list for scheduling failures (resources, taints)
- OOMKilled → check resource limits, suggest increasing memory
- ImagePullBackOff → verify image name, registry auth, network access
