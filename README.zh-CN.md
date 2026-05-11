# Agentdata

[English](README.md)

Agentdata 是一个本地优先的 CLI，用来检索、标准化和导出 AI 编程智能体的聊天记录。

它面向同时使用 Codex、Claude Code、Trae、Cursor、Windsurf 等工具的开发者，目标是把分散在不同产品里的 AI 工作记录整理成一个可迁移、可搜索、可审计的本地数据档案。

## 为什么做

AI 对话正在变成开发工作流的一部分。调试过程、设计决策、工具调用、失败尝试、最终修复，经常只存在于某一个 agent 产品里。

Agentdata 希望把这些记录变成用户自己可控的本地数据：

- 检索过去和 agent 的对话。
- 导出 session 为 JSONL 或 Markdown。
- 把不同工具的记录统一成同一种 schema。
- 默认把敏感数据留在本机。
- 为后续的授权、脱敏、打包和复用建立清晰的数据基础。

## 当前状态

Agentdata 还处在早期阶段。当前版本支持显式传入路径的 JSONL 输入。Codex、Claude Code、Trae、Cursor、Windsurf 等工具的自动发现和格式适配，需要更多真实机器验证后再加入。

当前已经支持：

- JSONL 解析
- 从 `~/.codex/sessions` 解析 Codex rollout
- 统一的 session/message 数据模型
- 本地消息搜索
- JSONL 导出
- Markdown 导出
- 基础 CLI 测试

暂未包含：

- 非 Codex 工具的自动发现
- SQLite/FTS 索引
- secret 脱敏
- 授权 manifest
- 云同步
- 数据市场或转售流程

## 安装

从源码安装：

```bash
go install github.com/LHYintheCode/agentdata/cmd/agentdata@latest
```

本地开发：

```bash
git clone https://github.com/LHYintheCode/agentdata.git
cd agentdata
go test ./...
go run ./cmd/agentdata version
```

## 快速开始

准备一个 JSONL 文件：

```jsonl
{"source":"codex","project":"/path/to/project","session_id":"s1","timestamp":"2026-05-11T01:02:03Z","role":"user","content":"Deploy the CLI"}
{"source":"codex","project":"/path/to/project","session_id":"s1","timestamp":"2026-05-11T01:03:04Z","role":"assistant","content":"Run go test ./..."}
```

扫描：

```bash
agentdata scan --path ./samples
```

搜索：

```bash
agentdata search --path ./samples "deploy"
```

导出：

```bash
agentdata export --path ./samples --format markdown
agentdata export --path ./samples --format jsonl
```

扫描本机 Codex sessions：

```bash
agentdata scan --source codex
agentdata search --source codex "deploy"
agentdata export --source codex --format markdown > codex-history.md
```

## 命令

```text
agentdata version
agentdata scan --path <file-or-directory>
agentdata scan --source codex [--path <codex-sessions-directory>]
agentdata search --path <file-or-directory> <query>
agentdata search --source codex [--path <codex-sessions-directory>] <query>
agentdata export --path <file-or-directory> --format jsonl|markdown
agentdata export --source codex [--path <codex-sessions-directory>] --format jsonl|markdown
```

## 数据模型

Agentdata 会把来源数据标准化成 session 和 message：

```json
{
  "id": "s1",
  "source": "codex",
  "project": "/path/to/project",
  "messages": [
    {
      "role": "user",
      "text": "Deploy the CLI",
      "timestamp": "2026-05-11T01:02:03Z"
    }
  ]
}
```

各 source adapter 负责把某个工具的内部格式转换成这个模型。搜索和导出层不应该直接依赖任何单一厂商的内部存储格式。

## 隐私

Agentdata 的设计基线是本地所有权：

- 只有在你传入路径时才读取本地文件。
- 不上传数据。
- 不采集遥测。
- 不修改原始聊天记录。
- 当前 MVP 阶段导出内容只写到 stdout。

未来如果加入数据分享能力，必须先有明确授权、脱敏和机器可读的 manifest，再允许任何数据离开用户机器。

## 路线图

- Claude Code source adapter
- Trae/Cursor/Windsurf 本地存储调研
- SQLite FTS 索引
- secret 和个人信息脱敏规则
- 带 manifest 的导出数据包
- 为不能直接调用 CLI 的 agent host 提供 MCP adapter

## 开发

```bash
go test ./...
go run ./cmd/agentdata version
```

## 许可证

许可证尚未确定。
