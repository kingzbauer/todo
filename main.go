package main

import (
	"fmt"
)

func main() {
	todos := List()
	for _, todo := range todos {
		fmt.Println(todo.Title, todo.DateCreated, todo.ID)
	}

	fmt.Println(GetID(todos[0]))

	todo := todos[3]
	Delete(todo)

	todos = List()
	for _, todo := range todos {
		fmt.Println(todo.Title, todo.DateCreated, todo.ID)
	}
}
