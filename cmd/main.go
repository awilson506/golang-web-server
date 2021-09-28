package main

import (
	"log"
	"net/http"

	"github.com/awilson506/golang-web-server/server"
)

func main() {

	//start the server
	s := server.NewServer()
	err := s.Start()

	if err != http.ErrServerClosed {
		log.Println(err)
	}
}
