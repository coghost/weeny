package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisStorage struct {
	// Address is the redis server address
	Address string
	// Password is the password for the redis server
	Password string
	// DB is the redis database. Default is 0
	DB int
	// Prefix is an optional string in the keys. It can be used
	// to use one redis database for independent scraping tasks.
	Prefix string
	// Client is the redis connection
	Client *redis.Client

	// Expiration time for Visited keys. After expiration pages
	// are to be visited again.
	Expires time.Duration
}

func MustNewRedisStorage(address, password string, db int, prefix string) *RedisStorage {
	st, err := NewRedisStorage(address, password, db, prefix)
	if err != nil {
		panic(err)
	}

	return st
}

var ErrPrefixMissing = errors.New("prefix is required")

func NewRedisStorage(address, password string, db int, prefix string) (*RedisStorage, error) {
	if prefix == "" {
		return nil, ErrPrefixMissing
	}

	return &RedisStorage{
		Address:  address,
		Password: password,
		DB:       db,
		Prefix:   fmt.Sprintf("W:%s", prefix),
	}, nil
}

// Init initializes the redis storage
func (s *RedisStorage) Init() error {
	if s.Client == nil {
		s.Client = redis.NewClient(&redis.Options{
			Addr:     s.Address,
			Password: s.Password,
			DB:       s.DB,
		})
	}

	_, err := s.Client.Ping().Result()
	if err != nil {
		return fmt.Errorf("redis connection error: %w", err)
	}

	return nil
}

// Clear removes all entries from the storage
func (s *RedisStorage) Clear() error {
	r2 := s.Client.Keys(s.Prefix + ":*")

	keys, err := r2.Result()
	if err != nil {
		return err
	}

	return s.Client.Del(keys...).Err()
}

func (s *RedisStorage) Visited(requestID string) error {
	return s.Client.Set(s.getIDStr(requestID), "1", s.Expires).Err()
}

func (s *RedisStorage) IsVisited(requestID string) (bool, error) {
	_, err := s.Client.Get(s.getIDStr(requestID)).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *RedisStorage) getIDStr(id string) string {
	return fmt.Sprintf("%s:%s", s.Prefix, id)
}
