package middlewares

import "github.com/gin-gonic/gin"

type prototype func(c *gin.Context)

const isAuth = true

func Tm(c *gin.Context, next prototype) {

	if isAuth {
		next(c)
		return
	}

	c.String(400, "unauth...")

}
