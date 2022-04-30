//Package microservice specifies the models (struct) and the business logic for the application (Order microservice).
//Abstracts the complexity from executables (inside /cmd) and keeps them as clean as possible.
package microservice

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//Constants to specify the order status when adding/updating in 'orders' table
//The statuses also correspond to a typical Order lifecycle.
const (
	received    = "RECEIVED"
	in_progress = "IN PROGRESS"
	shipped     = "SHIPPED"
	delivered   = "DELIVERED"
	cancelled   = "CANCELLED"
)

//ProductOfOrder saves the products in order
type ProductOfOrder struct {
	ProductId     int64 `json:"product_id"`
	ProductNumber int64 `json:"product_number"`
}

//NewJSONInfo to save JSON information that post/put from user to server
type NewJSONInfo struct {
	Email         string           `json:"email"`
	Password      string           `json:"password"`
	Cmd           string           `json:"cmd"`
	Role          string           `json:"role"`
	ProductName   string           `json:"product_name"`
	ProductNumber int64            `json:"product_number"`
	ProductId     int64            `json:"product_id"`
	Order_Id      int64            `json:"order_id"`
	Products      []ProductOfOrder `json:"products"`
}

//Connect to DB and get the DbMap
var dbmap = connectToDB()

//Ping is used to check the health of the ORDER service.
//If service is up and running, it returns a status of '200 OK'
func Ping(c *gin.Context) {}

//isNumber validates if an input string is a valid number or not.
//Returns true if it is a valid number, else false.
func isNumber(id string) bool {
	if _, err := strconv.Atoi(id); err == nil {
		return true
	} else {
		return false
	}
}

//isEmpty validates if a string is empty or not.
func isEmpty(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	} else {
		return false
	}
}

//checking login. User must login fist to do every other functions except creating new user. Return value is user_id in "users" DB
func LoginChecking(e string, p string, cmd string, c *gin.Context) int64 {
	var user User
	var query = "SELECT * FROM users WHERE email='" + e + "' ORDER BY email"
	fmt.Println(query)

	err := dbmap.SelectOne(&user, query)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Some errors orrur."})
		return 0
	} else if len(user.Email) == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "Invalid email."})
		return 0
	} else if p != user.Password {
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized, "error": "Invalid password."})
		return 0
	} else {
		if cmd == "login" {
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Login success!", "UserInfo": user})
		}
		return user.Id
	} //End of else-block
}

//check if email, password, cmd are empty. Return false and bad request code to user
func BaseInfoChecking(e string, p string, cmd string, c *gin.Context) bool {

	if isEmpty(e) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Email cannot be empty. Pass a valid string value and try again !"})
		return false
	}

	if isEmpty(p) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Password cannot be empty. Pass a valid string value and try again !"})
		return false
	}

	if isEmpty(cmd) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Cmd cannot be empty. Pass a valid string value and try again !"})
		return false
	}
	return true
}

//checking role of email. Return value is true if this user is vendor role.
func IsVendor(e string, c *gin.Context) bool {
	var user User
	var query = "SELECT * FROM users WHERE email='" + e + "' ORDER BY email"

	dbmap.SelectOne(&user, query)

	return (user.Role == "vendor")
}

