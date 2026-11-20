package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"

	"github.com/Protarium-Network/disney-infinity-nex-global/database"
	"github.com/Protarium-Network/disney-infinity-nex-global/globals"
	"github.com/Protarium-Network/disney-infinity-nex-global/version"
)

func init() {
	globals.Logger = plogger.NewLogger()
	common_globals.Logger = globals.Logger

	if err := godotenv.Load(); err != nil {
		globals.Logger.Warning("No .env file loaded, using the process environment only")
	}
}

// selectedGames returns the games to run: the comma-separated slugs of the
// GAMES variable, or every registered game when it is unset.
func selectedGames() ([]*version.Game, error) {
	value := strings.TrimSpace(os.Getenv("GAMES"))
	if value == "" {
		games := make([]*version.Game, 0, len(version.Games))
		for _, slug := range version.Slugs() {
			games = append(games, version.Games[slug])
		}
		return games, nil
	}

	var games []*version.Game
	for _, slug := range strings.Split(value, ",") {
		game, err := version.Lookup(strings.TrimSpace(slug))
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func requiredEnv(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("%s environment variable not set", name)
	}
	return value, nil
}

// portFromEnv returns the port in the environment variable, or def when the
// variable is unset.
func portFromEnv(name string, def int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return def, nil
	}

	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s is not a valid port: %s", name, value)
	}
	return port, nil
}

// newInstance reads the per-game configuration (<EnvPrefix>_* variables) and
// connects the game's database.
func newInstance(game *version.Game) (*globals.Instance, error) {
	inst := &globals.Instance{
		Game:             game,
		KerberosPassword: "password",
		SettingsPath:     strings.TrimSpace(os.Getenv(game.Env("SETTINGS_PATH"))),
		LocalAuthMode:    os.Getenv(game.Env("LOCAL_MODE")) == "1",
	}

	var err error
	if inst.AuthPort, err = portFromEnv(game.Env("AUTH_PORT"), game.Network.AuthPort); err != nil {
		return nil, err
	}
	if inst.SecurePort, err = portFromEnv(game.Env("SECURE_PORT"), game.Network.SecurePort); err != nil {
		return nil, err
	}

	if inst.SecureHost, err = requiredEnv(game.Env("SECURE_HOST")); err != nil {
		return nil, err
	}

	if kerberos := strings.TrimSpace(os.Getenv(game.Env("KERBEROS_PASSWORD"))); kerberos != "" {
		inst.KerberosPassword = kerberos
	} else {
		globals.Logger.Warningf("[%s] %s not set, using the default Kerberos password", game.Name, game.Env("KERBEROS_PASSWORD"))
	}

	if inst.LocalAuthMode {
		globals.Logger.Warningf("[%s] Local mode is enabled; use it only in an isolated preservation environment", game.Name)
	} else {
		keyName, secretName := game.Env("NEX_TOKEN_AES_KEY"), game.Env("NEX_PASSWORD_SECRET")

		value, err := requiredEnv(keyName)
		if err != nil {
			return nil, err
		}
		if inst.NEXTokenAESKey, err = hex.DecodeString(value); err != nil || len(inst.NEXTokenAESKey) != 32 {
			return nil, fmt.Errorf("%s must be a 64-character hexadecimal AES-256 key", keyName)
		}

		value, err = requiredEnv(secretName)
		if err != nil {
			return nil, err
		}
		if inst.NEXPasswordSecret, err = hex.DecodeString(value); err != nil || len(inst.NEXPasswordSecret) < 32 {
			return nil, fmt.Errorf("%s must be at least 32 bytes encoded as hexadecimal", secretName)
		}
	}

	inst.InitAccounts()

	// The database is optional: without it the game still serves login and
	// secure connection, but matchmaking and utility are disabled.
	if uri := strings.TrimSpace(os.Getenv(game.Env("POSTGRES_URI"))); uri != "" {
		if inst.DB, err = database.ConnectPostgres(uri); err != nil {
			return nil, fmt.Errorf("[%s] %w", game.Name, err)
		}
		globals.Logger.Successf("[%s] Connected to Postgres", game.Name)
	} else {
		globals.Logger.Warningf("[%s] %s not set: matchmaking and utility are disabled", game.Name, game.Env("POSTGRES_URI"))
	}

	return inst, nil
}

// checkPorts makes sure no two instances listen on the same port.
func checkPorts(instances []*globals.Instance) error {
	used := map[int]string{}
	for _, inst := range instances {
		for _, port := range []int{inst.AuthPort, inst.SecurePort} {
			if other, ok := used[port]; ok {
				return fmt.Errorf("UDP port %d is used by both %s and %s", port, other, inst.Game.Name)
			}
			used[port] = inst.Game.Name
		}
	}
	return nil
}
