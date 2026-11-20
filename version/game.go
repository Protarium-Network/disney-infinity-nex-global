// Package version describes the per-game configuration of the server: the
// NEX/PRUDP parameters and the set of protocols a game needs.
package version

import (
	"fmt"
	"sort"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

type ReportHandler func(pid types.PID, reportType types.UInt32, data types.QBuffer) error

type NetworkConfig struct {
	AccessKey      string
	LibraryVersion *nex.LibraryVersion
	AuthStreamID   uint8
	SecureStreamID uint8
	// SecureStreamType is the stream type advertised to clients in the
	// secure station URL.
	SecureStreamType constants.StreamType
	BuildName        string
	// AuthPort and SecurePort are the default UDP ports of the game. They
	// must differ between games running in the same process. They can be
	// overridden with <EnvPrefix>_AUTH_PORT / <EnvPrefix>_SECURE_PORT.
	AuthPort   int
	SecurePort int
}

type ProtocolRegistry struct {
	TicketGranting bool
	InsecureLogin  bool

	SecureConnection bool
	InsecureRegister bool
	ReportHandler    ReportHandler

	NATTraversal       bool
	MatchMaking        bool
	MatchMakingExt     bool
	MatchmakeExtension bool
	Ranking            bool
	StorageManager     bool
}

type Game struct {
	Slug      string
	Name      string
	EnvPrefix string
	TitleIDs  []string
	// GameServerIDs are the Pretendo game server IDs used to validate login
	// tokens against the account server.
	GameServerIDs []string

	Network   NetworkConfig
	Protocols ProtocolRegistry
}

type GameOption func(*Game)

func NewGame(slug, name, envPrefix string, opts ...GameOption) *Game {
	g := &Game{
		Slug:      slug,
		Name:      name,
		EnvPrefix: envPrefix,
		Network: NetworkConfig{
			AuthStreamID:     1,
			SecureStreamID:   1,
			SecureStreamType: constants.StreamTypeRVSecure,
		},
		Protocols: ProtocolRegistry{
			TicketGranting:   true,
			SecureConnection: true,
			ReportHandler: func(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
				return nil
			},
		},
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Env returns the name of a per-game environment variable, e.g. PN_DSI1_AUTH_PORT.
func (g *Game) Env(suffix string) string {
	return g.EnvPrefix + "_" + suffix
}

func (g *Game) UsesMatchmaking() bool {
	p := g.Protocols
	return p.MatchMaking || p.MatchMakingExt || p.MatchmakeExtension
}

// LegacyConnectionSignature reports whether the PRUDP connection signature is
// the legacy one, which only applies to NEX versions below 3.2.0. A wrong
// value shows up as an endless CONNECT retransmit on the client.
func (g *Game) LegacyConnectionSignature() bool {
	return g.Network.LibraryVersion != nil && !g.Network.LibraryVersion.GreaterOrEqual("3.2.0")
}

// UseStructureHeader reports whether structures carry a header (NEX 3.5.0+).
func (g *Game) UseStructureHeader() bool {
	return g.Network.LibraryVersion != nil && g.Network.LibraryVersion.GreaterOrEqual("3.5.0")
}

func WithNetwork(accessKey string, version *nex.LibraryVersion, buildName string) GameOption {
	return func(g *Game) {
		g.Network.AccessKey = accessKey
		g.Network.LibraryVersion = version
		g.Network.BuildName = buildName
	}
}

func WithPorts(auth, secure int) GameOption {
	return func(g *Game) {
		g.Network.AuthPort = auth
		g.Network.SecurePort = secure
	}
}

func WithTitleIDs(ids ...string) GameOption {
	return func(g *Game) { g.TitleIDs = ids }
}

func WithGameServerIDs(ids ...string) GameOption {
	return func(g *Game) { g.GameServerIDs = ids }
}

func WithInsecureMode(login, register bool) GameOption {
	return func(g *Game) {
		g.Protocols.InsecureLogin = login
		g.Protocols.InsecureRegister = register
	}
}

func WithAllMatchmaking() GameOption {
	return func(g *Game) {
		g.Protocols.MatchMaking = true
		g.Protocols.MatchMakingExt = true
		g.Protocols.MatchmakeExtension = true
	}
}

func WithStandardFeatures(nat, ranking, storage bool) GameOption {
	return func(g *Game) {
		g.Protocols.NATTraversal = nat
		g.Protocols.Ranking = ranking
		g.Protocols.StorageManager = storage
	}
}

// Games is the registry of supported games, keyed by slug.
var Games = map[string]*Game{}

func register(g *Game) *Game {
	Games[g.Slug] = g
	return g
}

// Lookup returns the game registered under slug.
func Lookup(slug string) (*Game, error) {
	if g, ok := Games[slug]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("unknown game %q (available: %v)", slug, Slugs())
}

func Slugs() []string {
	slugs := make([]string, 0, len(Games))
	for s := range Games {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}
