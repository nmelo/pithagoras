package db

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

const Sessions string = "Sessions"

var (
	db *bolt.DB
)

type Session struct {
	Key  string
	Date time.Time
}

func Connect() error {
	var err error
	db, err = bolt.Open("data.db", 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func PutInBucket(bucket []byte, key []byte, value []byte) error {
	if err := Connect(); err != nil {
		return err
	}
	defer db.Close()

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

func ListSessions() (sessions []Session, err error) {
	if err := Connect(); err != nil {
		return nil, err
	}
	defer db.Close()

	sessions = []Session{}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Sessions))
		if b == nil {
			return errors.New("bucket not found")
		}

		err := b.ForEach(func(k, v []byte) error {
			date := time.Time{}
			err := date.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			sessions = append(sessions, Session{Key: string(k), Date: date})
			return nil
		})
		return err
	})
	return sessions, err
}
