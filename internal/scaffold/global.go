package scaffold

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/roguenet/baker/internal/agents"
	"github.com/roguenet/baker/internal/blocks"
	"github.com/roguenet/baker/internal/config"
	"github.com/roguenet/baker/internal/recipe"
)

// GlobalInit creates ~/.baker with config, personal AGENTS.md and global memory.
// Existing files are never overwritten.
func GlobalInit() (*Plan, error) {
	p := NewPlan(config.Home(), false)
	p.NoLock = true
	cfg, err := toml.Marshal(config.DefaultGlobal())
	if err != nil {
		return nil, err
	}
	p.Once("config.toml", "# baker global config — see docs/configuration.md\n\n"+string(cfg))
	p.Once("AGENTS.md", recipe.Core("global-AGENTS.md"))
	p.Once("memory/MEMORY.md", recipe.Core("global-MEMORY.md"))
	p.Once("memory/lessons/.gitkeep", "")
	p.Once("memory/knowledge/.gitkeep", "")
	p.Once("skills/.gitkeep", "")
	p.Once("recipes/.gitkeep", "")
	return p, nil
}

var htmlComment = regexp.MustCompile(`(?s)<!--.*?-->`)

const codexHook = `notify = ["sh", "-c", "command -v baker >/dev/null 2>&1 && baker checkpoint --agent codex --quiet || true"]`

// GlobalSync pushes ~/.baker/AGENTS.md and ~/.baker/skills into each
// installed agent's user-level files. Plan root is the user's home dir.
func GlobalSync(force bool) (*Plan, error) {
	home := config.UserHome()
	p := NewPlan(home, force)
	g := config.LoadGlobal()

	rules, _ := os.ReadFile(filepath.Join(config.Home(), "AGENTS.md"))
	body := strings.TrimSpace(htmlComment.ReplaceAllString(string(rules), ""))
	body = regexp.MustCompile(`\n{3,}`).ReplaceAllString(body, "\n\n")
	memDir := filepath.Join(config.Home(), "memory")
	pointer := "## Cross-project memory\nLessons that hold across all my projects live in `" + memDir + "`. Read its `MEMORY.md` index when a task may benefit; open only the notes you need."
	inner := pointer
	if body != "" {
		inner = body + "\n\n" + pointer
	}

	skills := readDir(filepath.Join(config.Home(), "skills"))
	for _, id := range g.Agents {
		a, ok := agents.Get(id)
		if !ok || a.GlobalFile == "" {
			continue
		}
		// Only touch agents that are installed (their config dir exists).
		if _, err := os.Stat(filepath.Join(home, strings.Split(a.GlobalFile, "/")[0])); err != nil {
			continue
		}
		p.Block(a.GlobalFile, "global", "1", inner, "managed by baker; edit ~/.baker/AGENTS.md", blocks.Append)
		if a.GlobalSkillsDir != "" {
			for rel, content := range skills {
				p.Managed(path.Join(a.GlobalSkillsDir, rel), content)
			}
		}
	}

	// Codex has no per-turn hook; its `notify` program runs after each turn.
	if _, err := os.Stat(filepath.Join(home, ".codex")); err == nil && contains(g.Agents, "codex") {
		cur, _ := p.Read(".codex/config.toml")
		user := blocks.Strip(cur, blocks.Hash)
		if !regexp.MustCompile(`(?m)^notify\s*=`).MatchString(user) {
			p.Block(".codex/config.toml", "checkpoint", "1", codexHook, "", blocks.Prepend)
		} else {
			p.Warnings = append(p.Warnings, "~/.codex/config.toml already has a notify command; add `baker checkpoint --agent codex --quiet` to it yourself for automatic checkpoints")
		}
	}
	return p, nil
}

func readDir(dir string) map[string]string {
	out := map[string]string{}
	_ = filepath.WalkDir(dir, func(p string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || e.Name() == ".gitkeep" {
			return nil
		}
		b, err := os.ReadFile(p)
		if err == nil {
			rel, _ := filepath.Rel(dir, p)
			out[filepath.ToSlash(rel)] = string(b)
		}
		return nil
	})
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
