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
	dbPassword := flag.String("password", "", "Password of the database")
	dbIPHost := flag.String("host", "172.17.0.1", "Host IP of the database")
	dbPort := flag.String("dbport", "5432", "Port of the database")

	flag.Parse()

	var stgStore storage.StoreStorage
	var stgCustomer storage.CustomerStorage
	var stgCase storage.CaseStorage

	switch *dbDriver {
	case "postgres":
		db, err := postgres.InitPostgres(dbPassword, dbIPHost, dbPort)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			err = postgres.ClosePostgres(db)
			if err != nil {
				log.Fatal(err)
			}
		}()

		stgStore = &postgres.StoreDB{DB: db}
		stgCustomer = &postgres.CustomerDB{DB: db}
		stgCase = &postgres.CaseDB{DB: db}

	default:
		log.Fatalf("Unsupported driver %s", *dbDriver)
	}

	log.Printf("Listening on 0.0.0.0:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server.NewRouter(stgStore, stgCustomer, stgCase)))
}
