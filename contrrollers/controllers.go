package contrrollers

import (
	"Ecommerce-Backend/database"
	"Ecommerce-Backend/models"
	"Ecommerce-Backend/tokens"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProdCollection *mongo.Collection = database.UserData(database.Client, "Products")
var Validator = validator.New()

// pass encryption
func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

//VerifyPassword

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))

	if err != nil {
		return false, "Login or passwor is inavlid not matched"
	}
	return true, "here you Goooooooooo"
}

//signup

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.Users
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validator.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// UserCollection := database.UserData(database.Client, "users")
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exixtes with this email id"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()

		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists with this phone number"})
		}

		password := HashPassword(user.Password)
		user.Password = password
		user.Updated_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.User_Id = user.Id.Hex()

		token, refresh_toke, err := tokens.TokenGenerator(user.Email, user.First_Name, user.Last_Name, user.User_Id)
		user.Token = token
		user.Refresh_Token = refresh_toke
		user.UseCart = make([]models.ProductUser, 0)
		user.Order_Status = make([]models.Order, 0)
		user.Address_Details = make([]models.Address, 0)

		_, insertionErr := UserCollection.InsertOne(ctx, user)
		if insertionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user did not get created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "uer created Sucessfully")
	}

}

// login

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.Users
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//	UserCollection := database.UserData(database.Client, "users")
		var foundUser models.Users

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refresh_toke, err := tokens.TokenGenerator(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_Id)
		defer cancel()

		tokens.UpdateAllToken(token, refresh_toke, foundUser.User_Id)
		c.JSON(http.StatusFound, foundUser)
	}
}

// ProductViewrAdmin

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()

		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
		}

		products.Product_Id = primitive.NewObjectID()
		_, anyerr := ProdCollection.InsertOne(ctx, products)

		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong while creating new product")
			c.Abort()
		}

		c.JSON(http.StatusCreated, "product added sucess fully")

	}
}

// SearchProduct

func SearchProduct() gin.HandlerFunc {

	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cusror, err := ProdCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("something went wrong ,please try after sometime"))
			return
		}
		err = cusror.All(ctx, &productList)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("something went wrong ,please try after sometime"))
			return
		}

		if err := cusror.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}
		defer cusror.Close(ctx)
		c.IndentedJSON(http.StatusFound, productList)
	}
}

//SearchProductByQuery

func SearchProductByQuery() gin.HandlerFunc {

	return func(c *gin.Context) {
		var SearchProduct []models.Product

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Search Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchProductByQuery, err := ProdCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			log.Println(err)
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusInternalServerError, database.ErrCantGetItem)
			c.Abort()
			return
		}
		err = searchProductByQuery.All(ctx, &SearchProduct)

		if err != nil {
			log.Println(err)
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusInternalServerError, database.ErrCantGetItem)
			c.Abort()
			return
		}
		defer searchProductByQuery.Close(ctx)

		if err != nil {
			log.Panicln(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusFound, SearchProduct)
	}
}
