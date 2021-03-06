package postgres

import (
	"apigo/api"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type CustomerDB struct {
	DB *sql.DB
}
type Customer struct {
	id        int64
	firstName *string
	lastName  *string
	age       *int
	email     *string
	storeID   int64
	created   time.Time
	modified  time.Time
}

func (pdb *CustomerDB) GetCustomers() (*api.CustomersJSON, error) {
	rows, err := pdb.DB.Query("Select * From customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	customers := []*api.Customer{}
	for rows.Next() {
		var c Customer
		err = rows.Scan(&c.id, &c.firstName, &c.lastName, &c.age, &c.email, &c.storeID, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}

		age := ""
		if c.age != nil {
			age = strconv.Itoa(*c.age)
		}

		customers = append(
			customers,
			&api.Customer{
				ID:        strconv.FormatInt(c.id, 10),
				FirstName: c.firstName, LastName: c.lastName,
				Age:     &age,
				Email:   c.email,
				StoreID: strconv.FormatInt(c.storeID, 10),
			},
		)
	}
	return &api.CustomersJSON{Customers: customers}, nil
}

func (pdb *CustomerDB) GetCustomer(customerID string) (*api.Customer, error) {
	row := pdb.DB.QueryRow(fmt.Sprintf("Select * From customers WHERE id = '%s'", customerID))
	var c Customer
	err := row.Scan(&c.id, &c.firstName, &c.lastName, &c.age, &c.email, &c.storeID, &c.created, &c.modified)
	if err != nil {
		return nil, err
	}
	age := ""
	if c.age != nil {
		age = strconv.Itoa(*c.age)
	}
	return &api.Customer{ID: strconv.FormatInt(c.id, 10),
		FirstName: c.firstName,
		LastName:  c.lastName,
		Age:       &age,
		Email:     c.email,
		StoreID:   strconv.FormatInt(c.storeID, 10)}, nil
}

func (pdb *CustomerDB) DeleteCustomer(customerID string) (err error) {
	sqlStatement := `
		DELETE FROM customers
		WHERE id = $1
	`
	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)
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
	sqlStatement := fmt.Sprintf(`INSERT INTO customers (first_name, last_name, age, email, store_id, created_on, modified_on)
		VALUES ( '%s', '%s', '%s', '%s', '%s', '%s', '%s')
		RETURNING id`,
		*customerReq.FirstName, *customerReq.LastName, *customerReq.Age, *customerReq.Email, customerReq.StoreID, timeStr, timeStr)
	lastInsertid := int64(0)

	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)
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
		UPDATE customers
		SET first_name = $1, last_name = $2, age = $3, email = $4, store_id = $5, modified_on = $6
		WHERE id = $7
	`
	timeStr := time.Now().Format(time.RFC3339)

	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)

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
	resultsql, err := tx.ExecContext(
		ctx,
		sqlStatement,
		customerReq.FirstName,
		customerReq.LastName,
		customerReq.Age, customerReq.Email,
		customerReq.StoreID, timeStr, customerReq.ID)
	if err != nil {
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

func (pdb *CustomerDB) GetCustomersByStoreID(storeID string) (*api.CustomersJSON, error) {
	rows, err := pdb.DB.Query(fmt.Sprintf("Select * From customers Where store_id = '%s'", storeID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	customers := []*api.Customer{}
	for rows.Next() {
		var c Customer
		err = rows.Scan(&c.id, &c.firstName, &c.lastName, &c.age, &c.email, &c.storeID, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}

		age := ""
		if c.age != nil {
			age = strconv.Itoa(*c.age)
		}

		customers = append(
			customers,
			&api.Customer{
				ID:        strconv.FormatInt(c.id, 10),
				FirstName: c.firstName, LastName: c.lastName,
				Age:     &age,
				Email:   c.email,
				StoreID: strconv.FormatInt(c.storeID, 10),
			},
		)
	}
	return &api.CustomersJSON{Customers: customers}, nil
}
