package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/nmelo/pithagoras/pkg/bluetooth"
	"github.com/nmelo/pithagoras/pkg/session"
)

var (
	stop chan bool      = make(chan bool)
	wg   sync.WaitGroup = sync.WaitGroup{}
)

func installSignalHandler() {
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)

	go func(handler chan os.Signal) {
		for sig := range handler {
			if sig == os.Interrupt {
				stop <- true
			}
		}
	}(handler)
}

func main() {

	if err := session.Start(); err != nil {
		fmt.Println("Exiting: ", err)
		return
	}

	fmt.Println("Starting bluetooth service...")

	if err := bluetooth.Serve(&wg); err != nil {
		fmt.Println("Exiting: ", err)
		return
	} else {
		installSignalHandler()
		<-stop
		if err := session.End(); err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Println("Shutting down...")
	}
}
