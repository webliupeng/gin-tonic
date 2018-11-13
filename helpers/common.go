package helpers

import (
	"errors"
	"reflect"

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

func contain(target interface{}, obj interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	}

	return false, errors.New("not in array")
}
