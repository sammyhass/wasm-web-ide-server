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

	// check if already written
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": msgs,
		})
	}
}
