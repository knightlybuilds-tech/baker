package bakery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	root := t.TempDir()
	write(t, root, "memory/MEMORY.md", "# Memory index\n- [Codes](decisions/codes.md)\n")
	write(t, root, "memory/decisions/codes.md", "---\ntype: decision\nas_of: 2026-11-01\nlast_verified: 2026-11-01\nsource: tests/codes_test.go\n---\n# Account codes stay strings\n\n007 and 7 differ.\n\nSee [[retries]] and [[missing-note]].\n")
	write(t, root, "memory/lessons/retries.md", "---\ntype: lesson\nas_of: 2026-01-01\nlast_verified: 2026-01-01\n---\n# Retry only on 429\n\nBody.\n")
	write(t, root, "memory/hot.md", "# Goal\nprivate task state\n")
	return root
}

func TestReadBuildsNotesLinksAndFindings(t *testing.T) {
	now := time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)
	snap := Read(fixture(t), "demo", now)
	if !snap.HasMemory || snap.Mode != "local" || snap.Project != "demo" {
		t.Fatalf("unexpected snapshot header: %+v", snap)
	}
	if len(snap.Notes) != 2 {
		t.Fatalf("want 2 atomic notes (hot.md excluded), got %d", len(snap.Notes))
	}
	byID := map[string]Note{}
	for _, n := range snap.Notes {
		byID[n.ID] = n
	}
	codes := byID["decisions/codes.md"]
	if codes.Title != "Account codes stay strings" || codes.Kind != "decision" {
		t.Errorf("codes note parsed wrong: %+v", codes)
	}
	if strings.HasPrefix(codes.Body, "#") || !strings.HasPrefix(codes.Body, "007 and 7 differ.") {
		t.Errorf("title line should be stripped from body, got %q", codes.Body)
	}
	if len(codes.Links) != 1 || codes.Links[0] != "lessons/retries.md" {
		t.Errorf("links should resolve existing slugs to ids and drop missing ones, got %v", codes.Links)
	}
	retries := byID["lessons/retries.md"]
	if len(retries.Issues) == 0 {
		t.Errorf("stale note without source should report issues")
	}
	var warned bool
	for _, f := range snap.Findings {
		if f.Level == "warn" && strings.Contains(f.Msg, "missing-note") {
			warned = true
		}
	}
	if !warned {
		t.Errorf("memory check findings should include the broken link, got %+v", snap.Findings)
	}
}

func TestReadWithoutMemoryFolder(t *testing.T) {
	snap := Read(t.TempDir(), "empty", time.Now())
	if snap.HasMemory || len(snap.Notes) != 0 || snap.Notes == nil {
		t.Fatalf("want empty, non-nil notes and has_memory=false, got %+v", snap)
	}
}

func TestHandlerServesAPIAndUI(t *testing.T) {
	h := Handler(fixture(t), "demo")

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:4410/api/memory", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || !strings.Contains(rec.Header().Get("Content-Type"), "json") {
		t.Fatalf("api: %d %s", rec.Code, rec.Header().Get("Content-Type"))
	}
	var snap Snapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &snap); err != nil || len(snap.Notes) != 2 {
		t.Fatalf("api body: %v %s", err, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "private task state") {
		t.Error("hot.md must not be exposed")
	}

	req = httptest.NewRequest(http.MethodGet, "http://localhost:4410/", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "bakery") {
		t.Fatalf("ui index: %d", rec.Code)
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Error("missing CSP header")
	}
}

func TestHandlerRejectsNonLocalHostsAndWrites(t *testing.T) {
	h := Handler(fixture(t), "demo")
	for _, tc := range []struct {
		method, url string
		want        int
	}{
		{http.MethodGet, "http://evil.example:4410/api/memory", http.StatusForbidden},
		{http.MethodGet, "http://192.168.1.5:4410/api/memory", http.StatusForbidden},
		{http.MethodPost, "http://127.0.0.1:4410/api/memory", http.StatusMethodNotAllowed},
		{http.MethodGet, "http://[::1]:4410/api/memory", http.StatusOK},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.url, nil))
		if rec.Code != tc.want {
			t.Errorf("%s %s: got %d want %d", tc.method, tc.url, rec.Code, tc.want)
		}
	}
}
