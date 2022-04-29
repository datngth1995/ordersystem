package main

import (
	"log"
	"order/microservice"

	_ "github.com/lib/pq"
)

func main() {

	// initialize the DB and create required tables
	microservice.InitDB()
	log.Println("Initialized the DB successfully")
}
