// Package config loads and saves baker's TOML files:
// the project manifest (baker.toml), the lock (.baker/lock.toml),
// the global config (~/.baker/config.toml) and the project registry.
package config

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Standard is the version of the baker standard this binary writes.
const Standard = "0.1"

// DefaultAgents are targeted when nothing else is configured.
var DefaultAgents = []string{"claude", "codex", "gemini", "cursor", "copilot", "opencode"}

// MCPServer is one MCP server registered for every enabled agent.
type MCPServer struct {
	Name    string            `toml:"name"`
	Command string            `toml:"command"`
	Args    []string          `toml:"args,omitempty"`
	Env     map[string]string `toml:"env,omitempty"`
}

// Project is baker.toml at the project root.
type Project struct {
	Standard      string            `toml:"standard"`
	Name          string            `toml:"name"`
	Recipe        string            `toml:"recipe"`
	RecipeVersion int               `toml:"recipe_version"`
	Agents        []string          `toml:"agents"`
	Facts         map[string]string `toml:"facts,omitempty"`
	MCP           []MCPServer       `toml:"mcp,omitempty"`
}

const ProjectFile = "baker.toml"

func LoadProject(root string) (*Project, error) {
	var p Project
	if err := readTOML(filepath.Join(root, ProjectFile), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Project) Save(root string) error {
	return writeTOML(filepath.Join(root, ProjectFile), p)
}

// HasAgent reports whether an agent is enabled for this project.
func (p *Project) HasAgent(id string) bool {
	for _, a := range p.Agents {
		if a == id {
			return true
		}
	}
	return false
}

// FindRoot walks up from dir looking for baker.toml.
func FindRoot(dir string) (string, bool) {
	dir, _ = filepath.Abs(dir)
	for {
		if _, err := os.Stat(filepath.Join(dir, ProjectFile)); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

// Lock records hashes of everything baker generated, so later syncs can tell
// whether the user edited a file or block (and must not be overwritten).
type Lock struct {
	Files  map[string]string `toml:"files"`  // path -> sha256 of whole generated file
	Blocks map[string]string `toml:"blocks"` // "path#block-id" -> sha256 of block inner text
	MCP    []string          `toml:"mcp"`    // MCP server names baker registered
}

const LockFile = ".baker/lock.toml"

func LoadLock(root string) *Lock {
	l := &Lock{}
	_ = readTOML(filepath.Join(root, LockFile), l)
	if l.Files == nil {
		l.Files = map[string]string{}
	}
	if l.Blocks == nil {
		l.Blocks = map[string]string{}
	}
	return l
}

func (l *Lock) Save(root string) error {
	return writeTOML(filepath.Join(root, LockFile), l)
}

// Global is ~/.baker/config.toml.
type Global struct {
	Agents        []string            `toml:"agents"`
	DefaultRecipe string              `toml:"default_recipe"`
	Transcripts   map[string][]string `toml:"transcripts"` // agent -> config dirs (one per account)
}

// Home returns the global baker directory (~/.baker, or $BAKER_HOME).
func Home() string {
	if h := os.Getenv("BAKER_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".baker")
}

// UserHome returns the user's home dir; overridable for tests via $BAKER_USER_HOME.
func UserHome() string {
	if h := os.Getenv("BAKER_USER_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return home
}

// Expand turns a leading ~ into the user's home dir.
func Expand(p string) string {
	if p == "~" {
		return UserHome()
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(UserHome(), p[2:])
	}
	return p
}

func DefaultGlobal() *Global {
	return &Global{
		Agents:        append([]string(nil), DefaultAgents...),
		DefaultRecipe: "generic",
		Transcripts: map[string][]string{
			"claude": {"~/.claude"},
			"codex":  {"~/.codex"},
			"gemini": {"~/.gemini"},
		},
	}
}

// LoadGlobal returns the global config, or defaults when it doesn't exist yet.
func LoadGlobal() *Global {
	g := DefaultGlobal()
	_ = readTOML(filepath.Join(Home(), "config.toml"), g)
	if len(g.Agents) == 0 {
		g.Agents = append([]string(nil), DefaultAgents...)
	}
	return g
}

func (g *Global) Save() error {
	return writeTOML(filepath.Join(Home(), "config.toml"), g)
}

// Registry is ~/.baker/projects.toml: every project baker created or adopted.
type Registry struct {
	Projects []RegistryEntry `toml:"projects"`
}

type RegistryEntry struct {
	Path     string `toml:"path"`
	Name     string `toml:"name"`
	Recipe   string `toml:"recipe"`
	Standard string `toml:"standard"`
	Added    string `toml:"added"`
}

func registryPath() string { return filepath.Join(Home(), "projects.toml") }

func LoadRegistry() *Registry {
	r := &Registry{}
	_ = readTOML(registryPath(), r)
	return r
}

// Register adds or updates a project in the registry.
func Register(root string, p *Project) error {
	r := LoadRegistry()
	entry := RegistryEntry{Path: root, Name: p.Name, Recipe: p.Recipe, Standard: p.Standard, Added: time.Now().Format("2006-01-02")}
	found := false
	for i, e := range r.Projects {
		if e.Path == root {
			entry.Added = e.Added
			r.Projects[i] = entry
			found = true
		}
	}
	if !found {
		r.Projects = append(r.Projects, entry)
	}
	sort.Slice(r.Projects, func(i, j int) bool { return r.Projects[i].Path < r.Projects[j].Path })
	return writeTOML(registryPath(), r)
}

func readTOML(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(b, v)
}

func writeTOML(path string, v any) error {
	b, err := toml.Marshal(v)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// IsNotExist reports whether err means a file was missing.
func IsNotExist(err error) bool { return errors.Is(err, os.ErrNotExist) }
