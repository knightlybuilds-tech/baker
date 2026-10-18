package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bake(t *testing.T, root string, force bool) *Plan {
	t.Helper()
	proj, r, err := Prepare(root, "", []string{"claude", "codex", "gemini", "cursor"})
	if err != nil {
		t.Fatal(err)
	}
	p := NewPlan(root, force)
	if err := Build(p, proj, r); err != nil {
		t.Fatal(err)
	}
	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}
	return p
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBuildIsIdempotentAndProtectsEdits(t *testing.T) {
	t.Setenv("BAKER_HOME", t.TempDir())
	t.Setenv("BAKER_USER_HOME", t.TempDir())
	root := t.TempDir()

	bake(t, root, false)
	for _, f := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md", "baker.toml", "memory/MEMORY.md", "memory/hot.md",
		".agents/skills/session-close/SKILL.md", ".claude/skills/session-close/SKILL.md", ".claude/settings.json", ".gitignore"} {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("missing %s", f)
		}
	}
	if !strings.Contains(readFile(t, root, "CLAUDE.md"), "@AGENTS.md") {
		t.Error("CLAUDE.md doesn't import AGENTS.md")
	}

	// Second run changes nothing.
	proj, r, _ := Prepare(root, "", nil)
	p := NewPlan(root, false)
	_ = Build(p, proj, r)
	if p.HasWrites() {
		t.Errorf("second build wanted writes: %+v", p.Changes())
	}

	// User edits to notes survive; an edited managed block is skipped, then restored by --force.
	agents := readFile(t, root, "AGENTS.md")
	agents = strings.Replace(agents, "## Project notes", "## Project notes\nThis is a shop.", 1)
	agents = strings.Replace(agents, "Any agent can read and write it.", "EDITED", 1)
	_ = os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(agents), 0o644)

	bake(t, root, false)
	got := readFile(t, root, "AGENTS.md")
	if !strings.Contains(got, "This is a shop.") || !strings.Contains(got, "EDITED") {
		t.Error("user edits were overwritten without --force")
	}
	bake(t, root, true)
	got = readFile(t, root, "AGENTS.md")
	if !strings.Contains(got, "This is a shop.") || strings.Contains(got, "EDITED") {
		t.Error("--force should restore the block but keep user notes")
	}
}

func TestInitAdoptsExistingClaudeMD(t *testing.T) {
	t.Setenv("BAKER_HOME", t.TempDir())
	t.Setenv("BAKER_USER_HOME", t.TempDir())
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("Always use Pest.\n"), 0o644)

	bake(t, root, false)
	if !strings.Contains(readFile(t, root, "AGENTS.md"), "Always use Pest.") {
		t.Error("CLAUDE.md content not moved into AGENTS.md")
	}
	if strings.Contains(readFile(t, root, "CLAUDE.md"), "Always use Pest.") {
		t.Error("CLAUDE.md still holds the old content (would load twice)")
	}
}
