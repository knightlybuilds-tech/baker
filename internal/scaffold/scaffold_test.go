package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
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

func TestGeminiCheckpointHookMergesSettings(t *testing.T) {
	t.Setenv("BAKER_HOME", t.TempDir())
	t.Setenv("BAKER_USER_HOME", t.TempDir())
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".gemini"), 0o755)
	_ = os.WriteFile(filepath.Join(root, ".gemini/settings.json"), []byte(`{"theme": "Dracula"}`), 0o644)

	bake(t, root, false)
	var doc struct {
		Theme string `json:"theme"`
		Hooks map[string][]struct {
			Hooks []struct{ Type, Command string } `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(readFile(t, root, ".gemini/settings.json")), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Theme != "Dracula" {
		t.Error("user settings were dropped")
	}
	for _, event := range []string{"AfterAgent", "PreCompress"} {
		list := doc.Hooks[event]
		if len(list) != 1 || len(list[0].Hooks) != 1 || !strings.Contains(list[0].Hooks[0].Command, "baker checkpoint --agent gemini --quiet") {
			t.Errorf("%s hook: %+v", event, list)
		}
	}
	if strings.Contains(readFile(t, root, ".claude/settings.json"), "--agent gemini") {
		t.Error("Claude hook should name claude")
	}

	// A second build doesn't add the hook twice.
	bake(t, root, false)
	if n := strings.Count(readFile(t, root, ".gemini/settings.json"), "baker checkpoint"); n != 2 {
		t.Errorf("want 2 hooks after rebuild, got %d", n)
	}
}

func TestCodexNotifyHookIsValidTOML(t *testing.T) {
	var doc struct{ Notify []string }
	if err := toml.Unmarshal([]byte(codexHook), &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Notify) != 3 || doc.Notify[0] != "sh" || !strings.Contains(doc.Notify[2], `--quiet "$0"`) {
		t.Fatalf("notify: %q", doc.Notify)
	}
}
