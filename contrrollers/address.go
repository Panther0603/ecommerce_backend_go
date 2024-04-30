package contrrollers

import (
	"Ecommerce-Backend/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ivalid id"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.JSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var addressess models.Address
		addressess.Address_Id = primitive.NewObjectID()

		// filter
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}

		// unwind
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}

		// grouping
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		//aggretatting
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(500, "Internal server Error")
		}

		// now upadating tthe address
		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			log.Println(err)
			c.IndentedJSON(500, "internal server error")
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
			if size < 2 {

				filter := bson.D{primitive.E{Key: "_id", Value: address}}
				update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addressess}}}}
				_, err = UserCollection.UpdateOne(ctx, filter, update)
				if err != nil {
					log.Println(err)
				}
			} else {
				c.IndentedJSON(400, "Not allowed")
			}
			defer cancel()
			ctx.Done()
		}
	}
}

func EditHomeAddress() gin.HandlerFunc {

	return func(c *gin.Context) {

		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ivalid id"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			log.Panicln(err)
			c.JSON(500, "Internal server error")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var editAddress models.Address

		if err = c.BindJSON(&editAddress); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		// find the address
		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_number", Value: editAddress.House}, {Key: " address.0.street_name ", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pincode", Value: editAddress.Pincode}}}}

		// updatinng the fileterou out address with th update
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went worng"})
			c.Abort()
			return
		}

		defer cancel()
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "home address updated sucessfully"})
		ctx.Done()
	}
}

func EditOfficeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ivalid id"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			log.Panicln(err)
			c.JSON(500, "Internal server error")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var editAddress models.Address

		if err = c.BindJSON(editAddress); err != nil {
			c.JSON(400, err.Error())
			c.Abort()
			return
		}

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresss.1.house_number", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pincode", Value: editAddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusInternalServerError, "something went Wrong")
			c.Abort()
			return
		}
		defer cancel()
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, "office adress updated sucessfully")

		ctx.Done()
	}
}

func DeleteAddress() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
		}

		addresses := make([]models.Address, 0)

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "some issue with the id "})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.JSON(200, "Address sucessfully deleted")

	}
}
