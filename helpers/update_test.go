package helpers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/helpers"
)

func TestUpdate(t *testing.T) {
	item := &Item{}

	db.DB().Where("foo = ?", 1).Where("id > ?", 3).Find(&item)
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list/%d", item.ID), bytes.NewReader([]byte(`{
		"foo":"REPLY",
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 200, record.Code)
}

func TestUpdateWithCustomeResponse(t *testing.T) {

	item := &Item{}

	db.DB().Where("foo = ?", 1).Where("id > ?", 3).Find(&item)
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list2/%d", item.ID), bytes.NewReader([]byte(`{
		"foo":"REPLY",
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 200, record.Code)
}

func TestUpdateWithBadBody(t *testing.T) {

	item := &Item{}

	db.DB().Where("foo = ?", 1).Where("id > ?", 3).Find(&item)
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list/%d", item.ID), bytes.NewReader([]byte(`{
		"fooL
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, http.StatusBadRequest, record.Code)
}

func TestUpdateTypeMatch(t *testing.T) {

	item := &Item{}

	db.DB().Where("foo = ?", 1).Where("id > ?", 3).Find(&item)
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list/%d", item.ID), bytes.NewReader([]byte(`{
		"age": "asfas"
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, http.StatusUnprocessableEntity, record.Code)
}

func TestUpdateAUnupdatable(t *testing.T) {

	item := &Item{}

	db.DB().Where("foo = ?", 1).Where("id > ?", 3).Find(&item)
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list3/%d", item.ID), bytes.NewReader([]byte(`{
		"user_id": 1
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, http.StatusForbidden, record.Code)
}

type ItemDisallowUpdate struct {
	gorm.Model
	Foo    string
	Bar    string
	UserID int
	User   *User
}

func (i *ItemDisallowUpdate) TableName() string {
	return "items"
}

func init() {
	R = testRouter()

	R.PUT("/list/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.Update("item"),
	)

	R.PUT("/list2/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.Update("item"),
		func(c *gin.Context) {
			c.Get("updated")
		},
	)

	R.PUT("/list3/:id",
		helpers.FindOneByParam(&ItemDisallowUpdate{}, "id", "item"),
		helpers.Update("item"),
		func(c *gin.Context) {
			c.Get("updated")
		},
	)
}
