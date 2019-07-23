package storage

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"

	"github.com/Quard/gosh/internal/generator"
)

const dbLastIdKey = "gosh:lastid"

type RedisIdentifierStorage struct {
	lastIdentifier string
	db             *redis.Client
}

func NewRedisIdentifierStorage() (*RedisIdentifierStorage, error) {
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		addr = "localhost:6379"
	}
	db := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	return &RedisIdentifierStorage{db: db}, nil
}

func (storage *RedisIdentifierStorage) GetURL(identifier string) (string, error) {
	url, err := storage.db.Get(buildURLKey(identifier)).Result()
	if err == redis.Nil {
		return "", errors.New("identifier not found")
	} else if err != nil {
		log.Printf("[RedisIdentifierStorage.GetURL] redis get url error: %v", err)
		return "", errors.New("identifier not found")
	}

	return url, nil
}

func (storage *RedisIdentifierStorage) AddURL(url string) (string, error) {
	identifier, err := storage.getNextIdentifier()
	if err != nil {
		return "", errors.New("can't generate short url")
	}

	pipeline := storage.db.Pipeline()
	pipeline.Set(buildURLKey(identifier), url, 0)
	pipeline.Set(dbLastIdKey, identifier, 0)
	_, err = pipeline.Exec()
	if err != nil {
		log.Printf("[RedisIdentifierStorage.AddURL] redis set error: %v", err)
		return "", errors.New("can't generate short url")
	}

	return identifier, nil
}

func (storage *RedisIdentifierStorage) getNextIdentifier() (string, error) {
	var identifier string

	lastIdentifier, err := storage.db.Get(dbLastIdKey).Result()
	if err == redis.Nil {
		identifier = InitialIdentifier
	} else if err != nil {
		log.Printf("[RedisIdentifierStorage.getNextIdentifier] redis error: %v", err)
		return "", errors.New("can't get last identifier")
	} else {
		identifier, err = generator.GenerateNextSequence(lastIdentifier)
		if err != nil {
			log.Printf("[RedisIdentifierStorage.getNextIdentifier] sequence generation error: %v", err)
			return "", errors.New("can't generate new identifier")
		}
	}

	return identifier, nil
}

func buildURLKey(identifier string) string {
	return fmt.Sprintf("gosh:urls:%s", identifier)
}
