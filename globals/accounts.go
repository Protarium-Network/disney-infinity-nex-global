package globals

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// InitAccounts builds the two Quazal system accounts. Call after
// KerberosPassword is set.
func (i *Instance) InitAccounts() {
	i.AuthenticationServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", i.KerberosPassword, false)
	i.SecureServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", i.KerberosPassword, false)
}

func (i *Instance) AccountDetailsByPID(pid types.PID) (*nex.Account, *nex.Error) {
	if pid.Equals(i.AuthenticationServerAccount.PID) {
		return i.AuthenticationServerAccount, nil
	}

	if pid.Equals(i.SecureServerAccount.PID) {
		return i.SecureServerAccount, nil
	}

	password, errorCode := i.PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	return nex.NewAccount(pid, strconv.FormatUint(uint64(pid), 10), password, false), nil
}

func (i *Instance) AccountDetailsByUsername(username string) (*nex.Account, *nex.Error) {
	if username == i.AuthenticationServerAccount.Username {
		return i.AuthenticationServerAccount, nil
	}

	if username == i.SecureServerAccount.Username {
		return i.SecureServerAccount, nil
	}

	pidInt, err := strconv.ParseUint(username, 10, 32)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Invalid username")
	}

	pid := types.NewPID(pidInt)

	password, errorCode := i.PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	return nex.NewAccount(pid, username, password, false), nil
}
