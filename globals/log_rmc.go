package globals

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
)

// LogRMC prints every incoming RMC request, which is how unknown protocols
// and methods get discovered while reversing a game.
func LogRMC(tag string, endpoint *nex.PRUDPEndPoint) {
	endpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()
		if request == nil {
			return
		}

		fmt.Printf("[%s] PID=%d ProtocolID=%s (0x%02X) MethodID=%s (0x%02X) CallID=%d ParamsLen=%d\n",
			tag, uint64(packet.Sender().PID()), GetProtocolNameByID(request.ProtocolID), request.ProtocolID,
			GetMethodNameByID(request.ProtocolID, uint32(request.MethodID)), request.MethodID, request.CallID, len(request.Parameters))
	})

	endpoint.OnError(func(err *nex.Error) {
		Logger.Errorf("[%s] ResultCode=%d Message=%s", tag, err.ResultCode, err.Message)
	})
}
