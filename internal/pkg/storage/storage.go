package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	v string
	t Type
}

type Storage struct {
	inner  map[string]Value
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
		inner:  make(map[string]Value),
		logger: logger,
	}, nil
}

func (r Storage) Set(key, value string) {
	switch kind := getType(value); kind {
	case TypeInt:
		r.inner[key] = Value{v: value, t: kind}
	case TypeString:
		r.inner[key] = Value{v: value, t: kind}
	case TypeUndefined:
		r.logger.Error(
			"undefined value type",
			zap.String("key", key),
			zap.Any("value", value),
		)
	}
}

func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}

	return &res.v
}

type Type string

const (
	TypeInt       Type = "D"
	TypeString    Type = "S"
	TypeUndefined Type = "UN"
)

func getType(value string) Type {
	var val any

	val, err := strconv.Atoi(value)
	if err != nil {
		val = value
	}
	switch val.(type) {
	case int:
		return TypeInt
	case string:
		return TypeString
	default:
		return TypeUndefined
	}
}

func (r Storage) GetKind(key string) string {
	k, ok := r.inner[key]
	if !ok {
		return ""
	}
	_, err := strconv.Atoi(k.v)
	if err != nil {
		return "S"
	}
	return "D"
}
