package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
)

// Update - Generate a handler to handle update request
func Update(name string) gin.HandlerFunc {

	var retFunc gin.HandlerFunc
	retFunc = func(c *gin.Context) {
		currentHp := reflect.ValueOf(retFunc).Pointer()
		mainHp := reflect.ValueOf(c.Handler()).Pointer()
		instance, _ := c.Get(name)

		if updateable, ok := instance.(db.Updatable); ok {
			fields := updateable.UpdatableFields()
			msi := map[string]interface{}{}

			if err := c.ShouldBindBodyWith(&msi, binding.JSON); err != nil {
				ErrorResponse(c, http.StatusBadRequest, err.Error())
				return
			}

			updatedFields := map[string]interface{}{}
			for key := range msi {
				if ok, _ := contain(fields, key); ok {
					updatedFields[key] = msi[key]
				} else {
					fmt.Println("[warning]", key, "field does not allow updates")
				}
			}

			filterdData, _ := json.Marshal(updatedFields)

			c.Set(gin.BodyBytesKey, filterdData)
		} else {
			ErrorResponse(c, http.StatusForbidden, "Can not update this resource")
			return
		}

		if err := c.ShouldBindBodyWith(instance, binding.JSON); err == nil {
			if err := db.DB().Save(instance).Error; err == nil {
				if currentHp == mainHp {
					c.JSON(http.StatusOK, instance)
				} else {
					c.Set("updated", instance)
				}
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		} else {
			ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		}
	}
	return retFunc
}
