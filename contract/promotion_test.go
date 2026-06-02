package contract

import "testing"

func TestPromotionEligible(t *testing.T) {
	cases := []struct {
		name       string
		source     string
		importance float64
		want       bool
	}{
		{"user_saved always", "user_saved", 0, true},
		{"user_saved high", "user_saved", 0.99, true},
		{"agent_inferred above threshold", "agent_inferred", 0.8, true},
		{"agent_inferred at threshold", "agent_inferred", 0.7, true},
		{"agent_inferred below threshold", "agent_inferred", 0.69, false},
		{"agent_inferred zero", "agent_inferred", 0, false},
		{"system never", "system", 0.99, false},
		{"empty source never", "", 0.99, false},
		{"unknown source never", "imported", 0.99, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PromotionEligible(MemoryRecord{Source: tc.source, Importance: tc.importance})
			if got != tc.want {
				t.Fatalf("PromotionEligible(source=%q, importance=%v) = %v, want %v", tc.source, tc.importance, got, tc.want)
			}
		})
	}
}

func TestPromoteImportanceThreshold(t *testing.T) {
	if PromoteImportanceThreshold != 0.7 {
		t.Fatalf("PromoteImportanceThreshold = %v, want 0.7 (worker and gateway both depend on this exact value)", PromoteImportanceThreshold)
	}
}

func TestNormalizeDedupeContent(t *testing.T) {
	cases := []struct {
		name string
		rec  MemoryRecord
		want string
	}{
		{
			name: "uses normalized hash when present",
			rec:  MemoryRecord{NormalizedContentHash: "precomputed-hash", Content: "ignored"},
			want: "precomputed-hash",
		},
		{
			name: "trims the normalized hash",
			rec:  MemoryRecord{NormalizedContentHash: "  precomputed-hash  "},
			want: "precomputed-hash",
		},
		{
			name: "lowercases and collapses content when no hash",
			rec:  MemoryRecord{Content: "Hello   World"},
			want: "hello world",
		},
		{
			name: "trims leading and trailing whitespace of content",
			rec:  MemoryRecord{Content: "  Padded  Text \t"},
			want: "padded text",
		},
		{
			name: "blank hash falls through to content",
			rec:  MemoryRecord{NormalizedContentHash: "   ", Content: "Fallback"},
			want: "fallback",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeDedupeContent(tc.rec); got != tc.want {
				t.Fatalf("NormalizeDedupeContent = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDedupeKey_Deterministic(t *testing.T) {
	rec := MemoryRecord{TenantID: "t", UserID: "u", Category: "c", ProjectID: "p", Content: "hello world"}
	if DedupeKey(rec) != DedupeKey(rec) {
		t.Fatal("DedupeKey is not deterministic")
	}
}

// TestDedupeKey_Golden locks the exact key construction (part order, NUL
// separator, sha256 lowercase-hex) so a refactor cannot silently change the
// value — which would split the worker and gateway onto different dedupe
// identities. Golden values precomputed independently:
//
//	sha256("t\x00u\x00c\x00p\x00hello world")
//	sha256("tn\x00us\x00cat\x00proj\x00precomputed-hash")
func TestDedupeKey_Golden(t *testing.T) {
	const (
		goldenFromContent = "1d91863f3a5722a55a89e915a18379460cc78549db2a3e2f341243067e309bf4"
		goldenFromHash    = "222720a34cd9a87392c22f898359518a4e1d91191810e5806416e84eebe3f1b3"
	)
	if got := DedupeKey(MemoryRecord{TenantID: "t", UserID: "u", Category: "c", ProjectID: "p", Content: "Hello  World"}); got != goldenFromContent {
		t.Fatalf("DedupeKey(content) = %q, want golden %q", got, goldenFromContent)
	}
	if got := DedupeKey(MemoryRecord{TenantID: "tn", UserID: "us", Category: "cat", ProjectID: "proj", NormalizedContentHash: "precomputed-hash"}); got != goldenFromHash {
		t.Fatalf("DedupeKey(hash) = %q, want golden %q", got, goldenFromHash)
	}
}

func TestDedupeKey_FieldsAffectKey(t *testing.T) {
	base := MemoryRecord{TenantID: "t", UserID: "u", Category: "c", ProjectID: "p", Content: "same"}
	baseKey := DedupeKey(base)
	mutations := map[string]MemoryRecord{
		"tenant":  {TenantID: "t2", UserID: "u", Category: "c", ProjectID: "p", Content: "same"},
		"user":    {TenantID: "t", UserID: "u2", Category: "c", ProjectID: "p", Content: "same"},
		"category": {TenantID: "t", UserID: "u", Category: "c2", ProjectID: "p", Content: "same"},
		"project": {TenantID: "t", UserID: "u", Category: "c", ProjectID: "p2", Content: "same"},
		"content": {TenantID: "t", UserID: "u", Category: "c", ProjectID: "p", Content: "different"},
	}
	for field, rec := range mutations {
		if DedupeKey(rec) == baseKey {
			t.Fatalf("changing %s did not change the dedupe key", field)
		}
	}
}
