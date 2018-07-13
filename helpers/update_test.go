package helpers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/helpers"
)

func TestUpdate(t *testing.T) {

	item := &Item{}
	db.DB().First(&item)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/list/%d", item.ID), bytes.NewReader([]byte(`{
		"foo":"REPLY"
	}`)))
	record := httptest.NewRecorder()

	R.ServeHTTP(record, req)

	assert.Equal(t, 200, record.Code)
}

func init() {
	R = testRouter()

	R.PUT("/list/:id",
		helpers.FindOneByParam(&Item{}, "id", "item"),
		helpers.Update("item"),
	)
}
