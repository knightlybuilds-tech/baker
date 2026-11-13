// Package bakery serves a local, read-only Bakery view of a project's
// memory/ folder. It binds to loopback only, never writes files, and
// exposes nothing beyond the notes Baker already reads for agents.
package bakery

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/knightlybuilds-tech/baker/internal/memory"
)

// UI is the Bakery workspace front end, synced from bakery-workspace/dist
// with `make bakery-ui`.
//
//go:embed all:ui
var uiFS embed.FS

// Note is one memory note as the Bakery UI sees it.
type Note struct {
	ID           string   `json:"id"`
	Path         string   `json:"path"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Kind         string   `json:"kind"`
	Source       string   `json:"source"`
	AsOf         string   `json:"as_of"`
	LastVerified string   `json:"last_verified"`
	Tags         string   `json:"tags,omitempty"`
	Body         string   `json:"body"`
	Links        []string `json:"links"`
	Issues       []string `json:"issues"`
}

// Finding mirrors memory.Check output.
type Finding struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// Snapshot is the /api/memory response.
type Snapshot struct {
	Mode        string    `json:"mode"`
	Project     string    `json:"project"`
	Branch      string    `json:"branch,omitempty"`
	Head        string    `json:"head,omitempty"`
	HasMemory   bool      `json:"has_memory"`
	Notes       []Note    `json:"notes"`
	Findings    []Finding `json:"findings"`
	GeneratedAt time.Time `json:"generated_at"`
}

// Read builds a snapshot of root/memory. now is injectable for tests.
func Read(root, project string, now time.Time) Snapshot {
	dir := filepath.Join(root, "memory")
	snap := Snapshot{Mode: "local", Project: project, Notes: []Note{}, Findings: []Finding{}, GeneratedAt: now.UTC()}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return snap
	}
	snap.HasMemory = true
	snap.Branch, snap.Head = gitState(root)

	recalled := map[string]memory.RecalledNote{}
	for _, r := range memory.Recall(dir, "", now) {
		recalled[r.Path] = r
	}
	notes := memory.Notes(dir)
	bySlug := map[string]string{}
	for _, n := range notes {
		bySlug[n.Slug] = n.Path
	}
	for _, n := range notes {
		r := recalled[n.Path]
		links := []string{}
		for _, slug := range memory.Links(n.Body) {
			if id, ok := bySlug[slug]; ok && id != n.Path {
				links = append(links, id)
			}
		}
		issues := r.Issues
		if issues == nil {
			issues = []string{}
		}
		snap.Notes = append(snap.Notes, Note{
			ID: n.Path, Path: n.Path, Slug: n.Slug, Title: r.Title, Kind: n.Frontmatter["type"],
			Source: n.Frontmatter["source"], AsOf: n.Frontmatter["as_of"], LastVerified: n.Frontmatter["last_verified"],
			Tags: n.Frontmatter["tags"], Body: stripTitle(n.Body), Links: links, Issues: issues,
		})
	}
	levels := map[memory.Level]string{memory.OK: "ok", memory.Warn: "warn", memory.Error: "error"}
	for _, f := range memory.Check(dir, now) {
		snap.Findings = append(snap.Findings, Finding{Level: levels[f.Level], Msg: f.Msg})
	}
	return snap
}

// stripTitle drops the leading "# Title" line; the UI shows the title separately.
func stripTitle(body string) string {
	lines := strings.Split(strings.TrimLeft(body, "\n"), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		lines = lines[1:]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func gitState(root string) (branch, head string) {
	run := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
	return run("rev-parse", "--abbrev-ref", "HEAD"), run("rev-parse", "--short", "HEAD")
}

// Handler serves the UI and the read-only memory API for root.
func Handler(root, project string) http.Handler {
	ui, err := fs.Sub(uiFS, "ui")
	if err != nil {
		panic(err)
	}
	static := http.FileServer(http.FS(ui))
	mux := http.NewServeMux()
	mux.HandleFunc("/api/memory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(Read(root, project, time.Now()))
	})
	mux.Handle("/", static)
	return guard(mux)
}

// guard keeps the server local and read-only: GET/HEAD only, and a loopback
// Host header so other sites can't reach it through DNS rebinding.
func guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "read-only", http.StatusMethodNotAllowed)
			return
		}
		if !loopbackHost(r.Host) {
			http.Error(w, "local access only", http.StatusForbidden)
			return
		}
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
		next.ServeHTTP(w, r)
	})
}

func loopbackHost(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	host = strings.Trim(host, "[]")
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
