package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
)

type ModelInstanceCreator func(c *gin.Context) interface{}

// Create - Create a handler to handle create model instance. it will use Gin's BindJSON
func Create(instanceCreator ModelInstanceCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelInstance := instanceCreator(c)

		if err := c.ShouldBindJSON(modelInstance); err == nil {
			if err := db.DB().Save(modelInstance).Error; err == nil {
				c.JSON(http.StatusCreated, &modelInstance)
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		} else {
			ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
	}
}
