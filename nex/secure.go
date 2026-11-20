package nex

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	commonutility "github.com/PretendoNetwork/nex-protocols-common-go/v2/utility"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	utility "github.com/PretendoNetwork/nex-protocols-go/v2/utility"

	"github.com/Protarium-Network/disney-infinity-nex-global/globals"
)

// StartSecureServer builds and runs the secure server of a game. It blocks.
func StartSecureServer(inst *globals.Instance) {
	game := inst.Game
	protocols := game.Protocols
	tag := game.Name + " Secure"

	inst.SecureServer = newServer(inst)

	inst.SecureEndpoint = nex.NewPRUDPEndPoint(game.Network.SecureStreamID)
	inst.SecureEndpoint.IsSecureEndPoint = true
	inst.SecureEndpoint.ServerAccount = inst.SecureServerAccount
	inst.SecureEndpoint.AccountDetailsByPID = inst.AccountDetailsByPID
	inst.SecureEndpoint.AccountDetailsByUsername = inst.AccountDetailsByUsername
	inst.SecureServer.BindPRUDPEndPoint(inst.SecureEndpoint)

	globals.LogRMC(tag, inst.SecureEndpoint)

	inst.SecureEndpoint.OnConnectionEnded(func(connection *nex.PRUDPConnection) {
		fmt.Printf("[%s] PID=%d disconnected\n", tag, uint64(connection.PID()))
	})

	if protocols.SecureConnection {
		secureProtocol := secure.NewProtocol()
		inst.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
		commonSecure := commonsecure.NewCommonProtocol(secureProtocol)
		commonSecure.ValidateLoginData = inst.ValidateLoginData
		commonSecure.CreateReportDBRecord = protocols.ReportHandler

		if protocols.InsecureRegister {
			commonSecure.EnableInsecureRegister()
		}
	}

	if protocols.NATTraversal {
		natTraversalProtocol := nattraversal.NewProtocol()
		inst.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
		commonnattraversal.NewCommonProtocol(natTraversalProtocol)
	}

	hasDB := inst.DB != nil
	if !hasDB {
		globals.Logger.Warningf("[%s] No database: matchmaking and utility are not registered", tag)
	}

	if hasDB && game.UsesMatchmaking() {
		inst.MatchmakingManager = common_globals.NewMatchmakingManager(inst.SecureEndpoint, inst.DB)
		inst.MatchmakingManager.GetUserFriendPIDs = globals.GetUserFriendPIDs

		if protocols.MatchMaking {
			p := matchmaking.NewProtocol()
			inst.SecureEndpoint.RegisterServiceProtocol(p)
			commonmatchmaking.NewCommonProtocol(p).SetManager(inst.MatchmakingManager)
		}

		if protocols.MatchMakingExt {
			p := matchmakingext.NewProtocol()
			inst.SecureEndpoint.RegisterServiceProtocol(p)
			commonmatchmakingext.NewCommonProtocol(p).SetManager(inst.MatchmakingManager)
		}

		if protocols.MatchmakeExtension {
			p := matchmakeextension.NewProtocol()
			inst.SecureEndpoint.RegisterServiceProtocol(p)
			commonmatchmakeextension.NewCommonProtocol(p).SetManager(inst.MatchmakingManager)
		}
	}

	// The utility protocol is needed by most games (unique IDs) and is cheap.
	if hasDB {
		inst.UtilityManager = common_globals.NewUtilityManager(inst.SecureEndpoint, inst.DB)
		utilityProtocol := utility.NewProtocol()
		inst.SecureEndpoint.RegisterServiceProtocol(utilityProtocol)
		commonutility.NewCommonProtocol(utilityProtocol).SetManager(inst.UtilityManager)
	}

	if protocols.Ranking {
		globals.Logger.Warningf("[%s] Ranking is enabled in the game config but not implemented yet; requests will be refused", tag)
	}
	if protocols.StorageManager {
		globals.Logger.Warningf("[%s] StorageManager is enabled in the game config but not implemented yet; requests will be refused", tag)
	}

	globals.Logger.Successf("[%s] Secure server listening on UDP %d", game.Name, inst.SecurePort)
	inst.SecureServer.Listen(inst.SecurePort)
}
