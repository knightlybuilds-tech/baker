package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knightlybuilds-tech/baker/internal/agents"
	"github.com/knightlybuilds-tech/baker/internal/config"
	"github.com/knightlybuilds-tech/baker/internal/continuity"
	"github.com/spf13/cobra"
)

func checkpointCmd() *cobra.Command {
	var agent string
	var quiet bool
	cmd := &cobra.Command{
		Use:   "checkpoint [hook-payload]",
		Args:  cobra.ArbitraryArgs,
		Short: "Snapshot git state and the latest transcript (run by agent hooks)",
		Long: `Saves a mechanical snapshot to .baker/checkpoints/: branch, uncommitted
changes, and the tail of the latest agent transcript. No LLM, no network.
Agent hooks call this after every turn; with --quiet it never fails, so it
can't break your agent. Hook input is read as JSON from stdin (Claude Code,
Gemini CLI) or from the first argument (Codex notify).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCheckpoint(agent, quiet, args)
			if quiet {
				return nil
			}
			return err
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "which agent is calling (claude, codex, gemini...)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "print nothing and always exit 0 (for hooks)")
	return cmd
}

func runCheckpoint(agent string, quiet bool, args []string) error {
	var hook *continuity.HookInput
	var payload []byte
	if len(args) > 0 {
		payload = []byte(args[0])
	} else if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeCharDevice == 0 {
		payload, _ = io.ReadAll(io.LimitReader(os.Stdin, 1<<20))
	}
	var h continuity.HookInput
	if len(payload) > 0 && json.Unmarshal(payload, &h) == nil {
		hook = &h
	}
	wd, _ := os.Getwd()
	if hook != nil && hook.Cwd != "" {
		wd = hook.Cwd
	}
	root, ok := config.FindRoot(wd)
	if !ok {
		if quiet {
			return nil
		}
		return fmt.Errorf("not inside a baker project")
	}
	_, path, err := continuity.Take(root, agent, hook)
	if err == nil && !quiet {
		rel, _ := filepath.Rel(root, path)
		fmt.Println("✓ checkpoint saved:", rel)
	}
	return err
}

func handoffCmd() *cobra.Command {
	var from, out string
	var copy bool
	cmd := &cobra.Command{
		Use:   "handoff",
		Short: "Build a resume brief so any agent can continue the current task",
		Long: `Merges memory/hot.md, the latest checkpoint and the latest transcript from
any agent (Claude Code, Codex, Gemini CLI) into one short prompt. Paste it
into any agent or chat, or use ` + "`baker resume <agent>`" + ` to launch one with it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			brief, err := buildBrief(from)
			if err != nil {
				return err
			}
			switch {
			case out != "":
				if err := os.WriteFile(out, []byte(brief), 0o644); err != nil {
					return err
				}
				fmt.Fprintln(os.Stderr, "✓ brief written to", out)
			case copy:
				if err := clipboard(brief); err != nil {
					return err
				}
				fmt.Fprintln(os.Stderr, "✓ brief copied to clipboard — paste it into any agent")
			default:
				fmt.Print(brief)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "only read transcripts from this agent (claude, codex, gemini)")
	cmd.Flags().BoolVarP(&copy, "copy", "c", false, "copy the brief to the clipboard")
	cmd.Flags().StringVarP(&out, "out", "o", "", "write the brief to a file")
	return cmd
}

func resumeCmd() *cobra.Command {
	var from string
	cmd := &cobra.Command{
		Use:   "resume <agent>",
		Short: "Launch an agent with the resume brief as its first prompt",
		Example: `  baker resume codex             # Claude hit its limit → continue in Codex
  baker resume claude --from codex`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, ok := agents.Get(args[0])
			if !ok || a.Launch == nil {
				return fmt.Errorf("can't launch %q; use `baker handoff --copy` and paste into it", args[0])
			}
			brief, err := buildBrief(from)
			if err != nil {
				return err
			}
			argv := a.Launch(brief)
			bin, err := exec.LookPath(argv[0])
			if err != nil {
				_ = clipboard(brief)
				return fmt.Errorf("%s isn't installed or not on PATH (brief copied to clipboard instead)", argv[0])
			}
			fmt.Fprintf(os.Stderr, "→ starting %s with the resume brief…\n", a.Name)
			c := exec.Command(bin, argv[1:]...)
			c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
			return c.Run()
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "only read transcripts from this agent")
	return cmd
}

func buildBrief(from string) (string, error) {
	root, err := projectRoot()
	if err != nil {
		return "", err
	}
	name := filepath.Base(root)
	if p, err := config.LoadProject(root); err == nil && p.Name != "" {
		name = p.Name
	}
	return continuity.Brief(root, name, from)
}

func clipboard(s string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	}
	cmd.Stdin = strings.NewReader(s)
	return cmd.Run()
}
