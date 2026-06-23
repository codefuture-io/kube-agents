---
name: k8s-deploy
description: Manage Kubernetes deployments: rolling updates, rollbacks, scaling, and health verification.
---

# K8s Deploy

## Workflow

1. Use `deployment_get` to check current state and image version
2. Use `deployment_scale` to adjust replicas as needed
3. Use `pod_list` to verify pod rollout status after changes
4. Use `event_list` to watch for deployment-related events

## Safety Rules

- Always verify current state before scaling or updating
- Check available replicas before initiating a change
- Verify pod health after scaling operations
- Never scale to zero without explicit user confirmation
