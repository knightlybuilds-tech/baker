// Package recipe loads stack recipes (laravel, generic, ...) and renders
// their templates with facts detected from the project.
//
// A recipe is a folder:
//
//	recipe.toml   metadata, installer, detection, facts, MCP servers
//	agents.md     the managed recipe block for AGENTS.md (Go template)
//	notes.md      optional starting "Project notes" (written once, then user-owned)
//	skills/       optional stack skills, each skills/<name>/SKILL.md
//
// Lookup order: ./.baker/recipes → ~/.baker/recipes → built-in.
package recipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml/v2"
	"github.com/roguenet/baker/assets"
	"github.com/roguenet/baker/internal/config"
)

type Recipe struct {
	Name        string   `toml:"name"`
	Version     int      `toml:"version"`
	Description string   `toml:"description"`
	Install     Install  `toml:"install"`
	Detect      Detect   `toml:"detect"`
	Facts       []Fact   `toml:"facts"`
	MCP         []MCP    `toml:"mcp"`
	Gitignore   []string `toml:"gitignore"`

	Source string `toml:"-"` // "built-in" or a directory path
	files  fs.FS
}

type Install struct {
	Run string `toml:"run"` // template; runs in the parent dir and must create {{.Name}}
}

type Detect struct {
	File     string `toml:"file"`
	Contains string `toml:"contains"`
}

// Fact is a value read from the project after install, used in templates.
type Fact struct {
	Name       string `toml:"name"`
	File       string `toml:"file"`
	JSONPath   string `toml:"json_path"`   // dot path into a JSON file, e.g. "require.php"
	IfContains string `toml:"if_contains"` // or: does the file contain this string?
	Then       string `toml:"then"`
	Else       string `toml:"else"`
	Format     string `toml:"format"` // "version" strips ^~>=v, "major" keeps the first number
	Default    string `toml:"default"`
}

// MCP is a server the recipe registers, optionally only when a package is present.
type MCP struct {
	config.MCPServer
	WhenFile     string `toml:"when_file"`
	WhenContains string `toml:"when_contains"`
}

// Data is what templates see.
type Data struct {
	Name  string
	Facts map[string]string
}

// Load finds a recipe by name using the lookup order.
func Load(name, projectRoot string) (*Recipe, error) {
	dirs := []string{}
	if projectRoot != "" {
		dirs = append(dirs, filepath.Join(projectRoot, ".baker", "recipes", name))
	}
	dirs = append(dirs, filepath.Join(config.Home(), "recipes", name))
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(d, "recipe.toml")); err == nil {
			return loadFS(os.DirFS(d), d)
		}
	}
	sub, err := fs.Sub(assets.FS, "recipes/"+name)
	if err != nil {
		return nil, err
	}
	if _, err := fs.Stat(sub, "recipe.toml"); err != nil {
		return nil, fmt.Errorf("unknown recipe %q (try: %s)", name, strings.Join(Names(projectRoot), ", "))
	}
	return loadFS(sub, "built-in")
}

func loadFS(files fs.FS, source string) (*Recipe, error) {
	b, err := fs.ReadFile(files, "recipe.toml")
	if err != nil {
		return nil, err
	}
	r := &Recipe{}
	if err := toml.Unmarshal(b, r); err != nil {
		return nil, fmt.Errorf("%s/recipe.toml: %w", source, err)
	}
	r.Source, r.files = source, files
	return r, nil
}

