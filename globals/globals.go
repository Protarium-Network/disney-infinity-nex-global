package globals

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"

	"github.com/Protarium-Network/disney-infinity-nex-global/version"
)

var Logger *plogger.Logger

// Instance is the runtime state of one game. Several instances run side by
// side in the same process; each has its own ports, access key, accounts,
// database and PRUDP servers.
type Instance struct {
	Game *version.Game

	AuthPort   int
	SecurePort int
	// SecureHost is the host clients are told to reach the secure server on.
	SecureHost string

	// KerberosPassword is the shared secret for the Quazal Authentication /
	// Rendez-Vous system accounts (PIDs 1 and 2).
	KerberosPassword string

	// LocalAuthMode enables a self-contained preservation deployment:
	// accounts come from settings.json and the login token is not validated.
	// Keep it disabled for any shared/public deployment.
	LocalAuthMode bool
	SettingsPath  string

	// NEXTokenAESKey is the AES-256 key an external account server uses to
	// seal the NEX login token. Only used when LocalAuthMode is false.
	NEXTokenAESKey []byte
	// NEXPasswordSecret is the HMAC secret used to derive a stable per-PID NEX
	// password. Only used when LocalAuthMode is false; must match the account
	// server's secret.
	NEXPasswordSecret []byte

	AuthenticationServerAccount *nex.Account
	SecureServerAccount         *nex.Account

	AuthenticationServer   *nex.PRUDPServer
	AuthenticationEndpoint *nex.PRUDPEndPoint
	SecureServer           *nex.PRUDPServer
	SecureEndpoint         *nex.PRUDPEndPoint

	// DB backs matchmaking and utility. Use a separate database per game.
	DB                 *sql.DB
	MatchmakingManager *common_globals.MatchmakingManager
	UtilityManager     *common_globals.UtilityManager
}

// GetUserFriendPIDs is used by friend-scoped matchmaking searches. There is
// no friends system on this server, so it always returns an empty list.
func GetUserFriendPIDs(pid uint32) []uint32 {
	return []uint32{}
}
