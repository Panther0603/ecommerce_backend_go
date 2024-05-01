package routes

import (
	"Ecommerce-Backend/contrrollers"

	"github.com/gin-gonic/gin" // --> used for provinding the routes information
)

func UserRoutes(inccomigRoutes *gin.Engine) {
	inccomigRoutes.POST("/user/signup", contrrollers.SignUp())
	inccomigRoutes.POST("/login", contrrollers.Login())
	inccomigRoutes.POST("admin/addproduct", contrrollers.ProductViewerAdmin())
	inccomigRoutes.GET("/admin/viewproduct", contrrollers.SearchProduct())
	inccomigRoutes.GET("/users/search", contrrollers.SearchProductByQuery())
}
