// Package continuity lets work survive a session ending abruptly (usage
// limits, switching accounts or providers). It reads the local transcripts
// agents already keep, takes mechanical checkpoints, and builds a resume
// brief any agent can start from. No LLM is involved.
package continuity

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/roguenet/baker/internal/config"
)

type Message struct {
	Role string // "user" or "agent"
	Text string
}

type Todo struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

// Session is what baker could recover from one transcript.
type Session struct {
	Agent        string    `json:"agent"`
	File         string    `json:"file"`
	Updated      time.Time `json:"updated"`
	Messages     []Message `json:"messages"`
	Todos        []Todo    `json:"todos,omitempty"`
	FilesTouched []string  `json:"files_touched,omitempty"`
}

// LatestSession returns the most recent session for root across all agents,
// or only for agent `from` when set.
func LatestSession(root, from string) (*Session, error) {
	g := config.LoadGlobal()
	var best *Session
	for _, agent := range []string{"claude", "codex", "gemini"} {
		if from != "" && from != agent {
			continue
		}
		dirs := g.Transcripts[agent]
		if len(dirs) == 0 {
			dirs = config.DefaultGlobal().Transcripts[agent]
		}
		for _, d := range dirs {
			s, err := findSession(agent, config.Expand(d), root)
			if err != nil || s == nil {
				continue
			}
			if best == nil || s.Updated.After(best.Updated) {
				best = s
			}
		}
	}
	return best, nil
}

func findSession(agent, dir, root string) (*Session, error) {
	switch agent {
	case "claude":
		return claudeSession(dir, root)
	case "codex":
		return codexSession(dir, root)
	case "gemini":
		return geminiSession(dir, root)
	}
	return nil, nil
}

// ---- Claude Code: <dir>/projects/<path-with-dashes>/<session>.jsonl ----

var nonAlnum = regexp.MustCompile(`[^A-Za-z0-9]`)

