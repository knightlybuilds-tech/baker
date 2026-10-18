package scaffold

import (
	"path/filepath"

	"github.com/roguenet/baker/internal/config"
	"github.com/roguenet/baker/internal/recipe"
)

// Prepare loads (or creates) the project manifest for root, resolves the
// recipe, refreshes detected facts and merges recipe MCP servers.
// recipeName and agentIDs override what baker.toml says when non-empty.
func Prepare(root, recipeName string, agentIDs []string) (*config.Project, *recipe.Recipe, error) {
	proj, err := config.LoadProject(root)
	if err != nil {
		if !config.IsNotExist(err) {
			return nil, nil, err
		}
		g := config.LoadGlobal()
		proj = &config.Project{Name: filepath.Base(root), Agents: g.Agents}
	}
	proj.Standard = config.Standard
	if recipeName != "" {
		proj.Recipe = recipeName
	}
	if proj.Recipe == "" {
		proj.Recipe = recipe.DetectRecipe(root)
	}
	if len(agentIDs) > 0 {
		proj.Agents = agentIDs
	}
	r, err := recipe.Load(proj.Recipe, root)
	if err != nil {
		return nil, nil, err
	}
	proj.RecipeVersion = r.Version
	proj.Facts = r.ResolveFacts(root)
	if len(proj.Facts) == 0 {
		proj.Facts = nil
	}
	for _, s := range r.MCPFor(root) {
		if !hasServer(proj.MCP, s.Name) {
			proj.MCP = append(proj.MCP, s)
		}
	}
	return proj, r, nil
}

func hasServer(list []config.MCPServer, name string) bool {
	for _, s := range list {
		if s.Name == name {
			return true
		}
	}
	return false
}
