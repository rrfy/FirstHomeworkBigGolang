package main

import (
	"fmt"
	"homework1/internal/pkg/storage"
	"os"
)

func main() {
	ls := storage.NewListStorage()

	// Загрузка состояния с диска
	if err := ls.LoadFromDisk("db.json"); err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading state:", err)
	}

	// Примеры работы с базой данных
	ls.RPUSH("mylist", []string{"a", "b", "c"})
	ls.LSET("mylist", 1, "z")

	value, err := ls.LGET("mylist", 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Value at index 1:", value)
	}

	// Сохранение состояния на диск перед выходом
	if err := ls.SaveToDisk("db.json"); err != nil {
		fmt.Println("Error saving state:", err)
	}

}
