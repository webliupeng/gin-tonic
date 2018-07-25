package helpers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/helpers"
)

func TestDelete(t *testing.T) {
	item := &Item{}
	db.DB().Last(&item)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/list/%d", item.ID), nil)
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 204, record.Code)
}

func init() {
	R = testRouter()

	R.DELETE("/list/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.Delete("item"),
	)
}
