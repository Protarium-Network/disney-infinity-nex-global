package main

import (
	"testing"

	"github.com/Protarium-Network/disney-infinity-nex-global/globals"
	"github.com/Protarium-Network/disney-infinity-nex-global/version"
)

func TestCheckPorts(t *testing.T) {
	a := &globals.Instance{Game: version.DisneyInfinity1, AuthPort: 1, SecurePort: 2}
	b := &globals.Instance{Game: version.DisneyInfinity3, AuthPort: 3, SecurePort: 4}
	if err := checkPorts([]*globals.Instance{a, b}); err != nil {
		t.Fatalf("distinct ports rejected: %v", err)
	}

	b.SecurePort = 1
	if err := checkPorts([]*globals.Instance{a, b}); err == nil {
		t.Fatal("expected an error when two games share a port")
	}
}

func TestSelectedGamesDefaultsToAll(t *testing.T) {
	t.Setenv("GAMES", "")
	games, err := selectedGames()
	if err != nil || len(games) != len(version.Games) {
		t.Fatalf("got %d games, err=%v", len(games), err)
	}

	t.Setenv("GAMES", "disney-infinity-3.0")
	if games, err := selectedGames(); err != nil || len(games) != 1 || games[0] != version.DisneyInfinity3 {
		t.Fatalf("got %v, err=%v", games, err)
	}

	t.Setenv("GAMES", "nope")
	if _, err := selectedGames(); err == nil {
		t.Fatal("expected an error for an unknown game")
	}
}
