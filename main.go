package main

import (
	"fmt"
)

func main() {
	for _, todo := range List() {
		fmt.Println(todo.Title, todo.DateCreated, todo.ID)
	}
}
