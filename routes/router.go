package routes

import (
	"net/http"

	"github.com/NuthanReddy45/ELabourWebApp/middlewares"
	"github.com/gin-gonic/gin"
)

func Landing(c *gin.Context) {

	res, _ := c.Get("user")
	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": res,
	})

}

func Router() *gin.Engine {

	r := gin.New()
	r.Use(middlewares.CORSMiddleware())
	r.GET("/", middlewares.Auth, Landing)
	r.POST("/login", Login)
	r.POST("/register", Register)
	return r
}
