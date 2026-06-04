# Changelog

`github.com/costa92/llm-agent-memory-contract` 的所有重要变更都将记录在本文件中。

<!-- Keep a Changelog format: https://keepachangelog.com/en/1.1.0/ -->
<!-- Semver: https://semver.org/ -->

## [Unreleased]

## [0.2.0] - 2026-06-02

### Added

- **共享的晋升策略（M8 D3）。** `PromotionEligible(record)`、
  `PromoteImportanceThreshold = 0.7`、`DedupeKey(record)` 和
  `NormalizeDedupeContent(record)` —— 晋升规则与去重身份的单一可信来源，
  此前曾在固化工作进程和网关会话关闭路径中逐字节重复实现。纯属仅增量。

## [0.1.0] - 2026-05-26

### Added

- 初始的仅标准库、后端中立的持久记忆契约，逐字提取自
  `llm-agent-memory/memory/durable.go`。
- `MemoryRecord` 以及持久化聚合类型 `StoredEvent`、
  `OutboxMessage`、`IdempotencyEntry`。
- 全部 `*Input` / `*Result` DTO 以及 `DedupeAction`。
- 8 个存储端口接口：`RecordStore`、`Promoter`、`Deduper`、
  `AccessMarker`、`EventStore`、`IdempotencyStore`、`Outbox`、
  `MessagePublisher`。
- `RecordKind*` / `Dedupe*` 常量、`ErrInvalidRecordKind`，以及
  `NormalizeRecordKind` / `NormalizeWriteDefaults` / `SetWorkingDefault`
  辅助函数。

### Notes

- 设计上仅依赖标准库：没有任何第三方依赖，因此每个兄弟仓 module 都可以
  依赖它而不会引入一棵依赖树。
- 本 module 是一份 **持久化的 JSON schema**：这四个聚合类型使用默认的
  `encoding/json`（无 tag）直接序列化写入 Postgres，因此 wire 键名等于
  Go 字段名。重命名/改变字段类型、在值形式与指针形式之间切换、或添加
  `json` tag，都需要主版本号提升和数据库迁移计划。
