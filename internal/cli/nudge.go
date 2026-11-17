package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/knightlybuilds-tech/baker/internal/config"
	"github.com/knightlybuilds-tech/baker/internal/continuity"
	"github.com/spf13/cobra"
)

func nudgeCmd() *cobra.Command {
	var agent string
	cmd := &cobra.Command{
		Use:    "nudge",
		Args:   cobra.ArbitraryArgs,
		Hidden: true,
		Short:  "Ask the agent to save crumbs when a session did work but saved none (run by the Stop hook)",
		Long: `Run by Claude Code's Stop hook. If this session edited several files and
saved no crumbs, it replies with {"decision":"block","reason":...} so the
agent considers saving decisions and lessons before stopping. It fires at
most once per stretch of work, never while a previous nudge is still being
answered, and never fails: any problem means silence.

Turn it off with nudge = false in baker.toml (then baker sync), or set
BAKER_NO_NUDGE=1.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() { _ = recover() }()
			if out := runNudge(agent, os.Stdin); out != nil {
				_ = json.NewEncoder(os.Stdout).Encode(out)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "claude", "which agent is calling")
	return cmd
}

type nudgeInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
	StopHookActive bool   `json:"stop_hook_active"`
}

type nudgeReply struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

func runNudge(agent string, stdin io.Reader) *nudgeReply {
	if agent != "claude" || os.Getenv("BAKER_NO_NUDGE") != "" {
		return nil
	}
	payload, _ := io.ReadAll(io.LimitReader(stdin, 1<<20))
	var in nudgeInput
	if json.Unmarshal(payload, &in) != nil || in.TranscriptPath == "" || in.StopHookActive {
		return nil // stop_hook_active: we're already inside a nudge; never loop
	}
	wd := in.Cwd
	if wd == "" {
		wd, _ = os.Getwd()
	}
	root, ok := config.FindRoot(wd)
	if !ok {
		return nil
	}
	if p, err := config.LoadProject(root); err == nil && p.Nudge != nil && !*p.Nudge {
		return nil
	}
	session := in.SessionID
	if session == "" {
		session = strings.TrimSuffix(filepath.Base(in.TranscriptPath), filepath.Ext(in.TranscriptPath))
	}
	work := continuity.ClaudeWork(in.TranscriptPath, root)
	fire, save := continuity.ShouldNudge(work, continuity.LoadNudgeState(root, session))
	if save != nil {
		continuity.SaveNudgeState(root, session, *save)
	}
	if !fire {
		return nil
	}
	return &nudgeReply{Decision: "block", Reason: continuity.NudgeMessage}
}
