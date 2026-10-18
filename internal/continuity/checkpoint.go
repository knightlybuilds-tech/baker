package continuity

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	CheckpointDir  = ".baker/checkpoints"
	keepCheckpoint = 20
	keepMessages   = 8
	maxMessageLen  = 800
)

// GitState is a mechanical snapshot of the repo.
type GitState struct {
	Branch     string   `json:"branch,omitempty"`
	Head       string   `json:"head,omitempty"`
	Status     []string `json:"status,omitempty"`
	DiffStat   string   `json:"diffstat,omitempty"`
	IsRepo     bool     `json:"is_repo"`
	StatusMore int      `json:"status_more,omitempty"`
}

type Checkpoint struct {
	Time    time.Time `json:"time"`
	Agent   string    `json:"agent,omitempty"`
	Git     GitState  `json:"git"`
	Session *Session  `json:"session,omitempty"`
}

// HookInput is the JSON some agents (Claude Code) pass to hooks on stdin.
type HookInput struct {
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
}

// Take snapshots git state and the latest transcript and saves it.
func Take(root, agent string, hook *HookInput) (*Checkpoint, string, error) {
	cp := &Checkpoint{Time: time.Now(), Agent: agent, Git: Git(root)}
	var s *Session
	if hook != nil && hook.TranscriptPath != "" && agent == "claude" {
		s, _ = ParseClaude(hook.TranscriptPath)
	}
	if s == nil {
		s, _ = LatestSession(root, agent)
	}
	if s != nil {
		cp.Session = trimSession(s)
		if cp.Agent == "" {
			cp.Agent = s.Agent
		}
	}
	dir := filepath.Join(root, CheckpointDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, "", err
	}
	b, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(dir, cp.Time.Format("20060102-150405")+".json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return nil, "", err
	}
	prune(dir)
	return cp, path, nil
}

// LatestCheckpoint loads the newest checkpoint, if any.
func LatestCheckpoint(root string) *Checkpoint {
	files := newestFiles(filepath.Join(root, CheckpointDir), ".json", 1)
	if len(files) == 0 {
		return nil
	}
	b, err := os.ReadFile(files[0])
	if err != nil {
		return nil
	}
	var cp Checkpoint
	if json.Unmarshal(b, &cp) != nil {
		return nil
	}
	return &cp
}

// trimSession keeps only the tail of the conversation, redacted.
func trimSession(s *Session) *Session {
	out := *s
	msgs := s.Messages
	if len(msgs) > keepMessages {
		tail := msgs[len(msgs)-keepMessages:]
		// Always keep the user's latest request, even when the agent talked a lot after it.
		hasUser := false
		for _, m := range tail {
			hasUser = hasUser || m.Role == "user"
		}
		if !hasUser {
			for i := len(msgs) - keepMessages - 1; i >= 0; i-- {
				if msgs[i].Role == "user" {
					tail = append([]Message{msgs[i], {Role: "agent", Text: "…"}}, tail...)
					break
				}
			}
		}
		msgs = tail
	}
	out.Messages = make([]Message, len(msgs))
	for i, m := range msgs {
		t := Redact(m.Text)
		if len(t) > maxMessageLen {
			t = t[:maxMessageLen] + " …[truncated]"
		}
		out.Messages[i] = Message{m.Role, t}
	}
	for i, t := range out.Todos {
		out.Todos[i].Content = Redact(t.Content)
	}
	return &out
}

func prune(dir string) {
	entries, _ := os.ReadDir(dir)
	var names []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for len(names) > keepCheckpoint {
		_ = os.Remove(filepath.Join(dir, names[0]))
		names = names[1:]
	}
}

// Git reads branch, head, status and diffstat. Missing git is not an error.
func Git(root string) GitState {
	run := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimRight(string(out), "\n")
	}
	if run("rev-parse", "--is-inside-work-tree") != "true" {
		return GitState{}
	}
	g := GitState{IsRepo: true}
	g.Branch = run("branch", "--show-current")
	g.Head = run("log", "-1", "--format=%h %s")
	if st := run("status", "--short"); st != "" {
		lines := strings.Split(st, "\n")
		if len(lines) > 30 {
			g.StatusMore = len(lines) - 30
			lines = lines[:30]
		}
		g.Status = lines
	}
	g.DiffStat = run("diff", "--shortstat", "HEAD")
	return g
}
