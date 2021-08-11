package main

import (
	"apigo/internal/server"
	"apigo/internal/storage"
	"apigo/internal/storage/postgres"
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8000", "Port the server will be listening")
	dbDriver := flag.String("driver", "", "Database driver (postgres only supported now)")
	flag.Parse()

	var storage storage.DataBase
	var err error
	switch *dbDriver {
	case "postgres":
		storage, err = postgres.NewPostgresStorage()
		if err != nil {
			log.Fatal(err)
		}

		defer storage.(*postgres.PostgresDB).CloseDB()
	default:
		log.Fatalf("Unsupported driver %s", *dbDriver)
	}

	log.Printf("Listening on 0.0.0.0:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server.NewRouter(storage)))
}
