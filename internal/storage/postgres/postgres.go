package postgres

import (
	"apigo/api"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Db *sql.DB
}
type Store struct {
	id       string
	name     string
	created  time.Time
	modified time.Time
	deleted  bool
}

func (pdb *PostgresDB) InitDB() error {
	uri := "user=postgres password=secret host=172.17.0.1 port=5432 dbname=postgres connect_timeout=20 sslmode=disable"
	var err error
	pdb.Db, err = sql.Open("postgres", uri)
	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDB) CloseDB() error {
	pdb.Db.Close()
	return nil
}

func (pdb *PostgresDB) GetStores() ([]api.Store, error) {
	row, err := pdb.Db.Query("Select * From store WHERE deleted = FALSE")
	if err != nil {
		return nil, err
	}
	var stores []api.Store
	for row.Next() {
		var s Store
		err = row.Scan(&s.id, &s.name, &s.created, &s.modified, &s.deleted)
		if err != nil {
			return nil, err
		}
		stores = append(stores, api.Store{ID: s.id, Name: s.name})
	}
	return stores, nil
}

func (pdb *PostgresDB) GetStore(storeID string) (*api.Store, error) {
	row := pdb.Db.QueryRow("Select * From store WHERE id = '" + storeID + "' AND deleted = FALSE")
	s := &Store{}
	err := row.Scan(&s.id, &s.name, &s.created, &s.modified, &s.deleted)
	if err != nil {
		return nil, err
	}
	return &api.Store{ID: s.id, Name: s.name}, nil
}

func (pdb *PostgresDB) DeleteStore(storeID string) error {
	sqlStatement := `
		UPDATE store
		SET deleted = TRUE, modified_on = $1
		WHERE id = $2 AND deleted = FALSE
	`
	timeStr := time.Now().Format(time.RFC3339)
	_, err := pdb.Db.Exec(sqlStatement, timeStr, storeID)

	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDB) PostStore(storeReq *api.Store) error {
	sqlStatement := `
		INSERT INTO store (id, name, created_on, modified_on, deleted)
		VALUES ($1, $2, $3, $4, $5)
	`
	timeStr := time.Now().Format(time.RFC3339)
	_, err := pdb.Db.Exec(sqlStatement, storeReq.ID, storeReq.Name, timeStr, timeStr, false)

	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDB) PutStore(storeReq *api.Store) error {
	sqlStatement := `
		UPDATE store
		SET name = $1, modified_on = $2
		WHERE id = $3 AND deleted = FALSE
	`
	timeStr := time.Now().Format(time.RFC3339)
	_, err := pdb.Db.Exec(sqlStatement, storeReq.Name, timeStr, storeReq.ID)

	if err != nil {
		return err
	}
	return nil
}
