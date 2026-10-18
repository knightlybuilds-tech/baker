// Package blocks reads and writes baker-managed regions inside text files.
//
// A managed block looks like this (markdown style):
//
//	<!-- baker:begin core@0.1 -->
//	...content baker owns...
//	<!-- baker:end core -->
//
// Everything outside the markers belongs to the user and is never touched.
package blocks

import (
	"fmt"
	"regexp"
	"strings"
)

// Style is the comment syntax used for the markers.
type Style struct {
	Open  string
	Close string
}

var (
	Markdown = Style{Open: "<!-- ", Close: " -->"}
	Hash     = Style{Open: "# ", Close: ""} // yaml, toml, .gitignore
)

// StyleFor picks a marker style from a file name.
func StyleFor(path string) Style {
	if strings.HasSuffix(path, ".md") {
		return Markdown
	}
	return Hash
}

// Block is one managed region found in a file.
type Block struct {
	ID      string
	Version string
	Inner   string // content between the marker lines, without trailing newline
	Start   int    // byte offset of the begin marker
	End     int    // byte offset just past the end marker line
}

func beginRe(s Style, id string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(s.Open) + `baker:begin ` + regexp.QuoteMeta(id) + `@(\S+)[^\n]*\n`)
}

func endRe(s Style, id string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(s.Open) + `baker:end ` + regexp.QuoteMeta(id) + `(?:[ \t][^\n]*)?(\n|$)`)
}

// Find returns the block with the given id, if present.
func Find(content string, s Style, id string) (Block, bool) {
	bm := beginRe(s, id).FindStringSubmatchIndex(content)
	if bm == nil {
		return Block{}, false
	}
	rest := content[bm[1]:]
	em := endRe(s, id).FindStringIndex(rest)
	if em == nil {
		return Block{}, false
	}
	return Block{
		ID:      id,
		Version: content[bm[2]:bm[3]],
		Inner:   strings.TrimSuffix(rest[:em[0]], "\n"),
		Start:   bm[0],
		End:     bm[1] + em[1],
	}, true
}

// Render produces the full text of a block, including markers and a trailing newline.
func Render(s Style, id, version, inner string, note string) string {
	begin := fmt.Sprintf("%sbaker:begin %s@%s", s.Open, id, version)
	if note != "" {
		begin += " — " + note
	}
	begin += s.Close
	end := fmt.Sprintf("%sbaker:end %s%s", s.Open, id, s.Close)
	inner = strings.TrimRight(inner, "\n")
	return begin + "\n" + inner + "\n" + end + "\n"
}

// Where to put a block that doesn't exist yet.
type Position int

const (
	Append Position = iota
	Prepend
)

// Upsert replaces the block if it exists, or inserts it at pos.
func Upsert(content string, s Style, id, version, inner, note string, pos Position) string {
	rendered := Render(s, id, version, inner, note)
	if b, ok := Find(content, s, id); ok {
		return content[:b.Start] + rendered + content[b.End:]
	}
	if content == "" {
		return rendered
	}
	if pos == Prepend {
		return rendered + "\n" + content
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + rendered
}

// Strip removes every baker block from content, leaving only user text.
func Strip(content string, s Style) string {
	re := regexp.MustCompile(`(?ms)^` + regexp.QuoteMeta(s.Open) + `baker:begin (\S+?)@.*?^` + regexp.QuoteMeta(s.Open) + `baker:end [^\n]*(\n|$)`)
	return strings.TrimSpace(re.ReplaceAllString(content, ""))
}
