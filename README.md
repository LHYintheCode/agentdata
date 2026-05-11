# Agentdata

Agentdata is a local-first CLI for developers who want to search, normalize, and export their coding-agent chat history.

Agentdata 是一个本地优先的 CLI 工具，用于帮助开发者检索、标准化和导出自己与编程智能体的聊天记录。

## English

### What It Does

Agentdata starts with explicit-path JSONL ingestion and a small canonical schema. It is designed for future adapters for Codex, Claude Code, Trae, and other agent tools without tying the core data model to any one vendor format.

The long-term idea is to give users control over their own AI interaction data: search it locally, package it safely, and decide how it can be reused.

### Privacy Defaults

- Local-only by default.
- Read-only scanning.
- No telemetry.
- No cloud sync.
- Exports go to stdout in the initial MVP.
- Marketplace, resale, consent manifests, and redaction workflows are not enabled in the MVP.

### MVP Commands

```powershell
go run ./cmd/agentdata version
go run ./cmd/agentdata scan --path ./samples
go run ./cmd/agentdata search --path ./samples "deploy"
go run ./cmd/agentdata export --path ./samples --format markdown
```

### Current Status

Early-stage public CLI. The current implementation supports:

- JSONL parsing
- Canonical session/message model
- Local message search
- JSONL export
- Markdown export
- Basic CLI tests

Automatic discovery of Codex, Claude Code, Trae, Cursor, Windsurf, and other local storage paths will be added after more real-machine validation.

## 中文

### 这个项目做什么

Agentdata 的目标是做一个本地优先的开发者 CLI，用来整理不同 AI 编程工具里的聊天记录和 agent 轨迹。

第一版从显式传入路径的 JSONL 文件开始，把不同来源的数据统一成一个简单稳定的会话/消息模型。后续可以逐步加入 Codex、Claude Code、Trae、Cursor、Windsurf 等工具的本地数据适配器。

更长期的方向是让用户真正拥有自己的 AI 交互数据：可以本地检索，可以安全打包，也可以在明确授权后决定如何复用。

### 隐私默认策略

- 默认只在本地运行。
- 扫描是只读的。
- 不采集遥测数据。
- 不做云同步。
- MVP 阶段导出内容只输出到 stdout。
- MVP 阶段不包含数据市场、转售、授权清单、脱敏工作流。

### MVP 命令

```powershell
go run ./cmd/agentdata version
go run ./cmd/agentdata scan --path ./samples
go run ./cmd/agentdata search --path ./samples "deploy"
go run ./cmd/agentdata export --path ./samples --format markdown
```

### 当前状态

这是一个早期公开 CLI。目前已经支持：

- JSONL 解析
- 统一的 session/message 数据模型
- 本地消息搜索
- JSONL 导出
- Markdown 导出
- 基础 CLI 测试

Codex、Claude Code、Trae、Cursor、Windsurf 等工具的自动路径发现和格式适配，需要在更多真实机器上验证后再加入。