//POSThandler accepts the create new user, new product, new order in the DB.
func POSThandler(c *gin.Context) {
	var orderReq NewJSONInfo
	c.Bind(&orderReq)

	if BaseInfoChecking(orderReq.Email, orderReq.Password, orderReq.Cmd, c) == false {
		return
	}

	switch orderReq.Cmd {
	case "new_user":
		//Validate if the role is empty or not from request data. Role only can take "vendor" or "customer"
		if isEmpty(orderReq.Role) {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Role cannot be empty. Pass a valid string value and try again !"})
			return
		}
		if orderReq.Role != "customer" && orderReq.Role != "vendor" {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Role only has to be customer or vendor. Pass a valid string value and try again !"})
			return
		}

		//check if new email is already exist
		var user []User
		var query = "SELECT * FROM users WHERE email='" + orderReq.Email + "' ORDER BY email"
		fmt.Println(query)

		_, err := dbmap.Select(&user, query)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Some errors orrur."})
		} else if len(user) != 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Email is already exist."})
		} else {
			userInfo := &User{
				Email:    orderReq.Email,
				Password: orderReq.Password,
				Role:     orderReq.Role,
			}

			//Insert the new order single record in the "users" table.
			err := dbmap.Insert(userInfo)
			checkErr(err, "Add new customer failed in orders table")

			c.JSON(http.StatusCreated,
				gin.H{"status": http.StatusCreated, "message": "Create User is successful!", "email": userInfo.Email})
		}
	case "new_product":
		//checking login first to create new product
		if LoginChecking(orderReq.Email, orderReq.Password, orderReq.Cmd, c) == 0 {
			return
		}
		//checking role. Only vendor can create new product
		if IsVendor(orderReq.Email, c) == false {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Only vendors can modify products !"})
			return
		}
		//Validate if product name is empty
		if isEmpty(orderReq.ProductName) {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Product name cannot be empty. Pass a valid string value and try again !"})
			return
		}
		//Validate if product number is empty
		if orderReq.ProductNumber == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Product number cannot be empty. Pass a valid string value and try again !"})
			return
		}

		//Check if product name is already exist
		var product []Product
		var query = "SELECT * FROM products WHERE product_name='" + orderReq.ProductName + "' ORDER BY product_name"

		_, err := dbmap.Select(&product, query)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Some errors orrur."})
		} else if len(product) != 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Product is already exist."})
		} else {
			productInfo := &Product{
				Name:   orderReq.ProductName,
				Number: orderReq.ProductNumber,
			}
			//Insert the new product record in the "products" table.
			err := dbmap.Insert(productInfo)
			checkErr(err, "Add new product failed in orders table")

			c.JSON(http.StatusCreated,
				gin.H{"status": http.StatusCreated, "message": "Add new product is successfully!", "productName": productInfo.Name})
		}
	case "new_order":

		var user_id = LoginChecking(orderReq.Email, orderReq.Password, orderReq.Cmd, c)

		//checking login first to create new order
		if user_id == 0 {
			return
		}

		//Validate if the products array is empty or not.
		if len(orderReq.Products) == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Products cannot be empty. An order need to have at least 1 product. Add a product and try again !"})
			return
		}

		//order id is not auto increamental number. We set order id as the time created order
		order := &Order{
			Id:         time.Now().UnixNano(),
			CustomerId: user_id,
			Paied:      false,
		}

		//Insert the new order single record in the 'orders' table.
		err := dbmap.Insert(order)
		checkErr(err, "Add new order failed in orders table")

		//Iterate over the request data 'products array' and for each product
		//in this new order, add an entry into the 'orderproducts' table.
		for _, product := range orderReq.Products {
			orderProduct := &OrderProduct{
				OrderId:       order.Id,
				ProductId:     product.ProductId,
				ProductNumber: product.ProductNumber,
			}
			err := dbmap.Insert(orderProduct)
			checkErr(err, "Add new order_product mapping failed in order_products table")
		}

		c.JSON(http.StatusCreated,
			gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "resourceId": order.Id})
	}

}

//GET handler from DB based on the passed in path param. Execute get list of all products, get order information
func GETHandler(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	cmd := c.Query("cmd")

	//Valide if email, password, cmd are empty
	if BaseInfoChecking(email, password, cmd, c) == false {
		return
	}

	//Valide login success and return user_id
	user_id := LoginChecking(email, password, cmd, c)
	if user_id == 0 {
		return
	}

	switch cmd {
	//get list of all products
	case "get_product":
		var product []Product
		var query = "SELECT * FROM products;"
		_, err := dbmap.Select(&product, query)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Some errors orrur."})
		} else if len(product) == 0 {
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "There is no product at the moment"})
		} else {
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Get products success!", "Products": product})
		}
	//get list of order which belong to user_id. Other list of order will not disappear
	case "get_order":
		var order []Order
		var query = "SELECT * FROM orders WHERE user_id=" + strconv.Itoa(int(user_id))
		dbmap.Select(&order, query)
		for i, orderDetail := range order {
			var product []OrderProduct
			query = "SELECT * FROM orderproducts WHERE order_id=" + strconv.Itoa(int(order[i].Id))
			dbmap.Select(&product, query)
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "Order Info": orderDetail})
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "Product ordered Info": product})

		}
	}
}

