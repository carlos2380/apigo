# Apigo
Technical test for Aifi

## 1- How to run
### Prerequisites

To run the API you need to have Docker and Docker compose installed on the machine.
- Docker (Min version: 20.10.7): https://docs.docker.com/get-docker/
- Docker compose (Min version: 1.29.2): https://docs.docker.com/compose/install/
- AB apache, to do AB tests: https://en.wikipedia.org/wiki/ApacheBench

### Build and run
In the main folder of the project, where the file compose.yml is located. Execute:
```
# docker-compose up -d --build
```
This command builds the Dockerfile and pull postgres, adminer and swagger from DockerHub.
Once docker-compose has finished building and running the images. We can check that apigo is running
```
# docker-compose logs apigo
```
Now it's possible get the stores for example.
```
# curl http://localhost:8000/api/stores
```

#### Swagger
Swagger allows you to visualize and interact with the API.
```
http://localhost:8092
```
Initialize the browser on this URL and access the apigo API documentation to know the different requests and responses.


![swagger1](https://github.com/carlos2380/webCarlos2380/blob/master/swagger1.png)

You can interact with swagger and make requests and see the responses.

![swagger2](https://github.com/carlos2380/webCarlos2380/blob/master/swagger2.png)
![swagger3](https://github.com/carlos2380/webCarlos2380/blob/master/swagger3.png)

- The Apigo documentation that uses swagger to work is here: https://github.com/carlos2380/apigo/blob/main/swagger.yml

#### Adminer
It's a simple database manager.
To access to Adminer go to the next url:

```
http://localhost:8081
```

And set up configuration as:

![adminer](https://github.com/carlos2380/webCarlos2380/blob/master/adminer.png)
- The default password is: **secret**.


## 2- Performance
I tested the performance using the client and server on the same host. The results are different than in a real environment where the client and server do not share resources.

The environment to test the performance was an ubuntu virtual machine with 3 CPUs and 3GB of RAM.

### Own client
I have created a simple client to test server performance. 

https://github.com/carlos2380/apigo/blob/main/cmd/client/main.go

To run the client, after running the compose. Create the build:
```
# docker build -t client --target client .
```
and then, run the client.
```
# docker run -it client sh -c "/client -c 1 -nc 50000 -url http://172.17.0.1:8000/api/customers"
```
Where c is the number of threads, nc the number of transactions per thread and url the url to do the get.

#### Results
Executing concurrency 1 and 50000 transactions per thread we have a TPS (Transactions Per Second) of 3435 using 60% of a CPU
![own test1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1cresult.png)
![own top1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1ccpu.png)

Executing concurrency 2 and 50000 transactions per thread we have a TPS of 6020 using 80% of a CPU
![own test2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2cresult.png)
![own top2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2ccpu.png)

### AB Apache
Using AB testing it is easy to check the transactions per second specifying the number of threads and the number of total transactions 
```
# ab -k -c 2 -n 100000 http://172.17.0.1:8000/api/customers
```
Where -c is the number of concurrency and -n the number of total transactions.

#### Results
Executing concurrency 1 and 100000 transaction we have a TPS of 4207 using 73% of a CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1ccpu.png)

Executing concurrency 2 and 100000 transactions we have a TPS of 7527 using 92% of a CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2ccpu.png)

### Conclusion

Results obtained were similar, when the concurrency is doubled, the transactions are doubled, then GO is good and easy to parallelize the threads.
I got better results with AB because uses less CPU than my own client and both are sharing resources with apigo.

## 3- Documentation

### Database
The schema of the exercice is this:

| Customer  | Store     | Case           |
|-----------|-----------|----------------|
| ID        | ID        | ID             |
| FirstName | Name      | StartTimestamp |
| LastName  | Address   | EndTimestamp   |
| Age       | Customers | CustomerID     |
| Email     | Cases     | StoreID        |
| StoreID   |           |                |
| Cases     |           |                |

The UML that represents this scheme is:

![UML](https://github.com/carlos2380/webCarlos2380/blob/master/uml.png)

I have used Postgres as a rational database. This is the file with the implementation of the tables:

- https://github.com/carlos2380/apigo/blob/main/samples/initdb/create_tables.sql

I have enabled delete cascade to simplify the logic.

### Interface storage

To decouple the API with the database. I have created an interface for each model. This way I can have each model stored in different databases and also decouple the API.

The implementation of Storage now is this:

![postgresinterface](https://github.com/carlos2380/webCarlos2380/blob/master/postgresInter.png)

But it could be this:

![interinterface](https://github.com/carlos2380/webCarlos2380/blob/master/randInterface.png)

For this reason in functions like the following:

``` GO
func (sHandler *StorageHandler) GetStoreByCaseID(w http.ResponseWriter, r *http.Request) {
  //Code
  retCase, err := sHandler.StgCase.GetCase(id)
  //Code
  params["id"] = retCase.StoreID
  sHandler.GetStore(w, r)
  //Code
}
```

I make more than one call to the database because the information may be in different databases.

The SQL command for postgres would be this: (I supose that asks for case.id = 1 but can be any id to compare):

```SQL
SELECT stores
FROM cases
INNER JOIN stores ON cases.store_id = stores.id
WHERE cases.id = 1
```
### Handler

Because the functions in the router don't accept more parameters (only w http.ResponeWriter, r *http.Request)

I created a struct that has the information about the storage.

```GO
type StorageHandler struct {
	StgStore    storage.StoreStorage
	StgCustomer storage.CustomerStorage
	StgCase     storage.CaseStorage
}
```

then the function has the StorageHanler

```GO
func (sHandler *StorageHandler) GetCases(w http.ResponseWriter, r *http.Request)

```
and the function still doesn't have more parameters.

```GO
r.HandleFunc("/api/cases", stgHandler.GetCases).Methods(http.MethodGet, http.MethodOptions)

```

### Transactions

I added transactions on methods post, put and delete.

```GO
//I added err to return
func (pdb *CustomerDB) DeleteCustomer(customerID string) (err error) {
	// CODE
	//I get the context and the tx
	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	//At the end of the function, if all was correct I tried to commit and I returned the err
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()
	// More code 
}

```

### Close server properly

To prevent the server from shutting down while there are pending tasks. The server captures any shutdown signal and waits for pending tasks to finish.

The server has programmed a counter that after 5 seconds of receiving the shutdown signal. The server will shut down, avoiding stuck tasks and the server never shuts down.

```GO
func main() {
	//CODE
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

```

### Flags

I have added flags to start the application with different setups.

```GO
	port := flag.String("port", "8000", "Port the server will be listening")
	dbDriver := flag.String("driver", "", "Database driver (postgres only supported now)")
	dbPassword := flag.String("password", "", "Password of the database")
	dbIPHost := flag.String("host", "172.17.0.1", "Host IP of the database")
	dbPort := flag.String("dbport", "5432", "Port of the database")
	flag.Parse()
```

### TableTests

I use Table driven test to make the tests. In this way, there is a very simplified code and it works for all test cases.

- https://github.com/carlos2380/apigo/blob/main/internal/server/server_test.go

### CORS (Control Access HTTP)

I enabled CORS to be able to connect swagger with apigo.

```GO
func setHeaders(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}
```

```GO
	r.HandleFunc("/api/stores/{id}", stgHandler.GetStore).Methods(http.MethodGet, http.MethodOptions)
```
