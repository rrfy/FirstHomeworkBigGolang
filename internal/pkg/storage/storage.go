package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"go.uber.org/zap"
)

// Value хранит значение и его тип
type Value struct {
	v interface{}
	t Type
}

// Storage хранит данные с поддержкой различных типов
type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
}

// NewStorage создает новое хранилище
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

// Set добавляет значение в хранилище. Поддерживаются скалярные типы, списки и словари.
func (s *Storage) Set(key string, value interface{}) {
	switch kind := getType(value); kind {
	case TypeInt, TypeString, TypeBool, TypeFloat, TypeList, TypeMap:
		s.inner[key] = Value{v: value, t: kind}
	default:
		s.logger.Error(
			"undefined value type",
			zap.String("key", key),
			zap.Any("value", value),
		)
	}
}

// Get возвращает значение по ключу, если оно существует
func (s *Storage) Get(key string) (interface{}, error) {
	res, ok := s.inner[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return res.v, nil
}

// GetKind возвращает тип значения по ключу
func (s *Storage) GetKind(key string) (Type, error) {
	v, ok := s.inner[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v.t, nil
}

// SaveToDisk сохраняет данные на диск в формате JSON
func (s *Storage) SaveToDisk(filename string) error {
	data, err := json.Marshal(s.inner)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromDisk загружает данные с диска
func (s *Storage) LoadFromDisk(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.inner)
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

// getType определяет тип значения
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

// ListStorage расширяет функционал для работы со списками
type ListStorage struct {
	inner map[string][]interface{}
}

// NewListStorage создает хранилище для списков
func NewListStorage() *ListStorage {
	return &ListStorage{
		inner: make(map[string][]interface{}),
	}
}

// LPUSH добавляет элементы в начало списка
func (ls *ListStorage) LPUSH(key string, values []interface{}) {
	if _, err := ls.inner[key]; !err {
		ls.inner[key] = []interface{}{}
	}

	ls.inner[key] = append(values, ls.inner[key]...)
}

// RPUSH добавляет элементы в конец списка
func (ls *ListStorage) RPUSH(key string, values []interface{}) {
	if _, err := ls.inner[key]; !err {
		ls.inner[key] = []interface{}{}
	}

	ls.inner[key] = append(ls.inner[key], values...)
}

// LPOP удаляет и возвращает элементы с начала списка
func (ls *ListStorage) LPOP(key string, count ...int) ([]interface{}, error) {
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
	} else {
		return nil, errors.New("invalid count")
	}

	if end > len(ls.inner[key]) {
		end = len(ls.inner[key])
	}

	removed := ls.inner[key][start:end]
	ls.inner[key] = ls.inner[key][end:]

	return removed, nil
}

// DictionaryStorage хранит данные в формате словаря (ключ-значение)
type DictionaryStorage struct {
	inner map[string]map[string]interface{}
}

// NewDictionaryStorage создает новое хранилище для словарей
func NewDictionaryStorage() *DictionaryStorage {
	return &DictionaryStorage{
		inner: make(map[string]map[string]interface{}),
	}
}

// Set устанавливает значение для ключа в словаре
func (ds *DictionaryStorage) Set(key, field string, value interface{}) {
	if _, ok := ds.inner[key]; !ok {
		ds.inner[key] = make(map[string]interface{})
	}
	ds.inner[key][field] = value
}

// Get возвращает значение для ключа из словаря
func (ds *DictionaryStorage) Get(key, field string) (interface{}, error) {
	if _, ok := ds.inner[key]; !ok {
		return nil, errors.New("key not found")
	}
	value, ok := ds.inner[key][field]
	if !ok {
		return nil, errors.New("field not found")
	}
	return value, nil
}
