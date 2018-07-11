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
			"name": "tuice_test",
			"host":"10.10.4.22",
			"user": "xxx",
			"password": "C",
			"port": "3306"
		},
		"redis": {
			"host": "10.2.1.214",
			"port": "6379",
			"password": "xxx"
		},
		"open": {
			"url": "http://dev.open.admaster.co/"
		},
		"es": {
			"url": "http://10.10.3.13:9100/"
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
