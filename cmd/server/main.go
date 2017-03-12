package main

import (
	"fmt"

	"github.com/nmelo/pithagoras/pkg/blue"
	"github.com/nmelo/pithagoras/pkg/wifi"
)

func main() {
	fmt.Println("Getting list of wifi")
	wifi.PrintList()
	blue.Connect()
}
