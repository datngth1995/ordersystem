# A simple order system 
## Introduction
This service exposes the APIs for users to order products from vendors online that also tracks vendor inventory. The details below:
- Create new user with those informations: email, password, role(vendor or customer)
- Checking login
- Add/View/Remove products from inventory
- Create/Update an order

Tech stack includes:
- Golang
- RESTfull API
- Gin framework
- Postgresql

I only implement backend side. The REST endpoint can be invoked using any REST client application such as POSTMAN, etc.
## Installation
Make sure those packages are installed on your computer:
- Postgresql
- Golang
- Go Gin package dependencies: running these commands to get Gin package from **$GOPATH/src** 
```
go get github.com/gin-gonic/gin 
go get gopkg.in/gorp.v1 
go get -u github.com/lib/pq
```
## Running service
### Setup Postgresql
1. Login to Postgresql (default user is postgres)
2. Create new database dordersystem. You can create manually or run the scripts that I save in **$GOPATH/src/ordersystem/scripts/db**
3. Put correct database configuration. You might want to change **$GOPATH/src/ordersystem/configs/dbconf.json** to accomplish to your current configuration, for example password field. Our go application will use this file to connect to Postgresql.
```
{
  "host": "localhost",
  "port": 5432,
  "user": "postgres",
  "password": "123",
  "dbname": "dordersystem"
}
```
4. Initializing data base. Run db.go in **$GOPATH/src/ordersystem/cmd/db**. It also creates four tables in dordersystem database. The detail implementation of entity relationship refers to **$GOPATH/src/ordersystem/doc**
```
go run db.go
```
### Starting service
1. To start service, run main.go located at **$GOPATH/src/ordersystem/cmd/main**. Go application service locally
```
go run main.go
```
2. Using any REST client such as POSTMAN to execute add/update/get operation. The detail implementation and testing refer to **$GOPATH/src/ordersystem/doc**
## Limitation
Due to the given time constraints, it exposes only basic functions. There much things can improve to give better experience to user.
The limitations includes:
1. Validation and update product number after each paied order can be supported in next release.
2. Doesn't support export database into CSV
3. Front end implementation can be easy for testing
