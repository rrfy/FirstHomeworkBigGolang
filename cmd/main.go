package main

import (
	"fmt"
	"homework1/internal/pkg/storage"
)

func main() {
	s, _ := storage.NewStorage()

	s.Set("key2", "")

	res1 := s.Get("key2")
	res2 := s.Get("key3")

	fmt.Println(*res1, res2)

	s1, _ := storage.NewStorage()

	s1.Set("key2", "27")

	res3 := s.GetKind("key2")
	res4 := s1.GetKind("key2")

	fmt.Println(res3, res4)

}
