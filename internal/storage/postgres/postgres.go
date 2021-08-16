package postgres

import (
	"database/sql"
	"fmt"
)

func InitPostgres(dbPassword, dbIPHost, dbPort *string) (*sql.DB, error) {
	uri := fmt.Sprintf(`user=postgres password=%s 
		host=%s port=%s dbname=postgres connect_timeout=20 
		sslmode=disable`, *dbPassword, *dbIPHost, *dbPort)
	var err error
	postgresDB, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err
	}
	postgresDB.SetMaxIdleConns(128)
	postgresDB.SetMaxOpenConns(128)
	return postgresDB, nil
}
func ClosePostgres(db *sql.DB) error {
	return db.Close()
}
