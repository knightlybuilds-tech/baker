package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func note(verified, body string) string {
	return "---\ntype: lesson\nas_of: 2026-10-01\nlast_verified: " + verified + "\nsource: test\n---\n" + body
}

func TestCheckFindsProblems(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "MEMORY.md", "- [A](lessons/a.md) — ok\n- [Gone](lessons/gone.md) — broken\n")
	write(t, dir, "lessons/a.md", note("2026-10-01", "# A\nSee [[b]] and [[missing]]."))
	write(t, dir, "lessons/b.md", note("2025-01-01", "# B\n"+strings.Repeat("word ", 450)))
	write(t, dir, "knowledge/c.md", "# no frontmatter\n")

	now := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	got := ""
	for _, f := range Check(dir, now) {
		got += f.Msg + "\n"
	}
	for _, want := range []string{
		"links to lessons/gone.md",
		"lessons/b.md is not in MEMORY.md",
		"knowledge/c.md is missing frontmatter field \"type\"",
		"lessons/b.md has 452 words",
		"lessons/b.md last verified 2025-01-01",
		"links to [[missing]]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing finding %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "[[b]]") {
		t.Errorf("[[b]] exists and should resolve:\n%s", got)
	}
}

func TestHubNotesCountAsIndex(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "MEMORY.md", "- [Billing](knowledge/_billing.md) — hub\n")
	write(t, dir, "knowledge/_billing.md", note("2026-10-01", "# Billing\n- [[tenants-pay-monthly]]\n"))
	write(t, dir, "knowledge/tenants-pay-monthly.md", note("2026-10-01", "# Tenants pay monthly\n"))
	for _, f := range Check(dir, time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)) {
		if f.Level != OK {
			t.Errorf("unexpected finding: %s", f.Msg)
		}
	}
}
