# Apigo
Technical test for Aifi

## 1- How to run
### Prerequisites

To run the API you need to have Docker and Docker compose installed on the machine.
- Docker: https://docs.docker.com/get-docker/
- Docker compose https://docs.docker.com/compose/install/
- AB apache, to do AB tests: https://www.tutorialspoint.com/apache_bench/apache_bench_environment_setup.htm

### Build and run
In the main folder of the project. Where the file compose.yml is located. Execute:
```
# docker-compose up -d --build
```
This command build the Dockerfile and pull postgres, adminer and swagger from DockerHub.
Once docker-compose has finished building and running the images. We can check that apigo is running.
```
# docker-compose logs apigo
```
#### Swagger
```
http://localhost:8092
```
Init the browser in this url and access to the documentation of the API of api go to know the differents requests and responses.

![swagger1](https://github.com/carlos2380/webCarlos2380/blob/master/swagger1.png)

You can interact with swagger and do request and see the response

![swagger2](https://github.com/carlos2380/webCarlos2380/blob/master/swagger2.png)
![swagger3](https://github.com/carlos2380/webCarlos2380/blob/master/swagger3.png)

## 2- Performance
I check the performance using de client and the server in the same host. The results are worse than in a real enviroment where the client and the server don't share resources

The enviroment to test the performace were ubuntu virtual machine with 3 CPUs and 3GB RAM.

### Own client
I created a simple client to test the performace of a server. 

https://github.com/carlos2380/apigo/blob/main/cmd/client/main.go

To execute the client, affter to do the execution of compose. we create the build
```
# docker build -t client --target client .
```
and then we can run the client
```
# docker run -it client sh -c "/client -c 1 -nc 50000 --url=172.17.0.1:8000/api/customers"
```
Where c is the number of threads, nc the number of transactions per threads and url the url to do get.

#### Results
Executing concurrecy 1 and 50000 transactions per thread we have a TPS of 3435 using un 60% of CPU
![own test1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1cresult.png)
![own top1](https://github.com/carlos2380/webCarlos2380/blob/master/myclient1ccpu.png)

Executing concurrecy 2 and 50000 transactions per thread we have a TPS of 6020 using un 80% of CPU
![own test2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2cresult.png)
![own top2](https://github.com/carlos2380/webCarlos2380/blob/master/myclient2ccpu.png)

### AB Apache
Using AB testing is easy to check the transactions per second specifying the number of threads and the number total transactions 
```
# ab -k -c 2 -n 100000 http://172.17.0.1:8000/api/customers
```
Where -c is the number of cores and -n the number of total transactions.

#### Results
Executing concurrecy 1 and 100000 transaction we have a TPS of 4207 using un 73% of CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab1ccpu.png)

Executing concurrecy 2 and 100000 transactions we have a TPS of 7527 using un 92% of CPU
![ab test2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2cresult.png)
![ab top2](https://github.com/carlos2380/webCarlos2380/blob/master/ab2ccpu.png)

### Conclusion

The results are similar, whith the doble of the CPU we can duplicate the transactions, then GO is good and easy to paralelize the threads.
We get better results with AB because AB use less CPU than my own client and client and server are sharing resources.

## 3- Documentation

## 4- Next steps



