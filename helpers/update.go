package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
)

// Update - Generate a handler to handle update request
func Update(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance, _ := c.Get(name)

		if updateable, ok := instance.(db.Updatable); ok {
			fields := updateable.UpdatableFields()
			msi := map[string]interface{}{}

			c.ShouldBindBodyWith(&msi, binding.JSON)

			updatedFields := map[string]interface{}{}
			for _, val := range fields {
				if msi[val] != nil {
					updatedFields[val] = msi[val]
				}
			}

			filterdData, _ := json.Marshal(updatedFields)

			c.Set(gin.BodyBytesKey, filterdData)
		} else {
			ErrorResponse(c, http.StatusForbidden, "Can not update this resource")
		}

		if err := c.ShouldBindBodyWith(instance, binding.JSON); err == nil {
			if err := db.DB().Save(instance).Error; err == nil {
				c.JSON(http.StatusOK, instance)
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		} else {
			ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		}
	}
}
