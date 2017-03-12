package db

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func Connect() *bolt.DB {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Close() {
	db.Close()
}

func PutInBucket() {
	db := Connect()

	err := db.Update(func(tx *bolt.Tx) error {
		log.Println("Opening bucket...")
		b, err := tx.CreateBucketIfNotExists([]byte("Wifis"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("wifi1"), []byte("wifi1"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	db.Close()
}

func PrintBuckets() {
	db := Connect()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wifis"))

		err := b.ForEach(func(k, v []byte) error {
			log.Printf("Wifi: %s, %s", k, v)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	db.Close()
}
