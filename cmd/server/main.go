package main

import "github.com/lkeix/go-simple-mq/internal/server"

func main() {
	server.Run("tcp", "localhost:8080", 1024)
}
