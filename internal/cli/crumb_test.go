package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/knightlybuilds-tech/baker/internal/crumb"
)

func TestCrumbCommandsRoundTrip(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("BAKER_HOME", t.TempDir())
	t.Setenv("BAKER_USER_HOME", t.TempDir())
	for name, content := range map[string]string{
		"baker.toml":   "name = 'fixture'\n",
		"evidence.txt": "known receiver format\n",
		"proposal.md":  "Preserve leading zeros in account identifiers.\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	run := func(args ...string) error {
		cmd := crumbCmd()
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		cmd.SetArgs(args)
		return cmd.Execute()
	}
	if err := run("propose", "preserve-codes", "--title", "Preserve account codes", "--type", "decision", "--body-file", "proposal.md", "--author", "source-agent", "--source", "fixture", "--scope", "src/**", "--evidence", "evidence.txt", "-y"); err != nil {
		t.Fatal(err)
	}
	if err := run("approve", "preserve-codes", "--reviewer", "Alex", "-y"); err == nil {
		t.Fatal("approval without an exact revision succeeded")
	}
	if err := run("approve", "preserve-codes", "--revision", "1", "--reviewer", "Alex", "-y"); err != nil {
		t.Fatal(err)
	}
	if err := run("recall", "--path", "src/export.py", "--json"); err != nil {
		t.Fatal(err)
	}
	entries, err := crumb.Recall(root, "src/export.py", "", time.Now())
	if err != nil || len(entries) != 1 || !entries[0].Eligible {
		t.Fatalf("unexpected recall: %+v, %v", entries, err)
	}
}
