// Package crumb implements a local, Git-reviewable memory lifecycle. Approval
// records are attestations, not authentication; repository permissions govern trust.
package crumb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knightlybuilds-tech/baker/internal/blocks"
	"github.com/knightlybuilds-tech/baker/internal/scaffold"
)

const Schema = 1

var validID = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Evidence struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type Revision struct {
	Schema     int        `json:"schema"`
	ID         string     `json:"id"`
	Revision   int        `json:"revision"`
	Kind       string     `json:"kind"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Source     string     `json:"source"`
	Author     string     `json:"author"`
	Scopes     []string   `json:"scopes"`
	Evidence   []Evidence `json:"evidence"`
	Status     string     `json:"status"`
	Created    string     `json:"created"`
	Reviewer   string     `json:"reviewer,omitempty"`
	Approved   string     `json:"approved,omitempty"`
	BaseSHA256 string     `json:"base_sha256,omitempty"`
}

type Entry struct {
	Revision
	Eligible bool     `json:"eligible"`
	Issues   []string `json:"issues"`
}

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func recordPath(id string, rev int) string {
	return filepath.Join(".baker", "crumbs", id, fmt.Sprintf("%06d.json", rev))
}
func notePath(r Revision) string {
	dir := map[string]string{"decision": "decisions", "lesson": "lessons", "knowledge": "knowledge"}[r.Kind]
	return filepath.Join("memory", dir, r.ID+".md")
}

// localPath rejects escaping paths, including symlinked ancestors. It supports
// not-yet-created output paths by resolving the nearest existing ancestor.
func localPath(root, rel string) (string, error) {
	if rel == "" || filepath.IsAbs(rel) {
		return "", fmt.Errorf("expected a project-relative path: %q", rel)
	}
	clean := filepath.Clean(rel)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes project: %q", rel)
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	full := filepath.Join(base, clean)
	existing := full
	for {
		_, err = os.Lstat(existing)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return "", err
		}
		existing = parent
	}
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", err
	}
	within, err := filepath.Rel(base, resolved)
	if err != nil || within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("symlink escapes project: %q", rel)
	}
	return full, nil
}

func readLocal(root, rel string) ([]byte, error) {
	p, err := localPath(root, rel)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

func Records(root string) ([]Revision, error) {
	if _, err := localPath(root, ".baker/crumbs"); err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(root, ".baker", "crumbs", "*", "*.json"))
	if err != nil {
		return nil, err
	}
	out := []Revision{}
	for _, file := range files {
		rel, _ := filepath.Rel(root, file)
		b, err := readLocal(root, rel)
		if err != nil {
			return nil, err
		}
		var r Revision
		if err = json.Unmarshal(b, &r); err != nil {
			return nil, fmt.Errorf("%s: %w", rel, err)
		}
		if r.Schema != Schema || !validID.MatchString(r.ID) || r.Revision < 1 || filepath.Clean(rel) != recordPath(r.ID, r.Revision) {
			return nil, fmt.Errorf("invalid crumb record: %s", rel)
		}
		if _, ok := map[string]bool{"decision": true, "lesson": true, "knowledge": true}[r.Kind]; !ok {
			return nil, fmt.Errorf("invalid crumb kind in %s", rel)
		}
		if r.Status != "proposed" && r.Status != "approved" {
			return nil, fmt.Errorf("invalid status in %s", rel)
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ID == out[j].ID {
			return out[i].Revision < out[j].Revision
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func validate(r Revision) error {
	if !validID.MatchString(r.ID) || len(r.ID) > 100 {
		return fmt.Errorf("id must be a lowercase slug of at most 100 characters")
	}
	if r.Kind != "decision" && r.Kind != "lesson" && r.Kind != "knowledge" {
		return fmt.Errorf("kind must be decision, lesson, or knowledge")
	}
	for _, v := range []string{r.Title, r.Source, r.Author} {
		if strings.TrimSpace(v) == "" || strings.ContainsAny(v, "\r\n") {
			return fmt.Errorf("title, source and author must be nonempty single lines")
		}
	}
	if strings.TrimSpace(r.Body) == "" || len(strings.Fields(r.Body)) > 400 {
		return fmt.Errorf("body must contain 1 to 400 words")
	}
	if len(r.Evidence) == 0 {
		return fmt.Errorf("at least one local evidence file is required")
	}
	if len(r.Scopes) == 0 {
		return fmt.Errorf("at least one --scope is required (use ** for the project)")
	}
	for _, scope := range r.Scopes {
		if _, err := scopePattern(scope); err != nil {
			return err
		}
	}
	return nil
}

func plan(root string) *scaffold.Plan { p := scaffold.NewPlan(root, true); p.NoLock = true; return p }
func put(p *scaffold.Plan, path string, b []byte) error {
	if _, err := localPath(p.Root, path); err != nil {
		return err
	}
	p.Managed(path, string(b))
	return nil
}
func encode(r Revision) []byte { b, _ := json.MarshalIndent(r, "", "  "); return append(b, '\n') }

// Propose leaves the current approved note and index unchanged.
func Propose(root string, r Revision, now time.Time) (*scaffold.Plan, Revision, error) {
	if err := validate(r); err != nil {
		return nil, r, err
	}
	records, err := Records(root)
	if err != nil {
		return nil, r, err
	}
	r.Schema = Schema
	r.Revision = 1
	r.Status = "proposed"
	r.Created = now.UTC().Format(time.RFC3339Nano)
	r.Reviewer = ""
	r.Approved = ""
	r.BaseSHA256 = ""
	for _, old := range records {
		if old.ID != r.ID {
			continue
		}
		if old.Status == "proposed" {
			return nil, r, fmt.Errorf("%s already has an open proposal", r.ID)
		}
		r.Revision = old.Revision + 1
		if old.Kind != r.Kind {
			return nil, r, fmt.Errorf("cannot change a crumb's kind")
		}
	}
	current, err := readLocal(root, notePath(r))
	if err == nil {
		r.BaseSHA256 = digest(current)
	} else if !os.IsNotExist(err) {
		return nil, r, err
	}
	for i, e := range r.Evidence {
		b, err := readLocal(root, e.Path)
		if err != nil {
			return nil, r, fmt.Errorf("evidence %s: %w", e.Path, err)
		}
		r.Evidence[i] = Evidence{filepath.ToSlash(filepath.Clean(e.Path)), digest(b)}
	}
	p := plan(root)
	err = put(p, recordPath(r.ID, r.Revision), encode(r))
	return p, r, err
}

func markdown(r Revision) string {
	return fmt.Sprintf("---\nid: %s\nrevision: %d\ntype: %s\nstatus: approved\nas_of: %s\nlast_verified: %s\nsource: %s\nreviewer: %s\n---\n\n# %s\n\n%s\n", r.ID, r.Revision, r.Kind, r.Approved[:10], r.Approved[:10], strings.ReplaceAll(r.Source, ": ", " — "), r.Reviewer, r.Title, strings.TrimSpace(r.Body))
}

func evidenceIssues(root string, r Revision) []string {
	out := []string{}
	for _, e := range r.Evidence {
		b, err := readLocal(root, e.Path)
		if err != nil {
			out = append(out, "evidence unavailable: "+e.Path)
		} else if digest(b) != e.SHA256 {
			out = append(out, "evidence changed: "+e.Path)
		}
	}
	return out
}

// Approve checks both the previously reviewed note and the evidence snapshot.
// The reviewer is an explicit local attestation, not an authenticated identity.
func Approve(root, id string, rev int, reviewer string, now time.Time) (*scaffold.Plan, Revision, error) {
	var r Revision
	if !validID.MatchString(id) || rev < 1 || strings.TrimSpace(reviewer) == "" || strings.ContainsAny(reviewer, "\r\n") {
		return nil, r, fmt.Errorf("valid id, revision and single-line reviewer are required")
	}
	records, err := Records(root)
	if err != nil {
		return nil, r, err
	}
	found := false
	for _, x := range records {
		if x.ID == id && x.Revision == rev {
			r = x
			found = true
		}
		if x.ID == id && x.Revision > rev {
			return nil, r, fmt.Errorf("a newer revision exists")
		}
	}
	if !found || r.Status != "proposed" {
		return nil, r, fmt.Errorf("open proposal %s revision %d not found", id, rev)
	}
	if err := validate(r); err != nil {
		return nil, r, err
	}
	current, err := readLocal(root, notePath(r))
	if err != nil && !os.IsNotExist(err) {
		return nil, r, err
	}
	base := ""
	if err == nil {
		base = digest(current)
	}
	if base != r.BaseSHA256 {
		return nil, r, fmt.Errorf("approved note changed since proposal; reconcile before approval")
	}
	if issues := evidenceIssues(root, r); len(issues) > 0 {
		return nil, r, fmt.Errorf("cannot approve: %s", strings.Join(issues, "; "))
	}
	r.Status = "approved"
	r.Reviewer = reviewer
	r.Approved = now.UTC().Format(time.RFC3339Nano)
	p := plan(root)
	if err = put(p, recordPath(id, rev), encode(r)); err != nil {
		return nil, r, err
	}
	if err = put(p, notePath(r), []byte(markdown(r))); err != nil {
		return nil, r, err
	}
	if _, err = localPath(root, "memory/MEMORY.md"); err != nil {
		return nil, r, err
	}
	if _, exists := p.Read("memory/MEMORY.md"); !exists {
		p.Once("memory/MEMORY.md", "# Memory index\n")
	}
	p.Block("memory/MEMORY.md", "crumb:"+id, "1", fmt.Sprintf("- [%s](%s) — approved revision %d; use `baker experimental crumb recall` to check evidence", r.Title, strings.TrimPrefix(filepath.ToSlash(notePath(r)), "memory/"), rev), "", blocks.Append)
	return p, r, nil
}

func scopePattern(glob string) (*regexp.Regexp, error) {
	if glob == "" || strings.ContainsAny(glob, "\\\r\n") || strings.HasPrefix(glob, "/") {
		return nil, fmt.Errorf("scope must be a relative slash-separated glob")
	}
	for _, part := range strings.Split(glob, "/") {
		if part == ".." {
			return nil, fmt.Errorf("scope cannot escape project")
		}
	}
	p := regexp.QuoteMeta(glob)
	p = strings.ReplaceAll(p, `\*\*`, "@@DOUBLE@@")
	p = strings.ReplaceAll(p, `\*`, "[^/]*")
	p = strings.ReplaceAll(p, "@@DOUBLE@@", ".*")
	p = strings.ReplaceAll(p, `\?`, "[^/]")
	return regexp.Compile("^" + p + "$")
}

// Recall returns only the latest approved revision per id. Ineligible entries
// expose issues but omit the body, preventing accidental use of stale guidance.
func Recall(root, path, query string, now time.Time) ([]Entry, error) {
	records, err := Records(root)
	if err != nil {
		return nil, err
	}
	latest := map[string]Revision{}
	for _, r := range records {
		if r.Status == "approved" {
			latest[r.ID] = r
		}
	}
	out := []Entry{}
	for _, r := range latest {
		if _, err := validateApproved(r); err != nil {
			return nil, err
		}
		if path != "" {
			matches := false
			for _, g := range r.Scopes {
				re, err := scopePattern(g)
				if err != nil {
					return nil, err
				}
				matches = matches || re.MatchString(filepath.ToSlash(path))
			}
			if !matches {
				continue
			}
		}
		if query != "" && !strings.Contains(strings.ToLower(r.Title+" "+r.Body), strings.ToLower(query)) {
			continue
		}
		issues := evidenceIssues(root, r)
		b, err := readLocal(root, notePath(r))
		if err != nil || string(b) != markdown(r) {
			issues = append(issues, "approved note differs from review record")
		}
		approved, _ := time.Parse(time.RFC3339Nano, r.Approved)
		if now.Sub(approved) > 90*24*time.Hour {
			issues = append(issues, "approval is older than 90 days")
		}
		e := Entry{Revision: r, Eligible: len(issues) == 0, Issues: issues}
		if !e.Eligible {
			e.Body = ""
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func validateApproved(r Revision) (time.Time, error) {
	if err := validate(r); err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339Nano, r.Approved)
	if err != nil || r.Reviewer == "" {
		return t, fmt.Errorf("invalid approval for %s", r.ID)
	}
	return t, nil
}
