package main

import (
	"homework1/internal/pkg/server"
	"homework1/internal/pkg/storage"
)

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	s := server.New(":8090", &store)
	s.Start()
}
