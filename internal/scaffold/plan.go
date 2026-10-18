// Package scaffold computes every file baker wants to write as a Plan,
// shows it, and applies it. Nothing touches disk until Apply.
package scaffold

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knightlybuilds-tech/baker/internal/blocks"
	"github.com/knightlybuilds-tech/baker/internal/config"
)

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Skip   Action = "skip" // user edited it; left alone
)

type Change struct {
	Path   string
	Action Action
	Note   string
}

// Plan accumulates desired file contents relative to Root.
type Plan struct {
	Root     string
	Force    bool // overwrite files/blocks the user edited
	NoLock   bool // don't write .baker/lock.toml on Apply
	Lock     *config.Lock
	Warnings []string

	original map[string]*string // nil = file didn't exist
	pending  map[string]string
	notes    map[string][]string
	skipped  map[string]string
}

func NewPlan(root string, force bool) *Plan {
	return &Plan{
		Root: root, Force: force, Lock: config.LoadLock(root),
		original: map[string]*string{}, pending: map[string]string{},
		notes: map[string][]string{}, skipped: map[string]string{},
	}
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:8])
}

// Read returns the current planned content of path, and whether it exists.
func (p *Plan) Read(path string) (string, bool) {
	if c, ok := p.pending[path]; ok {
		return c, true
	}
	orig := p.load(path)
	if orig == nil {
		return "", false
	}
	return *orig, true
}

func (p *Plan) load(path string) *string {
	if o, ok := p.original[path]; ok {
		return o
	}
	b, err := os.ReadFile(filepath.Join(p.Root, path))
	if err != nil {
		p.original[path] = nil
		return nil
	}
	s := string(b)
	p.original[path] = &s
	return &s
}

func (p *Plan) set(path, content, note string) {
	p.load(path)
	p.pending[path] = content
	if note != "" {
		p.notes[path] = append(p.notes[path], note)
	}
}

func (p *Plan) skip(path, why string) {
	if prev, ok := p.skipped[path]; ok {
		why = prev + "; " + why
	}
	p.skipped[path] = why
}

// Once writes a file only if it doesn't exist. The user owns it afterwards.
func (p *Plan) Once(path, content string) {
	if _, ok := p.Read(path); !ok {
		p.set(path, content, "")
	}
}

// Managed writes a fully generated file. If the user edited it since baker
// last wrote it (hash differs from the lock), it is skipped unless Force.
func (p *Plan) Managed(path, content string) {
	cur, exists := p.Read(path)
	if exists && cur == content {
		p.Lock.Files[path] = hash(content)
		return
	}
	if exists && !p.Force {
		if h, ok := p.Lock.Files[path]; !ok || h != hash(cur) {
			p.skip(path, "changed by you since baker generated it")
			return
		}
	}
	p.set(path, content, "")
	p.Lock.Files[path] = hash(content)
}

// Block upserts a managed block. A block the user edited is skipped unless Force.
func (p *Plan) Block(path, id, version, inner, note string, pos blocks.Position) {
	style := blocks.StyleFor(path)
	cur, _ := p.Read(path)
	key := path + "#" + id
	inner = strings.TrimRight(inner, "\n")
	if b, ok := blocks.Find(cur, style, id); ok {
		if b.Inner == inner && b.Version == version {
			p.Lock.Blocks[key] = hash(inner)
			return
		}
		if !p.Force {
			if h, ok := p.Lock.Blocks[key]; ok && h != hash(b.Inner) {
				p.skip(path, fmt.Sprintf("block %q was edited by you", id))
				return
			}
		}
	}
	p.set(path, blocks.Upsert(cur, style, id, version, inner, note, pos), "")
	p.Lock.Blocks[key] = hash(inner)
}

// JSON loads a JSON object file (or {}), lets mutate change it, and writes it back if changed.
func (p *Plan) JSON(path string, mutate func(doc map[string]any)) error {
	cur, exists := p.Read(path)
	doc := map[string]any{}
	if exists && strings.TrimSpace(cur) != "" {
		if err := json.Unmarshal([]byte(cur), &doc); err != nil {
			p.Warnings = append(p.Warnings, fmt.Sprintf("%s: not valid JSON, left alone (%v)", path, err))
			return nil
		}
	}
	mutate(doc)
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false) // keep shell commands like `a && b > c` readable
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return err
	}
	out := buf.String()
	if exists && normalizeJSON(cur) == normalizeJSON(out) {
		return nil
	}
	if len(doc) == 0 && !exists {
		return nil
	}
	p.set(path, out, "")
	return nil
}

func normalizeJSON(s string) string {
	var v any
	if json.Unmarshal([]byte(s), &v) != nil {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func (p *Plan) dirHasFiles(dir string) bool {
	entries, err := os.ReadDir(filepath.Join(p.Root, dir))
	return err == nil && len(entries) > 0
}

// Note attaches an explanation to a path shown in the plan.
func (p *Plan) Note(path, note string) {
	p.notes[path] = append(p.notes[path], note)
}

// Changes lists what Apply would do, sorted by path.
func (p *Plan) Changes() []Change {
	var out []Change
	for path, content := range p.pending {
		orig := p.original[path]
		switch {
		case orig == nil:
			out = append(out, Change{path, Create, strings.Join(p.notes[path], "; ")})
		case *orig != content:
			out = append(out, Change{path, Update, strings.Join(p.notes[path], "; ")})
		}
	}
	for path, why := range p.skipped {
		if _, pending := p.pending[path]; !pending {
			out = append(out, Change{path, Skip, why + " — left alone, --force overwrites"})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// HasWrites reports whether Apply would write anything.
func (p *Plan) HasWrites() bool {
	for _, c := range p.Changes() {
		if c.Action != Skip {
			return true
		}
	}
	return false
}

// Print writes a human-readable plan.
func (p *Plan) Print(w io.Writer) {
	changes := p.Changes()
	if len(changes) == 0 {
		fmt.Fprintln(w, "  Nothing to change — everything is up to date.")
	}
	for _, c := range changes {
		sym := map[Action]string{Create: "+", Update: "~", Skip: "!"}[c.Action]
		line := fmt.Sprintf("  %s %-8s %s", sym, c.Action, c.Path)
		if c.Note != "" {
			line += "  (" + c.Note + ")"
		}
		fmt.Fprintln(w, line)
	}
	for path, why := range p.skipped {
		if _, pending := p.pending[path]; pending {
			fmt.Fprintf(w, "  ! %s: %s — left alone, --force overwrites\n", path, why)
		}
	}
	for _, warn := range p.Warnings {
		fmt.Fprintln(w, "  ! "+warn)
	}
}

// Apply writes every pending file and saves the lock.
func (p *Plan) Apply() error {
	for path, content := range p.pending {
		if o := p.original[path]; o != nil && *o == content {
			continue
		}
		full := filepath.Join(p.Root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}
	if p.NoLock {
		return nil
	}
	return p.Lock.Save(p.Root)
}
