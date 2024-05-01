package tokens

import (
	"Ecommerce-Backend/database"
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDeatils struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstName string, lastName string, uid string) (signedToken string, signedrefreshtoken string, err error) {

	claim := &SignedDeatils{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Duration(24)).Unix(),
		},
	}

	refreshcliams := &SignedDeatils{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claim).SignedString([]byte(SECET_KEY))

	if err != nil {
		return "", "", err
	}
	refresToken, err := jwt.NewWithClaims(jwt.SigningMethodES384, refreshcliams).SignedString([]byte(SECET_KEY))

	if err != nil {
		return "", "", err
	}
	return token, refresToken, err
}

func ValidateToken(signedtoken string) (claims *SignedDeatils, msg string) {

	token, err := jwt.ParseWithClaims(signedtoken, &SignedDeatils{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
	}
	claims, ok := token.Claims.(*SignedDeatils)

	if !ok {
		msg = "invalid token"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}

func UpdateAllToken(signedtoken string, signedrefreshToken, userId string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObject primitive.D
	updateObject = append(updateObject, bson.E{Key: "token", Value: signedtoken})
	updateObject = append(updateObject, bson.E{Key: "refresh_token", Value: signedrefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updatedAt", Value: updatedAt})

	upsert := true

	usertId, err := primitive.ObjectIDFromHex(userId)

	log.Println(usertId)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"user_id": userId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = UserData.UpdateOne(ctx, filter, bson.D{{
		Key: "$set", Value: updateObject,
	}}, &opt)

	if err != nil {
		upsert = false
		panic(upsert)
	}

	defer cancel()

}
