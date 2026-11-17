package continuity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTranscript(t *testing.T, lines ...string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "s.jsonl")
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func toolUse(name, key, value string) string {
	return fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":%q,"input":{%q:%q}}]}}`, name, key, value)
}

func TestClaudeWorkCountsEditsAndCrumbs(t *testing.T) {
	root := "/repo"
	tr := writeTranscript(t,
		toolUse("Edit", "file_path", "/repo/a.go"),
		toolUse("Write", "file_path", "/repo/a.go"),                      // same file counts once
		toolUse("MultiEdit", "file_path", "b.go"),                        // relative to root
		toolUse("Edit", "file_path", "/repo/memory/hot.md"),              // memory/ but not a crumb
		toolUse("Write", "file_path", "/repo/memory/lessons/x.md"),       // crumb by hand
		toolUse("Bash", "command", "baker memory add --kind lesson ..."), // crumb via CLI
		toolUse("Bash", "command", "go test ./..."),
		`{"type":"assistant","isSidechain":true,"message":{"content":[{"type":"tool_use","name":"Edit","input":{"file_path":"/repo/sub.go"}}]}}`,
		`not json`,
	)
	w := ClaudeWork(tr, root)
	if w.Edited != 2 || w.Crumbs != 2 {
		t.Errorf("got %+v, want 2 edited files and 2 crumbs", w)
	}
}

func TestShouldNudge(t *testing.T) {
	steps := []struct {
		name         string
		work         Work
		wantNudge    bool
		wantBaseline *NudgeState // nil = keep the stored state
	}{
		{"quiet start", Work{Edited: 1}, false, nil},
		{"enough work, no crumbs", Work{Edited: 3}, true, &NudgeState{3, 0}},
		{"same work again", Work{Edited: 3}, false, nil},
		{"a bit more work", Work{Edited: 10}, false, nil},
		{"agent saved a crumb", Work{Edited: 10, Crumbs: 1}, false, &NudgeState{10, 1}},
		{"lots more work, no new crumb", Work{Edited: 18, Crumbs: 1}, true, &NudgeState{18, 1}},
	}
	var state *NudgeState
	for _, s := range steps {
		nudge, save := ShouldNudge(s.work, state)
		if nudge != s.wantNudge {
			t.Errorf("%s: nudge=%v want %v", s.name, nudge, s.wantNudge)
		}
		if (save == nil) != (s.wantBaseline == nil) || (save != nil && *save != *s.wantBaseline) {
			t.Errorf("%s: save=%v want %v", s.name, save, s.wantBaseline)
		}
		if save != nil {
			state = save
		}
	}

	if nudge, save := ShouldNudge(Work{Edited: 5, Crumbs: 1}, nil); nudge || save == nil {
		t.Error("a session that already saved a crumb shouldn't be nudged, but its baseline should be stored")
	}
}

func TestNudgeStateRoundTrip(t *testing.T) {
	root := t.TempDir()
	if LoadNudgeState(root, "a/b") != nil {
		t.Fatal("want nil before saving")
	}
	SaveNudgeState(root, "a/b", NudgeState{Edited: 4, Crumbs: 1})
	if got := LoadNudgeState(root, "a/b"); got == nil || *got != (NudgeState{4, 1}) {
		t.Errorf("got %v", got)
	}
	if _, err := os.Stat(filepath.Join(root, ".baker", "nudge", "a_b.json")); err != nil {
		t.Errorf("session id should be sanitized into the filename: %v", err)
	}
}
