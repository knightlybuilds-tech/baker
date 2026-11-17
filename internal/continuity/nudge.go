package continuity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Work summarizes what a Claude Code session did, for the capture nudge.
type Work struct {
	Edited int // distinct files edited outside memory/
	Crumbs int // crumbs saved: note files written, or `baker memory add` calls
}

var memoryAddCall = regexp.MustCompile(`\bbaker\s+memory\s+add\b`)

// ClaudeWork scans a Claude Code transcript for file edits and crumb writes.
func ClaudeWork(transcript, root string) Work {
	edited := map[string]bool{}
	var w Work
	noteDirs := []string{}
	for _, d := range []string{"decisions", "lessons", "knowledge"} {
		noteDirs = append(noteDirs, filepath.Join(root, "memory", d)+string(filepath.Separator))
	}
	memDir := filepath.Join(root, "memory") + string(filepath.Separator)
	_ = eachJSONLine(transcript, func(line map[string]any) {
		if line["type"] != "assistant" || line["isSidechain"] == true {
			return
		}
		msg, _ := line["message"].(map[string]any)
		blocks, _ := msg["content"].([]any)
		for _, b := range blocks {
			blk, _ := b.(map[string]any)
			if blk["type"] != "tool_use" {
				continue
			}
			name, _ := blk["name"].(string)
			input, _ := blk["input"].(map[string]any)
			switch name {
			case "Edit", "Write", "MultiEdit", "NotebookEdit":
				fp, _ := input["file_path"].(string)
				if fp == "" {
					continue
				}
				if !filepath.IsAbs(fp) {
					fp = filepath.Join(root, fp)
				}
				fp = filepath.Clean(fp)
				isNote := false
				for _, d := range noteDirs {
					if strings.HasPrefix(fp, d) && strings.HasSuffix(fp, ".md") {
						isNote = true
					}
				}
				if isNote {
					w.Crumbs++
				} else if !strings.HasPrefix(fp, memDir) {
					edited[fp] = true
				}
			case "Bash":
				if cmd, _ := input["command"].(string); memoryAddCall.MatchString(cmd) {
					w.Crumbs++
				}
			}
		}
	})
	w.Edited = len(edited)
	return w
}

// NudgeState remembers, per session, the work level at the last nudge.
type NudgeState struct {
	Edited int `json:"edited"`
	Crumbs int `json:"crumbs"`
}

const (
	// FirstNudgeAt is how many edited files a session needs before the first nudge.
	FirstNudgeAt = 3
	// RenudgeAfter is how many more edited files before nudging again.
	RenudgeAfter = 8
)

// ShouldNudge decides whether to ask the agent to save crumbs now. save is
// the new baseline to store, or nil to leave the stored state as it is.
// The baseline moves only when a nudge fires or the agent saves a crumb, so
// quiet turns don't push the next nudge further away.
func ShouldNudge(w Work, last *NudgeState) (nudge bool, save *NudgeState) {
	now := &NudgeState{Edited: w.Edited, Crumbs: w.Crumbs}
	if last == nil {
		if w.Crumbs > 0 {
			return false, now
		}
		if w.Edited >= FirstNudgeAt {
			return true, now
		}
		return false, nil
	}
	if w.Crumbs > last.Crumbs {
		return false, now
	}
	if w.Edited-last.Edited >= RenudgeAfter {
		return true, now
	}
	return false, nil
}

// NudgeMessage is what the agent reads when the nudge fires.
const NudgeMessage = `Baker memory check before you stop: this session changed several files but saved no crumbs.

Did it produce anything the next agent or teammate would need that the code doesn't show? For example:
- a decision and why it was made
- a lesson learned the hard way (a bug cause, a trap, something that didn't work)
- a project fact (a constraint, an external system's behaviour, a convention)

If yes, save each one now, one claim per note:
  baker memory add --kind <decision|lesson|knowledge> --title "<the claim as one sentence>" --body "<detail and why it matters, under 200 words>" --source "<file, commit, issue, or session YYYY-MM-DD>"

If nothing is worth keeping, reply "No new crumbs." and stop. Don't save task progress here (that goes in memory/hot.md), and never save secrets.`

func nudgeStatePath(root, session string) string {
	safe := regexp.MustCompile(`[^A-Za-z0-9_-]+`).ReplaceAllString(session, "_")
	if safe == "" {
		safe = "unknown"
	}
	return filepath.Join(root, ".baker", "nudge", safe+".json")
}

// LoadNudgeState reads the stored state for a session, or nil.
func LoadNudgeState(root, session string) *NudgeState {
	b, err := os.ReadFile(nudgeStatePath(root, session))
	if err != nil {
		return nil
	}
	var s NudgeState
	if json.Unmarshal(b, &s) != nil {
		return nil
	}
	return &s
}

// SaveNudgeState stores the state for a session (best effort).
func SaveNudgeState(root, session string, s NudgeState) {
	p := nudgeStatePath(root, session)
	if os.MkdirAll(filepath.Dir(p), 0o755) != nil {
		return
	}
	b, _ := json.Marshal(s)
	_ = os.WriteFile(p, b, 0o644)
}
