package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	msgs := []string{}
	for _, err := range c.Errors {
		fmt.Println("Error: ", err.Error())
		msgs = append(msgs, err.Error())
	}

	if !c.Writer.Written() && len(msgs) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"info":  msgs,
		})
	}
}