//PUTHandler used to update product or update order. After fetching the old data, it updates with new data.
func PUTHandler(c *gin.Context) {

	email := c.Query("email")
	password := c.Query("password")
	cmd := c.Query("cmd")

	if BaseInfoChecking(email, password, cmd, c) == false {
		return
	}

	user_id := LoginChecking(email, password, cmd, c)
	if user_id == 0 {
		return
	}

	switch cmd {
	//Only vendor can update/remove product. Vendor must mention product id to update/remove product.
	//Name or number product is optional if having request to change
	case "update_product":
		if IsVendor(email, c) == false {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Only vendors can modify products !"})
			return
		}
		var product NewJSONInfo
		c.Bind(&product)
		if product.ProductId == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Product id cannot be empty. Pass a valid string value and try again !"})
			return
		}
		var query = "SELECT * FROM products WHERE product_id=" + strconv.Itoa(int(product.ProductId))
		var productLookup Product
		fmt.Println(query)

		err := dbmap.SelectOne(&productLookup, query)

		//Valide product id if exist then update it
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "No order with requested ID exists in the table. Invalid ID."})

		} else if len(productLookup.Name) == 0 {
			c.JSON(http.StatusNotFound,
				gin.H{"status": http.StatusNotFound, "error": "No product with requested ID exists in the table. Invalid ID."})
		} else {
			if product.ProductNumber != 0 {
				productLookup.Number = product.ProductNumber
			}
			if product.ProductName != "" {
				productLookup.Name = product.ProductName
			}
			dbmap.Update(&productLookup)
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Product Updated Successfully!", "product_id": productLookup})
		}
	//Only vendor can remove product. Vendor must mention product id to remove it
	case "remove_product":
		if IsVendor(email, c) == false {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Only vendors can modify products !"})
			return
		}
		var product NewJSONInfo
		c.Bind(&product)
		if product.ProductId == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Product id cannot be empty. Pass a valid string value and try again !"})
			return
		}
		var query = "SELECT * FROM products WHERE product_id=" + strconv.Itoa(int(product.ProductId))
		var productLookup []Product

		_, err := dbmap.Select(&productLookup, query)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "No order with requested ID exists in the table. Invalid ID."})

		} else if len(productLookup) == 0 {
			c.JSON(http.StatusNotFound,
				gin.H{"status": http.StatusNotFound, "error": "No product with requested ID exists in the table. Invalid ID."})
		} else {
			query = "DELETE FROM products WHERE product_id=" + strconv.Itoa(int(product.ProductId))
			dbmap.Exec(query)
			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Product deleted Successfully!", "product_id": product.ProductId})
		}
	//query exist order id, then update to orderproducts with the same order id
	case "update_order":
		var order NewJSONInfo
		c.Bind(&order)
		//Valide order id is empty
		if order.Order_Id == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Order id cannot be empty. Pass a valid string value and try again !"})
			return
		}
		var query = "SELECT * FROM orders WHERE order_id=" + strconv.Itoa(int(order.Order_Id)) + " AND user_id=" + strconv.Itoa(int(user_id))
		var order_query []Order
		dbmap.Select(&order_query, query)
		if len(order_query) == 0 {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "error": "Order id and user id dont match. Pass a valid string value and try again !"})
			return
		}
		query = "DELETE FROM orderproducts WHERE order_id=" + strconv.Itoa(int(order.Order_Id))
		dbmap.Exec(query)

		for _, product := range order.Products {
			orderProduct := &OrderProduct{
				OrderId:       order.Order_Id,
				ProductId:     product.ProductId,
				ProductNumber: product.ProductNumber,
			}
			err := dbmap.Insert(orderProduct)
			checkErr(err, "Add new order_product mapping failed in order_products table")
		}

		c.JSON(http.StatusCreated,
			gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "Order_Id": order.Order_Id})
	}
}
