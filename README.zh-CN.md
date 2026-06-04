[English](./README.md) | [简体中文](./README.zh-CN.md)

# llm-agent-memory-contract

为 llm-agent 生态提供的仅标准库、后端中立的持久记忆契约。逐字提取自
`llm-agent-memory/memory/durable.go`。

## 包含内容

- `MemoryRecord` 以及持久化聚合类型 `StoredEvent`、`OutboxMessage`、
  `IdempotencyEntry`。
- 全部 `*Input` / `*Result` DTO 以及 `DedupeAction`。
- 8 个存储端口接口：`RecordStore`、`Promoter`、`Deduper`、
  `AccessMarker`、`EventStore`、`IdempotencyStore`、`Outbox`、
  `MessagePublisher`。
- `RecordKind*` / `Dedupe*` 常量、`ErrInvalidRecordKind`，以及
  `NormalizeRecordKind` / `NormalizeWriteDefaults` / `SetWorkingDefault`
  辅助函数。

## 稳定性契约

本 module 是一份 **持久化的 JSON schema**，而不仅仅是一个 DTO 包。这四个
聚合类型使用默认的 `encoding/json`（无 tag）直接序列化写入 Postgres，因此
**wire 键名等于 Go 字段名**。

未经主版本号提升和数据库迁移计划，请勿：

- 重命名字段；
- 在值形式与指针形式之间切换字段；
- 修改字段的类型；
- 添加 `json` tag。

`golden_wire_test.go` 锁定了精确的 wire 字节。黄金测试变红意味着该改动会
损坏此前已持久化的行。

## 版本管理

独立 module；以 `llm-agent-memory-contract/vX.Y.Z` 形式打 tag。消费方必须在
一次协调发布波次中锚定 **相同的** 契约版本（`llm-agent-memory` 中的别名垫片
**并不能** 让混合版本的依赖图变得安全）。
