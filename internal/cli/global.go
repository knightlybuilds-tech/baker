package cli

import (
	"fmt"

	"github.com/roguenet/baker/internal/config"
	"github.com/roguenet/baker/internal/scaffold"
	"github.com/spf13/cobra"
)

func globalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Manage your personal layer in ~/.baker (rules, skills, cross-project memory)",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Create ~/.baker with config, your personal AGENTS.md and global memory",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := scaffold.GlobalInit()
			if err != nil {
				return err
			}
			if _, err := applyPlan(p, "→ baker global init: "+config.Home(), true); err != nil {
				return err
			}
			fmt.Printf("\nNext: write your personal rules in %s/AGENTS.md, then run `baker global sync`.\n", config.Home())
			return nil
		},
	})
	var force bool
	sync := &cobra.Command{
		Use:   "sync",
		Short: "Push your personal rules and skills into every installed agent's user-level files",
		Long: `Writes a managed block with ~/.baker/AGENTS.md (plus a pointer to
~/.baker/memory) into each installed agent's own user-level instructions:
~/.claude/CLAUDE.md, ~/.codex/AGENTS.md, ~/.gemini/GEMINI.md,
~/.config/opencode/AGENTS.md. Copies ~/.baker/skills to ~/.claude/skills.
Registers Codex's notify hook for automatic checkpoints.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := scaffold.GlobalSync(force)
			if err != nil {
				return err
			}
			_, err = applyPlan(p, "→ baker global sync (paths relative to your home folder)", true)
			return err
		},
	}
	sync.Flags().BoolVar(&force, "force", false, "overwrite blocks/files you edited")
	cmd.AddCommand(sync)
	return cmd
}
