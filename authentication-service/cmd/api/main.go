package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "80"

func main() {
	dbConn, err := openDB()
	if err != nil {
		log.Panic("error connecting to database")
	}

	app := Config{
		DB:     dbConn,
		Models: data.New(dbConn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: app.routes(),
	}

	log.Println("authentication service listening on port: ", PORT)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func openDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("connected to database")

	return db, nil
}
