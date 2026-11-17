package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/knightlybuilds-tech/baker/internal/continuity"
	"github.com/knightlybuilds-tech/baker/internal/memory"
	"github.com/knightlybuilds-tech/baker/internal/scaffold"
	"github.com/spf13/cobra"
)

func memoryAddCmd() *cobra.Command {
	var n memory.NewNote
	var bodyFile string
	var dryRun, asJSON bool
	cmd := &cobra.Command{
		Use:   "add",
		Args:  cobra.NoArgs,
		Short: "Save one crumb: validated, indexed in MEMORY.md, ready to commit",
		Long: `Save one atomic note (a decision, lesson or piece of knowledge) to memory/.

Baker checks the note before writing it: a one-sentence claim as the title,
a body under ~200 words, a source, no duplicate slug, no secrets. It then
writes memory/<kind>s/<slug>.md with frontmatter and adds one line to
memory/MEMORY.md. Ordinary notes need no approval.

Agents: use this instead of writing note files by hand.`,
		Example: `  baker memory add --kind lesson \
    --title "Retry payments only on HTTP 429 and 503" \
    --body "Other failures need a separate investigation; retrying them double-charged in #412." \
    --source "billing/retry.go; incident 2026-11-02" --tags payments,retries --link idempotency-keys

  echo "Longer body..." | baker memory add --kind decision --title "..." --source "..." --body-file -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := projectRoot()
			if err != nil {
				return err
			}
			if bodyFile != "" {
				var r io.Reader = os.Stdin
				if bodyFile != "-" {
					f, err := os.Open(bodyFile)
					if err != nil {
						return err
					}
					defer f.Close()
					r = f
				}
				b, err := io.ReadAll(io.LimitReader(r, 64<<10))
				if err != nil {
					return err
				}
				n.Body = string(b)
			}
			for _, field := range []string{n.Title, n.Body, n.Source, n.Hook} {
				if continuity.Redact(field) != field {
					return errors.New("this note looks like it contains a secret (token, key or password): remove it; memory is committed with the code")
				}
			}
			prep, err := memory.Prepare(filepath.Join(root, "memory"), n, time.Now())
			if err != nil {
				return err
			}
			if !dryRun {
				p := scaffold.NewPlan(root, false)
				p.NoLock = true
				p.Write(prep.Path, prep.Content)
				p.Write(prep.IndexPath, prep.Index)
				if err := p.Apply(); err != nil {
					return err
				}
			}
			if asJSON {
				return json.NewEncoder(os.Stdout).Encode(struct {
					*memory.Prepared
					Written bool `json:"written"`
				}{prep, !dryRun})
			}
			verb := "Saved"
			if dryRun {
				verb = "Would save"
				fmt.Print(prep.Content)
				fmt.Println()
			}
			fmt.Printf("✓ %s %s and indexed it in %s\n", verb, prep.Path, prep.IndexPath)
			for _, w := range prep.Warnings {
				fmt.Printf("  ! %s\n", w)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&n.Kind, "kind", "", "decision, lesson or knowledge (required)")
	f.StringVar(&n.Title, "title", "", "the claim as one full sentence (required)")
	f.StringVar(&n.Body, "body", "", "detail and why it matters, under ~200 words")
	f.StringVar(&bodyFile, "body-file", "", "read the body from a file, or - for stdin")
	f.StringVar(&n.Source, "source", "", "where this comes from: file path, commit, issue, doc, or \"session YYYY-MM-DD\" (required)")
	f.StringSliceVar(&n.Tags, "tags", nil, "comma-separated tags")
	f.StringSliceVar(&n.Links, "link", nil, "slug of a related note (repeatable)")
	f.StringVar(&n.Hook, "hook", "", "short phrase for the MEMORY.md line (default: the body's first sentence)")
	f.StringVar(&n.Confidence, "confidence", "", "high, medium or low (default medium)")
	f.StringVar(&n.AsOf, "as-of", "", "date the fact describes, YYYY-MM-DD (default today)")
	f.BoolVar(&dryRun, "dry-run", false, "validate and show the note without writing")
	f.BoolVar(&asJSON, "json", false, "print the result as JSON")
	return cmd
}
