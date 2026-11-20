package globals

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// PasswordFromPID returns the NEX (Kerberos) password for a player PID.
//
//   - LocalAuthMode: looked up in settings.json (same file format as
//     nex-viewer), for isolated preservation setups.
//   - otherwise: derived by HMAC-SHA256(NEXPasswordSecret, pid) so the
//     account server and this server agree on the secret without storing it.
func (i *Instance) PasswordFromPID(pid types.PID) (string, uint32) {
	if i.LocalAuthMode {
		return i.passwordFromPIDLocal(pid)
	}

	if len(i.NEXPasswordSecret) < 32 {
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	pidBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(pidBytes, uint64(pid))
	mac := hmac.New(sha256.New, i.NEXPasswordSecret)
	_, _ = mac.Write(pidBytes)

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), 0
}

type jsonAccount struct {
	Platform string  `json:"platform"`
	Username string  `json:"username"`
	Pid      float64 `json:"pid"`
	Password string  `json:"password"`
}

type settingsJSON struct {
	Accounts []jsonAccount `json:"accounts"`
}

func (i *Instance) passwordFromPIDLocal(pid types.PID) (string, uint32) {
	path := i.SettingsPath
	if path == "" {
		path = "settings.json"
	}

	file, err := os.ReadFile(path)
	if err != nil {
		Logger.Error(err.Error())
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	var data settingsJSON
	if err := json.Unmarshal(file, &data); err != nil {
		Logger.Error(err.Error())
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	want := strconv.FormatUint(uint64(pid), 10)
	for _, account := range data.Accounts {
		if account.Username == want {
			Logger.Infof("[%s] Using local account details for %v", i.Game.Name, account.Username)
			return account.Password, 0
		}
	}

	return "", nex.ResultCodes.RendezVous.InvalidUsername
}
