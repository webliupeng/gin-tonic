package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/webliupeng/gin-tonic/utils"
)

var db *gorm.DB

func init() {
	config := utils.GetConfig()
	var err error

	cs := os.Getenv("DB_URI")
	if cs == "" {
		cs = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true", config.Db.User, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.Name)
	}

	db, err = gorm.Open("mysql", cs)

	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(2)
	db.DB().SetConnMaxLifetime(time.Duration(10) * time.Minute)
	if err != nil {
		fmt.Println("bad connection", cs)
		panic(err)
	}
}

func DB() *gorm.DB {

	return db
}
