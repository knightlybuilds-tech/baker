package crumb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var clock = time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)

func fixture(t *testing.T) (string, Revision) {
	t.Helper()
	root := t.TempDir()
	write(t, root, "tests/retry.txt", "timeout accepted: do not retry\n")
	write(t, root, "memory/MEMORY.md", "# Team knowledge\n\nHuman notes stay here.\n")
	r := Revision{ID: "safe-retry", Kind: "decision", Title: "Retry only after checking the receipt", Body: "Check the receipt before retrying an ambiguous timeout.", Source: "INC-42", Author: "Priya via Claude", Scopes: []string{"src/billing/**"}, Evidence: []Evidence{{Path: "tests/retry.txt"}}}
	return root, r
}
func write(t *testing.T, root, path, body string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}
func propose(t *testing.T, root string, r Revision) Revision {
	t.Helper()
	p, r, err := Propose(root, r, clock)
	if err != nil {
		t.Fatal(err)
	}
	if err = p.Apply(); err != nil {
		t.Fatal(err)
	}
	return r
}
func approve(t *testing.T, root string, r Revision) Revision {
	t.Helper()
	p, r, err := Approve(root, r.ID, r.Revision, "Alex", clock)
	if err != nil {
		t.Fatal(err)
	}
	if err = p.Apply(); err != nil {
		t.Fatal(err)
	}
	return r
}
func recall(t *testing.T, root, path string) []Entry {
	t.Helper()
	out, err := Recall(root, path, "", clock)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestLifecycleAndSupersession(t *testing.T) {
	root, r := fixture(t)
	p, _, err := Propose(root, r, clock)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(root, ".baker/crumbs")); !os.IsNotExist(err) {
		t.Fatal("planning wrote files")
	}
	_ = p
	r = propose(t, root, r)
	if got := recall(t, root, ""); len(got) != 0 {
		t.Fatal("proposal leaked into recall")
	}
	if _, err = os.Stat(filepath.Join(root, notePath(r))); !os.IsNotExist(err) {
		t.Fatal("proposing created approved note")
	}
	r = approve(t, root, r)
	got := recall(t, root, "src/billing/retry.go")
	if len(got) != 1 || !got[0].Eligible || got[0].Revision.Revision != 1 {
		t.Fatalf("unexpected recall: %+v", got)
	}
	if got := recall(t, root, "src/auth/login.go"); len(got) != 0 {
		t.Fatal("scope leakage")
	}
	r.Body = "Revised policy: keep the idempotency key when checking the receipt."
	r = propose(t, root, r)
	if got := recall(t, root, ""); len(got) != 1 || got[0].Revision.Revision != 1 {
		t.Fatal("open update replaced approved revision")
	}
	approve(t, root, r)
	got = recall(t, root, "")
	if len(got) != 1 || got[0].Revision.Revision != 2 || !strings.Contains(got[0].Body, "idempotency") {
		t.Fatalf("supersession failed: %+v", got)
	}
	index, _ := os.ReadFile(filepath.Join(root, "memory/MEMORY.md"))
	if !strings.Contains(string(index), "Human notes stay here.") || strings.Count(string(index), "](") != 1 {
		t.Fatal("index was overwritten or duplicated")
	}
	if _, _, err := Approve(root, r.ID, 1, "Alex", clock); err == nil {
		t.Fatal("obsolete revision accepted")
	}
}

func TestChangedEvidenceBlocksApprovalAndRecall(t *testing.T) {
	root, r := fixture(t)
	r = propose(t, root, r)
	write(t, root, "tests/retry.txt", "changed evidence")
	if _, _, err := Approve(root, r.ID, 1, "Alex", clock); err == nil {
		t.Fatal("changed evidence approved")
	}
	write(t, root, "tests/retry.txt", "timeout accepted: do not retry\n")
	approve(t, root, r)
	write(t, root, "tests/retry.txt", "changed again")
	got := recall(t, root, "")
	if len(got) != 1 || got[0].Eligible || got[0].Body != "" || len(got[0].Issues) == 0 {
		t.Fatalf("stale guidance exposed: %+v", got)
	}
}

func TestManualEditsAndOldApprovalsAreNotSilentlyTrusted(t *testing.T) {
	root, r := fixture(t)
	r = approve(t, root, propose(t, root, r))
	r.Body = "New proposal"
	update := propose(t, root, r)
	write(t, root, notePath(r), "human changed the note")
	if _, _, err := Approve(root, update.ID, 2, "Alex", clock); err == nil {
		t.Fatal("concurrent edit overwritten")
	}
	got := recall(t, root, "")
	if got[0].Eligible {
		t.Fatal("manual note edit trusted")
	}
	write(t, root, notePath(r), markdown(r))
	aged, err := Recall(root, "", "", clock.Add(91*24*time.Hour))
	if err != nil || aged[0].Eligible {
		t.Fatal("old approval trusted")
	}
}

func TestProjectIsolationAndPathBoundaries(t *testing.T) {
	root, r := fixture(t)
	approve(t, root, propose(t, root, r))
	other, _ := fixture(t)
	if got := recall(t, other, ""); len(got) != 0 {
		t.Fatal("cross-project leakage")
	}
	for _, id := range []string{"../outside", "UPPER", "a/b"} {
		r.ID = id
		if _, _, err := Propose(other, r, clock); err == nil {
			t.Fatalf("accepted %s", id)
		}
	}
	r.ID = "safe-retry"
	r.Evidence = []Evidence{{Path: "../secret"}}
	if _, _, err := Propose(other, r, clock); err == nil {
		t.Fatal("escaping evidence accepted")
	}
	outside := t.TempDir()
	write(t, outside, "secret.txt", "secret")
	if err := os.Symlink(outside, filepath.Join(other, "external")); err != nil {
		t.Skip("symlinks unavailable")
	}
	r.Evidence = []Evidence{{Path: "external/secret.txt"}}
	if _, _, err := Propose(other, r, clock); err == nil {
		t.Fatal("symlink evidence accepted")
	}
}

func TestMissingAndCorruptRecordsFailClosed(t *testing.T) {
	root, r := fixture(t)
	r = approve(t, root, propose(t, root, r))
	if err := os.Remove(filepath.Join(root, "tests/retry.txt")); err != nil {
		t.Fatal(err)
	}
	if recall(t, root, "")[0].Eligible {
		t.Fatal("missing evidence trusted")
	}
	write(t, root, recordPath(r.ID, r.Revision), `{"schema":999}`)
	if _, err := Records(root); err == nil {
		t.Fatal("unknown schema accepted")
	}
}

func TestDuplicateProposalAndInvalidScope(t *testing.T) {
	root, r := fixture(t)
	propose(t, root, r)
	if _, _, err := Propose(root, r, clock); err == nil {
		t.Fatal("duplicate proposal accepted")
	}
	for _, scope := range []string{"../**", "/etc/**", ""} {
		other, r := fixture(t)
		r.Scopes = []string{scope}
		if _, _, err := Propose(other, r, clock); err == nil {
			t.Fatalf("invalid scope accepted: %s", scope)
		}
	}
}
