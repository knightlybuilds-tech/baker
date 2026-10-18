package continuity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/roguenet/baker/internal/memory"
)

// Brief builds the resume prompt from the three layers:
// memory/hot.md (written by the agent), the newest checkpoint (written by
// hooks) and the newest transcript (recovered from the agent's own logs).
func Brief(root, name, from string) (string, error) {
	var b strings.Builder
	w := func(f string, a ...any) { fmt.Fprintf(&b, f, a...) }

	session, _ := LatestSession(root, from)
	cp := LatestCheckpoint(root)
	// Prefer a checkpoint if it's newer than any transcript we can read
	// (e.g. the previous agent's transcript format isn't supported).
	if cp != nil && cp.Session != nil && (session == nil || cp.Time.After(session.Updated.Add(time.Minute))) {
		session = cp.Session
	}
	if session != nil {
		session = trimSession(session)
	}

	w("# Continuing work on %s\n\n", name)
	if session != nil {
		w("You're picking up a task that **%s** was working on (last activity %s). ", session.Agent, ago(session.Updated))
	} else {
		w("You're picking up a task another agent was working on. ")
	}
	w("Its session ended — likely a usage limit or a switch of account or provider.\n\n")
	w("Read `AGENTS.md` first. Sources are labeled: **hot.md** was written on purpose by the previous agent; **transcript** lines are raw and may stop mid-thought.\n\n")

	hot, _ := os.ReadFile(filepath.Join(root, "memory", "hot.md"))
	fm, body := memory.ParseFrontmatter(string(hot))
	w("## Handoff note (hot.md")
	if fm["updated"] != "" {
		w(", updated %s", fm["updated"])
	}
	w(")\n\n")
	if hotIsEmpty(body) {
		w("_Empty — the previous agent didn't write one. Rely on the repo state and transcript below._\n\n")
	} else {
		w("%s\n\n", strings.TrimSpace(Redact(stripHeadingLevel(body))))
	}

	g := Git(root)
	w("## Repo state (live)\n\n")
	if !g.IsRepo {
		w("_Not a git repository._\n\n")
	} else {
		w("- Branch: `%s`\n- Last commit: %s\n", orNone(g.Branch), orNone(g.Head))
		if len(g.Status) == 0 {
			w("- Working tree clean\n\n")
		} else {
			if g.DiffStat != "" {
				w("- Uncommitted: %s\n", strings.TrimSpace(g.DiffStat))
			}
			w("\n```\n%s\n", strings.Join(g.Status, "\n"))
			if g.StatusMore > 0 {
				w("… and %d more\n", g.StatusMore)
			}
			w("```\n\n")
		}
	}

	if session != nil {
		if len(session.Todos) > 0 {
			w("## Previous agent's todo list (transcript)\n\n")
			for _, t := range session.Todos {
				box := " "
				switch t.Status {
				case "completed":
					box = "x"
				case "in_progress":
					box = "~"
				}
				w("- [%s] %s\n", box, t.Content)
			}
			w("\n")
		}
		if len(session.FilesTouched) > 0 {
			w("## Files the previous agent edited (transcript)\n\n")
			for _, f := range session.FilesTouched {
				if rel, err := filepath.Rel(root, f); err == nil && !strings.HasPrefix(rel, "..") {
					f = rel
				}
				w("- `%s`\n", f)
			}
			w("\n")
		}
		if len(session.Messages) > 0 {
			w("## Last exchange before the session ended (transcript)\n\n")
			for _, m := range session.Messages {
				who := "Agent"
				if m.Role == "user" {
					who = "User"
				}
				w("**%s:** %s\n\n", who, quoteBlock(m.Text))
			}
		}
	}

	w("## Your first move\n\n")
	w("1. Read `AGENTS.md`, `memory/hot.md` and the files above.\n")
	w("2. Tell the user in one or two lines what you think the task is and what you'll do next. Wait for their go.\n")
	w("3. Continue from Next. Keep `memory/hot.md` updated after every finished step.\n")
	return b.String(), nil
}

// hotIsEmpty is true when hot.md only has the template headings.
func hotIsEmpty(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "#") && !strings.HasPrefix(t, "<!--") {
			return false
		}
	}
	return true
}

// stripHeadingLevel demotes "# X" to "### X" so hot.md nests under the brief's headings.
func stripHeadingLevel(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "#") {
			lines[i] = "##" + l
		}
	}
	return strings.Join(lines, "\n")
}

func quoteBlock(s string) string {
	s = strings.TrimSpace(s)
	if !strings.Contains(s, "\n") {
		return s
	}
	return "\n> " + strings.ReplaceAll(s, "\n", "\n> ")
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

func ago(t time.Time) string {
	if t.IsZero() {
		return "unknown time"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return t.Format("2006-01-02 15:04")
}
