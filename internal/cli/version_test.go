package cli

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	info := func(v string, ok bool) func() (*debug.BuildInfo, bool) {
		return func() (*debug.BuildInfo, bool) {
			if !ok {
				return nil, false
			}
			return &debug.BuildInfo{Main: debug.Module{Path: "github.com/knightlybuilds-tech/baker", Version: v}}, true
		}
	}
	for _, tc := range []struct {
		name, ldflags, module string
		ok                    bool
		want                  string
	}{
		{"ldflags wins", "v0.2.0", "v0.1.0", true, "v0.2.0"},
		{"go install @tag", "dev", "v0.1.0", true, "v0.1.0"},
		{"empty ldflags", "", "v0.1.0", true, "v0.1.0"},
		{"local build without VCS", "dev", "(devel)", true, "dev"},
		{"no build info", "dev", "", false, "dev"},
	} {
		if got := resolveVersion(tc.ldflags, info(tc.module, tc.ok)); got != tc.want {
			t.Errorf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}
