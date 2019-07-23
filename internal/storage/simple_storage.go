package storage

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/boltdb/bolt"

	"github.com/Quard/gosh/internal/generator"
)

const dbPath = "gosh.bolt.db"

var dbBucketCommon = []byte("common")
var dbBucketUrl = []byte("urls")
var dbLastIdKeyName = []byte("lastIdentifier")

type SimpleIdentifierStorage struct {
	lastIdentifier string
	db             *bolt.DB
}

func NewSimpleIdentifierStorage() (*SimpleIdentifierStorage, error) {
	db, err := bolt.Open(dbPath, 0644, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open database, %v", err)
	}

	return &SimpleIdentifierStorage{db: db}, nil
}

func (storage *SimpleIdentifierStorage) GetURL(identifier string) (string, error) {
	var url []byte

	err := storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbBucketUrl)
		if bucket == nil {
			log.Printf("[SimpleIdentifierStorage.GetURL] bucket get/create error")
			return errors.New("URLs bucket not created")
		}

		url = bucket.Get([]byte(identifier))

		return nil
	})

	if err != nil {
		return "", err
	}

	return string(url), nil
}

func (storage *SimpleIdentifierStorage) AddURL(url string) (string, error) {
	identifier, err := storage.getNextIdentifier()
	if err != nil {
		return "", errors.New("can't generate short url")
	}

	err = storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(dbBucketUrl)
		if err != nil {
			log.Printf("[SimpleIdentifierStorage.AddURL] bucket get/create error: %v", err)
			return err
		}

		err = bucket.Put([]byte(identifier), []byte(url))
		if err != nil {
			log.Printf("[SimpleIdentifierStorage.AddURL] save id/url pair error: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return "", errors.New("can't create short url")
	}

	return identifier, nil
}

func (storage *SimpleIdentifierStorage) getNextIdentifier() (string, error) {
	var identifier string

	err := storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(dbBucketCommon)
		if err != nil {
			log.Printf("[SimpleIdentifierStorage.getNextIdentifier] bucket get/create error: %v", err)
			return err
		}

		if len(storage.lastIdentifier) == 0 {
			storage.lastIdentifier = string(bucket.Get(dbLastIdKeyName))
		}

		if len(storage.lastIdentifier) == 0 {
			storage.lastIdentifier = InitialIdentifier
			identifier = InitialIdentifier
		} else {
			identifier, err = generator.GenerateNextSequence(storage.lastIdentifier)
			storage.lastIdentifier = identifier
			if err != nil {
				log.Printf("[SimpleIdentifierStorage.getNextIdentifier] generate next id error: %v", err)
				return err
			}
		}

		err = bucket.Put(dbLastIdKeyName, []byte(identifier))
		if err != nil {
			log.Printf("[SimpleIdentifierStorage.getNextIdentifier] save last id error: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("can't generate next identifier, %v", err)
		return "", errors.New("can't generate short url")
	}

	return identifier, nil
}

func assertSequenceEqual(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
