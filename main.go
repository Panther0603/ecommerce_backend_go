package main

import (
	"Ecommerce-Backend/contrrollers"
	"Ecommerce-Backend/database"
	"Ecommerce-Backend/middleware"
	"Ecommerce-Backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := contrrollers.NewApplication(database.ProductData(database.Client, "Prodicts"), database.UserData(database.Client, ":users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Aunthication())

	// some new routes

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.removeItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/intantbuy", app.IntantBuy())

	// running on a particular port
	log.Fatal(router.Run(":" + port))

}
