package memory

import (
	"strings"
	"time"
)

// RecalledNote is portable project context, not an approval or authority claim.
type RecalledNote struct {
	Path         string   `json:"path"`
	Title        string   `json:"title"`
	Kind         string   `json:"kind"`
	Source       string   `json:"source"`
	LastVerified string   `json:"last_verified"`
	Body         string   `json:"body,omitempty"`
	Issues       []string `json:"issues"`
}

// Recall reads ordinary Markdown notes. Missing metadata is reported, not an
// approval gate. Stale notes remain available as historical leads.
func Recall(dir, query string, now time.Time) []RecalledNote {
	out := []RecalledNote{}
	for _, note := range Notes(dir) {
		if query != "" && !strings.Contains(strings.ToLower(note.Body), strings.ToLower(query)) {
			continue
		}
		title := note.Slug
		for _, line := range strings.Split(note.Body, "\n") {
			if strings.HasPrefix(line, "# ") {
				title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
				break
			}
		}
		n := RecalledNote{Path: note.Path, Title: title, Kind: note.Frontmatter["type"], Source: note.Frontmatter["source"], LastVerified: note.Frontmatter["last_verified"], Body: note.Body, Issues: []string{}}
		for _, field := range RequiredFields {
			if note.Frontmatter[field] == "" {
				n.Issues = append(n.Issues, "missing "+field)
			}
		}
		if n.LastVerified != "" {
			t, err := time.Parse("2006-01-02", n.LastVerified)
			if err != nil {
				n.Issues = append(n.Issues, "invalid last_verified date")
			} else if now.Sub(t) > StaleDays*24*time.Hour {
				n.Issues = append(n.Issues, "last verification is older than 90 days")
			} else if t.After(now) {
				n.Issues = append(n.Issues, "last verification date is in the future")
			}
		}
		out = append(out, n)
	}
	return out
}
