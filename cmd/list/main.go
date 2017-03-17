package main

import (
	"fmt"

	"github.com/nmelo/pithagoras/pkg/db"
	"github.com/nmelo/pithagoras/pkg/session"
)

func main() {
	db.Connect()

	err := db.PrintBucket(session.Bucket)
	if err != nil {
		fmt.Printf("Exiting: %s\n", err)
	}

	defer db.Close()
}
