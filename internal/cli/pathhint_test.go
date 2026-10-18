package cli

import (
	"strings"
	"testing"
)

func TestPathAdvicePicksShell(t *testing.T) {
	cases := map[string]struct{ goos, shell, want string }{
		"zsh":     {"darwin", "/bin/zsh", "~/.zshrc"},
		"bash":    {"linux", "/bin/bash", "~/.bashrc"},
		"fish":    {"linux", "/usr/bin/fish", "fish_add_path"},
		"default": {"darwin", "", "~/.zshrc"},
		"windows": {"windows", "", "Environment Variables"},
	}
	for name, c := range cases {
		got := pathAdvice("/opt/go/bin", c.goos, c.shell)
		if !strings.Contains(got, c.want) || !strings.Contains(got, "/opt/go/bin") {
			t.Errorf("%s: want %q in:\n%s", name, c.want, got)
		}
	}
}

func TestPathHintSilentWhenDisabled(t *testing.T) {
	t.Setenv("BAKER_NO_PATH_HINT", "1")
	if h := pathHint(); h != "" {
		t.Errorf("expected no hint, got %q", h)
	}
}
