package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/helpers"
)

func TestFindOne(t *testing.T) {
	req, _ := http.NewRequest("GET", "/mylist/123123123", nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 404, record.Code)
}

func init() {
	R = testRouter()
	R.GET("/mylist/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.ServeJSONFromContext("item"),
	)

}
