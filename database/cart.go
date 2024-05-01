package database

import (
	"Ecommerce-Backend/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// here defing some of teh general error

	ErrorCantFindProduct   = errors.New("can't find the producrs ")
	ErrorCantDecodeProduct = errors.New("can't able to decode the product")
	ErrUserIdIsNotValid    = errors.New("this user is not valid ")
	ErrCantUpdateUser      = errors.New("not able to update user")
	ErrCantRemoveCartItem  = errors.New("cannot able to remove this item from cart")
	ErrCantGetItem         = errors.New("was unable tp get the item, please try after sometime")
	ErrCantBuyCartItem     = errors.New("can not update purache order ")
)

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) (err error) {

	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productId})

	if err != nil {
		return ErrorCantFindProduct
	}
	var productCart []models.Product

	err = searchFromDb.All(ctx, &productCart)

	if err != nil {
		return ErrorCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, proodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productId}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, proodCollection *mongo.Collection, userCollection *mongo.Collection, userID string) error {

	// things will be done in this function  because these thing happen during the making buy

	// fetch the cart of user
	// find cart total
	//added order to user collection
	// added item in the cart to order list
	// empty the user cart

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.Users
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}

	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "usercart.price"}}}}}}
	pointCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	if err != nil {
		panic(err)
	}

	var getusercart []bson.M

	err = pointCursor.All(ctx, &getusercart)

	if err != nil {
		panic(err)
	}

	var tot_price int32
	for _, user_item := range getusercart {

		price := user_item["total"]
		tot_price = price.(int32)
	}

	orderCart.Price = uint64(tot_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Panic(err)
		return ErrCantBuyCartItem
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Panic(err)
		return ErrCantBuyCartItem
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order.$[].order_list": bson.M{"$each": getCartItems.UseCart}}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		log.Panic(err)
		return ErrCantBuyCartItem
	}

	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)

	if err != nil {
		log.Panic(err)
		return ErrCantBuyCartItem
	}

	return nil

}

func IntantBuy(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Panic(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.Product
	var order_details models.Order

	order_details.Order_ID = primitive.NewObjectID()
	order_details.Ordered_At = time.Now()
	order_details.Order_Cart = make([]models.ProductUser, 0)
	order_details.Payment_Method.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&product_details)

	if err != nil {
		log.Panic(err)
		return ErrorCantFindProduct
	}

	order_details.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_details}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return ErrCantUpdateUser
	}
	filter1 := bson.D{primitive.E{Key: "_id", Value: id}}
	update1 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}

	_, err = userCollection.UpdateOne(ctx, filter1, update1)

	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}
