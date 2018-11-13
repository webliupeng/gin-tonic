package utils_test

import (
	"testing"

	"github.com/webliupeng/gin-tonic/utils"
)

func TestRedis(t *testing.T) {

	r := utils.Redis()

	r.Set("a", "b", 1000).Result()
}
