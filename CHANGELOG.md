# Changelog

All notable changes to `github.com/costa92/llm-agent-memory-contract` will be
documented in this file.

<!-- Keep a Changelog format: https://keepachangelog.com/en/1.1.0/ -->
<!-- Semver: https://semver.org/ -->

## [Unreleased]

## [0.2.0] - 2026-06-02

### Added

- **Shared promotion policy (M8 D3).** `PromotionEligible(record)`,
  `PromoteImportanceThreshold = 0.7`, `DedupeKey(record)` and
  `NormalizeDedupeContent(record)` — the single source of truth for the
  promotion rule and dedupe identity, previously duplicated byte-for-byte in
  the consolidation worker and the gateway session-close path. Purely additive.

## [0.1.0] - 2026-05-26

### Added

- Initial stdlib-only, backend-neutral durable-memory contract, extracted
  verbatim from `llm-agent-memory/memory/durable.go`.
- `MemoryRecord` plus the persisted aggregates `StoredEvent`,
  `OutboxMessage`, `IdempotencyEntry`.
- All `*Input` / `*Result` DTOs and `DedupeAction`.
- The 8 storage-port interfaces: `RecordStore`, `Promoter`, `Deduper`,
  `AccessMarker`, `EventStore`, `IdempotencyStore`, `Outbox`,
  `MessagePublisher`.
- `RecordKind*` / `Dedupe*` constants, `ErrInvalidRecordKind`, and the
  `NormalizeRecordKind` / `NormalizeWriteDefaults` / `SetWorkingDefault`
  helpers.

### Notes

- Stdlib-only by design: no third-party dependencies, so every sibling
  module can depend on it without pulling a dependency tree.
- This module is a **persisted JSON schema**: the four aggregate types are
  marshaled with the default `encoding/json` (no tags) straight into
  Postgres, so wire keys equal Go field names. Renaming/retyping a field,
  switching value/pointer form, or adding a `json` tag requires a major
  version bump and a DB migration plan.
