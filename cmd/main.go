package main

import (
	"flights/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()
	if err := srv.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
