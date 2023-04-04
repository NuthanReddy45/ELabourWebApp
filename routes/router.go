package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/NuthanReddy45/ELabourWebApp/models"
	"github.com/NuthanReddy45/ELabourWebApp/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func Landing(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": "Vokay lets roll bimches ",
	})

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

	curData := models.User{Email: temp.Email, PhnNo: temp.PhnNo, Password: temp.Password1}

	inserted, err := DB.InsertOne(context.Background(), curData, nil)

	if err != nil {
		log.Fatal(err)
		c.String(500, "sry ..")
		return
	}

	c.String(200, "registration successful  %d", inserted.InsertedID)
}

func Login(c *gin.Context) {

	var temp LoginData
	err := c.ShouldBindJSON(&temp)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}

	DB = util.Db

	filter := bson.D{{Key: "email", Value: temp.Email}}

	var res models.User

	err = DB.FindOne(context.TODO(), filter, nil).Decode(&res)

	fmt.Println("temp= ", res)
	if err == nil {

		if res.Password == temp.Password {
			c.String(200, "Logged in ")
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

func Router() *gin.Engine {

	r := gin.Default()
	r.GET("/", Landing)
	r.POST("/login", Login)
	r.POST("/register", Register)
	return r
}
