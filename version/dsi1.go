package version

import "github.com/PretendoNetwork/nex-go/v2"

var DisneyInfinity1 = register(NewGame(
	"disney-infinity-1.0",
	"Disney Infinity 1.0",
	"PN_DSI1",
	WithTitleIDs("00039E00"),
	WithGameServerIDs("101AC900"),
	WithNetwork("2ef57176", nex.NewLibraryVersion(2, 0, 0), ""),
	WithPorts(29100, 29101),
	WithInsecureMode(false, true),
	WithAllMatchmaking(),
	WithStandardFeatures(true, true, true),
))
