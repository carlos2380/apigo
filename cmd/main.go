package main

import (
	"apigo/internal/server"
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8000", "Port the server will be listening")
	flag.Parse()
	server := server.NewServer()
	log.Printf("Listening on 0.0.0.0:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, server.Router))
}
