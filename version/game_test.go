package version

import (
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
)

func TestLookup(t *testing.T) {
	g, err := Lookup("disney-infinity-1.0")
	if err != nil || g != DisneyInfinity1 {
		t.Fatalf("Lookup = %v, %v", g, err)
	}
	if _, err := Lookup("nope"); err == nil {
		t.Fatal("expected an error for an unknown game")
	}
}

func TestEnv(t *testing.T) {
	if got := DisneyInfinity1.Env("AUTH_PORT"); got != "PN_DSI1_AUTH_PORT" {
		t.Fatalf("Env = %q", got)
	}
}

// Games running in the same process must not share a port, an env prefix or
// an access key, otherwise clients get confused between them.
func TestGamesDoNotCollide(t *testing.T) {
	ports := map[int]string{}
	prefixes := map[string]string{}
	keys := map[string]string{}

	for slug, g := range Games {
		for _, p := range []int{g.Network.AuthPort, g.Network.SecurePort} {
			if p < 1 || p > 65535 {
				t.Errorf("%s: invalid port %d", slug, p)
			}
			if other, ok := ports[p]; ok {
				t.Errorf("port %d used by both %s and %s", p, other, slug)
			}
			ports[p] = slug
		}
		if other, ok := prefixes[g.EnvPrefix]; ok {
			t.Errorf("env prefix %s used by both %s and %s", g.EnvPrefix, other, slug)
		}
		prefixes[g.EnvPrefix] = slug
		if other, ok := keys[g.Network.AccessKey]; ok {
			t.Errorf("access key %s used by both %s and %s", g.Network.AccessKey, other, slug)
		}
		keys[g.Network.AccessKey] = slug
	}

	if len(Games) < 2 {
		t.Fatalf("expected at least 2 games, got %d", len(Games))
	}
}

func TestNEXVersionFlags(t *testing.T) {
	cases := []struct {
		v      *nex.LibraryVersion
		legacy bool
		header bool
	}{
		{nex.NewLibraryVersion(2, 0, 0), true, false},
		{nex.NewLibraryVersion(3, 2, 0), false, false},
		{nex.NewLibraryVersion(3, 5, 2), false, true},
	}
	for _, c := range cases {
		g := NewGame("x", "x", "X", WithNetwork("k", c.v, ""))
		if g.LegacyConnectionSignature() != c.legacy || g.UseStructureHeader() != c.header {
			t.Errorf("NEX %d.%d.%d: legacy=%v header=%v", c.v.Major, c.v.Minor, c.v.Patch,
				g.LegacyConnectionSignature(), g.UseStructureHeader())
		}
	}
}
