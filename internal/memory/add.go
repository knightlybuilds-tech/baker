package memory

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Kinds maps a note kind to its folder under memory/ and its MEMORY.md heading.
var Kinds = map[string]struct{ Dir, Heading string }{
	"decision":  {"decisions", "## Decisions"},
	"lesson":    {"lessons", "## Lessons"},
	"knowledge": {"knowledge", "## Knowledge"},
}

// SoftWords is the size an atomic note should stay under; MaxWords is the hard cap.
const SoftWords = 200

// NewNote is one crumb to add.
type NewNote struct {
	Kind       string
	Title      string
	Body       string
	Source     string
	Tags       []string
	Links      []string // slugs of related notes
	Hook       string   // short index hook; defaults to the body's first sentence
	Confidence string   // high | medium | low
	AsOf       string   // YYYY-MM-DD; defaults to today
}

// Prepared is a validated note ready to write. Paths are relative to the project root.
type Prepared struct {
	Slug      string   `json:"slug"`
	Path      string   `json:"path"`
	Content   string   `json:"-"`
	IndexPath string   `json:"index_path"`
	Index     string   `json:"-"`
	Warnings  []string `json:"warnings"`
}

var (
	slugStrip = regexp.MustCompile(`[^a-z0-9]+`)
	dateRe    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	stopWords = map[string]bool{"a": true, "an": true, "the": true, "is": true, "are": true, "be": true, "to": true, "of": true, "and": true, "or": true, "in": true, "on": true, "for": true, "with": true, "must": true, "should": true, "we": true, "it": true}
)

// Slugify turns a claim into a filename slug of at most 60 characters.
func Slugify(title string) string {
	s := strings.Trim(slugStrip.ReplaceAllString(strings.ToLower(title), "-"), "-")
	if len(s) > 60 {
		s = s[:60]
		if i := strings.LastIndex(s, "-"); i > 30 {
			s = s[:i]
		}
		// Don't end a cut slug on a dangling little word ("...-nudged-by-a").
		parts := strings.Split(strings.Trim(s, "-"), "-")
		for len(parts) > 3 && (stopWords[parts[len(parts)-1]] || len(parts[len(parts)-1]) <= 2) {
			parts = parts[:len(parts)-1]
		}
		s = strings.Join(parts, "-")
	}
	return strings.Trim(s, "-")
}

