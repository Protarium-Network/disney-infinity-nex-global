// Package nex wires the PRUDP authentication and secure servers.
package nex

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/v2/ticket-granting"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/v2/ticket-granting"

	"github.com/Protarium-Network/disney-infinity-nex-global/globals"
)

// newServer builds a PRUDP server with the parameters (access key, NEX
// version) of the instance's game.
func newServer(inst *globals.Instance) *nex.PRUDPServer {
	game := inst.Game

	server := nex.NewPRUDPServer()
	server.AccessKey = game.Network.AccessKey
	server.LibraryVersions.SetDefault(game.Network.LibraryVersion)
	server.PRUDPV0Settings.LegacyConnectionSignature = game.LegacyConnectionSignature()
	server.PRUDPV1Settings.LegacyConnectionSignature = game.LegacyConnectionSignature()
	server.ByteStreamSettings.UseStructureHeader = game.UseStructureHeader()

	return server
}

// StartAuthenticationServer builds and runs the authentication server of a
// game. It blocks.
func StartAuthenticationServer(inst *globals.Instance) {
	game := inst.Game

	inst.AuthenticationServer = newServer(inst)

	inst.AuthenticationEndpoint = nex.NewPRUDPEndPoint(game.Network.AuthStreamID)
	inst.AuthenticationEndpoint.ServerAccount = inst.AuthenticationServerAccount
	inst.AuthenticationEndpoint.AccountDetailsByPID = inst.AccountDetailsByPID
	inst.AuthenticationEndpoint.AccountDetailsByUsername = inst.AccountDetailsByUsername
	inst.AuthenticationServer.BindPRUDPEndPoint(inst.AuthenticationEndpoint)

	globals.LogRMC(game.Name+" Authentication", inst.AuthenticationEndpoint)

	ticketGrantingProtocol := ticket_granting.NewProtocol()
	inst.AuthenticationEndpoint.RegisterServiceProtocol(ticketGrantingProtocol)
	commonTicketGranting := common_ticket_granting.NewCommonProtocol(ticketGrantingProtocol)

	if game.Protocols.InsecureLogin {
		commonTicketGranting.EnableInsecureLogin()
	}

	secureStationURL := types.NewStationURL("")
	secureStationURL.SetURLType(constants.StationURLPRUDPS)
	secureStationURL.SetAddress(inst.SecureHost)
	secureStationURL.SetPortNumber(uint16(inst.SecurePort))
	secureStationURL.SetConnectionID(1)
	secureStationURL.SetPrincipalID(types.NewPID(2))
	secureStationURL.SetStreamID(game.Network.SecureStreamID)
	secureStationURL.SetStreamType(game.Network.SecureStreamType)
	secureStationURL.SetType(uint8(constants.StationURLFlagPublic))

	commonTicketGranting.ValidateLoginData = inst.ValidateLoginData
	commonTicketGranting.SecureStationURL = secureStationURL
	commonTicketGranting.BuildName = types.NewString(game.Network.BuildName)
	commonTicketGranting.SecureServerAccount = inst.SecureServerAccount

	globals.Logger.Successf("[%s] Authentication server listening on UDP %d", game.Name, inst.AuthPort)
	inst.AuthenticationServer.Listen(inst.AuthPort)
}
