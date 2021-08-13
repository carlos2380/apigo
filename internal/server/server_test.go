package server_test

import (
	"apigo/internal/server"
	"apigo/internal/storage/postgres"
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func exec(db *sql.DB, comand string) *sql.DB {
	var err error
	_, err = db.Exec(comand)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestServer(t *testing.T) {
	password := "secret"
	ip := "172.17.0.1"
	port := "5432"
	db, err := postgres.InitPostgres(&password, &ip, &port)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = postgres.ClosePostgres(db)
		if err != nil {
			log.Fatal(err)
		}
	}()

	db = exec(db, "DELETE FROM cases")
	db = exec(db, "DELETE FROM customers")
	db = exec(db, "DELETE FROM stores")
	db = exec(db, "ALTER SEQUENCE cases_id_seq RESTART")
	db = exec(db, "ALTER SEQUENCE customers_id_seq RESTART")
	db = exec(db, "ALTER SEQUENCE stores_id_seq RESTART")

	store := &postgres.StoreDB{DB: db}
	customer := &postgres.CustomerDB{DB: db}
	cs := &postgres.CaseDB{DB: db}

	srv := httptest.NewServer(server.NewRouter(store, customer, cs))
	defer srv.Close()

	tableTest := []struct {
		desc               string
		httpMethod         string
		url                string
		bodyReqStr         string
		bodyRespStr        string
		expectedStatusCode int
	}{
		{
			"Test 1: Get a stores that doesn't exist",
			http.MethodGet,
			srv.URL + "/api/stores/1000",
			"",
			`{"error":"Item Not Found"}`,
			http.StatusNotFound,
		},
		{
			"Test 2: Create a store",
			http.MethodPost,
			srv.URL + "/api/stores",
			`{"id":"1","name":"store1","address":"address1"}`,
			`{"id":"1"}`,
			http.StatusOK,
		},
		{
			"Test 3: Get the store 1",
			http.MethodGet,
			srv.URL + "/api/stores/1",
			"",
			`{"id":"1","name":"store1","address":"address1"}`,
			http.StatusOK,
		},
		{
			"Test 4: Delete a store that doesn't exist",
			http.MethodDelete,
			srv.URL + "/api/stores/100",
			"",
			`{"error":"Item Not Found"}`,
			http.StatusNotFound,
		},
		{
			"Test 5: Create a store",
			http.MethodPost,
			srv.URL + "/api/stores",
			`{"id":"2","name":"store2","address":"address2"}`,
			`{"id":"2"}`,
			http.StatusOK,
		},
		{
			"Test 6: Get all stores",
			http.MethodGet,
			srv.URL + "/api/stores",
			"",
			`{"stores":[{"id":"1","name":"store1","address":"address1"},{"id":"2","name":"store2","address":"address2"}]}`,
			http.StatusOK,
		},
		{
			"Test 7: Delete store",
			http.MethodDelete,
			srv.URL + "/api/stores/1",
			"",
			"",
			http.StatusOK,
		},
		{
			"Test 8: Get store that is deleted",
			http.MethodGet,
			srv.URL + "/api/stores/1",
			"",
			`{"error":"Item Not Found"}`,
			http.StatusNotFound,
		},
		{
			"Test 9: Create a store with wrong parameters",
			http.MethodPost,
			srv.URL + "/api/stores",
			`{"names":"store2","adxress":"address2"}`,
			`{"error":"Bad Request"}`,
			http.StatusBadRequest,
		},
		{
			"Test 10: Do a delete in stores",
			http.MethodDelete,
			srv.URL + "/api/stores",
			"",
			"",
			http.StatusMethodNotAllowed,
		},
		{
			"Test 11: Create customer",
			http.MethodPost,
			srv.URL + "/api/customers",
			`{"first_name":"jhon","last_name":"connor","age":"30","email":"jconnor@gmail.com"}`,
			`{"id":"1"}`,
			http.StatusOK,
		},
		{
			"Test 12: Get customer by Id",
			http.MethodGet,
			srv.URL + "/api/customers/1",
			``,
			`{"id":"1","first_name":"jhon","last_name":"connor","age":"30","email":"jconnor@gmail.com"}`,
			http.StatusOK,
		},
		{
			"Test 13: Put on customer 1",
			http.MethodPut,
			srv.URL + "/api/customers/1",
			`{"id":"1","first_name":"mick","last_name":"addams","age":"20","email":"maddams@gmail.com"}`,
			``,
			http.StatusOK,
		},
		{
			"Test 14: Get customers",
			http.MethodGet,
			srv.URL + "/api/customers",
			``,
			`{"customers":[{"id":"1","first_name":"mick","last_name":"addams","age":"20","email":"maddams@gmail.com"}]}`,
			http.StatusOK,
		},
		{
			"Test 15: Delete by Id",
			http.MethodDelete,
			srv.URL + "/api/customers/1",
			``,
			``,
			http.StatusOK,
		},
		{
			"Test : Create customer",
			http.MethodPost,
			srv.URL + "/api/customers",
			`{"id":"2","first_name":"jhon","last_name":"connor","age":"30","email":"jconnor@gmail.com"}`,
			`{"id":"2"}`,
			http.StatusOK,
		},
		{
			"Test 16: Create case",
			http.MethodPost,
			srv.URL + "/api/cases",
			`{"start_time_stamp":"2021-08-12 04:35:36","end_time_stamp":"2021-08-12 04:35:36","customer_id":"2","store_id":"2"}`,
			`{"id":"1"}`,
			http.StatusOK,
		},
		{
			"Test 17: Create case",
			http.MethodPost,
			srv.URL + "/api/cases",
			`{"start_time_stamp":"2021-08-14 04:35:36","end_time_stamp":"2021-08-14 04:35:36","customer_id":"2","store_id":"2"}`,
			`{"id":"2"}`,
			http.StatusOK,
		},
		{
			"Test 14: Get cases",
			http.MethodGet,
			srv.URL + "/api/cases",
			``,
			`{"cases":[{"id":"1","start_time_stamp":"2021-08-12 04:35:36 +0000 +0000","end_time_stamp":"2021-08-12 04:35:36 +0000 +0000","customer_id":"2","store_id":"2"},{"id":"2","start_time_stamp":"2021-08-14 04:35:36 +0000 +0000","end_time_stamp":"2021-08-14 04:35:36 +0000 +0000","customer_id":"2","store_id":"2"}]}`,
			http.StatusOK,
		},
		{
			"Test 19: Delete by Id",
			http.MethodDelete,
			srv.URL + "/api/cases/1",
			``,
			``,
			http.StatusOK,
		},
		{
			"Test: Delete by Id",
			http.MethodDelete,
			srv.URL + "/api/customers/2",
			``,
			``,
			http.StatusOK,
		},
		{
			"Test 20: Get cases",
			http.MethodGet,
			srv.URL + "/api/cases",
			``,
			`{"cases":[]}`,
			http.StatusOK,
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, tt.httpMethod, tt.url, strings.NewReader(tt.bodyReqStr))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Fatalf("Expected Status Code %d but found %d", tt.expectedStatusCode, resp.StatusCode)
			}

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			respStr := string(respBytes)
			if strings.TrimRight(respStr, "\n") != tt.bodyRespStr {
				t.Fatalf("Expected body %s but found %s", tt.bodyRespStr, respStr)
			}
		})
	}
}
