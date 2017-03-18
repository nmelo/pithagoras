package main

import (
	"fmt"

	"github.com/nmelo/pithagoras/pkg/db"
)

func main() {
	fmt.Println("Printing bucket...")
	objects, err := db.ListSessions()
	if err != nil {
		fmt.Printf("Exiting: %s\n", err)
	}

	fmt.Println("Sessions:")
	for _, v := range objects {
		fmt.Printf("\t%s, %s\n", v.Key, v.Date)
	}
}
