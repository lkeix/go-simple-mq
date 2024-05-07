package main

import (
	"fmt"
	"time"

	gosimplemq "github.com/lkeix/go-simple-mq"
)

func main() {
	c, err := gosimplemq.NewClient("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	c.Send([]byte("Hello, world!"))
	time.Sleep(1 * time.Second)
	msg, err := c.Receive()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg))
}
