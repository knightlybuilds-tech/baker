package blocks

import (
	"strings"
	"testing"
)

func TestUpsertInsertsThenReplaces(t *testing.T) {
	doc := "# my-app\n\nUser notes.\n"
	doc = Upsert(doc, Markdown, "core", "0.1", "first", "", Append)
	if !strings.Contains(doc, "<!-- baker:begin core@0.1 -->\nfirst\n<!-- baker:end core -->") {
		t.Fatalf("block not inserted:\n%s", doc)
	}
	doc = Upsert(doc, Markdown, "core", "0.2", "second", "", Append)
	if strings.Contains(doc, "first") || !strings.Contains(doc, "core@0.2 -->\nsecond") {
		t.Fatalf("block not replaced:\n%s", doc)
	}
	if !strings.HasPrefix(doc, "# my-app\n\nUser notes.\n") {
		t.Fatalf("user text changed:\n%s", doc)
	}
}

func TestFindReadsVersionAndInner(t *testing.T) {
	doc := "a\n# baker:begin baker@1\nx\ny\n# baker:end baker\nb\n"
	b, ok := Find(doc, Hash, "baker")
	if !ok || b.Version != "1" || b.Inner != "x\ny" {
		t.Fatalf("got %+v ok=%v", b, ok)
	}
}

func TestIDsDoNotCollide(t *testing.T) {
	doc := Render(Markdown, "recipe", "1", "R", "") + Render(Markdown, "recipe:laravel", "1", "L", "")
	b, ok := Find(doc, Markdown, "recipe")
	if !ok || b.Inner != "R" {
		t.Fatalf("got %+v", b)
	}
	doc = Render(Markdown, "recipe:laravel", "1", "L", "") + Render(Markdown, "recipe", "1", "R", "")
	b, ok = Find(doc, Markdown, "recipe")
	if !ok || b.Inner != "R" {
		t.Fatalf("got %+v", b)
	}
}

func TestStrip(t *testing.T) {
	doc := "keep me\n" + Render(Markdown, "bridge", "1", "@AGENTS.md", "generated") + "and me\n"
	if got := Strip(doc, Markdown); got != "keep me\nand me" {
		t.Fatalf("got %q", got)
	}
}
