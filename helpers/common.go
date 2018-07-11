package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// ErrorResponse - Handle erro message
func ErrorResponse(c *gin.Context, code int, msg string) {

	c.Error(errors.New(msg))
	c.JSON(code, gin.H{
		"message": msg,
	})
	c.Abort()
}
