package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalMemoryDoesNotDependOnReviewRecords(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	files := map[string]string{
		"baker.toml":                       "name = 'standalone'\n",
		"memory/lessons/codes.md":          "# Preserve account codes\nKeep identifiers as strings.\n",
		".baker/crumbs/broken/000001.json": "malformed experimental data",
	}
	for path, body := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, mode := range []string{"list", "recall"} {
		cmd := memoryCmd()
		cmd.SetArgs([]string{mode, "--query", "codes", "--json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("local %s depended on experimental review state: %v", mode, err)
		}
	}
}
