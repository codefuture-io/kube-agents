# tRPC-Agent-Go —— Go 语言多智能体 AI 开发框架

> **tRPC-Agent-Go** 是腾讯 tRPC 团队开源的 Go 语言智能体（Agent）开发框架，构建在 tRPC 生态之上，提供去中心化、可协作、能涌现的多 Agent 能力。支持 LLM 推理、工具调用、记忆持久化、知识检索、图工作流编排、人机协作和全链路可观测性，专为生产级 AI 应用设计。

- **GitHub**: [github.com/trpc-group/trpc-agent-go](https://github.com/trpc-group/trpc-agent-go)
- **官方文档**: [trpc-group.github.io/trpc-agent-go/zh/](https://trpc-group.github.io/trpc-agent-go/zh/)
- **API 参考**: [pkg.go.dev/trpc.group/trpc-go/trpc-agent-go](https://pkg.go.dev/trpc.group/trpc-go/trpc-agent-go)
- **License**: Apache 2.0 | **最新版本**: v1.9.1 | **Go**: ≥ 1.21
- **生态配套**: [tRPC-A2A-Go](https://github.com/trpc-group/trpc-a2a-go) | [tRPC-MCP-Go](https://github.com/trpc-group/trpc-mcp-go)

---

## 目录

- [一、入门篇](#一入门篇)
  - [1.1 项目背景与定位](#11-项目背景与定位)
  - [1.2 架构设计](#12-架构设计)
  - [1.3 核心设计理念](#13-核心设计理念)
  - [1.4 快速开始](#14-快速开始)
  - [1.5 示例全景](#15-示例全景)
- [二、Agent 篇](#二agent-篇)
  - [2.1 Agent 概述](#21-agent-概述)
  - [2.2 LLMAgent —— 大模型智能体](#22-llmagent--大模型智能体)
  - [2.3 消息可见性控制](#23-消息可见性控制)
  - [2.4 推理内容模式](#24-推理内容模式)
  - [2.5 Agent 动态创建](#25-agent-动态创建)
  - [2.6 Agent 作为 Tool](#26-agent-作为-tool)
  - [2.7 运行时取消](#27-运行时取消)
- [三、多 Agent 协作篇](#三多-agent-协作篇)
  - [3.1 多 Agent 协作概述](#31-多-agent-协作概述)
  - [3.2 ChainAgent —— 链式顺序执行](#32-chainagent--链式顺序执行)
  - [3.3 ParallelAgent —— 并行执行与结果合并](#33-parallelagent--并行执行与结果合并)
  - [3.4 CycleAgent —— 循环迭代执行](#34-cycleagent--循环迭代执行)
  - [3.5 组合嵌套模式](#35-组合嵌套模式)
- [四、Graph 工作流篇](#四graph-工作流篇)
  - [4.1 GraphAgent 概述](#41-graphagent-概述)
  - [4.2 核心概念：Graph / Node / State / Edge](#42-核心概念graph--node--state--edge)
  - [4.3 BSP 与 DAG 执行引擎](#43-bsp-与-dag-执行引擎)
  - [4.4 条件路由与多条件并行路由](#44-条件路由与多条件并行路由)
  - [4.5 CheckPoint 与时间旅行](#45-checkpoint-与时间旅行)
  - [4.6 Human-in-the-Loop（HITL）](#46-human-in-the-loophitl)
  - [4.7 子 Agent 节点](#47-子-agent-节点)
  - [4.8 GraphAgent 完整示例](#48-graphagent-完整示例)
- [五、模型与工具篇](#五模型与工具篇)
  - [5.1 Model —— 多模型支持](#51-model--多模型支持)
  - [5.2 Function Tool —— 函数工具](#52-function-tool--函数工具)
  - [5.3 MCP Tool —— 模型上下文协议](#53-mcp-tool--模型上下文协议)
  - [5.4 内置工具集](#54-内置工具集)
  - [5.5 工具组合与调用策略](#55-工具组合与调用策略)
- [六、记忆与会话篇](#六记忆与会话篇)
  - [6.1 Session —— 会话管理](#61-session--会话管理)
  - [6.2 Memory —— 长期记忆](#62-memory--长期记忆)
  - [6.3 Memory Tools 集成](#63-memory-tools-集成)
  - [6.4 占位符变量 —— 会话状态注入](#64-占位符变量--会话状态注入)
- [七、知识检索篇（RAG）](#七知识检索篇rag)
  - [7.1 Knowledge 系统概述](#71-knowledge-系统概述)
  - [7.2 基本使用流程](#72-基本使用流程)
  - [7.3 智能过滤与多知识库](#73-智能过滤与多知识库)
  - [7.4 存储后端与 Embedder](#74-存储后端与-embedder)
- [八、技能与制品篇](#八技能与制品篇)
  - [8.1 Agent Skills —— 可复用工作流](#81-agent-skills--可复用工作流)
  - [8.2 技能发现与动态安装](#82-技能发现与动态安装)
  - [8.3 Artifacts —— 版本化文件存储](#83-artifacts--版本化文件存储)
- [九、可观测性与评估篇](#九可观测性与评估篇)
  - [9.1 OpenTelemetry 全链路追踪](#91-opentelemetry-全链路追踪)
  - [9.2 Callback 机制](#92-callback-机制)
  - [9.3 Evaluation —— 评估与基准测试](#93-evaluation--评估与基准测试)
- [十、服务与集成篇](#十服务与集成篇)
  - [10.1 AG-UI —— Agent-用户交互协议](#101-ag-ui--agent-用户交互协议)
  - [10.2 A2A —— Agent-to-Agent 互操作](#102-a2a--agent-to-agent-互操作)
  - [10.3 Gateway Server](#103-gateway-server)
- [十一、进阶篇](#十一进阶篇)
  - [11.1 Runner 深入理解](#111-runner-深入理解)
  - [11.2 Planner —— 规划与推理](#112-planner--规划与推理)
  - [11.3 代码执行](#113-代码执行)
  - [11.4 Prompt Caching](#114-prompt-caching)
  - [11.5 生产部署建议](#115-生产部署建议)
  - [11.6 生态与社区](#116-生态与社区)

---

# 一、入门篇

## 1.1 项目背景与定位

当前主流 Agent 框架（AutoGen、CrewAI、Agno、ADK 等）多数基于 Python，而 Go 语言在微服务、高并发与部署方面有天然优势。tRPC-Agent-Go 填补了 Go 生态中 **"去中心化、可协作、能涌现"的自主多 Agent 框架** 的空白。

框架直接利用 Go 的高并发能力与 tRPC 生态，将 LLM 的推理、协商和自适应性引入 Go 场景，满足复杂业务对 **"智能 + 性能"** 的双重需求。

**生产验证**：已在腾讯元宝、腾讯视频、腾讯新闻、IMA、QQ 音乐等业务中落地打磨。

## 1.2 架构设计

tRPC-Agent-Go 采用模块化架构，所有组件都可插拔，通过事件驱动机制实现组件间解耦通信：

```
┌──────────────────────────────────────────────────────────┐
│                        Runner                             │
│        (执行器：串联 Session/Memory/Planner 等能力)        │
├──────────────────────────────────────────────────────────┤
│   Agent（核心执行单元）    Planner（规划与推理）             │
│   LLMAgent · ChainAgent · ParallelAgent · CycleAgent     │
│   GraphAgent（图工作流）· 自定义 Agent                    │
├──────────────────────────────────────────────────────────┤
│   Model           Tool            Knowledge              │
│   OpenAI/DeepSeek Function/MCP    RAG 检索增强            │
├──────────────────────────────────────────────────────────┤
│   Session         Memory          Skill       Artifact   │
│   会话状态管理     长期记忆        可复用技能   版本化文件   │
├──────────────────────────────────────────────────────────┤
│   Telemetry       Evaluation      Server                 │
│   全链路追踪       评估基准        AG-UI / A2A / Gateway  │
└──────────────────────────────────────────────────────────┘
```

**核心包职责**：

| 包 | 职责 |
|---|---|
| `agent` | 核心执行单元，处理用户输入并生成响应 |
| `runner` | Agent 执行器，管理执行流程，串联 Session/Memory Service |
| `model` | LLM 模型抽象，支持 OpenAI、DeepSeek 等 |
| `tool` | 工具抽象，支持 Function Tool、MCP Tool、DuckDuckGo 等 |
| `session` | 会话状态管理与事件存储 |
| `memory` | 长期记忆与个性化信息记录 |
| `knowledge` | RAG 知识检索能力 |
| `planner` | Agent 规划与推理能力 |
| `artifact` | 版本化的文件存储与检索（图片、报告等） |
| `skill` | 基于 SKILL.md 的可复用 Agent 技能 |
| `graph` | 图工作流引擎（StateGraph、条件路由、BSP/DAG） |
| `evaluation` | Agent 评估框架 |
| `server` | HTTP 服务（AG-UI、A2A、Gateway） |
| `telemetry` | OpenTelemetry 追踪与指标 |

## 1.3 核心设计理念

| 原则 | 说明 |
|------|------|
| **模块化可插拔** | 每个组件独立抽象，可按需替换实现 |
| **事件驱动解耦** | 组件间通过事件流通信，支持 Callback 注入 |
| **多 Agent 涌现** | 不止于编排式工作流，支持真正的去中心化多 Agent 协作 |
| **类型安全** | Graph State Schema + Go 泛型确保编译时安全 |
| **生产就绪** | 内建追踪、会话持久化、评估体系，开箱即用 |

## 1.4 快速开始

### 安装

```bash
go get trpc.group/trpc-go/trpc-agent-go
```

依赖 Go 1.21 及以上版本。

### 最简 LLMAgent

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
    "trpc.group/trpc-go/trpc-agent-go/tool"
    "trpc.group/trpc-go/trpc-agent-go/tool/function"
)

func main() {
    // 创建模型
    modelInstance := openai.New("deepseek-chat")

    // 创建工具
    calculatorTool := function.NewFunctionTool(
        calculator,
        function.WithName("calculator"),
        function.WithDescription("执行加减乘除。参数：a、b 为数值，op 取值 add/sub/mul/div"),
    )

    // 启用流式输出
    genConfig := model.GenerationConfig{Stream: true}

    // 创建 Agent
    agent := llmagent.New("assistant",
        llmagent.WithModel(modelInstance),
        llmagent.WithTools([]tool.Tool{calculatorTool}),
        llmagent.WithGenerationConfig(genConfig),
    )

    // 创建 Runner 并执行
    r := runner.NewRunner("calculator-app", agent)
    ctx := context.Background()
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("计算 2+3 等于多少"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 处理流式事件
    for event := range events {
        if event.Object == "chat.completion.chunk" {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
    }
}
```

### 多 Agent 链式协作

```go
// 创建分析 Agent
analysisAgent := llmagent.New("analyzer",
    llmagent.WithModel(modelInstance),
    llmagent.WithInstruction("分析给定问题并拆解为子任务"),
)

// 创建执行 Agent
executionAgent := llmagent.New("executor",
    llmagent.WithModel(modelInstance),
    llmagent.WithTools([]tool.Tool{calculatorTool, searchTool}),
)

// 组合为链式 Agent
chain := chainagent.New("problem-solver",
    chainagent.WithSubAgents([]agent.Agent{analysisAgent, executionAgent}),
)

// 运行
r := runner.NewRunner("multi-agent-app", chain)
events, _ := r.Run(ctx, "user-001", "session-001", message)
```

### 运行示例

```bash
git clone https://github.com/trpc-group/trpc-agent-go.git
cd trpc-agent-go
export OPENAI_API_KEY="your-api-key"
cd examples/runner
go run . -model="gpt-4o-mini" -streaming=true
```

## 1.5 示例全景

`examples/` 目录覆盖了框架全部核心能力：

| 类别 | 示例目录 | 说明 |
|------|---------|------|
| **Agent 基础** | `examples/llmagent` | LLMAgent 完整配置与流式输出 |
| **多 Agent** | `examples/multiagent` | ChainAgent / ParallelAgent / CycleAgent |
| **GraphAgent** | `examples/graph` | 图工作流、条件路由、状态管理、HITL |
| **工具** | `examples/multitools` | 多工具编排 |
| | `examples/duckduckgo` | 网络搜索工具 |
| | `examples/filetoolset` | 文件系统工具集 |
| | `examples/agenttool` | Agent 作为可调用工具 |
| **记忆** | `examples/memory` | 内存/Redis 记忆服务 |
| **知识** | `examples/knowledge` | RAG 检索、向量存储 |
| **MCP** | `examples/mcptool` | MCP 协议工具集成 |
| **技能** | `examples/skillrun` | Agent Skills (SKILL.md) |
| | `examples/skillfind` | 技能发现与动态安装 |
| **可观测** | `examples/telemetry` | OpenTelemetry 全链路追踪 |
| **评估** | `examples/evaluation` | Agent 评估基准 |
| **HITL** | `examples/humaninloop` | 人机协作 |
| **代码执行** | `examples/codeexecution` | 安全代码执行沙箱 |
| **制品** | `examples/artifact` | 版本化文件存储 |
| **UI** | `examples/agui` | AG-UI SSE 服务（CopilotKit 客户端） |
| **A2A** | `examples/a2aadk` | Agent-to-Agent 跨运行时互操作 |
| **Gateway** | `openclaw` | 类 OpenClaw 网关服务 |

---

# 二、Agent 篇

## 2.1 Agent 概述

Agent 是 tRPC-Agent-Go 的核心执行单元，每个 Agent 实现统一接口：

```go
type Agent interface {
    Run(ctx context.Context, invocation *Invocation) (<-chan *event.Event, error)
    Tools() []tool.Tool
    Info() Info
    SubAgents() []Agent
    FindSubAgent(name string) Agent
}
```

框架内置五种 Agent 类型：

| Agent | 特点 | 适用场景 |
|-------|------|---------|
| **LLMAgent** | 大模型智能体，支持工具调用和推理 | 对话助手、工具调用 |
| **ChainAgent** | 链式顺序执行子 Agent | 流水线任务 |
| **ParallelAgent** | 并行执行子 Agent，合并结果 | 多专家并发分析 |
| **CycleAgent** | 循环迭代直到终止条件 | 自我优化、迭代改进 |
| **GraphAgent** | 图工作流，条件路由，状态管理 | 复杂业务流程 |

## 2.2 LLMAgent —— 大模型智能体

LLMAgent 是最常用的 Agent 类型，封装了 LLM 调用、工具使用和推理循环。

### 基础配置

```go
llmAgent := llmagent.New("agent-name",
    // 模型（必填）
    llmagent.WithModel(modelInstance),

    // 描述与指令
    llmagent.WithDescription("Agent 功能描述"),
    llmagent.WithInstruction("You are a helpful assistant. Focus on {topic}."),

    // 生成参数
    llmagent.WithGenerationConfig(model.GenerationConfig{
        MaxTokens:   ptrOf(2000),
        Temperature: ptrOf(0.7),
        Stream:      true,
    }),

    // 工具
    llmagent.WithTools([]tool.Tool{calculatorTool, searchTool}),
)
```

### 全部配置选项

| 选项 | 说明 |
|------|------|
| `WithModel(m)` | 设置 LLM 模型 |
| `WithDescription(d)` | Agent 能力描述（作为 SubAgent 时影响父 Agent 的调用决策） |
| `WithInstruction(i)` | System Prompt，支持 `{key}` 占比符变量 |
| `WithGenerationConfig(c)` | 生成参数（MaxTokens、Temperature、Stream、Stop 等） |
| `WithTools(t)` | 工具列表 |
| `WithToolSets(ts)` | 工具集（批量注册一组工具） |
| `WithSubAgents(a)` | 子 Agent 列表 |
| `WithPlanner(p)` | 规划器 |
| `WithMemory(m)` | 记忆服务 |
| `WithKnowledge(k)` | 知识库（自动添加 `knowledge_search` 工具） |
| `WithCodeExecutor(e)` | 代码执行器 |
| `WithOutputKey(k)` | 将 Agent 输出写入 Session 的键名 |
| `WithMessageFilterMode(m)` | 消息可见性控制模式 |

## 2.3 消息可见性控制

在多 Agent 场景中，控制每个 Agent 能看到哪些历史消息至关重要。

**简化模式**（推荐）：

```go
llmagent.WithMessageFilterMode(llmagent.FullContext)        // 全部消息可见
llmagent.WithMessageFilterMode(llmagent.RequestContext)     // 仅当前请求内消息
llmagent.WithMessageFilterMode(llmagent.IsolatedRequest)    // 仅当前 Agent 当前请求
llmagent.WithMessageFilterMode(llmagent.IsolatedInvocation) // 仅当前 invocation
```

**精细控制模式**：

```go
llmAgent := llmagent.New("agent",
    // 时间维度过滤
    llmagent.WithMessageTimelineFilterMode(llmagent.TimelineFilterCurrentRequest),
    // 分支维度过滤
    llmagent.WithMessageBranchFilterMode(llmagent.BranchFilterModePrefix),
)
```

| 时间维度 | 可见范围 |
|---------|---------|
| `TimelineFilterAll` | 包含全部历史消息 |
| `TimelineFilterCurrentRequest` | 仅包含当前 `runner.Run` 内的消息 |
| `TimelineFilterCurrentInvocation` | 仅包含当前 invocation 内的消息 |

| 分支维度 | 匹配规则 |
|---------|---------|
| `BranchFilterModePrefix` | 层级匹配（祖先/自己/子孙都算） — **默认值** |
| `BranchFilterModeSubtree` | 仅当前 key 及其子孙（不含父级） |
| `BranchFilterModeExact` | `Event.FilterKey == Invocation.eventFilterKey` |
| `BranchFilterModeAll` | 忽略 FilterKey，包含全部消息 |

最终传给模型的消息需同时满足时间维度与分支维度两个条件。

## 2.4 推理内容模式

支持 DeepSeek 等模型的思考过程展示：

```go
llmagent.WithReasoningContentMode(llmagent.ReasoningContentModeKeepAll)
// ReasoningContentModeKeepAll     - 保留所有思考内容
// ReasoningContentModeDiscardAll  - 丢弃所有思考内容
```

## 2.5 Agent 动态创建

按请求动态创建 Agent（不同 Prompt、模型、工具）：

```go
r := runner.NewRunnerWithAgentFactory("app", "agent-name",
    func(ctx context.Context, ro agent.RunOptions) (agent.Agent, error) {
        a := llmagent.New("assistant",
            llmagent.WithInstruction(ro.Instruction),
            // 根据 ro 动态配置
        )
        return a, nil
    },
)

events, _ := r.Run(ctx, "user-001", "session-001",
    model.NewUserMessage("Hello"),
    agent.WithInstruction("You are a helpful assistant."),
)
```

## 2.6 Agent 作为 Tool

将一个 Agent 包装为可被其他 Agent 调用的 Tool：

```go
// 创建专业 Agent
researchAgent := llmagent.New("researcher", ...)

// 包装为 Tool
researchTool := agenttool.NewTool(researchAgent,
    agenttool.WithSkipStreaming(false), // 支持流式
)

// 父 Agent 使用
parentAgent := llmagent.New("coordinator",
    llmagent.WithTools([]tool.Tool{researchTool}),
)
```

## 2.7 运行时取消

框架提供三种取消方式：

**方式一：Context 取消（推荐）**

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

events, _ := r.Run(ctx, "user-001", "session-001", message)
go func() {
    time.Sleep(2 * time.Second)
    cancel()  // 2 秒后取消
}()
for range events {} // 持续消费直到通道关闭
```

**方式二：Ctrl+C（命令行程序）**

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
events, _ := r.Run(ctx, "user-001", "session-001", message)
```

**方式三：按 RequestID 取消（服务端场景）**

```go
events, _ := r.Run(ctx, "user-001", "session-001", message,
    agent.WithRequestID("req-123"),
)
mr := r.(runner.ManagedRunner)
mr.Cancel("req-123")
```

重要：取消后必须持续消费事件通道直到关闭，否则 Agent 协程可能阻塞在通道写入上。

---

# 三、多 Agent 协作篇

## 3.1 多 Agent 协作概述

tRPC-Agent-Go 提供三种内置多 Agent 协作模式，可自由嵌套组合：

```
ChainAgent         ParallelAgent         CycleAgent
 A → B → C         A ↘               A（循环）
 顺序串联           B → Merge         ↓ ↑
 流水线式           C ↗              B（执行）
                   并发协作          迭代优化
```

## 3.2 ChainAgent —— 链式顺序执行

子 Agent 按顺序执行，前一个的输出成为后一个的上下文：

```go
chain := chainagent.New("pipeline",
    chainagent.WithSubAgents([]agent.Agent{
        analyzer,     // 分析阶段
        processor,    // 处理阶段
        reporter,     // 报告阶段
    }),
)

r := runner.NewRunner("app", chain)
events, _ := r.Run(ctx, "user-001", "session-001", message)
```

**适用场景**：多步骤流水线（分析 → 处理 → 汇总）、翻译 + 润色链。

## 3.3 ParallelAgent —— 并行执行与结果合并

多个子 Agent 同时执行，结果自动合并：

```go
parallel := parallelagent.New("concurrent",
    parallelagent.WithSubAgents([]agent.Agent{
        sentimentAgent,   // 情感分析
        keywordsAgent,    // 关键词提取
        summaryAgent,     // 摘要生成
    }),
)
```

**适用场景**：多专家从不同角度并发分析同一内容。

## 3.4 CycleAgent —— 循环迭代执行

循环执行子 Agent，直到满足终止条件：

```go
cycle := cycleagent.New("optimizer",
    cycleagent.WithSubAgents([]agent.Agent{planner, executor}),
    cycleagent.WithMaxIterations(5),
    // WithStopCondition: 基于事件判断终止
)
```

**适用场景**：自我优化、反复改进生成质量。

## 3.5 组合嵌套模式

三种模式可自由嵌套，构建复杂的多 Agent 体系：

```go
// 并行分析组
analyzeGroup := parallelagent.New("analysis",
    parallelagent.WithSubAgents([]agent.Agent{
        marketAnalyzer, techAnalyzer, riskAnalyzer,
    }),
)

// 迭代改进组
refinementLoop := cycleagent.New("refinement",
    cycleagent.WithSubAgents([]agent.Agent{writer, reviewer}),
    cycleagent.WithMaxIterations(3),
)

// 链式串联全部
pipeline := chainagent.New("full-pipeline",
    chainagent.WithSubAgents([]agent.Agent{
        analyzeGroup,      // 并行分析
        planner,            // 规划
        refinementLoop,     // 迭代改进
        publisher,          // 发布
    }),
)
```

---

# 四、Graph 工作流篇

## 4.1 GraphAgent 概述

GraphAgent 是 tRPC-Agent-Go 的图工作流引擎，功能等价于 LangGraph for Go。将可控的工作流编排与可扩展的 Agent 能力结合。

**核心特性**：

- **类型安全的状态管理**：Schema 定义状态结构，支持自定义 Reducer（Overwrite / Append / Merge）
- **条件路由**：基于状态动态选择执行路径，支持多条件并行路由（BSP 引擎）
- **两种执行引擎**：BSP（默认超步式确定性并行）与 DAG（即时调度无全局屏障）
- **LLM / Tool / Agent 节点**：内置大模型、工具调用和子 Agent 节点
- **流式执行**：节点通过 EventEmitter 发射实时事件和进度
- **CheckPoint**：基于检查点的时间旅行，浏览历史并恢复状态
- **HITL**：中断与恢复的交互式工作流
- **子 Agent 调度**：通过 `AddAgentNode` 委托执行子 Agent

## 4.2 核心概念：Graph / Node / State / Edge

### Graph — 工作流容器

```go
import "trpc.group/trpc-go/trpc-agent-go/graph"

schema := graph.NewStateSchema()
sg := graph.NewStateGraph(schema)
```

### Node — 处理步骤

```go
type NodeFunc func(ctx context.Context, state graph.State) (any, error)

// 添加普通节点
sg.AddNode("process", processFunc)

// 添加 LLM 节点
sg.AddLLMNode("ask_llm", modelInstance, systemPrompt, tools)

// 添加工具执行节点
sg.AddToolsNode("execute_tools", tools)

// 添加子 Agent 节点
sg.AddAgentNode("sub_agent",
    graph.WithAgentNodeName("researcher"),
)
```

### State — 数据容器与 Schema

```go
// 使用预设的消息状态 Schema
schema := graph.MessagesStateSchema()
// 预定义键：messages / user_input / last_response / system_prompt

// 自定义 Schema
schema := graph.NewStateSchema().
    AddField("counter", graph.StateField{
        DefaultValue: 0,
        Reducer:      graph.OverwriteReducer, // 覆盖策略
    }).
    AddField("results", graph.StateField{
        DefaultValue: []string{},
        Reducer:      graph.AppendReducer,    // 追加策略
    })
```

### Edge — 数据流与控制流

```go
// 普通边
sg.AddEdge("node_a", "node_b")

// 虚拟节点：Start 和 End
sg.SetEntryPoint("prepare")   // Start 自动连接到指定节点
sg.SetFinishPoint("finish")   // 指定节点自动连接到 End

// 条件边
sg.AddConditionalEdges("router", conditionFunc, pathMap)

// 多条件并行边（BSP 引擎独有）
sg.AddMultiConditionalEdges("router", multiConditionFunc, pathMap)
```

## 4.3 BSP 与 DAG 执行引擎

| 引擎 | 模式 | 特点 |
|------|------|------|
| **BSP**（默认） | Bulk Synchronous Parallel | 超步式：计划/执行/合并，确定性并行，支持多条件路由 |
| **DAG** | Directed Acyclic Graph | 即时调度：无全局屏障，节点就绪即执行 |

```go
// 使用 DAG 引擎
dagGraph, _ := sg.Compile(graph.WithDAGMode())
```

## 4.4 条件路由与多条件并行路由

**条件路由** — 基于状态动态选择单一路径：

```go
sg.AddConditionalEdges("llm_node",
    func(ctx context.Context, s graph.State) (string, error) {
        msgs := s[graph.StateKeyMessages].([]model.Message)
        last := msgs[len(msgs)-1]
        if len(last.ToolCalls) > 0 {
            return "use_tools", nil
        }
        return "finish", nil
    },
    map[string]string{
        "use_tools": "execute_tools",
        "finish":    graph.End,
    },
)
```

**多条件并行路由** — 一个节点同时路由到多个目标：

```go
sg.AddMultiConditionalEdges("router",
    func(ctx context.Context, s graph.State) ([]string, error) {
        // 返回多个键 → 对应节点并行执行
        return []string{"goA", "goB", "goC"}, nil
    },
    map[string]string{
        "goA": "node_a", "goB": "node_b", "goC": "node_c",
    },
)
sg.SetFinishPoint("node_a").SetFinishPoint("node_b").SetFinishPoint("node_c")
```

## 4.5 CheckPoint 与时间旅行

完整的状态快照与恢复机制：

```go
graphAgent, _ := graphagent.New("workflow", compiledGraph,
    graphagent.WithCheckpointSaver(saver),
)

// 从检查点恢复
events, _ := runner.Run(ctx, userID, sessionID, message,
    agent.WithResume(true),
    agent.WithCheckpointID("checkpoint-xxx"),
)
```

**核心能力**：
- **原子检查点**：状态与待写入数据原子化存储
- **检查点谱系**：父子关系追踪，支持执行历史浏览
- **时间旅行**：在任意历史检查点恢复执行

## 4.6 Human-in-the-Loop（HITL）

```go
// 节点中触发中断
func approvalNode(ctx context.Context, s graph.State) (any, error) {
    if !s["approved"].(bool) {
        return nil, graph.NewInterruptError("需要审批确认")
    }
    return graph.State{"status": "approved"}, nil
}

// 恢复执行
events, _ := runner.Run(ctx, userID, sessionID,
    graph.Command{ResumeMap: map[string]any{
        "approval": true,
    }},
    agent.WithResume(true),
)
```

## 4.7 子 Agent 节点

在 Graph 流程中委托子 Agent：

```go
graphAgent, _ := graphagent.New("main-workflow", compiledGraph,
    graphagent.WithSubAgents([]agent.Agent{
        researchAgent,
        codeAgent,
    }),
)

sg.AddAgentNode("delegate_research",
    graph.WithAgentNodeName("researchAgent"),
)
```

## 4.8 GraphAgent 完整示例

经典的 "prepare → ask LLM → 可能调用工具" ReAct 循环：

```
START → prepare → ask_llm ──(tool_calls)──→ tools → ask_llm
                         └──(no_tools)───→ fallback → END
```

```go
schema := graph.NewStateSchema().AddField(graph.StateKeyMessages, graph.StateField{
    Reducer: graph.AppendReducer,
})

sg := graph.NewStateGraph(schema)

sg.AddNode("prepare", prepareInput)
sg.AddLLMNode("ask_llm", modelInstance, systemPrompt, tools)
sg.AddToolsNode("tools", tools)
sg.AddNode("fallback", formatResponse)

sg.SetEntryPoint("prepare")
sg.AddEdge("prepare", "ask_llm")

sg.AddConditionalEdges("ask_llm",
    func(ctx context.Context, s graph.State) (string, error) {
        msgs := s[graph.StateKeyMessages].([]model.Message)
        if len(msgs[len(msgs)-1].ToolCalls) > 0 {
            return "use_tools", nil
        }
        return "finish", nil
    },
    map[string]string{"use_tools": "tools", "finish": "fallback"},
)

sg.AddEdge("tools", "ask_llm")  // 工具结果反馈给 LLM，形成循环
sg.AddEdge("fallback", graph.End)

compiled, _ := sg.Compile()
ag, _ := graphagent.New("react-agent", compiled)
r := runner.NewRunner("react-app", ag)
```

---

# 五、模型与工具篇

## 5.1 Model —— 多模型支持

框架通过统一接口支持多种 LLM：

```go
type Model interface {
    Generate(ctx context.Context, messages []Message, ...) (*Response, error)
    Stream(ctx context.Context, messages []Message, ...) (<-chan *StreamChunk, error)
    Info() ModelInfo
}
```

**内置模型**：

| 模型 | 用法 |
|------|------|
| OpenAI 兼容 | `openai.New("gpt-4o-mini")` |
| DeepSeek | `openai.New("deepseek-chat", openai.WithVariant(openai.VariantDeepSeek))` |

```bash
# 环境变量配置
export OPENAI_API_KEY="your-api-key"
export OPENAI_BASE_URL="https://api.deepseek.com"  # 任意兼容 API
```

## 5.2 Function Tool —— 函数工具

任意函数注册为 Tool，框架从结构体 tag 自动推断参数 Schema：

```go
type CalculatorReq struct {
    A  float64 `json:"a" jsonschema:"description=第一个操作数,required"`
    B  float64 `json:"b" jsonschema:"description=第二个操作数,required"`
    Op string  `json:"op" jsonschema:"description=运算,enum=add,enum=sub,enum=mul,enum=div,required"`
}
type CalculatorRsp struct {
    Result float64 `json:"result"`
}

tool := function.NewFunctionTool(calculatorFunc,
    function.WithName("calculator"),
    function.WithDescription("执行加减乘除"),
)
```

## 5.3 MCP Tool —— 模型上下文协议

通过 MCP 协议集成外部工具服务器：

```go
import "trpc.group/trpc-go/trpc-agent-go/tool/mcptool"

mcpToolSet := mcptool.NewMCPToolSet(serverConn)
tools := mcpToolSet.Tools() // 从 MCP 服务获取所有工具
```

基于 [tRPC-MCP-Go](https://github.com/trpc-group/trpc-mcp-go)，完整遵循 MCP 规范。

## 5.4 内置工具集

| 工具 | 包 | 功能 |
|------|-----|------|
| DuckDuckGo | `tool/duckduckgo` | 网络搜索 |
| FileToolSet | `tool/file` | 文件读/写/列表/编辑/搜索 |
| AgentTool | `tool/agenttool` | 将 Agent 包装为可调用工具 |
| Memory Tools | `memory` 包 | 记忆增/删/改/查 |
| Knowledge Search | `knowledge` 包 | 知识库语义搜索 |
| Skill Tools | `skill` 包 | 技能加载与运行 |

## 5.5 工具组合与调用策略

```go
agent := llmagent.New("assistant",
    llmagent.WithTools([]tool.Tool{
        calculatorTool, searchTool, weatherTool,
    }),
    llmagent.WithToolSets([]tool.ToolSet{
        fileToolSet,   // 文件工具集
        mcpToolSet,    // MCP 工具集
    }),
)
```

---

# 六、记忆与会话篇

## 6.1 Session —— 会话管理

管理用户会话状态和事件历史：

```go
import "trpc.group/trpc-go/trpc-agent-go/session"

// 开发：内存存储
sessionService := session.NewInMemorySessionService()

// 生产：Redis 存储
sessionService := session.NewRedisSessionService(redisClient)

runner := runner.NewRunner("app", agent,
    runner.WithSessionService(sessionService),
)
```

**状态操作**：

```go
// 用户级别状态
sessionService.UpdateUserState(ctx, session.UserKey{
    AppName: "my-app", UserID: "user-001",
}, session.StateMap{"topics": []byte("AI, ML")})

// 应用级别状态
sessionService.UpdateAppState(ctx, "my-app", session.StateMap{
    "banner": []byte("Research Mode"),
})

// 会话初始化
sessionService.CreateSession(ctx, session.Key{
    AppName: "my-app", UserID: "user-001", SessionID: "sess-001",
}, session.StateMap{"topic": []byte("AI agents")})
```

## 6.2 Memory —— 长期记忆

记录用户的长期记忆和个性化信息：

```go
import (
    "trpc.group/trpc-go/trpc-agent-go/memory"
    memorysvc "trpc.group/trpc-go/trpc-agent-go/memory/service"
)

memoryService := memorysvc.NewInMemoryService()

// 接入 Agent
agent := llmagent.New("assistant",
    llmagent.WithModel(modelInstance),
    llmagent.WithTools(memoryService.Tools()),
)

// 接入 Runner
r := runner.NewRunner("app", agent,
    runner.WithMemoryService(memoryService),
)
```

## 6.3 Memory Tools 集成

Agent 可在对话中自主管理记忆：

| 工具 | 功能 |
|------|------|
| `memory_add` | 添加新记忆 |
| `memory_update` | 更新已有记忆 |
| `memory_delete` | 删除记忆 |
| `memory_search` | 语义搜索记忆 |
| `memory_load` | 加载记忆列表 |

## 6.4 占位符变量 —— 会话状态注入

Instruction 和 SystemPrompt 中的占位符自动从会话状态注入：

```go
llm := llmagent.New("research-agent",
    llmagent.WithModel(modelInstance),
    llmagent.WithInstruction(
        "You are a research assistant. Focus: {research_topics}. " +
        "User interests: {user:topics?}. " +   // 可选，不存在时替换为空串
        "App mode: {app:banner?}." +
        "Case: {invocation:case}",
    ),
)
```

| 语法 | 来源 | 说明 |
|------|------|------|
| `{key}` | Session.State | 会话级键值 |
| `{user:key}` | User State | 用户命名空间 |
| `{app:key}` | App State | 应用命名空间 |
| `{temp:key}` | Temp State | 临时命名空间 |
| `{invocation:key}` | Invocation.State | 当前 invocation |
| `{key?}` | 可选 | 不存在时替换为空串，不保留原样 |

---

# 七、知识检索篇（RAG）

## 7.1 Knowledge 系统概述

提供完整的检索增强生成（RAG）能力：

- **智能检索**：向量相似度语义搜索
- **多源支持**：文件、目录、URL
- **文档提取**：PDF、HTML 自动转 Markdown/文本
- **灵活存储**：内存、PostgreSQL（pgvector）、TcVector
- **知识过滤**：元数据静态过滤 + Agent 智能过滤
- **并发处理**：批量加载与高性能索引

## 7.2 基本使用流程

```go
import (
    "trpc.group/trpc-go/trpc-agent-go/knowledge"
    "trpc.group/trpc-go/trpc-agent-go/knowledge/source"
)

// 1. 创建 Knowledge
kb := knowledge.New(
    knowledge.WithVectorStore(vectorStore),
    knowledge.WithEmbedder(embedder),
    knowledge.WithSources([]source.Source{
        source.NewFileSource("./docs"),
        source.NewURLSource("https://example.com/doc.md"),
    }),
)

// 2. 加载文档
kb.Load(ctx)

// 3. 创建搜索工具
searchTool := knowledge.NewKnowledgeSearchTool(kb,
    knowledge.WithToolName("knowledge_search"),
    knowledge.WithToolDescription("搜索知识库获取相关信息"),
)

// 4. 集成
agent := llmagent.New("rag-agent",
    llmagent.WithModel(modelInstance),
    llmagent.WithTools([]tool.Tool{searchTool}),
)
```

**简便集成**（单一知识库，自动添加 `knowledge_search` 工具）：

```go
agent := llmagent.New("agent",
    llmagent.WithKnowledge(kb),
)
```

## 7.3 智能过滤与多知识库

**智能过滤**：Agent 自动决定知识库过滤条件：

```go
filterTool := knowledge.NewAgenticFilterSearchTool(kb,
    knowledge.WithToolName("smart_search"),
)
```

**多知识库**：手动创建多个不同名的搜索工具：

```go
productTool := knowledge.NewKnowledgeSearchTool(productKB,
    knowledge.WithToolName("search_products"),
)
docTool := knowledge.NewKnowledgeSearchTool(docKB,
    knowledge.WithToolName("search_docs"),
)

agent := llmagent.New("agent",
    llmagent.WithTools([]tool.Tool{productTool, docTool}),
)
```

## 7.4 存储后端与 Embedder

| 后端 | 适用场景 |
|------|---------|
| InMemory | 开发测试 |
| PostgreSQL + pgvector | 通用生产环境 |
| TcVector | 腾讯云环境 |

---

# 八、技能与制品篇

## 8.1 Agent Skills —— 可复用工作流

Skills 是基于 `SKILL.md` 规范的可复用工作流。每个 Skill 是一个包含描述文件和可选文档/脚本的目录：

```go
import (
    "trpc.group/trpc-go/trpc-agent-go/skill"
    "trpc.group/trpc-go/trpc-agent-go/tool/skilltool"
)

// 从本地目录加载
repo, _ := skill.NewFSRepository("./skills")

// 从远程加载（支持 .zip / .tar.gz，自动缓存）
repo, _ := skill.NewFSRepository("https://example.com/skills.zip")

// 多源组合
repo, _ := skill.NewFSRepository("./shared-skills", "./private-skills")

// 长时间运行中刷新
repo.Refresh()

// 创建技能工具
tools := []tool.Tool{
    skilltool.NewLoadTool(repo),
    skilltool.NewRunTool(repo, localexec.New()),
    skilltool.NewListDocsTool(repo),
    skilltool.NewSelectDocsTool(repo),
}

agent := llmagent.New("agent",
    llmagent.WithTools(tools),
    llmagent.WithCodeExecutor(exec),
    llmagent.WithEnableCodeExecutionResponseProcessor(false),
)
```

## 8.2 技能发现与动态安装

`examples/skillfind` 展示了完整的技能发现闭环：

1. Agent 使用内置 `skill-find` 技能搜索公开可用技能
2. 从 GitHub 安装技能到用户私有目录
3. 调用 `repo.Refresh()` 刷新技能列表
4. 在同一对话中立即使用新技能

## 8.3 Artifacts —— 版本化文件存储

存储 Agent 和工具生成的图片、文件、报告：

```go
import "trpc.group/trpc-go/trpc-agent-go/artifact"

artifactService := artifact.NewInMemoryService()

// 保存
artifactService.Save(ctx, userKey, artifactID, &artifact.Artifact{
    Name:    "report.pdf",
    Content: reportBytes,
    Version: 1,
})

r := runner.NewRunner("app", agent,
    runner.WithArtifactService(artifactService),
)
```

支持后端：InMemory / S3 / COS（腾讯云对象存储）。

---

# 九、可观测性与评估篇

## 9.1 OpenTelemetry 全链路追踪

内建 OpenTelemetry，覆盖 Model、Tool、Runner 各层：

```go
// 启动 Langfuse 集成
clean, _ := langfuse.Start(ctx)
defer clean(ctx)

r := runner.NewRunner("app", agent)

events, _ := r.Run(ctx, "user-001", "session-001", message,
    agent.WithSpanAttributes(
        attribute.String("langfuse.user.id", "user-001"),
        attribute.String("langfuse.session.id", "session-001"),
    ),
)
```

追踪覆盖：Model 调用与 Token 消耗、Tool 调用参数、Agent 推理过程、Runner 生命周期。

## 9.2 Callback 机制

Agent / Model / Tool 三层 Callback 注入点：

```go
agent := llmagent.New("agent",
    // Agent 级
    llmagent.WithAgentCallbacks(&agent.Callbacks{
        BeforeAgent: func(...) { /* 前置处理 */ },
        AfterAgent:  func(...) { /* 后置处理 */ },
    }),
    // Model 级
    llmagent.WithModelCallbacks(&model.Callbacks{
        BeforeModel: func(...) { /* 模型调用前 */ },
        AfterModel:  func(...) { /* 模型调用后 */ },
    }),
    // Tool 级
    llmagent.WithToolCallbacks(&tool.Callbacks{
        BeforeTool: func(...) { /* 工具调用前 */ },
        AfterTool:  func(...) { /* 工具调用后 */ },
    }),
)
```

## 9.3 Evaluation —— 评估与基准测试

```go
import "trpc.group/trpc-go/trpc-agent-go/evaluation"

evaluator, _ := evaluation.New("app", runner,
    evaluation.WithNumRuns(3),       // 多次运行取平均
)
defer evaluator.Close()

result, _ := evaluator.Evaluate(ctx, "math-basic")

// 结果包含
result.OverallStatus   // 整体状态
result.Metrics         // 各项指标
result.Runs            // 每轮详细结果
```

支持本地文件存储和内存存储。

---

# 十、服务与集成篇

## 10.1 AG-UI —— Agent-用户交互协议

基于 SSE 的 Agent-用户交互协议，提供标准化的前后端接口：

```go
import "trpc.group/trpc-go/trpc-agent-go/server/agui"

server := agui.NewServer(runner,
    agui.WithPath("/api/agent/chat"),
)
server.Start(":8080")
```

客户端示例包括 CopilotKit、TDesign Chat。

## 10.2 A2A —— Agent-to-Agent 互操作

跨运行时、跨语言的 Agent 互操作协议：

```go
import "trpc.group/trpc-go/trpc-agent-go/server/a2a"

server := a2a.NewServer(agent,
    a2a.WithPath("/a2a"),
)
```

基于 [tRPC-A2A-Go](https://github.com/trpc-group/trpc-a2a-go)。`examples/a2aadk` 展示了与 Python ADK A2A 服务的流式对话和代码执行。

## 10.3 Gateway Server

类 OpenClaw 的网关服务，提供：

- 稳定会话 ID 与会话级序列化
- 安全控制：白名单（allowlist）+ 提及门控（mention gating）
- 多平台接入（Telegram 等）

---

# 十一、进阶篇

## 11.1 Runner 深入理解

Runner 是框架的编排中心：

```go
type Runner interface {
    Run(ctx context.Context, userID, sessionID string,
        message model.Message, opts ...agent.RunOption) (<-chan *event.Event, error)
    Close() error
}
```

**Runner 五大职责**：

1. **会话管理**：自动创建/查找 Session，注入 SessionService
2. **记忆集成**：串联 MemoryService，管理长期记忆
3. **上下文注入**：将 Session 状态注入 Agent 的 Instruction 占位符
4. **事件流管理**：统一事件通道，桥接 Agent 输出与消费方
5. **生命周期**：取消、超时、资源清理

## 11.2 Planner —— 规划与推理

```go
import "trpc.group/trpc-go/trpc-agent-go/planner"

planner := planner.NewBuiltinPlanner(
    planner.WithReasoningEffort("high"),
)

agent := llmagent.New("planner-agent",
    llmagent.WithPlanner(planner),
)
```

Planner 在 Agent 推理前进行层级规划，决定执行策略和工具选择。

## 11.3 代码执行

安全的代码执行沙箱：

```go
import "trpc.group/trpc-go/trpc-agent-go/codeexecutor"

executor := codeexecutor.NewLocalCodeExecutor(
    codeexecutor.WithWorkDir("/tmp/sandbox"),
)

agent := llmagent.New("coder-agent",
    llmagent.WithCodeExecutor(executor),
)
```

## 11.4 Prompt Caching

框架自动优化：System Prompt 和工具定义等稳定内容自动标记为可缓存，节省约 90% 的重复 Token 成本。

## 11.5 生产部署建议

### 存储选择

| 组件 | 开发 | 生产 |
|------|------|------|
| SessionService | InMemory | Redis |
| MemoryService | InMemory | Redis |
| VectorStore | InMemory | PostgreSQL + pgvector / TcVector |
| ArtifactService | InMemory | S3 / COS |
| CheckPointStore | InMemory | Redis / 数据库 |

### 高可用

- Runner 无状态 → 水平扩展
- `NewRunnerWithAgentFactory` 按请求隔离 Agent
- Redis SessionService 支持多实例共享

### 监控

- 接入 OpenTelemetry + Langfuse
- 定期运行 Evaluation 跟踪 Agent 质量
- 结构化日志记录关键决策

## 11.6 生态与社区

### 项目地址

| 资源 | 链接 |
|------|------|
| tRPC-Agent-Go | [github.com/trpc-group/trpc-agent-go](https://github.com/trpc-group/trpc-agent-go) |
| tRPC-A2A-Go | [github.com/trpc-group/trpc-a2a-go](https://github.com/trpc-group/trpc-a2a-go) |
| tRPC-MCP-Go | [github.com/trpc-group/trpc-mcp-go](https://github.com/trpc-group/trpc-mcp-go) |
| 文档 | [trpc-group.github.io/trpc-agent-go/](https://trpc-group.github.io/trpc-agent-go/) |

### 社区

- **Issues**: [github.com/trpc-group/trpc-agent-go/issues](https://github.com/trpc-group/trpc-agent-go/issues)
- **贡献指南**: `CONTRIBUTING.md`

### 版本与成熟度

- 最新版本：**v1.9.1**（2026-05-11）
- Go 要求：**≥ 1.21**
- 发布版本：**114+**
- 贡献者：**40+**
- Stars：**1.2K+**
- 生产验证：腾讯元宝、腾讯视频、腾讯新闻、IMA、QQ 音乐等

### 致谢

感谢 ADK、Agno、CrewAI、AutoGen 等开源框架的启发。

---

> **tRPC-Agent-Go** —— 用 Go 语言，构建去中心化、可协作、能涌现的下一代智能 Agent 系统。
