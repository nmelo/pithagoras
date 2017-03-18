package wifi

import (
	"github.com/skycoin/skycoin/src/aether/wifi"
)

func Scan() (networks []network.WifiNetwork, err error) {
	networks = []network.WifiNetwork{}
	ifaces, err := network.NewWifiInterfaces()
	if err != nil {
		return
	}
	for _, v := range ifaces {
		nw, err := v.Scan()
		if err != nil {
			return nil, err
		}
		networks = append(networks, nw...)
	}
	return
}
