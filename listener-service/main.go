package main

import (
	"listener/event"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	connection, err := connect()
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	log.Println("Connection established with rabbitMQ successfully")

	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {

	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		return nil, err
	}

	return connection, nil
}
