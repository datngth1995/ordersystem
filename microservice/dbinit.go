package microservice

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

//Configuration type corresponds to the DB Configs JSON from the 'dbconf.json' file
type Configuration struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

//Order type corresponds to the 'orders' table in the DB
type Order struct {
	Id         int64 `db:"order_id"`
	CustomerId int64 `db:"user_id"`
	Paied      bool  `db:"paied"`
}

//OrderProduct type corresponds to the 'orderproducts' table in the DB
type OrderProduct struct {
	Id            int64 `db:"order_product_id"`
	OrderId       int64 `db:"order_id"`
	ProductId     int64 `db:"product_id"`
	ProductNumber int64 `db:"product_number"`
}

//Customer type corresponds to the 'customers' table in the DB
type User struct {
	Id       int64  `db:"user_id"`
	Role     string `db:"role"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

//Product type corresponds to the 'products' table in the DB
type Product struct {
	Id     int64  `db:"product_id" json:"product_id"`
	Name   string `db:"product_name"`
	Number int64  `db:"product_number"`
}

func InitDB() {

	//Invoke method to connect to DB and add table definitions to dbmap
	dbmap := connectToDB()
	defer dbmap.Db.Close()

	// Create the tables. In a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err := dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
}

//connectToDB establishes a connection to the DB, then create TABLE users, products, orders. orderproducts
func connectToDB() *gorp.DbMap {

	//get configuration of server from dbconf.json
	config := loadDBConfigs("../../configs/dbconf.json")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	// connect to db
	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err, "sql.Open failed")

	err = db.Ping()
	checkErr(err, "Connection to DB or ping failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	//initialize all tables
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Product{}, "products").SetKeys(true, "Id")
	dbmap.AddTableWithName(Order{}, "orders").SetKeys(false, "Id")
	dbmap.AddTableWithName(OrderProduct{}, "orderproducts").SetKeys(true, "Id")
	return dbmap
}

//loadDBConfigs loads configuration of server. Then return it under json format
func loadDBConfigs(filepath string) Configuration {

	configFile, err := os.Open(filepath)
	defer configFile.Close()
	checkErr(err, "Error reading DB configs from JSON file")
	jsonParser := json.NewDecoder(configFile)
	config := Configuration{}
	jsonParser.Decode(&config)
	return config
}

//checkErr checks for error and logs when present.
func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
