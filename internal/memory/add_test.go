package memory

import (
	"strings"
	"testing"
	"time"
)

var addNow = time.Date(2026, 11, 17, 9, 0, 0, 0, time.UTC)

func addFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, "MEMORY.md", "# Memory index\n\n<!-- comment -->\n\n## Decisions\n- [Old](decisions/old.md) — x\n\n## Lessons\n\n## Knowledge\n")
	write(t, dir, "decisions/old.md", "---\ntype: decision\n---\n# Account codes must remain strings\n\nBody.\n")
	return dir
}

func TestPrepareBuildsNoteAndIndexesIt(t *testing.T) {
	dir := addFixture(t)
	p, err := Prepare(dir, NewNote{
		Kind: "Lesson", Title: "Retry payments only on HTTP 429 and 503.", Body: "Other failures need investigating. More detail here.",
		Source: "billing/retry.go", Tags: []string{"Payments, retries", "payments"}, Links: []string{"old", "[[missing]]"},
	}, addNow)
	if err != nil {
		t.Fatal(err)
	}
	if p.Path != "memory/lessons/retry-payments-only-on-http-429-and-503.md" {
		t.Errorf("path: %s", p.Path)
	}
	for _, want := range []string{"type: lesson\n", "tags: [payments, retries]\n", "as_of: 2026-11-17\n", "last_verified: 2026-11-17\n",
		"source: billing/retry.go\n", "confidence: medium\n", "# Retry payments only on HTTP 429 and 503\n", "Related: [[old]], [[missing]]\n"} {
		if !strings.Contains(p.Content, want) {
			t.Errorf("content missing %q:\n%s", want, p.Content)
		}
	}
	fm, _ := ParseFrontmatter(p.Content)
	for _, f := range RequiredFields {
		if fm[f] == "" {
			t.Errorf("required field %s empty", f)
		}
	}
	wantLine := "- [Retry payments only on HTTP 429 and 503](lessons/retry-payments-only-on-http-429-and-503.md) — Other failures need investigating."
	lessons := strings.Index(p.Index, "## Lessons")
	knowledge := strings.Index(p.Index, "## Knowledge")
	if i := strings.Index(p.Index, wantLine); i < lessons || i > knowledge {
		t.Errorf("index line not under Lessons:\n%s", p.Index)
	}
	if !strings.Contains(strings.Join(p.Warnings, "\n"), "[[missing]]") {
		t.Errorf("want a warning for the missing link, got %v", p.Warnings)
	}
}

func TestPrepareRejectsBadNotes(t *testing.T) {
	dir := addFixture(t)
	good := NewNote{Kind: "decision", Title: "A new claim here", Body: "b", Source: "s"}
	for name, mutate := range map[string]func(n *NewNote){
		"bad kind":       func(n *NewNote) { n.Kind = "idea" },
		"no title":       func(n *NewNote) { n.Title = " " },
		"multiline":      func(n *NewNote) { n.Title = "a\nb" },
		"no body":        func(n *NewNote) { n.Body = "" },
		"no source":      func(n *NewNote) { n.Source = "" },
		"bad date":       func(n *NewNote) { n.AsOf = "yesterday" },
		"bad confidence": func(n *NewNote) { n.Confidence = "sure" },
		"too long":       func(n *NewNote) { n.Body = strings.Repeat("word ", MaxWords+1) },
		"duplicate slug": func(n *NewNote) { n.Title = "Old" },
	} {
		n := good
		mutate(&n)
		if _, err := Prepare(dir, n, addNow); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
	if _, err := Prepare(dir, good, addNow); err != nil {
		t.Errorf("baseline note should pass: %v", err)
	}
}

func TestPrepareWarnsOnSimilarTitle(t *testing.T) {
	p, err := Prepare(addFixture(t), NewNote{Kind: "decision", Title: "Account codes remain strings always", Body: "b", Source: "s"}, addNow)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(p.Warnings, "\n"), "similar to memory/decisions/old.md") {
		t.Errorf("want a similarity warning, got %v", p.Warnings)
	}
}

func TestPrepareNeedsIndex(t *testing.T) {
	dir := t.TempDir()
	if _, err := Prepare(dir, NewNote{Kind: "lesson", Title: "x y z", Body: "b", Source: "s"}, addNow); err == nil {
		t.Error("want an error without MEMORY.md")
	}
}

func TestInsertUnder(t *testing.T) {
	cases := []struct{ index, heading, want string }{
		{"## A\n- 1\n\n## B\n", "## A", "## A\n- 1\n- new\n\n## B\n"},
		{"## A\n\n## B\n", "## A", "## A\n- new\n\n## B\n"},
		{"## A\n- 1\n## B\n", "## B", "## A\n- 1\n## B\n- new\n"},
		{"# Index\n", "## Knowledge", "# Index\n\n## Knowledge\n- new\n"},
	}
	for _, c := range cases {
		if got := insertUnder(c.index, c.heading, "- new"); got != c.want {
			t.Errorf("insertUnder(%q, %q):\n got %q\nwant %q", c.index, c.heading, got, c.want)
		}
	}
}

func TestSlugify(t *testing.T) {
	if got := Slugify("Retry payments only on HTTP 429 & 503!"); got != "retry-payments-only-on-http-429-503" {
		t.Errorf("got %q", got)
	}
	if got := Slugify("Agents save crumbs through baker memory add, nudged by a Claude Stop hook"); got != "agents-save-crumbs-through-baker-memory-add-nudged" {
		t.Errorf("cut slug should drop dangling short words, got %q", got)
	}
	long := Slugify(strings.Repeat("abcdef ", 20))
	if len(long) > 60 || strings.HasSuffix(long, "-") {
		t.Errorf("long slug %q", long)
	}
}
