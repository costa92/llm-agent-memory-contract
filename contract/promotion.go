package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// PromoteImportanceThreshold is the minimum Importance an agent_inferred record
// must reach to be eligible for promotion out of the working tier.
const PromoteImportanceThreshold = 0.7

// PromotionEligible reports whether a record qualifies for promotion out of the
// working tier, by source: user_saved always qualifies; agent_inferred
// qualifies at or above PromoteImportanceThreshold; every other source does not.
//
// This is the single source of truth shared by the consolidation worker
// (async, outbox-driven) and the gateway session-close path (synchronous), so
// the two promotion triggers agree on eligibility. Per-caller concerns — the
// promotion Reason string and the idempotency-key salt — intentionally stay in
// each caller and are NOT part of this shared contract.
func PromotionEligible(record MemoryRecord) bool {
	switch record.Source {
	case "user_saved":
		return true
	case "agent_inferred":
		return record.Importance >= PromoteImportanceThreshold
	default:
		return false
	}
}

// DedupeKey derives the dedupe identity for a record: a stable lowercase-hex
// sha256 over tenant, user, category, project, and the normalized content.
// Two records that hash to the same key are dedupe candidates.
//
// The key includes project_id so identical content in different projects does
// not collide; it deliberately excludes session_id so the same content across
// sessions IS a valid dedupe candidate. Single source of truth shared by the
// worker and the gateway session-close path (see PromotionEligible).
func DedupeKey(record MemoryRecord) string {
	return dedupeHashParts(
		record.TenantID,
		record.UserID,
		record.Category,
		record.ProjectID,
		NormalizeDedupeContent(record),
	)
}

// NormalizeDedupeContent returns the content component of the dedupe key: the
// trimmed NormalizedContentHash when present, otherwise the record Content
// lowercased with surrounding and interior whitespace collapsed to single
// spaces.
func NormalizeDedupeContent(record MemoryRecord) string {
	if h := strings.TrimSpace(record.NormalizedContentHash); h != "" {
		return h
	}
	return strings.Join(strings.Fields(strings.ToLower(record.Content)), " ")
}

func dedupeHashParts(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
}
