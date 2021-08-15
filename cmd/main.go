package main

import (
	"apigo/internal/server"
	"apigo/internal/storage"
	"apigo/internal/storage/postgres"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
				log.Println(err)
			}
		}()

		stgStore = &postgres.StoreDB{DB: db}
		stgCustomer = &postgres.CustomerDB{DB: db}
		stgCase = &postgres.CaseDB{DB: db}

	default:
		log.Fatalf("Unsupported driver %s", *dbDriver)
	}

	router := server.NewRouter(stgStore, stgCustomer, stgCase)

	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	log.Printf("Listening on 0.0.0.0:%s", *port)

	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
