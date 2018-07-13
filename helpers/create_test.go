package helpers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/helpers"
)

func TestCreate(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list", bytes.NewReader([]byte(`{
		"foo":"REPLY"
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 201, record.Code)
}

func TestCreateWithBadJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list", bytes.NewReader([]byte(`{
		"foo":"REPLY",
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 400, record.Code)
}

func init() {
	R = testRouter()

	item := &Item{}
	item.Bar = "test_abc"
	db.DB().Save(&item)

	R.POST("/list", helpers.Should(func(c *gin.Context) bool {
		return true
	}), helpers.Create(func(c *gin.Context) interface{} {
		item := &Item{}
		item.Bar = "haha"
		return item
	}))
}
