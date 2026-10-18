// Package cli wires baker's commands.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/knightlybuilds-tech/baker/internal/config"
	"github.com/knightlybuilds-tech/baker/internal/memory"
	"github.com/knightlybuilds-tech/baker/internal/scaffold"
	"github.com/spf13/cobra"
)

// Version is set at build time with -ldflags "-X github.com/knightlybuilds-tech/baker/internal/cli.Version=v0.1.0".
var Version = "dev"

func Execute() error {
	root := &cobra.Command{
		Use:   "baker",
		Short: "Scaffold AI-ready projects: one AGENTS.md, memory, skills and continuity for every agent",
		Long: `baker bakes projects that every AI coding agent understands on day one.

It writes one source of truth (AGENTS.md), bridges it to each agent
(Claude Code, Codex, Gemini CLI, Cursor, Copilot, OpenCode, Aider...),
sets up plain-markdown memory, shared skills, MCP config, and lets you
continue a task in another agent or account when a session ends.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
	}
	root.AddCommand(newCmd(), initCmd(), syncCmd(), doctorCmd(), globalCmd(),
		checkpointCmd(), handoffCmd(), resumeCmd(), listCmd(), recipesCmd())
	// After a plain `go install`, the first run is usually by full path:
	// tell the user how to put baker on PATH. Hooks (checkpoint) stay silent.
	if len(os.Args) < 2 || os.Args[1] != "checkpoint" {
		printPathHint(os.Stderr)
	}
	return root.Execute()
}

// projectRoot finds the baker project containing the working directory.
func projectRoot() (string, error) {
	wd, _ := os.Getwd()
	root, ok := config.FindRoot(wd)
	if !ok {
		return "", fmt.Errorf("not inside a baker project (no baker.toml found) — run `baker init` first")
	}
	return root, nil
}

// applyPlan prints the plan, asks for confirmation unless yes, and applies it.
func applyPlan(p *scaffold.Plan, title string, yes bool) (bool, error) {
	fmt.Println(title)
	p.Print(os.Stdout)
	if !p.HasWrites() {
		return false, nil
	}
	if !yes && !confirm("Apply these changes?") {
		fmt.Println("Nothing written.")
		return false, nil
	}
	if err := p.Apply(); err != nil {
		return false, err
	}
	return true, nil
}

func confirm(q string) bool {
	fmt.Printf("%s [y/N] ", q)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

func printFindings(fs []memory.Finding) (errors int) {
	for _, f := range fs {
		sym := map[memory.Level]string{memory.OK: "✓", memory.Warn: "!", memory.Error: "✗"}[f.Level]
		fmt.Printf("  %s %s\n", sym, f.Msg)
		if f.Level == memory.Error {
			errors++
		}
	}
	return errors
}

func splitList(s string) []string {
	var out []string
	for _, x := range strings.Split(s, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}
