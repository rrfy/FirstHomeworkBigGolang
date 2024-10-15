package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

type ListStorage struct {
	inner map[string][]string
}

func NewListStorage() *ListStorage {
	return &ListStorage{
		inner: make(map[string][]string),
	}
}

func (ls *ListStorage) LPUSH(key string, values []string) {
	if _, err := ls.inner[key]; !err {
		ls.inner[key] = []string{}
	}

	ls.inner[key] = append(values, ls.inner[key]...)
}

func (ls *ListStorage) RPUSH(key string, values []string) {
	if _, err := ls.inner[key]; !err {
		ls.inner[key] = []string{}
	}

	ls.inner[key] = append(ls.inner[key], values...)
}

func (ls *ListStorage) RADDToSet(key string, values []string) {
	if _, err := ls.inner[key]; !err {
		ls.inner[key] = []string{}
	}

	for _, value := range values {
		if !contains(ls.inner[key], value) {
			ls.inner[key] = append(ls.inner[key], value)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (ls *ListStorage) LPOP(key string, count ...int) ([]string, error) {
	if _, exist := ls.inner[key]; !exist {
		return nil, errors.New("no such key")
	}

	if len(ls.inner[key]) == 0 {
		return nil, errors.New("list is empty")
	}

	var start, end int
	if len(count) == 0 {
		start, end = 0, 1
	} else if len(count) == 1 {
		start, end = 0, count[0]
	} else if len(count) == 2 {
		start, end = count[0], count[1]
	} else {
		return nil, errors.New("invalid count")
	}

	if start < 0 {
		start = len(ls.inner[key]) + start
	}
	if end < 0 {
		end = len(ls.inner[key]) + end + 1
	}
	if end > len(ls.inner[key]) {
		end = len(ls.inner[key])
	}

	if start > end {
		return nil, errors.New("invalid range")
	}

	removed := ls.inner[key][start:end]

	ls.inner[key] = append(ls.inner[key][:start], ls.inner[key][end:]...)

	return removed, nil
}

func (ls *ListStorage) RPOP(key string, count ...int) ([]string, error) {
	if _, exists := ls.inner[key]; !exists {
		return nil, errors.New("no such key")
	}
	if len(ls.inner[key]) == 0 {
		return nil, errors.New("list is empty")
	}

	var start, end int
	length := len(ls.inner[key])
	if len(count) == 0 {
		start, end = length-1, length
	} else if len(count) == 1 {
		start, end = length-count[0], length
	} else if len(count) == 2 {
		start, end = count[0], count[1]
	} else {
		return nil, errors.New("invalid count")
	}

	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end + 1
	}
	if end > length {
		end = length
	}

	if start > end {
		return nil, errors.New("invalid range")
	}

	removed := ls.inner[key][start:end]
	ls.inner[key] = append(ls.inner[key][:start], ls.inner[key][end:]...)

	return removed, nil
}

func (ls *ListStorage) LSET(key string, index int, element string) error {
	if _, exists := ls.inner[key]; !exists {
		return errors.New("no such key")
	}
	if index < 0 {
		index += len(ls.inner[key])
	}
	if index < 0 || index >= len(ls.inner[key]) {
		return errors.New("index out of range")
	}

	ls.inner[key][index] = element
	return nil
}

func (ls *ListStorage) LGET(key string, index int) (string, error) {
	if _, exists := ls.inner[key]; !exists {
		return "", errors.New("no such key")
	}
	if index < 0 {
		index += len(ls.inner[key])
	}
	if index < 0 || index >= len(ls.inner[key]) {
		return "", errors.New("index out of range")
	}

	return ls.inner[key][index], nil
}

func (ls *ListStorage) SaveToDisk(filename string) error {
	data, err := json.Marshal(ls.inner)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (ls *ListStorage) LoadFromDisk(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &ls.inner)
}
