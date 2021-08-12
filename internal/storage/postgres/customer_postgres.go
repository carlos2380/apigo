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

type CustomerDB struct {
	Db *sql.DB
}
type Customer struct {
	id        int64
	firstName string
	lastName  string
	age       int
	email     string
	created   time.Time
	modified  time.Time
}

func NewPostgresCustomer() (storage.CustomerStorage, error) {
	uri := "user=postgres password=secret host=172.17.0.1 port=5432 dbname=postgres connect_timeout=20 sslmode=disable"
	var err error
	postgresDB, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err

	}
	return &CustomerDB{Db: postgresDB}, nil
}

func (pdb *CustomerDB) CloseDB() error {
	return pdb.Db.Close()
}

func (pdb *CustomerDB) GetCustomers() ([]*api.Customer, error) {
	rows, err := pdb.Db.Query("Select * From customer")
	if err != nil {
		return nil, err
	}
	var customers []*api.Customer
	for rows.Next() {
		var c Customer
		err = rows.Scan(&c.id, &c.firstName, &c.lastName, &c.age, &c.email, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &api.Customer{Id: strconv.FormatInt(c.id, 10), FirstName: c.firstName, LastName: c.lastName, Age: strconv.Itoa(c.age), Email: c.email})
	}
	return customers, nil
}

func (pdb *CustomerDB) GetCustomer(customerID string) (*api.Customer, error) {
	row := pdb.Db.QueryRow("Select * From customer WHERE id = '" + customerID + "'")
	var c Customer
	err := row.Scan(&c.id, &c.firstName, &c.lastName, &c.age, &c.email, &c.created, &c.modified)
	if err != nil {
		return nil, err
	}
	return &api.Customer{Id: strconv.FormatInt(c.id, 10), FirstName: c.firstName, LastName: c.lastName, Age: strconv.Itoa(c.age), Email: c.email}, nil
}

func (pdb *CustomerDB) DeleteCustomer(customerID string) (err error) {
	sqlStatement := `
		DELETE FROM customer
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

	resultsql, err := tx.ExecContext(ctx, sqlStatement, customerID)
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

func (pdb *CustomerDB) PostCustomer(customerReq *api.Customer) (_ string, err error) {

	timeStr := time.Now().Format(time.RFC3339)
	sqlStatement := fmt.Sprintf(`INSERT INTO customer (first_name, last_name, age, email, created_on, modified_on)
		VALUES ( '%s', '%s', '%s', '%s', '%s', '%s')
		RETURNING id`,
		customerReq.FirstName, customerReq.LastName, customerReq.Age, customerReq.Email, timeStr, timeStr)
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

func (pdb *CustomerDB) PutCustomer(customerReq *api.Customer) (err error) {
	sqlStatement := `
		UPDATE customer
		SET first_name = $1, last_name = $2, age = $3, email = $4, modified_on = $5
		WHERE id = $6
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
	_, err = tx.ExecContext(ctx, sqlStatement, customerReq.FirstName, customerReq.LastName, customerReq.Age, customerReq.Email, timeStr, customerReq.Id)
	if err != nil {
		return err
	}

	return nil
}
