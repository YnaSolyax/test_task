package main

import (
	"test_task/server"
)

func main() {
	store := server.NewTodoStore()
	server.StartServ()
	server.Handler(store)
}
