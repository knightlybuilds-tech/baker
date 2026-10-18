// Package assets holds everything baker writes into projects: the core
// AGENTS.md block, built-in skills, memory templates and built-in recipes.
package assets

import "embed"

//go:embed all:core all:recipes all:skills
var FS embed.FS
