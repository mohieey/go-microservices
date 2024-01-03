package main

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = "80"

func main() {
	rabbitMQConn, err := connect()
	if err != nil {
		panic(err)
	}
	defer rabbitMQConn.Close()
	log.Println("rabbitMQConn established with rabbitMQ successfully")

	// rabbitMQConsumer, err := event.NewConsumer(rabbitMQConn)
	// if err != nil {
	// 	panic(err)
	// }

	app := Config{
		RabbitMQ: rabbitMQConn,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: app.routes(),
	}

	log.Println("message broker listening on port: ", PORT)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connect() (*amqp.Connection, error) {

	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		return nil, err
	}

	return connection, nil
}
