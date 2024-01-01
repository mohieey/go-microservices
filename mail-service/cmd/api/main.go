package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const PORT = "80"

func main() {
	app := Config{
		Mailer: *createMailer(),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: app.routes(),
	}
	log.Println("mail service listening on port: ", PORT)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func createMailer() *Mail {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	return &Mail{
		Domain:     os.Getenv("MAIL_DOMAIN"),
		Host:       os.Getenv("MAIL_HOST"),
		Port:       port,
		Username:   os.Getenv("MAIL_USERNAME"),
		Password:   os.Getenv("MAIL_PSSWORD"),
		Encryption: os.Getenv("MAIL_PASSWORD"),
		From:       os.Getenv("FROM_ADDRESS"),
		FromName:   os.Getenv("FROM_NAME"),
	}
}
