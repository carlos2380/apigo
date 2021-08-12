package postgres

import (
	"apigo/api"
	"apigo/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type StoreDB struct {
	Db *sql.DB
}
type Store struct {
	id       int64
	name     string
	address  string
	created  time.Time
	modified time.Time
}

func NewPostgresStorage() (storage.StoreStorage, error) {
	uri := "user=postgres password=secret host=172.17.0.1 port=5432 dbname=postgres connect_timeout=20 sslmode=disable"
	var err error
	postgresDB, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err

	}
	return &StoreDB{Db: postgresDB}, nil
}
func (pdb *StoreDB) CloseDB() error {
	return pdb.Db.Close()
}

func (pdb *StoreDB) GetStores() ([]*api.Store, error) {
	rows, err := pdb.Db.Query("Select * From store")
	if err != nil {
		return nil, err
	}
	var stores []*api.Store
	for rows.Next() {
		var s Store
		err = rows.Scan(&s.id, &s.name, &s.address, &s.created, &s.modified)
		if err != nil {
			return nil, err
		}
		stores = append(stores, &api.Store{Id: strconv.FormatInt(s.id, 10), Name: s.name, Address: s.address})
	}
	return stores, nil
}

func (pdb *StoreDB) GetStore(storeID string) (*api.Store, error) {
	row := pdb.Db.QueryRow("Select * From store WHERE id = '" + storeID + "'")
	s := &Store{}
	err := row.Scan(&s.id, &s.name, &s.address, &s.created, &s.modified)
	if err != nil {
		return nil, err
	}
	return &api.Store{Id: strconv.FormatInt(s.id, 10), Name: s.name, Address: s.address}, nil
}

func (pdb *StoreDB) DeleteStore(storeID string) (err error) {
	sqlStatement := `
		DELETE FROM store
		WHERE id = $1
	`
	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	resultsql, err := tx.ExecContext(ctx, sqlStatement, storeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowsAffected, err := resultsql.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		err = errors.New("item no found")
		return
	}

	return nil
}

func (pdb *StoreDB) PostStore(storeReq *api.Store) (_ string, err error) {

	timeStr := time.Now().Format(time.RFC3339)
	sqlStatement := fmt.Sprintf(`INSERT INTO store (name, address, created_on, modified_on)
		VALUES ( '%s', '%s', '%s', '%s')
		RETURNING id`,
		storeReq.Name, storeReq.Address, timeStr, timeStr)
	lastInsertid := int64(0)

	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	err = tx.QueryRowContext(ctx, sqlStatement).Scan(&lastInsertid)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(lastInsertid, 10), nil
}

func (pdb *StoreDB) PutStore(storeReq *api.Store) (err error) {
	sqlStatement := `
		UPDATE store
		SET name = $1, address = $2, modified_on = $3
		WHERE id = $4
	`
	timeStr := time.Now().Format(time.RFC3339)

	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlStatement, storeReq.Name, storeReq.Address, timeStr, storeReq.Id)
	if err != nil {
		return err
	}

	return nil
}
