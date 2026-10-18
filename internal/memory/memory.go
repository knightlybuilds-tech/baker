// Package memory checks a memory/ folder against the atomic-note rules:
// every note indexed, links resolve, notes stay small, frontmatter is
// present and facts are re-verified within the freshness window.
package memory

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	MaxWords  = 400
	StaleDays = 90
)

// NoteDirs hold atomic notes.
var NoteDirs = []string{"decisions", "lessons", "knowledge"}

var RequiredFields = []string{"type", "as_of", "last_verified", "source"}

type Level int

const (
	OK Level = iota
	Warn
	Error
)

type Finding struct {
	Level Level
	Msg   string
}

type Note struct {
	Path        string // relative to the memory dir, e.g. lessons/foo.md
	Slug        string
	Frontmatter map[string]string
	Body        string
}

var (
	mdLink   = regexp.MustCompile(`\]\(([^)\s]+\.md)\)`)
	wikiLink = regexp.MustCompile(`\[\[([^\]|#]+)(?:[#|][^\]]*)?\]\]`)
)

// Notes loads every atomic note under dir.
func Notes(dir string) []Note {
	var out []Note
	for _, sub := range NoteDirs {
		_ = filepath.WalkDir(filepath.Join(dir, sub), func(p string, e os.DirEntry, err error) error {
			if err != nil || e.IsDir() || !strings.HasSuffix(p, ".md") {
				return nil
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			rel, _ := filepath.Rel(dir, p)
			fm, body := ParseFrontmatter(string(b))
			out = append(out, Note{Path: filepath.ToSlash(rel), Slug: strings.TrimSuffix(filepath.Base(p), ".md"), Frontmatter: fm, Body: body})
			return nil
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// ParseFrontmatter splits simple `key: value` YAML frontmatter from the body.
func ParseFrontmatter(s string) (map[string]string, string) {
	fm := map[string]string{}
	if !strings.HasPrefix(s, "---\n") {
		return fm, s
	}
	end := strings.Index(s[4:], "\n---")
	if end < 0 {
		return fm, s
	}
	sc := bufio.NewScanner(strings.NewReader(s[4 : 4+end]))
	for sc.Scan() {
		k, v, ok := strings.Cut(sc.Text(), ":")
		if ok && !strings.HasPrefix(k, " ") {
			fm[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	body := s[4+end+4:]
	return fm, strings.TrimPrefix(body, "\n")
}

// Check runs every memory rule on dir (usually <project>/memory).
func Check(dir string, now time.Time) []Finding {
	var out []Finding
	add := func(l Level, f string, a ...any) { out = append(out, Finding{l, fmt.Sprintf(f, a...)}) }

	index, err := os.ReadFile(filepath.Join(dir, "MEMORY.md"))
	if err != nil {
		add(Error, "memory/MEMORY.md is missing — agents have no index to start from")
		return out
	}
	notes := Notes(dir)
	bySlug := map[string]bool{}
	byPath := map[string]bool{}
	for _, n := range notes {
		bySlug[n.Slug] = true
		byPath[n.Path] = true
	}

	// Indexed = linked from MEMORY.md, or from a hub note (filename starts with "_") that is itself linked.
	indexed := map[string]bool{}
	for _, m := range mdLink.FindAllStringSubmatch(string(index), -1) {
		target := filepath.ToSlash(filepath.Clean(m[1]))
		indexed[target] = true
		if _, err := os.Stat(filepath.Join(dir, target)); err != nil {
			add(Error, "MEMORY.md links to %s, which doesn't exist", target)
		}
	}
	for _, n := range notes {
		if strings.HasPrefix(n.Slug, "_") && indexed[n.Path] {
			for _, m := range mdLink.FindAllStringSubmatch(n.Body, -1) {
				indexed[filepath.ToSlash(filepath.Join(filepath.Dir(n.Path), m[1]))] = true
			}
			for _, m := range wikiLink.FindAllStringSubmatch(n.Body, -1) {
				for _, o := range notes {
					if o.Slug == strings.TrimSpace(m[1]) {
						indexed[o.Path] = true
					}
				}
			}
		}
	}

	for _, n := range notes {
		where := "memory/" + n.Path
		if !indexed[n.Path] {
			add(Warn, "%s is not in MEMORY.md (or a linked hub note) — agents won't find it", where)
		}
		for _, f := range RequiredFields {
			if n.Frontmatter[f] == "" {
				add(Warn, "%s is missing frontmatter field %q", where, f)
			}
		}
		if words := len(strings.Fields(n.Body)); words > MaxWords && !strings.HasPrefix(n.Slug, "_") {
			add(Warn, "%s has %d words (max %d) — split it into atomic notes", where, words, MaxWords)
		}
		if lv := n.Frontmatter["last_verified"]; lv != "" {
			if t, err := time.Parse("2006-01-02", lv); err == nil && now.Sub(t) > StaleDays*24*time.Hour {
				add(Warn, "%s last verified %s (over %d days) — re-check it", where, lv, StaleDays)
			}
		}
		for _, m := range wikiLink.FindAllStringSubmatch(n.Body, -1) {
			if slug := strings.TrimSpace(m[1]); !bySlug[slug] {
				add(Warn, "%s links to [[%s]], which doesn't exist", where, slug)
			}
		}
	}
	if len(out) == 0 {
		add(OK, "memory: %d notes, all indexed, linked and fresh", len(notes))
	}
	return out
}
