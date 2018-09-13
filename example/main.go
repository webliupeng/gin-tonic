package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/webliupeng/gin-tonic/utils"
)

func main() {

	viper.SetEnvPrefix("EXA")
	utils.GetConfig()

	e := gin.Default()

	utils.Redis().Get("ab")

	e.Run(":9911")
}

func init() {
	viper.SetEnvPrefix("EXA")
}
