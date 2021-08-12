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

	var stgStore storage.StoreStorage
	var stgCustomer storage.CustomerStorage

	var err error
	switch *dbDriver {
	case "postgres":
		stgStore, err = postgres.NewPostgresStorage()
		if err != nil {
			log.Fatal(err)
		}
		defer stgStore.(*postgres.StoreDB).CloseDB()

		stgCustomer, err = postgres.NewPostgresCustomer()
		if err != nil {
			log.Fatal(err)
		}
		defer stgCustomer.(*postgres.CustomerDB).CloseDB()

		/*stgStore, err = postgres.NewPostgresStorage()
		if err != nil {
			log.Fatal(err)
		}
		defer stgStore.(*postgres.StoreDB).CloseDB()*/
	default:
		log.Fatalf("Unsupported driver %s", *dbDriver)
	}

	log.Printf("Listening on 0.0.0.0:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server.NewRouter(stgStore, stgCustomer)))
}
