package helpers

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/webliupeng/gin-tonic/db"
)

// ServeJSONFromContext serve json from context by name
func ServeJSONFromContext(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if s, e := c.Get(name); e {
			c.JSON(http.StatusOK, s)
		}
	}
}

// FindOneByParam is find a row by specified model type and primary key
func FindOneByParam(modelIns interface{}, paramName string, contextName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if param, ex := c.Params.Get(paramName); ex {
			theType := reflect.TypeOf(modelIns).Elem()
			modelInstance := reflect.New(theType).Elem().Addr().Interface()

			if err := db.DB().Find(modelInstance, "id=?", param).Error; err == nil {
				c.Set(contextName, modelInstance)
			} else if gorm.IsRecordNotFoundError(err) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
		} else {
			c.AbortWithError(http.StatusUnprocessableEntity, errors.New("need a param to query"))
		}
	}
}

// Should -
func Should(args ...func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := false

		for _, checker := range args {
			if checker(c) {
				pass = true
				return
			}
		}
		if !pass {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		}
	}
}
