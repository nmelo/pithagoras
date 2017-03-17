package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func Connect() error {
	var err error
	db, err = bolt.Open("data.db", 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	return db.Close()
}

func PutInBucket(bucket []byte, key []byte, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = b.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func PrintBucket(bucket string) error {

	fmt.Println("Printing bucket...")
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket not found")
		}

		err := b.ForEach(func(k, v []byte) error {
			date := time.Time{}
			err := date.UnmarshalBinary(v)
			if err != nil {
				return err
			}
			log.Printf("Sessions: %s, %s", k, date)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
