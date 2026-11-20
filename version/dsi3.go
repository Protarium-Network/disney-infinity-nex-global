package version

import "github.com/PretendoNetwork/nex-go/v2"

var DisneyInfinity3 = register(NewGame(
	"disney-infinity-3.0",
	"Disney Infinity 3.0",
	"PN_DSI3",
	WithGameServerIDs("101AC900"),
	WithNetwork("83a601e5", nex.NewLibraryVersion(3, 9, 1), ""),
	WithPorts(29200, 29201),
	WithInsecureMode(false, true),
	WithAllMatchmaking(),
	WithStandardFeatures(true, false, false),
))
