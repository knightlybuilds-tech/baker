package continuity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRedact(t *testing.T) {
	in := `key sk-ant-api03-abcdefghijklmnopqrstuv and STRIPE_SECRET=whsec_123456 and ghp_abcdefghijklmnopqrstuvwxyz123 and password: "hunter2"`
	out := Redact(in)
	for _, leak := range []string{"abcdefghijklmnopqrstuv", "whsec_123456", "ghp_abc", "hunter2"} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %q in %q", leak, out)
		}
	}
	if !strings.Contains(out, "STRIPE_SECRET=[REDACTED]") {
		t.Errorf("key name should stay readable: %q", out)
	}
}

func TestParseClaude(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "s.jsonl")
	lines := []string{
		`{"type":"user","message":{"role":"user","content":"add billing<system-reminder>noise</system-reminder>"}}`,
		`{"type":"user","isMeta":true,"message":{"role":"user","content":"meta"}}`,
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"On it."},{"type":"tool_use","name":"Edit","input":{"file_path":"/p/app/Billing.php"}},{"type":"tool_use","name":"TodoWrite","input":{"todos":[{"content":"write job","status":"in_progress"}]}}]}}`,
		`{"type":"user","message":{"role":"user","content":[{"type":"tool_result","content":"ok"}]}}`,
		`not json`,
	}
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := ParseClaude(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Messages) != 2 || s.Messages[0].Text != "add billing" || s.Messages[1].Text != "On it." {
		t.Fatalf("messages: %+v", s.Messages)
	}
	if len(s.Todos) != 1 || s.Todos[0].Status != "in_progress" {
		t.Fatalf("todos: %+v", s.Todos)
	}
	if len(s.FilesTouched) != 1 || s.FilesTouched[0] != "/p/app/Billing.php" {
		t.Fatalf("files: %+v", s.FilesTouched)
	}
}

func TestParseCodex(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "rollout.jsonl")
	lines := []string{
		`{"type":"session_meta","payload":{"cwd":"/p"}}`,
		`{"type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"<environment_context>x</environment_context>"}]}}`,
		`{"type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"fix the test"}]}}`,
		`{"type":"response_item","payload":{"type":"message","role":"assistant","content":[{"type":"output_text","text":"Fixed."}]}}`,
	}
	_ = os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644)
	if !fileMentionsCwd(p, "/p") {
		t.Fatal("cwd not detected")
	}
	s, _ := ParseCodex(p)
	if len(s.Messages) != 2 || s.Messages[0].Text != "fix the test" || s.Messages[1].Role != "agent" {
		t.Fatalf("messages: %+v", s.Messages)
	}
}

func TestHotIsEmpty(t *testing.T) {
	if !hotIsEmpty("# Goal\n\n# Done\n\n# Next\n") {
		t.Error("template should count as empty")
	}
	if hotIsEmpty("# Goal\nShip billing\n") {
		t.Error("filled hot.md should not be empty")
	}
}

func TestTrimKeepsLatestUserMessage(t *testing.T) {
	s := &Session{Messages: []Message{{"user", "build the thing"}}}
	for i := 0; i < 20; i++ {
		s.Messages = append(s.Messages, Message{"agent", "working"})
	}
	got := trimSession(s).Messages
	if got[0].Role != "user" || got[0].Text != "build the thing" {
		t.Fatalf("latest user request dropped: %+v", got[:2])
	}
}
