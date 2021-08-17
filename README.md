# Apigo
Technical test for Aifi

## 1- How to run
### Prerequisites

To run the API you need to have Docker and Docker compose installed on the machine.
- Docker (Min version: 20.10.7): https://docs.docker.com/get-docker/
- Docker compose (Min version: 1.29.2): https://docs.docker.com/compose/install/
- AB apache, to do AB tests: https://www.tutorialspoint.com/apache_bench/apache_bench_environment_setup.htm

### Build and run
In the main folder of the project. Where the file compose.yml is located. Execute:
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
Swagger allows to visualize and interact with the API.
```
http://localhost:8092
```
Initialize the browser on this url and access the apigo API documentation to know the different requests and responses.


![swagger1](https://github.com/carlos2380/webCarlos2380/blob/master/swagger1.png)

You can interact with swagger and make requests and see the responses.

![swagger2](https://github.com/carlos2380/webCarlos2380/blob/master/swagger2.png)
![swagger3](https://github.com/carlos2380/webCarlos2380/blob/master/swagger3.png)

- The Apigo documentation that uses swagger to work is here: https://github.com/carlos2380/apigo/blob/main/swagger.yml

#### Adminer
Is a simple database manager.
To access to Adminer go to the next url:
```
http://localhost:8081
```
And set up configuration as:
![adminer](https://github.com/carlos2380/webCarlos2380/blob/master/adminer.png)
- Password is **secret** by default.


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
Executing concurrnecy 1 and 50000 transactions per thread we have a TPS (Transactions Per Second) of 3435 using un 60% of CPU
![own test1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1cresult.png)
![own top1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1ccpu.png)

Executing concurrnecy 2 and 50000 transactions per thread we have a TPS of 6020 using un 80% of CPU
![own test2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2cresult.png)
![own top2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2ccpu.png)

### AB Apache
Using AB testing is easy to check the transactions per second specifying the number of threads and the number of total transactions 
```
# ab -k -c 2 -n 100000 http://172.17.0.1:8000/api/customers
```
Where -c is the number of cores and -n the number of total transactions.

#### Results
Executing concurrnecy 1 and 100000 transaction we have a TPS of 4207 using un 73% of CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1ccpu.png)

Executing concurrnecy 2 and 100000 transactions we have a TPS of 7527 using un 92% of CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2ccpu.png)

### Conclusion

The results are similar, with double the CPU we can double the transactions, then GO is good and easy to parallelize the threads.
We get better results with AB because AB uses less CPU than my own client and the client and server are sharing resources.

## 3- Documentation

## 4- Next steps



