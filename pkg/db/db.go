package db

import (
	"errors"
	"time"

	"encoding/json"

	"bytes"
	"encoding/binary"

	"github.com/boltdb/bolt"
)

const Sessions string = "Sessions"
const Messages string = "Messages"

var (
	db *bolt.DB
)

type Session struct {
	Key  string
	Date time.Time
}

type Message struct {
	ID       string
	Username string
	Text     string
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

func ClearSessions() (err error) {
	if err := Connect(); err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Sessions))
		if b == nil {
			return errors.New("bucket not found")
		}

		err := b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
		return err
	})
	return
}

func AddMessage(userID, username, text string) error {
	if err := Connect(); err != nil {
		return err
	}
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(Messages))
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		val, err := json.Marshal(&Message{
			ID:       userID,
			Username: username,
			Text:     text,
		})
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, id)
		if err != nil {
			return err
		}

		err = b.Put(buf.Bytes(), val)
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

func ListMessages() (messages []Message, err error) {
	if err := Connect(); err != nil {
		return nil, err
	}
	defer db.Close()

	messages = []Message{}
	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte(Messages))
		if err != nil {
			return err
		}

		err = b.ForEach(func(k, v []byte) error {
			m := &Message{}
			err := json.Unmarshal(v, m)
			if err != nil {
				return err
			}
			messages = append(messages, *m)
			return nil
		})
		return err
	})
	return messages, err
}

func ClearMessages() (err error) {
	if err := Connect(); err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Messages))
		if b == nil {
			return errors.New("bucket not found")
		}

		err := b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
		return err
	})
	return
}