func claudeSession(dir, root string) (*Session, error) {
	projects := filepath.Join(dir, "projects")
	candidates := newestFiles(filepath.Join(projects, nonAlnum.ReplaceAllString(root, "-")), ".jsonl", 1)
	if len(candidates) == 0 {
		// Fall back to scanning for a transcript whose cwd is root.
		for _, f := range newestFiles(projects, ".jsonl", 200) {
			if fileMentionsCwd(f, root) {
				candidates = []string{f}
				break
			}
		}
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	return ParseClaude(candidates[0])
}

var noiseTags = regexp.MustCompile(`(?s)<(system-reminder|command-[a-z-]+|local-command-[a-z-]+|ide_[a-z_]+)[^>]*>.*?</[^>]+>`)

// ParseClaude reads a Claude Code transcript.
func ParseClaude(path string) (*Session, error) {
	s := &Session{Agent: "claude", File: path}
	touched := map[string]bool{}
	err := eachJSONLine(path, func(line map[string]any) {
		if line["isMeta"] == true || line["isSidechain"] == true {
			return
		}
		typ, _ := line["type"].(string)
		if typ != "user" && typ != "assistant" {
			return
		}
		msg, _ := line["message"].(map[string]any)
		var texts []string
		switch c := msg["content"].(type) {
		case string:
			texts = append(texts, c)
		case []any:
			for _, b := range c {
				blk, _ := b.(map[string]any)
				switch blk["type"] {
				case "text":
					t, _ := blk["text"].(string)
					texts = append(texts, t)
				case "tool_use":
					name, _ := blk["name"].(string)
					input, _ := blk["input"].(map[string]any)
					if name == "TodoWrite" {
						s.Todos = parseTodos(input["todos"])
					}
					if fp, ok := input["file_path"].(string); ok && (name == "Edit" || name == "Write" || name == "MultiEdit" || name == "NotebookEdit") {
						touched[fp] = true
					}
				}
			}
		}
		text := strings.TrimSpace(noiseTags.ReplaceAllString(strings.Join(texts, "\n"), ""))
		if text == "" {
			return
		}
		role := "agent"
		if typ == "user" {
			role = "user"
		}
		s.Messages = append(s.Messages, Message{role, text})
	})
	s.FilesTouched = sortedKeys(touched)
	s.Updated = modTime(path)
	return s, err
}

// ---- Codex: <dir>/sessions/YYYY/MM/DD/rollout-*.jsonl ----

func codexSession(dir, root string) (*Session, error) {
	for _, f := range newestFiles(filepath.Join(dir, "sessions"), ".jsonl", 100) {
		if fileMentionsCwd(f, root) {
			return ParseCodex(f)
		}
	}
	return nil, nil
}

// ParseCodex reads a Codex CLI rollout file.
func ParseCodex(path string) (*Session, error) {
	s := &Session{Agent: "codex", File: path}
	touched := map[string]bool{}
	err := eachJSONLine(path, func(line map[string]any) {
		item := line
		if p, ok := line["payload"].(map[string]any); ok {
			item = p
		}
		switch item["type"] {
		case "function_call", "custom_tool_call":
			codexToolCall(s, touched, item)
			return
		case "message":
		default:
			return
		}
		role, _ := item["role"].(string)
		if role != "user" && role != "assistant" {
			return
		}
		var texts []string
		content, _ := item["content"].([]any)
		for _, b := range content {
			blk, _ := b.(map[string]any)
			if t, ok := blk["text"].(string); ok {
				texts = append(texts, t)
			}
		}
		text := strings.TrimSpace(strings.Join(texts, "\n"))
		if text == "" || strings.HasPrefix(text, "<environment_context>") || strings.HasPrefix(text, "<user_instructions>") || strings.HasPrefix(text, "# AGENTS.md") {
			return
		}
		if role == "assistant" {
			role = "agent"
		}
		s.Messages = append(s.Messages, Message{role, text})
	})
	s.FilesTouched = sortedKeys(touched)
	s.Updated = modTime(path)
	return s, err
}

var patchFile = regexp.MustCompile(`(?m)^\*\*\* (?:Add|Update|Delete) File: (.+)$`)

// codexToolCall picks up the plan (update_plan) and edited files (apply_patch).
// apply_patch is a custom tool with the patch as raw input; older versions sent
// it as a function call with JSON arguments.
func codexToolCall(s *Session, touched map[string]bool, item map[string]any) {
	name, _ := item["name"].(string)
	raw, _ := item["input"].(string)
	if args, ok := item["arguments"].(string); ok {
		raw = args
	}
	switch name {
	case "update_plan":
		var args struct {
			Plan []struct {
				Step   string `json:"step"`
				Status string `json:"status"`
			} `json:"plan"`
		}
		if json.Unmarshal([]byte(raw), &args) == nil {
			s.Todos = nil
			for _, p := range args.Plan {
				s.Todos = append(s.Todos, Todo{Content: p.Step, Status: p.Status})
			}
		}
	case "apply_patch":
		var args struct {
			Input string `json:"input"`
		}
		if json.Unmarshal([]byte(raw), &args) == nil && args.Input != "" {
			raw = args.Input
		}
		for _, m := range patchFile.FindAllStringSubmatch(raw, -1) {
			touched[strings.TrimSpace(m[1])] = true
		}
	}
}

// ---- Gemini CLI: <dir>/tmp/<project-id>/chats/session-*.jsonl ----
//
// Current versions name the project dir with a short slug recorded in
// <dir>/projects.json and in <dir>/tmp/<slug>/.project_root. Older versions
// used sha256(root) and wrote one session-*.json document per chat.

func geminiSession(dir, root string) (*Session, error) {
	var files []string
	for _, d := range geminiProjectDirs(dir, root) {
		for _, f := range newestFiles(filepath.Join(d, "chats"), "", 50) {
			// Subagent chats live in subdirectories and aren't named session-*.
			base := filepath.Base(f)
			if strings.HasPrefix(base, "session-") && (strings.HasSuffix(base, ".jsonl") || strings.HasSuffix(base, ".json")) {
				files = append(files, f)
			}
		}
	}
	if len(files) == 0 {
		return nil, nil
	}
	sort.Slice(files, func(i, j int) bool { return modTime(files[i]).After(modTime(files[j])) })
	return ParseGemini(files[0])
}

// geminiProjectDirs lists the temp dirs Gemini may have used for root.
func geminiProjectDirs(dir, root string) []string {
	tmp := filepath.Join(dir, "tmp")
	seen := map[string]bool{}
	var out []string
	add := func(name string) {
		if name != "" && !seen[name] {
			seen[name] = true
			out = append(out, filepath.Join(tmp, name))
		}
	}
	var reg struct {
		Projects map[string]string `json:"projects"`
	}
	if b, err := os.ReadFile(filepath.Join(dir, "projects.json")); err == nil && json.Unmarshal(b, &reg) == nil {
		add(reg.Projects[root])
	}
	entries, _ := os.ReadDir(tmp)
	for _, e := range entries {
		if b, err := os.ReadFile(filepath.Join(tmp, e.Name(), ".project_root")); err == nil && strings.TrimSpace(string(b)) == root {
			add(e.Name())
		}
	}
	sum := sha256.Sum256([]byte(root))
	add(hex.EncodeToString(sum[:]))
	return out
}

type geminiMessage struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content any    `json:"content"`
}

