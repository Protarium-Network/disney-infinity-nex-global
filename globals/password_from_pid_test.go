package globals

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/plogger-go"

	"github.com/Protarium-Network/disney-infinity-nex-global/version"
)

func testInstance(secret []byte) *Instance {
	Logger = plogger.NewLogger()
	return &Instance{Game: version.DisneyInfinity1, NEXPasswordSecret: secret, KerberosPassword: "k"}
}

func TestPasswordFromPIDDerivedIsStableAndKeyed(t *testing.T) {
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(i)
	}
	inst := testInstance(secret)

	a, code := inst.PasswordFromPID(types.NewPID(1234))
	if code != 0 || a == "" {
		t.Fatalf("expected a password, got code=%d value=%q", code, a)
	}
	if b, _ := inst.PasswordFromPID(types.NewPID(1234)); b != a {
		t.Fatalf("derivation not stable: %q != %q", a, b)
	}
	if c, _ := inst.PasswordFromPID(types.NewPID(1235)); c == a {
		t.Fatal("different PIDs produced the same password")
	}

	inst.NEXPasswordSecret = make([]byte, 16) // too short
	if _, code := inst.PasswordFromPID(types.NewPID(1234)); code == 0 {
		t.Fatal("expected failure with a short secret")
	}
}

func TestPasswordsDifferBetweenGames(t *testing.T) {
	a := testInstance(make([]byte, 32))
	secret := make([]byte, 32)
	secret[0] = 1
	b := testInstance(secret)

	pa, _ := a.PasswordFromPID(types.NewPID(7))
	pb, _ := b.PasswordFromPID(types.NewPID(7))
	if pa == pb {
		t.Fatal("instances with different secrets derived the same password")
	}
}

func TestPasswordFromPIDLocalReadsSettings(t *testing.T) {
	inst := testInstance(nil)
	inst.LocalAuthMode = true

	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"accounts":[{"username":"42","password":"hunter2"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	inst.SettingsPath = path

	if got, code := inst.PasswordFromPID(types.NewPID(42)); code != 0 || got != "hunter2" {
		t.Fatalf("got (%q, %d), want (\"hunter2\", 0)", got, code)
	}
	if _, code := inst.PasswordFromPID(types.NewPID(99)); code == 0 {
		t.Fatal("expected failure for a PID missing from settings.json")
	}
}
