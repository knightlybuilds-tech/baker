package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/knightlybuilds-tech/baker/internal/agents"
	"github.com/knightlybuilds-tech/baker/internal/config"
	"github.com/knightlybuilds-tech/baker/internal/recipe"
	"github.com/knightlybuilds-tech/baker/internal/scaffold"
	"github.com/spf13/cobra"
)

var validName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func newCmd() *cobra.Command {
	var recipeName, agentList string
	var bare, noGit bool
	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new project folder with the full AI setup",
		Example: `  baker new shop --recipe laravel
  baker new notes --recipe generic --agents claude,codex
  baker new api --recipe laravel --bare   # skip the stack installer`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !validName.MatchString(name) {
				return fmt.Errorf("invalid project name %q (letters, numbers, . _ - only)", name)
			}
			wd, _ := os.Getwd()
			root := filepath.Join(wd, name)
			if entries, err := os.ReadDir(root); err == nil && len(entries) > 0 {
				return fmt.Errorf("%s already exists and isn't empty — use `baker init` inside it instead", root)
			}
			if recipeName == "" {
				recipeName = config.LoadGlobal().DefaultRecipe
			}
			r, err := recipe.Load(recipeName, "")
			if err != nil {
				return err
			}

			install, err := r.InstallCommand(recipe.Data{Name: name})
			if err != nil {
				return err
			}
			if install != "" && !bare {
				fmt.Printf("→ Running the %s installer: %s\n", r.Name, install)
				if err := shell(wd, install); err != nil {
					return fmt.Errorf("installer failed: %w", err)
				}
			}
			if err := os.MkdirAll(root, 0o755); err != nil {
				return err
			}

			proj, r, err := scaffold.Prepare(root, recipeName, splitList(agentList))
			if err != nil {
				return err
			}
			p := scaffold.NewPlan(root, false)
			if err := scaffold.Build(p, proj, r); err != nil {
				return err
			}
			if _, err := applyPlan(p, fmt.Sprintf("→ Baking %s (recipe: %s)", name, r.Name), true); err != nil {
				return err
			}
			if !noGit {
				gitInit(root)
			}
			_ = config.Register(root, proj)
			fmt.Printf("\n✓ %s is ready. Next:\n  cd %s\n  edit AGENTS.md → Project notes (what this project is)\n  start any agent: claude, codex, gemini, cursor…\n", name, name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&recipeName, "recipe", "r", "", "stack recipe (see `baker recipes`)")
	cmd.Flags().StringVar(&agentList, "agents", "", "comma-separated agents to target (default from ~/.baker/config.toml)")
	cmd.Flags().BoolVar(&bare, "bare", false, "skip the recipe's installer; only add the AI layer")
	cmd.Flags().BoolVar(&noGit, "no-git", false, "don't run git init")
	return cmd
}

func initCmd() *cobra.Command {
	var recipeName, agentList string
	var yes, force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Add the AI setup to an existing project (shows a plan first)",
		Long: `Adds baker to the current folder. Detects the stack, shows every file it
will create or change, and asks before writing. Existing CLAUDE.md / GEMINI.md
instructions are moved into AGENTS.md, never thrown away.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := os.Getwd()
			proj, r, err := scaffold.Prepare(root, recipeName, splitList(agentList))
			if err != nil {
				return err
			}
			p := scaffold.NewPlan(root, force)
			if err := scaffold.Build(p, proj, r); err != nil {
				return err
			}
			applied, err := applyPlan(p, fmt.Sprintf("→ baker init (recipe: %s, agents: %s)", r.Name, strings.Join(proj.Agents, ", ")), yes)
			if err != nil || !applied {
				return err
			}
			_ = config.Register(root, proj)
			fmt.Println("\n✓ Done. Run `baker doctor` any time to check the setup.")
			return nil
		},
	}
	cmd.Flags().StringVarP(&recipeName, "recipe", "r", "", "stack recipe (default: detected)")
	cmd.Flags().StringVar(&agentList, "agents", "", "comma-separated agents to target")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "apply without asking")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite files and blocks you edited")
	return cmd
}

func syncCmd() *cobra.Command {
	var dryRun, force bool
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Regenerate bridges, skills, MCP config and managed blocks",
		Long: `Re-renders everything baker manages from AGENTS.md and baker.toml:
bridge files, skill mirrors, MCP configs, hooks, and the core/recipe blocks
in AGENTS.md. Files or blocks you edited by hand are left alone (--force to
overwrite). Run it after editing baker.toml or upgrading baker.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := projectRoot()
			if err != nil {
				return err
			}
			proj, r, err := scaffold.Prepare(root, "", nil)
			if err != nil {
				return err
			}
			p := scaffold.NewPlan(root, force)
			if err := scaffold.Build(p, proj, r); err != nil {
				return err
			}
			if dryRun {
				fmt.Println("→ baker sync (dry run)")
				p.Print(os.Stdout)
				return nil
			}
			_, err = applyPlan(p, "→ baker sync", true)
			return err
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would change without writing")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite files and blocks you edited")
	return cmd
}

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the project against the baker standard",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := projectRoot()
			if err != nil {
				return err
			}
			fmt.Println("→ baker doctor:", root)
			if n := printFindings(scaffold.Doctor(root)); n > 0 {
				return fmt.Errorf("%d problem(s) found", n)
			}
			return nil
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List projects baker created or adopted",
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := config.LoadRegistry()
			if len(reg.Projects) == 0 {
				fmt.Println("No projects yet. Create one with `baker new <name>`.")
				return nil
			}
			for _, e := range reg.Projects {
				status := ""
				if _, err := os.Stat(e.Path); err != nil {
					status = "  (missing)"
				}
				fmt.Printf("  %-24s %-10s std %-5s %s%s\n", e.Name, e.Recipe, e.Standard, e.Path, status)
			}
			return nil
		},
	}
}

func recipesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recipes",
		Short: "List available stack recipes",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			root, _ := config.FindRoot(wd)
			for _, n := range recipe.Names(root) {
				r, err := recipe.Load(n, root)
				if err != nil {
					continue
				}
				fmt.Printf("  %-12s %s  [%s]\n", n, r.Description, r.Source)
			}
			fmt.Printf("\nAgents: %s\n", strings.Join(agents.IDs(), ", "))
			return nil
		},
	}
}

func shell(dir, command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = dir, os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func gitInit(root string) {
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		return
	}
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	_ = cmd.Run()
}
