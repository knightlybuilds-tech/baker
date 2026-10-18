package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/roguenet/baker/internal/agents"
	"github.com/roguenet/baker/internal/blocks"
	"github.com/roguenet/baker/internal/config"
	"github.com/roguenet/baker/internal/memory"
)

// Doctor checks a project against the baker standard.
func Doctor(root string) []memory.Finding {
	var out []memory.Finding
	add := func(l memory.Level, f string, a ...any) {
		out = append(out, memory.Finding{Level: l, Msg: fmt.Sprintf(f, a...)})
	}
	read := func(p string) (string, bool) {
		b, err := os.ReadFile(filepath.Join(root, p))
		return string(b), err == nil
	}

	proj, err := config.LoadProject(root)
	if err != nil {
		add(memory.Error, "baker.toml missing or invalid — run `baker init`")
		return out
	}
	add(memory.OK, "baker.toml: recipe %s, standard %s, agents %s", proj.Recipe, proj.Standard, strings.Join(proj.Agents, ", "))
	if proj.Standard != config.Standard {
		add(memory.Warn, "project uses standard %s, this baker writes %s — run `baker sync`", proj.Standard, config.Standard)
	}

	agentsMD, ok := read("AGENTS.md")
	switch {
	case !ok:
		add(memory.Error, "AGENTS.md is missing — run `baker sync`")
	case len(agentsMD) > 32*1024:
		add(memory.Warn, "AGENTS.md is %d KiB; Codex truncates past 32 KiB — move detail into memory notes", len(agentsMD)/1024)
	default:
		if _, ok := blocks.Find(agentsMD, blocks.Markdown, "core"); !ok {
			add(memory.Warn, "AGENTS.md has no baker core block — run `baker sync`")
		} else {
			add(memory.OK, "AGENTS.md: %d KiB, core block present", (len(agentsMD)+1023)/1024)
		}
	}

	for _, a := range agents.For(proj) {
		if a.BridgeFile == "" {
			continue
		}
		c, ok := read(a.BridgeFile)
		b, found := blocks.Find(c, blocks.StyleFor(a.BridgeFile), "bridge")
		switch {
		case !ok:
			add(memory.Error, "%s is missing, so %s won't read AGENTS.md — run `baker sync`", a.BridgeFile, a.Name)
		case !found || !strings.Contains(b.Inner, "AGENTS.md"):
			add(memory.Error, "%s doesn't import AGENTS.md — run `baker sync`", a.BridgeFile)
		default:
			add(memory.OK, "%s bridges %s to AGENTS.md", a.BridgeFile, a.Name)
		}
		if a.SkillsDir != "" {
			if drift := skillDrift(root, a.SkillsDir); drift != "" {
				add(memory.Warn, "%s is out of date with .agents/skills (%s) — run `baker sync`", a.SkillsDir, drift)
			}
		}
	}

	for _, f := range agents.LegacyFiles {
		if info, err := os.Stat(filepath.Join(root, f)); err == nil && !info.IsDir() {
			add(memory.Warn, "%s is a legacy rules file; some tools (Zed) read it instead of AGENTS.md — move its content into AGENTS.md and delete it", f)
		}
	}

	if proj.HasAgent("claude") {
		if c, ok := read(".claude/settings.json"); !ok || !strings.Contains(c, "baker checkpoint") {
			add(memory.Warn, "Claude Code checkpoint hook not installed — run `baker sync`")
		}
		if _, err := exec.LookPath("baker"); err != nil {
			add(memory.Warn, "`baker` is not on your PATH, so agent hooks can't take checkpoints")
		}
	}

	out = append(out, memory.Check(filepath.Join(root, "memory"), time.Now())...)
	return out
}

// skillDrift reports the first file that differs between .agents/skills and a mirror.
func skillDrift(root, mirror string) string {
	src := filepath.Join(root, ".agents", "skills")
	drift := ""
	_ = filepath.WalkDir(src, func(p string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || drift != "" {
			return nil
		}
		rel, _ := filepath.Rel(src, p)
		a, _ := os.ReadFile(p)
		b, err := os.ReadFile(filepath.Join(root, mirror, rel))
		if err != nil || string(a) != string(b) {
			drift = rel
		}
		return nil
	})
	return drift
}
