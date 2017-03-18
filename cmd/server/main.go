package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/nmelo/pithagoras/pkg/session"
	"github.com/nmelo/pithagoras/pkg/ui"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func waitSignal() {
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)
	go func(handler chan os.Signal) {
		for sig := range handler {
			if sig == os.Interrupt {
				cancel()
			}
		}
	}(handler)
}

func main() {

	ctx, cancel = context.WithCancel(context.Background())

	if err := session.Start(); err != nil {
		fmt.Println("Exiting: ", err)
		return
	}

	//go bluetooth.Serve(ctx)

	go ui.Serve(ctx)

	waitSignal()

	select {
	case <-ctx.Done():
		if err := session.End(); err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Println("Shutting down...")
	}
}