// Names lists every recipe available from any source.
func Names(projectRoot string) []string {
	seen := map[string]bool{}
	add := func(dir fs.FS) {
		entries, _ := fs.ReadDir(dir, ".")
		for _, e := range entries {
			if e.IsDir() {
				seen[e.Name()] = true
			}
		}
	}
	if sub, err := fs.Sub(assets.FS, "recipes"); err == nil {
		add(sub)
	}
	add(os.DirFS(filepath.Join(config.Home(), "recipes")))
	if projectRoot != "" {
		add(os.DirFS(filepath.Join(projectRoot, ".baker", "recipes")))
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// DetectRecipe picks the first recipe whose detect rule matches root, else "generic".
func DetectRecipe(root string) string {
	for _, n := range Names(root) {
		if n == "generic" {
			continue
		}
		r, err := Load(n, root)
		if err != nil || r.Detect.File == "" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, r.Detect.File))
		if err != nil {
			continue
		}
		if r.Detect.Contains == "" || bytes.Contains(b, []byte(r.Detect.Contains)) {
			return n
		}
	}
	return "generic"
}

// ResolveFacts reads every fact from the project files.
func (r *Recipe) ResolveFacts(root string) map[string]string {
	out := map[string]string{}
	for _, f := range r.Facts {
		v := f.resolve(root)
		if v == "" {
			v = f.Default
		}
		out[f.Name] = v
	}
	return out
}

func (f Fact) resolve(root string) string {
	b, err := os.ReadFile(filepath.Join(root, f.File))
	if err != nil {
		return ""
	}
	var v string
	switch {
	case f.IfContains != "":
		if bytes.Contains(b, []byte(f.IfContains)) {
			v = f.Then
		} else {
			v = f.Else
		}
	case f.JSONPath != "":
		var doc any
		if json.Unmarshal(b, &doc) != nil {
			return ""
		}
		for _, key := range strings.Split(f.JSONPath, ".") {
			m, ok := doc.(map[string]any)
			if !ok {
				return ""
			}
			doc = m[key]
		}
		if s, ok := doc.(string); ok {
			v = s
		}
	}
	return format(v, f.Format)
}

var versionNum = regexp.MustCompile(`\d+(\.\d+)*`)

func format(v, how string) string {
	switch how {
	case "version":
		return versionNum.FindString(v)
	case "major":
		return strings.SplitN(versionNum.FindString(v), ".", 2)[0]
	}
	return v
}

// MCPFor returns the MCP servers whose conditions hold in root.
func (r *Recipe) MCPFor(root string) []config.MCPServer {
	var out []config.MCPServer
	for _, m := range r.MCP {
		if m.WhenFile != "" {
			b, err := os.ReadFile(filepath.Join(root, m.WhenFile))
			if err != nil || !bytes.Contains(b, []byte(m.WhenContains)) {
				continue
			}
		}
		out = append(out, m.MCPServer)
	}
	return out
}

// Render executes one of the recipe's template files. Missing files render empty.
func (r *Recipe) Render(file string, d Data) (string, error) {
	b, err := fs.ReadFile(r.files, file)
	if err != nil {
		return "", nil
	}
	return Exec(file, string(b), d)
}

// InstallCommand renders the installer command, or "" when the recipe has none.
func (r *Recipe) InstallCommand(d Data) (string, error) {
	if r.Install.Run == "" {
		return "", nil
	}
	return Exec("install", r.Install.Run, d)
}

// Skills returns the recipe's skill files as relative path -> content.
func (r *Recipe) Skills() map[string]string {
	return readTree(r.files, "skills")
}

// Exec renders a Go template string.
func Exec(name, text string, d Data) (string, error) {
	t, err := template.New(name).Option("missingkey=zero").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// readTree reads every file under dir in fsys as relative path -> content.
func readTree(fsys fs.FS, dir string) map[string]string {
	out := map[string]string{}
	_ = fs.WalkDir(fsys, dir, func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(fsys, p)
		if err == nil {
			out[strings.TrimPrefix(p, dir+"/")] = string(b)
		}
		return nil
	})
	return out
}

// CoreSkills returns baker's built-in skills (session-close, consolidate, ingest).
func CoreSkills() map[string]string {
	return readTree(assets.FS, "skills")
}

// Core reads a file from assets/core.
func Core(name string) string {
	b, _ := fs.ReadFile(assets.FS, "core/"+name)
	return string(b)
}
