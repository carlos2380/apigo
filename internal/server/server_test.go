package server_test

import (
	"apigo/internal/server"
	"apigo/internal/storage/postgres"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	storage, err := postgres.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	defer storage.(*postgres.PostgresDB).CloseDB()
	srv := httptest.NewServer(server.NewRouter(storage))
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
			`{"name":"store1","address":"address1"}`,
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
			"",
			http.StatusOK,
		},
		{
			"Test 5: Create a store",
			http.MethodPost,
			srv.URL + "/api/stores",
			`{"name":"store2","address":"address2"}`,
			`{"id":"2"}`,
			http.StatusOK,
		},
		{
			"Test 6: Get all stores",
			http.MethodGet,
			srv.URL + "/api/stores",
			"",
			`[{"id":"1","name":"store1","address":"address1"},{"id":"2","name":"store2","address":"address2"}]`,
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
	}

	for _, tt := range tableTest {
		t.Run(tt.desc, func(t *testing.T) {
			req, err := http.NewRequest(tt.httpMethod, tt.url, strings.NewReader(tt.bodyReqStr))
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
