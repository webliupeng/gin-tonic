package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
)

func Delete(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _ := c.Get(name)
		if err := db.DB().Delete(s).Error; err == nil {
			c.Status(204)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
