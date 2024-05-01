package contrrollers

import (
	"Ecommerce-Backend/database"
	"Ecommerce-Backend/models"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	userCollection *mongo.Collection
	prodCollection *mongo.Collection
}

func NewApplication(userCollection *mongo.Collection, prodCollection *mongo.Collection) *Application {
	return &Application{
		userCollection: userCollection,
		prodCollection: prodCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryId := c.Query("id")

		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New(""))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "issuse while mapping"})
		}

		c.JSON(200, "sucessfully added to cart")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		productQueryId := c.Query("id")

		if productQueryId == "" {
			c.AbortWithError(http.StatusInternalServerError, errors.New("product id is empty"))
			return
		}
		userId := c.Query("userId")

		if userId == "" {
			c.AbortWithError(http.StatusInternalServerError, errors.New("user id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("error while creating prodid"))
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "sucessfully item removed fromcart")
	}

}

func (app *Application) GetItemFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			log.Panicln("user is is empty")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("user id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.prodCollection, app.userCollection, userID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "order placed success fully")
	}
}

func BuyFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "id is invalid"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.Users

		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "user not found")
			return
		}

		// here we will use the aggregation concept and the stage to achive the aggregation

		// first filter
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}

		// unwind

		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}

		// grouping

		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "usercart.price"}}}}}}

		// after doing all the  stage we need to run the aggregation stages in a pipeline
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil {
			log.Println(err)
			c.JSON(500, "Error")
			return
		}
		var listing []bson.M

		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithError(500, errors.New("Error"))
			return
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UseCart)
		}
		ctx.Done()

		//c.JSON(http.StatusOK, listing)
	}
}

func (app *Application) IntantBuy() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryId := c.Query("id")

		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New(""))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.IntantBuy(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(200, "sucessfully sucess fully placed the order")

	}
}
