package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "80"

func main() {
	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: app.routes(),
	}

	log.Println("message broker listening on port: ", PORT)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
