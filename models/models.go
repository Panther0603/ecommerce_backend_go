package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id              primitive.ObjectID `json:"id" bson:"_id"`
	First_Name      string             `json:"first_name" bson:"first_name" validate:"required,min=3,max=30"`
	Last_Name       string             `json:"last_name" bson:"last_name" validate:"required,min=3,max=30" `
	Password        string             `json:"password" bson:"password" valiadate:"required,min=6"`
	Email           string             `json:"email" bson:"email" validate:"required"`
	Phone           string             `json:"phone" bson:"phone" validate:"required"`
	Token           string             `json:"token" bson:"token"`
	Refresh_Token   string             `json:"refresh_token" bson:"refresh_token"`
	User_Id         string             `json:"user_id" bson:"uswr_id"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_Date    time.Time          `json:"updated_at" bson:"updated_at"`
	UseCart         []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

type Address struct {
	Address_Id primitive.ObjectID `json:"id" bson:"_id"`
	House      string             `json:"house_number" bson:"house_number"`
	Street     string             `json:"street_name" bson:"street_name"`
	City       string             `json:"city_name" bson:"city_name"`
	Pincode    string             `json:"pincode" bson:"pincode "`
}

type Product struct {
	Product_Id   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name" bson:"product_name"`
	Price        uint64             `json:"price" bson:"price"`
	Rating       uint8              `json:"rating" bson:"rating"`
	Image        string             `json:"image" bson:"image"`
}

type ProductUser struct {
	Product_Id   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name" bson:"product_name"`
	Price        uint64             `json:"price" bson:"price"`
	Rating       uint8              `json:"rating" bson:"rating"`
	Image        string             `json:"image" bson:"iamge"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"productuser" bson:"productuser"`
	Price          uint64             `json:"price" bosn:"price"`
	Discount       uint32             `json:"discount" bson:"discount"`
	Payment_Method Payemnt            `json:"payment_method" bson:"payment_method"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
}

type Payemnt struct {
	Digital bool
	COD     bool
}
