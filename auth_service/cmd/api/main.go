package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8081"

type Application struct {
	DB      *sql.DB
	Models  data.Models
}

func main() {
	log.Printf("Starting authentication Service")

	db := connectToDB()
	if db == nil {
		log.Fatal("Could not connect to Postgres")
	}

	app := Application{
		DB:      db,
		Models:  data.New(db),
	}
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting authentication Service on port :%s\n", webPort)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for range 10 {
		db, err := sql.Open("pgx", dsn)

		if err != nil && db.Ping() == nil {
			log.Println("Could not connect to Postgres")
			return db
		}

		log.Println("Waiting fpr Postgres")
		time.Sleep(time.Second * 2)
	}

	return nil
}
