package main

import (
	"apigo/internal/server"
	"apigo/internal/storage/postgres"
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8000", "Port the server will be listening")
	flag.Parse()
	ctrlDB := &postgres.PostgresDB{}
	server := server.NewServer(ctrlDB)
	log.Printf("Listening on 0.0.0.0:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server.Router))
}
