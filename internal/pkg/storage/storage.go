package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Storage struct {
	inner  map[string]string
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()
	logger.Info("created new storage")

	return Storage{
		inner:  make(map[string]string),
		logger: logger,
	}, nil
}

func (r Storage) Set(key, value string) {
	r.inner[key] = value
}

func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}

	return &res
}

func (r Storage) GetKind(key string) string {
	k, ok := r.inner[key]
	if !ok {
		return ""
	}
	_, err := strconv.Atoi(k)
	if err != nil {
		return "S"
	}
	return "D"
}
