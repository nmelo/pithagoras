package session

import (
	"time"

	"fmt"

	"crypto/rand"

	"github.com/nmelo/pithagoras/pkg/db"
)

const Bucket string = "Sessions"

var (
	SessionID string
)

func Start() error {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	SessionID = fmt.Sprintf("%X", randBytes)

	fmt.Println("Saving session start...")

	start := time.Now()
	startBytes, err := start.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshall start date: %s", err)
	}
	return db.PutInBucket([]byte(Bucket), []byte(fmt.Sprintf("start-%s", SessionID)), startBytes)
}

func End() error {
	fmt.Println("Saving session end...")

	end := time.Now()
	endBytes, err := end.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshall end date: %s", err)
	}
	return db.PutInBucket([]byte(Bucket), []byte(fmt.Sprintf("end-%s", SessionID)), endBytes)
}
