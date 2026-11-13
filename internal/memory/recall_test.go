package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecallOrdinaryNotesWithoutApproval(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "lessons"), 0755); err != nil {
		t.Fatal(err)
	}
	body := "---\ntype: lesson\nas_of: 2026-11-05\nlast_verified: 2026-11-05\nsource: local investigation\n---\n# Preserve leading zeros\n\nIdentifiers are strings.\n"
	if err := os.WriteFile(filepath.Join(dir, "lessons", "codes.md"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	got := Recall(dir, "LEADING ZEROS", now)
	if len(got) != 1 || len(got[0].Issues) != 0 || got[0].Source != "local investigation" || !strings.Contains(got[0].Body, "Identifiers") {
		t.Fatalf("ordinary note not recalled: %+v", got)
	}
	if len(Recall(dir, "unrelated", now)) != 0 {
		t.Fatal("query did not filter")
	}
	got = Recall(dir, "", now.Add(91*24*time.Hour))
	if len(got) != 1 || len(got[0].Issues) == 0 || got[0].Body == "" {
		t.Fatal("stale notes should be flagged and remain historical leads")
	}
}

func TestRecallFlagsMissingMetadata(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "knowledge"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "knowledge", "old.md"), []byte("# Legacy knowledge\nStill useful as a lead."), 0644); err != nil {
		t.Fatal(err)
	}
	got := Recall(dir, "", time.Now())
	if len(got) != 1 || len(got[0].Issues) != len(RequiredFields) {
		t.Fatalf("missing metadata not reported: %+v", got)
	}
}
