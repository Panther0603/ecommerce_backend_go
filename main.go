package main

import (
	"Ecommerce-Backend/contrrollers"
	"Ecommerce-Backend/database"
	"Ecommerce-Backend/middleware"
	"Ecommerce-Backend/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "6000"

	app := contrrollers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, ":users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// some new routes

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.GetItemFromCart())
	router.GET("/intantbuy", app.IntantBuy())

	// running on a particular port
	fmt.Print("connecting to server")
	log.Fatal(router.Run(":" + port))
	fmt.Print("code is live at server " + port)

}
