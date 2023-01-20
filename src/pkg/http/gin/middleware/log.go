package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LogErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		errs := c.Errors

		for _, err := range errs {
			log.Println(err.Error())
		}
	}
}
