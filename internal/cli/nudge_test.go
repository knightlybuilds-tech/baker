package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func nudgeFixture(t *testing.T, manifest string) (root, transcript string) {
	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "baker.toml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, f := range []string{"a.go", "b.go", "c.go"} {
		lines = append(lines, fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Edit","input":{"file_path":%q}}]}}`, filepath.Join(root, f)))
	}
	transcript = filepath.Join(root, "t.jsonl")
	if err := os.WriteFile(transcript, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, transcript
}

func hookInput(root, transcript string, active bool) *strings.Reader {
	return strings.NewReader(fmt.Sprintf(`{"session_id":"s1","transcript_path":%q,"cwd":%q,"stop_hook_active":%v}`, transcript, root, active))
}

func TestRunNudge(t *testing.T) {
	t.Setenv("BAKER_NO_NUDGE", "")
	root, tr := nudgeFixture(t, "name = 'x'\n")
	if out := runNudge("claude", hookInput(root, tr, true)); out != nil {
		t.Error("never nudge while stop_hook_active (loop guard)")
	}
	if out := runNudge("codex", hookInput(root, tr, false)); out != nil {
		t.Error("only claude supports the nudge reply")
	}
	out := runNudge("claude", hookInput(root, tr, false))
	if out == nil || out.Decision != "block" || !strings.Contains(out.Reason, "baker memory add") {
		t.Fatalf("want a block decision with instructions, got %+v", out)
	}
	if again := runNudge("claude", hookInput(root, tr, false)); again != nil {
		t.Error("don't nudge twice for the same work")
	}
	if out := runNudge("claude", strings.NewReader("garbage")); out != nil {
		t.Error("bad input must stay silent")
	}

	off, tr2 := nudgeFixture(t, "name = 'x'\nnudge = false\n")
	if out := runNudge("claude", hookInput(off, tr2, false)); out != nil {
		t.Error("nudge = false in baker.toml must silence it")
	}
}
