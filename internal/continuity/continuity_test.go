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

func TestParseCodexToolCalls(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "rollout.jsonl")
	lines := []string{
		`{"type":"session_meta","payload":{"cwd":"/p","cli_version":"0.115.0"}}`,
		`{"type":"response_item","payload":{"type":"message","role":"developer","content":[{"type":"input_text","text":"sandbox rules"}]}}`,
		`{"type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"# AGENTS.md instructions for /p"}]}}`,
		`{"type":"response_item","payload":{"type":"function_call","name":"update_plan","arguments":"{\"plan\":[{\"step\":\"write job\",\"status\":\"in_progress\"},{\"step\":\"test\",\"status\":\"pending\"}]}","call_id":"1"}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call","name":"apply_patch","input":"*** Begin Patch\n*** Update File: app/Billing.php\n@@\n*** Add File: tests/BillingTest.php\n*** End Patch","call_id":"2"}}`,
		`{"type":"response_item","payload":{"type":"function_call","name":"apply_patch","arguments":"{\"input\":\"*** Begin Patch\\n*** Delete File: old.php\\n*** End Patch\"}","call_id":"3"}}`,
	}
	_ = os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644)
	s, err := ParseCodex(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Messages) != 0 {
		t.Errorf("developer and AGENTS.md messages should be skipped: %+v", s.Messages)
	}
	if len(s.Todos) != 2 || s.Todos[0].Content != "write job" || s.Todos[0].Status != "in_progress" {
		t.Errorf("todos: %+v", s.Todos)
	}
	want := []string{"app/Billing.php", "old.php", "tests/BillingTest.php"}
	if strings.Join(s.FilesTouched, ",") != strings.Join(want, ",") {
		t.Errorf("files: %+v", s.FilesTouched)
	}
}

func TestParseGeminiJSONL(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "session-2026-10-18T10-00-abcd1234.jsonl")
	lines := []string{
		`{"sessionId":"abcd1234","projectHash":"x","startTime":"2026-10-18T10:00:00Z","kind":"main"}`,
		`{"id":"u1","timestamp":"t","type":"user","content":[{"text":"add billing"}]}`,
		`{"id":"g1","timestamp":"t","type":"gemini","content":"Draft one."}`,
		`{"$set":{"lastUpdated":"t"}}`,
		`{"id":"g1","timestamp":"t","type":"gemini","content":"On it."}`,
		`{"id":"i1","timestamp":"t","type":"info","content":"Model switched"}`,
		`{"id":"u2","timestamp":"t","type":"user","content":"never mind"}`,
		`{"$rewindTo":"u2"}`,
	}
	_ = os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644)
	s, err := ParseGemini(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Messages) != 2 || s.Messages[0] != (Message{"user", "add billing"}) || s.Messages[1] != (Message{"agent", "On it."}) {
		t.Fatalf("messages: %+v", s.Messages)
	}
}

func TestParseGeminiLegacyJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "session-old.json")
	_ = os.WriteFile(p, []byte(`{"sessionId":"x","messages":[{"id":"1","type":"user","content":"hi"},{"id":"2","type":"gemini","content":[{"text":"hello"}]}]}`), 0o644)
	s, err := ParseGemini(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Messages) != 2 || s.Messages[1].Text != "hello" {
		t.Fatalf("messages: %+v", s.Messages)
	}
}

func TestGeminiSessionFindsProjectDir(t *testing.T) {
	write := func(p, body string) {
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(body), 0o644)
	}
	chat := `{"sessionId":"s"}` + "\n" + `{"id":"1","type":"user","content":"from %s"}`

	// Current layout: slug dir with a .project_root marker; subagent chats are ignored.
	dir := t.TempDir()
	write(filepath.Join(dir, "tmp/shop/.project_root"), "/p/shop\n")
	write(filepath.Join(dir, "tmp/shop/chats/session-2026-10-18T10-00-aaaa.jsonl"), strings.Replace(chat, "%s", "marker", 1))
	write(filepath.Join(dir, "tmp/shop/chats/parent-id/sub.jsonl"), strings.Replace(chat, "%s", "subagent", 1))
	if s, _ := geminiSession(dir, "/p/shop"); s == nil || s.Messages[0].Text != "from marker" {
		t.Fatalf("marker lookup: %+v", s)
	}

	// Registry only.
	dir = t.TempDir()
	write(filepath.Join(dir, "projects.json"), `{"projects":{"/p/shop":"shop-1"}}`)
	write(filepath.Join(dir, "tmp/shop-1/chats/session-x.jsonl"), strings.Replace(chat, "%s", "registry", 1))
	if s, _ := geminiSession(dir, "/p/shop"); s == nil || s.Messages[0].Text != "from registry" {
		t.Fatalf("registry lookup: %+v", s)
	}

	// Legacy sha256 dir; another project's chats aren't picked up.
	dir = t.TempDir()
	write(filepath.Join(dir, "tmp/a7f0d9b5f8b4a1e3c8e1d8f0b0f6a7f6c0c2f1d8a6e5b4c3d2e1f0a9b8c7d6e5/chats/session-x.json"), `{"messages":[]}`)
	if s, _ := geminiSession(dir, "/p/shop"); s != nil {
		t.Fatalf("unrelated dir matched: %+v", s)
	}
	legacy := geminiProjectDirs(dir, "/p/shop")
	write(filepath.Join(legacy[len(legacy)-1], "chats/session-x.json"), `{"messages":[{"type":"user","content":"from hash"}]}`)
	if s, _ := geminiSession(dir, "/p/shop"); s == nil || s.Messages[0].Text != "from hash" {
		t.Fatalf("legacy lookup: %+v", s)
	}
}

func TestTakeUsesGeminiTranscriptPath(t *testing.T) {
	t.Setenv("BAKER_HOME", t.TempDir())
	t.Setenv("BAKER_USER_HOME", t.TempDir())
	root := t.TempDir()
	p := filepath.Join(t.TempDir(), "session-x.jsonl")
	_ = os.WriteFile(p, []byte(`{"id":"1","type":"user","content":"ship it"}`), 0o644)
	cp, _, err := Take(root, "gemini", &HookInput{TranscriptPath: p, Cwd: root})
	if err != nil {
		t.Fatal(err)
	}
	if cp.Session == nil || cp.Session.Agent != "gemini" || cp.Session.Messages[0].Text != "ship it" {
		t.Fatalf("session: %+v", cp.Session)
	}
}
