# kube-agents

AI-powered Kubernetes operations assistant built on [tRPC-Agent-Go](https://github.com/trpc-group/trpc-agent-go) and DeepSeek.

Users describe K8s resource operations in natural language — the agent executes them via the Kubernetes API.

```
> List all failing pods in production namespace

I found 3 pods not in Running state:

1. payment-svc-7d4f8b9c-x2k  → CrashLoopBackOff (restarts: 12)
2. cache-worker-5b6c9f-4mz    → Pending (waiting for CPU)
3. api-gateway-8a3d2e-1qk    → ImagePullBackOff

Let me investigate further...
```

## Features

- **Natural Language K8s Operations** — list, describe, create, delete pods, deployments, services, and more
- **Multi-Model Support** — DeepSeek (default), OpenAI-compatible models
- **Authentication & Authorization** — K8s ServiceAccount TokenReview + SubjectAccessReview RBAC
- **Agent-to-Agent (A2A)** — inter-agent communication protocol for multi-agent collaboration
- **HTTP API** — OpenAI-compatible `/v1/chat/completions` endpoint with SSE streaming
- **gRPC API** — streaming chat and tool listing (proto defined)
- **CLI Interface** — `chat`, `serve`, `version` commands via cobra
- **Plugin System** — built-in audit, rate-limiting, and logging plugins; extensible via `plugin.Plugin` interface
- **MCP Tools** — Model Context Protocol integration for external tool servers (stdio/sse/streamable)
- **Skills** — reusable SKILL.md-based workflow modules (diagnose, deploy, security audit)
- **RAG Knowledge** — K8s documentation retrieval for informed responses
- **Session & Memory** — persistent conversation state and long-term user memory
- **Graceful Shutdown** — HTTP server drains in-flight requests before exit (configurable timeout)
- **Structured Logging** — `log/slog` with JSON/text formats, log levels, file rotation via lumberjack

## Quick Start

### Prerequisites

- Go 1.21+
- DeepSeek API key (or any OpenAI-compatible provider)
- Kubernetes cluster (optional — K8s tools auto-disable if unreachable)

### Install

```bash
git clone https://github.com/codefuture-io/kube-agents.git
cd kube-agents
make build
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DEEPSEEK_API_KEY` | Yes | DeepSeek API key |
| `OPENAI_API_KEY` | No | Fallback if `DEEPSEEK_API_KEY` not set |
| `OPENAI_BASE_URL` | No | Custom API base URL (default: `https://api.deepseek.com`) |
| `KUBECONFIG` | No | Path to kubeconfig (default: `~/.kube/config`) |

### Run

```bash
# Interactive chat
./bin/kube-agents chat --api-key=sk-your-key

# Start HTTP API server
./bin/kube-agents serve --api-key=sk-your-key

# With config file + plugins
./bin/kube-agents serve --config=config/kube-agents.yaml --api-key=sk-your-key

# Print version
./bin/kube-agents version
```

## CLI Reference

### Commands

| Command | Description |
|---------|-------------|
| `kube-agents chat` | Interactive terminal chat session |
| `kube-agents serve` | Start HTTP + gRPC + A2A API servers |
| `kube-agents version` | Print version table |

### Flags

| Flag | Default | Env Fallback | Description |
|------|---------|-------------|-------------|
| `--api-key` | — | `DEEPSEEK_API_KEY` → `OPENAI_API_KEY` | API key |
| `--base-url` | — | `OPENAI_BASE_URL` | API base URL |
| `--model` | `deepseek-chat` | — | Model name |
| `--config` | `config/kube-agents.yaml` | — | Config file path |
| `--log-level` | `info` | — | Log level: debug, info, warn, error |
| `--log-format` | `text` | — | Log format: text, json |
| `--log-add-source` | `false` | — | Include source file and line in log output |
| `--log-file-output` | `false` | — | Write logs to file instead of stderr |
| `--log-file-path` | — | — | Log file path |
| `--log-max-size` | `100` | — | Max log file size in MB before rotation |
| `--log-max-backups` | `10` | — | Max old log files to retain |
| `--log-max-age` | `30` | — | Max days to retain old log files |

### Logging

All server and plugin output uses Go's native `log/slog` structured logging, with optional JSON format and file output with rotation (via `lumberjack`).

```bash
# JSON format with debug level
./bin/kube-agents serve --api-key=sk-xxx --log-format=json --log-level=debug

# File output with log rotation
./bin/kube-agents serve --api-key=sk-xxx \
  --log-file-output --log-file-path=/var/log/kube-agents.log \
  --log-max-size=500 --log-max-backups=30
```

## API

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/chat/completions` | OpenAI-compatible chat (SSE streaming) |

### Non-Streaming

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-chat",
    "messages": [{"role": "user", "content": "List all namespaces"}],
    "stream": false
  }'
```

### Streaming (SSE)

```bash
curl -N -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-chat",
    "messages": [{"role": "user", "content": "Get logs from pod payment-svc"}],
    "stream": true
  }'
```

### Response

```json
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion",
  "model": "gpt-3.5-turbo",
  "choices": [{
    "index": 0,
    "message": {"role": "assistant", "content": "..."},
    "finish_reason": "stop"
  }],
  "usage": {"prompt_tokens": 3019, "completion_tokens": 47, "total_tokens": 3066}
}
```

## K8s Tools

When connected to a Kubernetes cluster, the agent has access to these tools:

| Tool | Resource | Operations |
|------|----------|------------|
| `pod_list` | Pods | List with label selector, namespace filter |
| `pod_get` | Pods | Detail: containers, conditions, events, labels |
| `pod_logs` | Pods | Fetch recent logs (configurable tail + container) |
| `pod_delete` | Pods | Delete by name |
| `deployment_list` | Deployments | List with ready/available status |
| `deployment_get` | Deployments | Detail: replicas, image, strategy, conditions |
| `deployment_scale` | Deployments | Scale to N replicas |
| `service_list` | Services | List with type, cluster IP, ports |
| `service_get` | Services | Detail: ports, selector, external IP |
| `namespace_list` | Namespaces | List all namespaces |
| `namespace_get` | Namespaces | Detail: labels |
| `namespace_set` | Namespaces | Change active namespace |
| `event_list` | Events | Filter by namespace and type (Normal/Warning) |
| `ingress_list` | Ingress | List with hosts and ingress class |
| `ingress_get` | Ingress | Detail: hosts, TLS, routing rules, backend services |
| `hpa_list` | HPA | List with target, min/max replicas, current CPU |
| `hpa_get` | HPA | Detail: metrics, current load, conditions |
| `configmap_list` | ConfigMap | List with key count |
| `configmap_get` | ConfigMap | Detail: keys, data content (>256 chars truncated) |
| `secret_list` | Secret | List with type and key count (values never exposed) |
| `secret_get` | Secret | Metadata: type, key names, labels (values never exposed) |
| `resource_list` | Generic | Any K8s/CRD resource via dynamic client |
| `resource_get` | Generic | Get specific resource by name + type |
| `cluster_info` | Cluster | Version, API groups, node/namespace count |

## Skills

Reusable workflow skills located in `skills/`:

| Skill | Description |
|-------|-------------|
| `k8s-diagnose` | Diagnose pod crashes, resource bottlenecks, scheduling failures |
| `k8s-deploy` | Manage rolling updates, rollbacks, scaling, health checks |
| `k8s-security` | Audit RBAC, pod security contexts, network policies, secrets |

## Plugins

Built-in plugins registered by name in the config:

| Plugin | Hook Points | Description |
|--------|------------|-------------|
| `audit` | `AfterTool` | Logs all tool invocations with timestamps |
| `ratelimit` | `BeforeAgent` | Per-session request rate limiting (100 req/session) |
| `logging` | `BeforeAgent`, `AfterAgent` | Agent lifecycle logging |

Custom plugins implement the `plugin.Plugin` interface and register via `runner.WithPlugins()`.

## Configuration

Full reference — `config/kube-agents.yaml`:

```yaml
server:
  http:                         # OpenAI-compatible HTTP API
    enabled: true
    port: 8080
  grpc:                         # gRPC API (proto: api/proto/)
    enabled: false
    port: 9090
  a2a:                          # Agent-to-Agent protocol
    enabled: false
    host: "kube-agents.default.svc.cluster.local:8080"

model:
  provider: deepseek
  name: deepseek-chat
  apiKeyEnv: DEEPSEEK_API_KEY
  baseUrl: ""                   # Override API base URL

auth:
  mode: serviceaccount          # serviceaccount | jwt
  jwt:
    issuer: kube-agents
    secretEnv: JWT_SECRET

session:
  backend: memory               # memory | redis
  redis:
    addr: "redis:6379"
    passwordEnv: REDIS_PASSWORD
    db: 0

memory:
  backend: memory               # memory | redis
  redis:
    addr: "redis:6379"
    passwordEnv: REDIS_PASSWORD
    db: 1

knowledge:                      # RAG document sources
  sources:
    - type: file
      path: ./docs/k8s-reference/
    - type: url
      path: https://kubernetes.io/docs/

mcpServers:                     # External MCP tool servers
  - name: kubectl
    transport: stdio
    command: kubectl
    args: ["mcp"]

plugins:                        # Built-in plugin names
  - audit
  - ratelimit
  - logging

skillsDir: ./skills/            # SKILL.md search path
```

## Architecture

```
HTTP/gRPC/CLI/A2A
       │
       ▼
┌──────────────────┐
│   Auth Plugin    │  ← TokenReview + SubjectAccessReview
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│     Runner       │  ← Session / Memory / Plugin
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  K8s LLMAgent    │  ← System prompt + 16 K8s tools + memory + knowledge
└──────┬───────────┘
       │
       ├──→ K8s tools (pod, deploy, svc, event, …)
       ├──→ MCP tools (external tool servers)
       ├──→ Skills (k8s-diagnose, k8s-deploy, k8s-security)
       └──→ Knowledge search (RAG)
                  │
                  ▼
           client-go → K8s API Server
```

## Deployment

### Docker

```bash
make docker-build IMG=your-registry/kube-agents:0.1.0
make docker-push  IMG=your-registry/kube-agents:0.1.0
```

Multi-arch builds:

```bash
make docker-buildx IMG=your-registry/kube-agents:0.1.0
```

### Kubernetes

```bash
kubectl apply -f deployments/
```

This creates:
- `ServiceAccount` + `ClusterRole` + `ClusterRoleBinding`
- `ConfigMap` (config + skills)
- `Deployment` (single replica, 128Mi-512Mi memory, distroless image)
- `Service` (port 8080)

### RBAC

The ServiceAccount requires these permissions (see `deployments/rbac.yaml`):

| API Group | Resources | Verbs |
|-----------|-----------|-------|
| `""` (core) | pods, services, namespaces, events, nodes | get, list, watch |
| `""` (core) | pods, services | create, update, delete |
| `apps` | deployments, deployments/scale, replicasets | get, list, watch, update |
| `authentication.k8s.io` | tokenreviews | create |
| `authorization.k8s.io` | subjectaccessreviews | create |
| `networking.k8s.io`, `batch`, `apiextensions.k8s.io` | `*` | get, list |

## Development

```bash
make fmt          # Format code
make vet          # Static analysis
make build        # Build binary → bin/kube-agents
make run          # Run from source
make test         # Run tests with race detection
make test-cover   # Run tests with coverage report
```

## License

Apache License 2.0 — see LICENSE
