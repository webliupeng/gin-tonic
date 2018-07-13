package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtConfig(t *testing.T) {

	defer os.RemoveAll("./configs")
	_ = os.Mkdir("./configs", 0777)

	text := `{
		"app": {
			"port": "8001"
		},
		"db": {
			"name": "test",
			"host":"localhost",
			"user": "xxx",
			"password": "C",
			"port": "3306"
		},
		"redis": {
			"host": "localhost",
			"port": "6379",
			"password": "xxx"
		},
		"es": {
			"url": "http://localhost/"
		},
		"ext": {
			"foo": "bar",
			"num": 1
		}
	}`

	if err := ioutil.WriteFile("./configs/config.test.json", []byte(text), 0644); err != nil {
		t.Error(err)
	}

	config := GetConfig()
	assert.Equal(t, "bar", config.GetExt("foo").Str())
	assert.Equal(t, float64(1), config.GetExt("num").Float64())
}
