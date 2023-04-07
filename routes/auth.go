package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NuthanReddy45/ELabourWebApp/models"
	"github.com/NuthanReddy45/ELabourWebApp/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var DB *mongo.Collection

type LoginData struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerData struct {
	Email     string `json:"email,omitempty"`
	PhnNo     string `json:"phnNo,omitempty"`
	Password1 string `json:"password1,omitempty"`
	Password2 string `json:"password2,omitempty"`
}

func Login(c *gin.Context) {

	var temp LoginData
	err := c.ShouldBindJSON(&temp)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}

	DB = util.Db

	filter := bson.D{{Key: "email", Value: temp.Email}}

	var res models.User

	err = DB.FindOne(context.TODO(), filter, nil).Decode(&res)

	if err == nil {

		err1 := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(temp.Password))

		if err1 == nil {

			fmt.Println("in login id = ", res.ID)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": res.ID,
				"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			sec := []byte(os.Getenv("SECRET_KEY"))

			tokenString, err2 := token.SignedString(sec)

			if err2 != nil {
				c.String(500, "error signing Jwt token")
				return
			}

			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Auth", tokenString, 3600*24*3, "", "", false, true)

			c.IndentedJSON(http.StatusOK, nil)
			return
		}
		c.String(400, "Invalid credits ")
		return
	}

	if err != mongo.ErrNoDocuments {
		c.String(500, "Server Error ")
		return
	}

	c.String(400, "Invalid credits 2 ")
}

func Register(c *gin.Context) {

	var temp registerData
	err := c.ShouldBindJSON(&temp)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}

	if temp.Password1 != temp.Password2 {
		c.String(400, "Passwords do not match")
		return
	}

	DB = util.Db
	var res models.User

	filter := bson.D{{Key: "email", Value: temp.Email}}
	err = DB.FindOne(context.Background(), filter).Decode(&res)

	if err == nil {
		c.String(400, "user already exsits")
		return
	}
	if err != mongo.ErrNoDocuments {
		c.String(500, "Server Error ")
		return
	}

	// hash the password

	hash, errHasing := bcrypt.GenerateFromPassword([]byte(temp.Password1), 10)

	if errHasing != nil {
		c.String(500, "error in hashing password")
		return
	}

	curData := models.User{Email: temp.Email, PhnNo: temp.PhnNo, Password: string(hash)}

	inserted, err := DB.InsertOne(context.Background(), curData, nil)

	if err != nil {
		log.Fatal(err)
		c.String(500, "sry ..")
		return
	}
	fmt.Println("id in generating tkn ", inserted.InsertedID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": inserted.InsertedID,
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	sec := []byte(os.Getenv("SECRET_KEY"))

	tokenString, err2 := token.SignedString(sec)

	if err2 != nil {
		c.String(500, "error signing Jwt token")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 3600*24*3, "", "", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{
		"Inserted_id": inserted,
	})
}