// Prepare validates n against the notes in dir (<root>/memory) and builds the
// note file and the updated MEMORY.md. It writes nothing.
func Prepare(dir string, n NewNote, now time.Time) (*Prepared, error) {
	n.Kind = strings.ToLower(strings.TrimSpace(n.Kind))
	n.Title = strings.TrimSuffix(strings.TrimSpace(n.Title), ".")
	n.Body = strings.TrimSpace(n.Body)
	n.Source = strings.TrimSpace(n.Source)
	kind, ok := Kinds[n.Kind]
	switch {
	case !ok:
		return nil, fmt.Errorf("kind must be decision, lesson or knowledge (got %q)", n.Kind)
	case n.Title == "":
		return nil, errors.New("a title is required: the claim as one full sentence")
	case strings.ContainsAny(n.Title, "\n\r"):
		return nil, errors.New("the title must be a single line")
	case len(n.Title) > 140:
		return nil, errors.New("the title is over 140 characters: state the claim, put detail in the body")
	case n.Body == "":
		return nil, errors.New("a body is required: the detail, and why it matters here")
	case n.Source == "":
		return nil, errors.New("a source is required: a file path, commit, issue, doc or \"session YYYY-MM-DD\"")
	}
	if n.AsOf == "" {
		n.AsOf = now.Format("2006-01-02")
	} else if !dateRe.MatchString(n.AsOf) {
		return nil, fmt.Errorf("as-of must be YYYY-MM-DD (got %q)", n.AsOf)
	}
	if n.Confidence == "" {
		n.Confidence = "medium"
	} else if n.Confidence != "high" && n.Confidence != "medium" && n.Confidence != "low" {
		return nil, fmt.Errorf("confidence must be high, medium or low (got %q)", n.Confidence)
	}
	words := len(strings.Fields(n.Body))
	if words > MaxWords {
		return nil, fmt.Errorf("the body has %d words (max %d): split it into several notes", words, MaxWords)
	}

	slug := Slugify(n.Title)
	if slug == "" {
		return nil, errors.New("the title has no letters or digits to make a filename from")
	}
	p := &Prepared{Slug: slug, Path: filepath.ToSlash(filepath.Join("memory", kind.Dir, slug+".md")), IndexPath: "memory/MEMORY.md", Warnings: []string{}}
	if words > SoftWords {
		p.Warnings = append(p.Warnings, fmt.Sprintf("body has %d words; atomic notes stay under %d", words, SoftWords))
	}

	existing := Notes(dir)
	bySlug := map[string]Note{}
	for _, e := range existing {
		bySlug[e.Slug] = e
	}
	if e, dup := bySlug[slug]; dup {
		return nil, fmt.Errorf("memory/%s already exists: edit that note instead, or choose a different claim", e.Path)
	}
	for _, e := range existing {
		if similar(n.Title, titleOf(e)) {
			p.Warnings = append(p.Warnings, fmt.Sprintf("similar to memory/%s (%q): check it isn't a duplicate", e.Path, titleOf(e)))
		}
	}
	linked := map[string]bool{}
	for _, l := range Links(n.Body) {
		linked[l] = true
	}
	var extra []string
	for _, l := range n.Links {
		l = strings.TrimSuffix(strings.Trim(strings.TrimSpace(l), "[]"), ".md")
		if l == "" || linked[l] || l == slug {
			continue
		}
		linked[l] = true
		extra = append(extra, "[["+l+"]]")
	}
	for l := range linked {
		if _, ok := bySlug[l]; !ok {
			p.Warnings = append(p.Warnings, fmt.Sprintf("links to [[%s]], which doesn't exist yet", l))
		}
	}

	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "type: %s\nscope: project\n", n.Kind)
	fmt.Fprintf(&b, "tags: [%s]\n", strings.Join(cleanTags(n.Tags), ", "))
	fmt.Fprintf(&b, "as_of: %s\nlast_verified: %s\n", n.AsOf, now.Format("2006-01-02"))
	fmt.Fprintf(&b, "source: %s\nconfidence: %s\n", oneLine(n.Source), n.Confidence)
	b.WriteString("---\n")
	fmt.Fprintf(&b, "# %s\n\n%s\n", n.Title, n.Body)
	if len(extra) > 0 {
		fmt.Fprintf(&b, "\nRelated: %s\n", strings.Join(extra, ", "))
	}
	p.Content = b.String()

	index, err := os.ReadFile(filepath.Join(dir, "MEMORY.md"))
	if err != nil {
		return nil, errors.New("memory/MEMORY.md is missing: run `baker init` first")
	}
	hook := strings.TrimSpace(n.Hook)
	if hook == "" {
		hook = firstSentence(n.Body, 80)
	}
	line := fmt.Sprintf("- [%s](%s/%s.md) — %s", n.Title, kind.Dir, slug, oneLine(hook))
	p.Index = insertUnder(string(index), kind.Heading, line)
	return p, nil
}

func titleOf(n Note) string {
	for _, line := range strings.Split(n.Body, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return n.Slug
}

// similar reports whether two titles share most of their meaningful words.
func similar(a, b string) bool {
	ta, tb := tokens(a), tokens(b)
	if len(ta) < 3 || len(tb) < 3 {
		return false
	}
	inter := 0
	for w := range ta {
		if tb[w] {
			inter++
		}
	}
	union := len(ta) + len(tb) - inter
	return float64(inter)/float64(union) >= 0.6
}

func tokens(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.Fields(slugStrip.ReplaceAllString(strings.ToLower(s), " ")) {
		if !stopWords[w] {
			out[w] = true
		}
	}
	return out
}

func cleanTags(tags []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, t := range tags {
		for _, part := range strings.Split(t, ",") {
			part = Slugify(part)
			if part != "" && !seen[part] {
				seen[part] = true
				out = append(out, part)
			}
		}
	}
	return out
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

func firstSentence(body string, max int) string {
	s := oneLine(body)
	for _, sep := range []string{". ", "? ", "! "} {
		if i := strings.Index(s, sep); i > 0 {
			s = s[:i+1]
			break
		}
	}
	if len(s) > max {
		s = strings.TrimSpace(s[:max])
		if i := strings.LastIndex(s, " "); i > max/2 {
			s = s[:i]
		}
		s += "…"
	}
	return s
}

// insertUnder adds line at the end of the list under heading, creating the
// heading at the end of the file if it's missing.
func insertUnder(index, heading, line string) string {
	lines := strings.Split(strings.TrimRight(index, "\n"), "\n")
	at := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == heading {
			at = i
			break
		}
	}
	if at < 0 {
		return strings.TrimRight(index, "\n") + "\n\n" + heading + "\n" + line + "\n"
	}
	end := at + 1
	for end < len(lines) && !strings.HasPrefix(lines[end], "## ") {
		end++
	}
	pos := end
	for pos > at+1 && strings.TrimSpace(lines[pos-1]) == "" {
		pos--
	}
	out := append([]string{}, lines[:pos]...)
	out = append(out, line)
	if pos < len(lines) && strings.HasPrefix(lines[pos], "## ") {
		out = append(out, "")
	}
	out = append(out, lines[pos:]...)
	return strings.Join(out, "\n") + "\n"
}
