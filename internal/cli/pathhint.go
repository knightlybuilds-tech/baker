package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// goBinDirs lists where `go install` puts binaries: $GOBIN, $GOPATH/bin, ~/go/bin.
func goBinDirs() []string {
	var dirs []string
	if d := os.Getenv("GOBIN"); d != "" {
		dirs = append(dirs, d)
	}
	if p := os.Getenv("GOPATH"); p != "" {
		for _, part := range filepath.SplitList(p) {
			dirs = append(dirs, filepath.Join(part, "bin"))
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "go", "bin"))
	}
	return dirs
}

// pathHint returns advice when baker was installed with `go install` but
// its folder isn't on PATH (so typing `baker` fails and agent hooks can't
// run it). It returns "" when there's nothing to say.
func pathHint() string {
	if os.Getenv("BAKER_NO_PATH_HINT") != "" {
		return ""
	}
	if _, err := exec.LookPath("baker"); err == nil {
		return ""
	}
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	dir := filepath.Dir(exe)
	for _, d := range goBinDirs() {
		if sameDir(d, dir) {
			return pathAdvice(dir, runtime.GOOS, os.Getenv("SHELL"))
		}
	}
	return ""
}

func sameDir(a, b string) bool {
	ra, err1 := filepath.EvalSymlinks(a)
	rb, err2 := filepath.EvalSymlinks(b)
	if err1 == nil && err2 == nil {
		return ra == rb
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// pathAdvice builds the message for a folder, OS and login shell.
func pathAdvice(dir, goos, shell string) string {
	shown := dir
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(dir, home) {
		shown = "$HOME" + strings.TrimPrefix(dir, home)
	}
	var fix string
	switch {
	case goos == "windows":
		fix = fmt.Sprintf("  Add %s to your Path: Settings → System → About → Advanced system settings → Environment Variables.\n  Then open a new terminal.", dir)
	case strings.HasSuffix(shell, "fish"):
		fix = fmt.Sprintf("  fish_add_path %s", shown)
	case strings.HasSuffix(shell, "bash"):
		fix = fmt.Sprintf("  echo 'export PATH=\"%s:$PATH\"' >> ~/.bashrc && source ~/.bashrc", shown)
	default: // zsh is the macOS default
		fix = fmt.Sprintf("  echo 'export PATH=\"%s:$PATH\"' >> ~/.zshrc && source ~/.zshrc", shown)
	}
	return fmt.Sprintf("baker is installed in %s, but that folder isn't on your PATH.\n"+
		"Typing `baker` won't work, and agent hooks can't save checkpoints. Fix it with:\n\n%s\n", dir, fix)
}

func printPathHint(w io.Writer) {
	if hint := pathHint(); hint != "" {
		fmt.Fprintln(w, hint)
	}
}
