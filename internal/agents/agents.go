// Package agents describes every AI coding agent baker supports: which files
// bridge it to AGENTS.md, where its MCP config and user-level instructions
// live, how to launch it, and where it keeps transcripts.
//
// Paths marked "verify" in docs/agents.md come from vendor docs as of
// October 2026 and should be re-checked when a tool changes.
package agents

import (
	"sort"

	"github.com/roguenet/baker/internal/config"
)

// MCPFormat says how a tool's MCP config file is shaped.
type MCPFormat int

const (
	MCPNone      MCPFormat = iota
	MCPServers             // JSON {"mcpServers": {name: {command, args, env}}}  — Claude, Cursor, Gemini
	MCPVSCode              // JSON {"servers": {name: {type: "stdio", command, args, env}}} — Copilot
	MCPOpenCode            // JSON {"mcp": {name: {type: "local", command: [...], environment}}}
	MCPCodexTOML           // TOML [mcp_servers.name] tables
)

type Agent struct {
	ID   string
	Name string

	// ReadsAgentsMD: the tool reads AGENTS.md natively, so it needs no bridge file.
	ReadsAgentsMD bool

	// Bridge file pointing at AGENTS.md (when ReadsAgentsMD is false).
	BridgeFile    string
	BridgeContent string

	// SkillsDir: if set, skills are mirrored here because the tool doesn't read .agents/skills.
	SkillsDir string

	MCPFile   string
	MCPFormat MCPFormat

	// GlobalFile: user-level instructions file, relative to the user's home.
	GlobalFile string
	// GlobalSkillsDir: user-level skills dir, relative to the user's home.
	GlobalSkillsDir string

	// Launch builds the command line that starts the tool with a first prompt.
	Launch func(prompt string) []string
}

var all = []Agent{
	{
		ID: "claude", Name: "Claude Code",
		BridgeFile: "CLAUDE.md", BridgeContent: "@AGENTS.md",
		SkillsDir: ".claude/skills",
		MCPFile:   ".mcp.json", MCPFormat: MCPServers,
		GlobalFile: ".claude/CLAUDE.md", GlobalSkillsDir: ".claude/skills",
		Launch: func(p string) []string { return []string{"claude", p} },
	},
	{
		ID: "codex", Name: "OpenAI Codex", ReadsAgentsMD: true,
		MCPFile: ".codex/config.toml", MCPFormat: MCPCodexTOML,
		GlobalFile: ".codex/AGENTS.md",
		Launch:     func(p string) []string { return []string{"codex", p} },
	},
	{
		ID: "gemini", Name: "Gemini CLI",
		BridgeFile: "GEMINI.md", BridgeContent: "@AGENTS.md",
		MCPFile: ".gemini/settings.json", MCPFormat: MCPServers,
		GlobalFile: ".gemini/GEMINI.md",
		Launch:     func(p string) []string { return []string{"gemini", "-i", p} },
	},
	{
		ID: "cursor", Name: "Cursor", ReadsAgentsMD: true,
		MCPFile: ".cursor/mcp.json", MCPFormat: MCPServers,
	},
	{
		ID: "copilot", Name: "GitHub Copilot", ReadsAgentsMD: true,
		MCPFile: ".vscode/mcp.json", MCPFormat: MCPVSCode,
	},
	{
		ID: "opencode", Name: "OpenCode", ReadsAgentsMD: true,
		MCPFile: "opencode.json", MCPFormat: MCPOpenCode,
		GlobalFile: ".config/opencode/AGENTS.md",
		Launch:     func(p string) []string { return []string{"opencode", "--prompt", p} },
	},
	{
		ID: "aider", Name: "Aider",
		BridgeFile: ".aider.conf.yml", BridgeContent: "read: AGENTS.md",
		Launch: func(p string) []string { return []string{"aider", "--message", p} },
	},
	{ID: "windsurf", Name: "Windsurf", ReadsAgentsMD: true},
	{ID: "amp", Name: "Amp", ReadsAgentsMD: true},
	{ID: "zed", Name: "Zed", ReadsAgentsMD: true},
	{ID: "cline", Name: "Cline", ReadsAgentsMD: true},
}

// Get returns an agent by id.
func Get(id string) (Agent, bool) {
	for _, a := range all {
		if a.ID == id {
			return a, true
		}
	}
	return Agent{}, false
}

// IDs lists all supported agent ids.
func IDs() []string {
	out := make([]string, 0, len(all))
	for _, a := range all {
		out = append(out, a.ID)
	}
	sort.Strings(out)
	return out
}

// For returns the agents enabled in a project, skipping unknown ids.
func For(p *config.Project) []Agent {
	var out []Agent
	for _, id := range p.Agents {
		if a, ok := Get(id); ok {
			out = append(out, a)
		}
	}
	return out
}

// LegacyFiles hide AGENTS.md from some tools (Zed reads the first match only).
var LegacyFiles = []string{".cursorrules", ".windsurfrules", ".rules", ".clinerules"}
