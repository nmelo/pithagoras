package wifi

import (
	"fmt"

	"github.com/skycoin/skycoin/src/aether/wifi"
)

func PrintList() {
	ifaces, err := network.NewWifiInterfaces()
	if err != nil {
		fmt.Printf("getting interfaces: %s", err)
	}
	for _, v := range ifaces {
		fmt.Println(v.Name)
	}
}
