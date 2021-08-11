package postgres

import (
	"apigo/api"
	"apigo/internal/storage"
	"context"
	"database/sql"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Db *sql.DB
}
type Store struct {
	id       int64
	name     string
	address  string
	created  time.Time
	modified time.Time
}

func NewPostgresStorage() (storage.DataBase, error) {
	uri := "user=postgres password=secret host=172.17.0.1 port=5432 dbname=postgres connect_timeout=20 sslmode=disable"
	var err error
	postgresDB, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err

	}
	return &PostgresDB{Db: postgresDB}, nil
}
func (pdb *PostgresDB) CloseDB() error {
	return pdb.Db.Close()
}

func (pdb *PostgresDB) GetStores() ([]api.Store, error) {
	rows, err := pdb.Db.Query("Select * From store")
	if err != nil {
		return nil, err
	}
	var stores []api.Store
	for rows.Next() {
		var s Store
		err = rows.Scan(&s.id, &s.name, &s.address, &s.created, &s.modified)
		if err != nil {
			return nil, err
		}
		stores = append(stores, api.Store{Id: strconv.FormatInt(s.id, 10), Name: s.name, Address: s.address})
	}
	return stores, nil
}

func (pdb *PostgresDB) GetStore(storeID string) (*api.Store, error) {
	row := pdb.Db.QueryRow("Select * From store WHERE id = '" + storeID + "'")
	s := &Store{}
	err := row.Scan(&s.id, &s.name, &s.address, &s.created, &s.modified)
	if err != nil {
		return nil, err
	}
	return &api.Store{Id: strconv.FormatInt(s.id, 10), Name: s.name, Address: s.address}, nil
}

func (pdb *PostgresDB) DeleteStore(storeID string) (string, error) {
	sqlStatement := `
		DELETE FROM store
		WHERE id = $1
	`
	contxt := context.Background()
	tx, err := pdb.Db.BeginTx(contxt, nil)
	if err != nil {
		return "", err
	}

	resultsql, err := tx.ExecContext(contxt, sqlStatement, storeID)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	rowsAffected, err := resultsql.RowsAffected()
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(rowsAffected, 10), nil
}

func (pdb *PostgresDB) PostStore(storeReq *api.Store) (string, error) {

	timeStr := time.Now().Format(time.RFC3339)
	sqlStatement := "INSERT INTO store (name, address, created_on, modified_on)" +
		"VALUES ( '" + storeReq.Name + "', '" + storeReq.Address + "', '" + timeStr + "', '" + timeStr + "')" +
		"RETURNING id"
	lastInsertid := int64(0)

	contxt := context.Background()
	tx, err := pdb.Db.BeginTx(contxt, nil)
	if err != nil {
		return "", err
	}

	err = tx.QueryRowContext(contxt, sqlStatement).Scan(&lastInsertid)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(lastInsertid, 10), nil
}

func (pdb *PostgresDB) PutStore(storeReq *api.Store) error {
	sqlStatement := `
		UPDATE store
		SET name = $1, address = $2, modified_on = $3
		WHERE id = $4
	`
	timeStr := time.Now().Format(time.RFC3339)

	contxt := context.Background()
	tx, err := pdb.Db.BeginTx(contxt, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(contxt, sqlStatement, storeReq.Name, storeReq.Address, timeStr, storeReq.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
