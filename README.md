# mw-backend-test

Simple E-Commerce RESTful API using 3 services on mono repository(User, Product, Transaction)

## Building and Running The App

Prerequisites:

1. Golang 1.14
2. Golang migrate (https://github.com/golang-migrate/migrate)

The App is separate into three services, User service, Product service, Transaction service but in one App

**Step 1 Checkout**

```bash
$ git clone https://github.com/arieffian/mw-backend-test.git
$ cd mw-backend-test
```

**Step 2 Start MySQL Service**

```bash
$ cd deployments
$ docker-compose up mysql
```

**Step 2 Run Migration**

```bash
$ migrate -database mysql://mw-backend:mw-backend@/mw-backend -path ./sql up

```

**Step 2 Build and Run**

```bash
$ go build -ldflags "-s -w" -o build/mw-backend-test.app cmd/main.go
```

***Running the app as user service***

```bash
$ mw-backend-test.app user
```

***Running the app as product service***

```bash
$ mw-backend-test.app user
```

***Running the app as transaction service***

```bash
$ mw-backend-test.app transaction
```

## Running using Docker Container

Prerequisites:

1. Docker
2. Docker Compose

For simple deployment or testing you can just run the docker compose file

```bash
$ cd deployments
$ docker-compose up 
```

## Testing App

```bash
$ go test
``` 

## Calling API##
```bash
$ curl -X POST -H 'content-type: application/json' --data '{"name": "acer"}' http://localhost:8080/brand
``` 

```bash
$ curl http://localhost:8080/product?id=1
``` 

```bash
$ curl -X POST -H 'content-type: application/json' --data '{"brand_id": 4, "name": "predator", "qty": 3, "price": 1050}' http://localhost:8080/product
``` 

```bash
$ curl http://localhost:8080/product/brand?id=1
``` 