// ParseGemini reads a Gemini CLI chat: the current JSONL log (a metadata line,
// message records, and $set / $rewindTo updates; $patch edits are ignored) or
// a legacy JSON document.
func ParseGemini(path string) (*Session, error) {
	s := &Session{Agent: "gemini", File: path, Updated: modTime(path)}
	var msgs []geminiMessage
	if strings.HasSuffix(path, ".json") {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var doc struct {
			Messages []geminiMessage `json:"messages"`
		}
		if err := json.Unmarshal(b, &doc); err != nil {
			return nil, err
		}
		msgs = doc.Messages
	} else {
		err := eachJSONLine(path, func(line map[string]any) {
			if set, ok := line["$set"].(map[string]any); ok {
				if list, ok := set["messages"].([]any); ok {
					msgs = nil
					for _, m := range list {
						msgs = append(msgs, toGeminiMessage(m))
					}
				}
				return
			}
			if id, ok := line["$rewindTo"].(string); ok {
				for i, m := range msgs {
					if m.ID == id {
						msgs = msgs[:i]
						break
					}
				}
				return
			}
			if _, ok := line["type"]; ok {
				m := toGeminiMessage(line)
				// A record with a known id replaces the earlier version.
				for i := range msgs {
					if m.ID != "" && msgs[i].ID == m.ID {
						msgs[i] = m
						return
					}
				}
				msgs = append(msgs, m)
			}
		})
		if err != nil {
			return nil, err
		}
	}
	for _, m := range msgs {
		if m.Type != "user" && m.Type != "gemini" {
			continue // info, error, warning
		}
		text := strings.TrimSpace(geminiText(m.Content))
		if text == "" {
			continue
		}
		role := "agent"
		if m.Type == "user" {
			role = "user"
		}
		s.Messages = append(s.Messages, Message{role, text})
	}
	return s, nil
}

func toGeminiMessage(v any) geminiMessage {
	var m geminiMessage
	b, _ := json.Marshal(v)
	_ = json.Unmarshal(b, &m)
	return m
}

// geminiText flattens a PartListUnion: a string, a part, or a list of either.
func geminiText(v any) string {
	switch c := v.(type) {
	case string:
		return c
	case map[string]any:
		t, _ := c["text"].(string)
		return t
	case []any:
		var parts []string
		for _, p := range c {
			if t := geminiText(p); t != "" {
				parts = append(parts, t)
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

// ---- helpers ----

func eachJSONLine(path string, fn func(map[string]any)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for sc.Scan() {
		var line map[string]any
		if json.Unmarshal(sc.Bytes(), &line) == nil {
			fn(line)
		}
	}
	return sc.Err()
}

// fileMentionsCwd checks the first lines of a transcript for "cwd": root.
func fileMentionsCwd(path, root string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	needle, _ := json.Marshal(root)
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for i := 0; i < 50 && sc.Scan(); i++ {
		line := sc.Text()
		if strings.Contains(line, `"cwd":`+string(needle)) || strings.Contains(line, `"cwd": `+string(needle)) {
			return true
		}
	}
	return false
}

// newestFiles returns up to n files with ext under dir, newest first.
func newestFiles(dir, ext string, n int) []string {
	type fm struct {
		path string
		mod  time.Time
	}
	var all []fm
	_ = filepath.WalkDir(dir, func(p string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(p, ext) {
			return nil
		}
		if info, err := e.Info(); err == nil {
			all = append(all, fm{p, info.ModTime()})
		}
		return nil
	})
	sort.Slice(all, func(i, j int) bool { return all[i].mod.After(all[j].mod) })
	out := []string{}
	for i := 0; i < len(all) && i < n; i++ {
		out = append(out, all[i].path)
	}
	return out
}

func parseTodos(v any) []Todo {
	b, _ := json.Marshal(v)
	var t []Todo
	_ = json.Unmarshal(b, &t)
	return t
}

func modTime(p string) time.Time {
	if info, err := os.Stat(p); err == nil {
		return info.ModTime()
	}
	return time.Time{}
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
