package cli

import "runtime/debug"

// Version is set at build time with -ldflags "-X github.com/knightlybuilds-tech/baker/internal/cli.Version=v0.1.0".
// Release builds and `make build` set it. When it's unset, resolveVersion falls
// back to the module version Go records in the binary, so
// `go install github.com/knightlybuilds-tech/baker/cmd/baker@v0.1.0` reports v0.1.0.
var Version = "dev"

// resolveVersion prefers an ldflags version, then the module version from
// build info, then "dev". "(devel)" means a local build with no VCS info.
func resolveVersion(ldflags string, info func() (*debug.BuildInfo, bool)) string {
	if ldflags != "" && ldflags != "dev" {
		return ldflags
	}
	if bi, ok := info(); ok && bi != nil {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}
