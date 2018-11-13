package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/webliupeng/gin-tonic/db"
)

type ModelInstanceCreator func(c *gin.Context) interface{}

// Create - Create a handler to handle create model instance. it will use Gin's BindJSON
func Create(instanceCreator ModelInstanceCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelInstance := instanceCreator(c)
		if modelInstance == nil {
			return
		}
		if insertable, ok := modelInstance.(db.Insertable); !ok {
			ErrorResponse(c, http.StatusForbidden, "Can not write this resource")
			return
		} else {
			fields := insertable.InsertableFields()
			msi := map[string]interface{}{}

			if err := c.ShouldBindBodyWith(&msi, binding.JSON); err != nil {
				ErrorResponse(c, http.StatusBadRequest, err.Error())
				return
			}

			insertFields := map[string]interface{}{}
			for key := range msi {
				if ok, _ := contain(fields, key); ok {
					insertFields[key] = msi[key]
				} else {
					fmt.Println("[warning]", key, "field does not allow inserts")
				}
			}
			filterdData, _ := json.Marshal(insertFields)

			c.Set(gin.BodyBytesKey, filterdData)
		}

		if err := c.ShouldBindBodyWith(modelInstance, binding.JSON); err == nil {
			if err := db.DB().Save(modelInstance).Error; err == nil {
				c.JSON(http.StatusCreated, &modelInstance)
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		}
	}
}
