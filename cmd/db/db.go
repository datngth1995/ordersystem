package main

import (
	"log"
	"ordersystem/microservice"

	_ "github.com/lib/pq"
)

func main() {

	// initialize the DB and create required tables
	microservice.InitDB()
	log.Println("Initialized the DB successfully")
}
