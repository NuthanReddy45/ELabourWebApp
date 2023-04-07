package middlewares

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth(c *gin.Context) {

	tokenString, err := c.Cookie("Auth")

	if err != nil {
		c.String(400, "UnAuthorised.. (no tkn) ")
		return
	}

	token, err1 := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{}
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})

	if err1 != nil {
		c.String(400, " Errror in parsing jwt tkn  ")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check if tkn is expired

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// find the user  (userId = claims["sub"])

		var CurUser models.User

		str := fmt.Sprint(claims["sub"])

		req, err3 := primitive.ObjectIDFromHex(str)

		if err3 != nil {
			log.Fatal("not a object id (??) ")
			c.AbortWithStatus(400)
			return
		}

		filter := bson.D{{Key: "_id", Value: req}}
		err = util.Db.FindOne(context.Background(), filter).Decode(&CurUser)

		if err != nil {
			fmt.Println("error ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", CurUser)

		c.Next()
	} else {
		c.String(400, "UnAuthorised... ")
		return
	}

}
