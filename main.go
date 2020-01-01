package main

func main() {
	todo := Todo{Title: "Test todo"}
	_, err := InsertOne(todo)
	if err != nil {
		panic(err)
	}
}
