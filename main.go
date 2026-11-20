package main

import (
	"os"
	"sync"

	"github.com/Protarium-Network/disney-infinity-nex-global/globals"
	"github.com/Protarium-Network/disney-infinity-nex-global/nex"
)

func fatal(err error) {
	globals.Logger.Critical(err.Error())
	os.Exit(1)
}

func main() {
	games, err := selectedGames()
	if err != nil {
		fatal(err)
	}

	instances := make([]*globals.Instance, 0, len(games))
	for _, game := range games {
		inst, err := newInstance(game)
		if err != nil {
			fatal(err)
		}
		instances = append(instances, inst)
	}

	if err := checkPorts(instances); err != nil {
		fatal(err)
	}

	var wg sync.WaitGroup
	for _, inst := range instances {
		wg.Add(2)

		go func() {
			defer wg.Done()
			nex.StartAuthenticationServer(inst)
		}()
		go func() {
			defer wg.Done()
			nex.StartSecureServer(inst)
		}()
	}

	wg.Wait()
}
