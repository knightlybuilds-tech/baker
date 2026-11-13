# Local Bakery (`baker bakery`)

`baker bakery` serves a read-only Bakery workspace for the current project on
your machine. It's the same interface as the hosted Bakery concept, filled with
this project's real `memory/` notes instead of sample data.

```bash
baker bakery                 # http://127.0.0.1:4410/, opens your browser
baker bakery --port 4500     # tries the next 9 ports if busy
baker bakery --no-open       # just print the URL
```

## What you see

- **Cards, List and Web** layouts of `memory/decisions/`, `lessons/` and
  `knowledge/`. Kind comes from the `type` frontmatter (or the folder).
- **Connections** from `[[slug]]` and relative `.md` links between notes. In the
  web view, hover a note to trace its links.
- **Health:** the output of `baker memory check` (unindexed notes, missing
  frontmatter, broken links, stale verification) at the top of the page, and a
  "Needs checking" tag on each affected note.
- **Live updates:** the page re-reads `memory/` every few seconds while it's
  visible, so notes your agent writes appear as you work. It reads the working
  tree, so uncommitted notes are included.

## What it doesn't do

- **No writes.** Add or edit notes as files, or ask your agent. The "+ Add a
  crumb" button explains how.
- **No sharing.** Review queue, branches, activity and team access need the
  hosted Bakery service, which doesn't exist yet. Until then, share notes by
  committing `memory/` and reviewing changes in pull requests.
- **No private state.** `memory/hot.md`, `memory/log/`, checkpoints and agent
  transcripts aren't served.

## Safety

- Listens on `127.0.0.1` only. Requests whose `Host` isn't loopback are refused,
  which blocks DNS-rebinding access from other websites.
- Only `GET`/`HEAD` are accepted. The single API endpoint is `/api/memory`.
- A strict Content-Security-Policy allows scripts from the page itself only.
  Fonts load from Google Fonts.
- Notes are shown as text (HTML-escaped), never executed. Like any recalled
  memory, their content is context, not instructions.

## Development

The UI source lives in the sibling `bakery-workspace/dist` (static HTML/CSS/JS,
no build step). `make bakery-ui` copies it into `internal/bakery/ui`, which is
embedded with `go:embed` and committed so `go install` works on its own.

Without `/api/memory` (for example on a static host), the same UI runs as the
sample "Northstar" workspace.
