package main

import (
	"fmt"

	"github.com/nmelo/pithagoras/pkg/wifi"
)

func main() {
	fmt.Println("Getting list of wifis")
	wifi.PrintList()
}
