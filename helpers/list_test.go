package helpers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/objx"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/helpers"
)

type Item struct {
	gorm.Model
	Foo    string
	Bar    string
	UserID string
}

func (i Item) UpdatableFields() []string {
	return []string{"foo", "bar"}
}

func (i Item) FilterableFields() []string {
	return []string{"foo", "bar", "id"}
}

var R *gin.Engine

func TestList(t *testing.T) {

	req, _ := http.NewRequest("GET", "/list?.offset=10&.maxResults=10", nil)

	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	//println("limited", record.Body.String())
	obj, _ := objx.FromJSON(record.Body.String())

	assert.Equal(t, 10, len(obj.Get("data").InterSlice()))
	assert.Equal(t, record.Code, 200)
}

func TestIncludes(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/1/list?.includes=account", nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, record.Code, 200)
}

func TestLike(t *testing.T) {
	req, _ := http.NewRequest("GET", "/list?bar_like=test_*", nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, record.Code, 200)
}

func TestGetOne(t *testing.T) {

	item := &Item{}
	item.Foo = "hah"
	db.DB().Save(&item)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/list/%v", item.ID), nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, record.Code, 200)
}

func TestGranterThan(t *testing.T) {
	req, _ := http.NewRequest("GET", "/list?id_gt=2", nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, record.Code, 200)
}

func TestLessThan(t *testing.T) {
	req, _ := http.NewRequest("GET", "/list?id_lt=2&id_gt=3", nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, record.Code, 200)
}

func testRouter() *gin.Engine {
	if R == nil {
		R = gin.Default()
	}
	return R
}
func init() {
	R = testRouter()
	db.DB().DropTable(&Item{})
	db.DB().AutoMigrate(&Item{})

	for i := 0; i < 50; i++ {
		item := &Item{}
		item.Bar = fmt.Sprintf("bar%v", i)
		db.DB().Save(&item)
	}

	item := &Item{}
	item.Bar = "test_abc"
	db.DB().Save(&item)
	R.GET("/list", helpers.List(&Item{}))
	R.GET("/list/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.ServeJSONFromContext("item"),
	)
	R.GET("/users/:userid/list",
		helpers.List(&Item{}, helpers.CriteriaByParam("userId")),
	)
}
