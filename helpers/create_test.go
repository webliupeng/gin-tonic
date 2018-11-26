package helpers_test

import (
	"bytes"
	"encoding/json"
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

func TestCreateUninsterable(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list", bytes.NewReader([]byte(`{
		"foo":"test",
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	result := map[string]interface{}{}
	json.Unmarshal(record.Body.Bytes(), &result)

	if val, ok := result["UserID"].(int); ok {
		assert.Equal(t, 0, val)
	}
	assert.Equal(t, 201, record.Code)
}

func TestCreateEmptyFieldValueRequired(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list", bytes.NewReader([]byte(`{
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	result := map[string]interface{}{}
	json.Unmarshal(record.Body.Bytes(), &result)

	if val, ok := result["UserID"].(int); ok {
		assert.Equal(t, 0, val)
	}

	assert.Equal(t, 400, record.Code)
}

func TestCreateWithoutWritable(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list2", bytes.NewReader([]byte(`{
		"foo":"REPLY"
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 403, record.Code)
}

func TestCreateWithBadJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list", bytes.NewReader([]byte(`{
		"foo":"REPLY",
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 400, record.Code)
}

func TestShould(t *testing.T) {
	req, _ := http.NewRequest("POST", "/list3", bytes.NewReader([]byte(`{
		"foo":"REPLY",
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 403, record.Code)
}

type ReadOnlyItem struct {
	Bar string
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

	R.POST("/list2", helpers.Create(func(c *gin.Context) interface{} {
		item := &ReadOnlyItem{}
		item.Bar = "haha"
		return item
	}))

	R.POST("/list3", helpers.Should(func(c *gin.Context) bool {
		return false
	}), helpers.Create(func(c *gin.Context) interface{} {
		item := &ReadOnlyItem{}
		item.Bar = "haha"
		return item
	}))
}
