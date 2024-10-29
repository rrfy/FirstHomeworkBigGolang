package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Value struct {
	v         interface{}
	t         Type
	expiresAt time.Time
}

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}

	logger.Info("created new storage")

	return Storage{
		inner:  make(map[string]Value),
		logger: logger,
	}, nil
}

func (s *Storage) Set(key string, value interface{}, ttl ...time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiration := time.Time{}
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0])
	}

	s.inner[key] = Value{
		v:         value,
		t:         getType(value),
		expiresAt: expiration,
	}
}

func (s *Storage) Get(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.inner[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	if value.expiresAt.IsZero() || value.expiresAt.After(time.Now()) {
		return value.v, nil
	}

	delete(s.inner, key)
	return nil, errors.New("key expired")
}

func (s *Storage) ActiveExpiration(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for k, v := range s.inner {
			if !v.expiresAt.IsZero() && v.expiresAt.Before(time.Now()) {
				delete(s.inner, k)
				s.logger.Info("expired key removed", zap.String("key", k))
			}
		}
		s.mu.Unlock()
	}
}

func (s *Storage) SaveToDisk(filename string) error {
	data, err := json.Marshal(s.inner)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (s *Storage) LoadFromDisk(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.inner)
}

func getType(value interface{}) Type {
	switch value.(type) {
	case int:
		return TypeInt
	case string:
		return TypeString
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case []interface{}:
		return TypeList
	case map[string]interface{}:
		return TypeMap
	default:
		return TypeUndefined
	}
}

type Type string

const (
	TypeInt       Type = "INT"
	TypeString    Type = "STRING"
	TypeBool      Type = "BOOL"
	TypeFloat     Type = "FLOAT"
	TypeList      Type = "LIST"
	TypeMap       Type = "MAP"
	TypeUndefined Type = "UNDEFINED"
)
