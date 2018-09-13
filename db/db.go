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

}

// DB - DB will run a gorm.DB pointer
func DB() *gorm.DB {
	if db == nil {
		config := utils.GetConfig()
		var err error

		cs := os.Getenv("DB_URI")
		if cs == "" {
			cs = fmt.Sprintf(
				"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true",
				config.GetString("db.user"),
				config.GetString("db.password"),
				config.GetString("db.host"),
				config.GetString("db.port"),
				config.GetString("db.name"))
		}

		db, err = gorm.Open("mysql", cs)

		db.DB().SetMaxOpenConns(100)
		db.DB().SetMaxIdleConns(20)
		db.DB().SetConnMaxLifetime(time.Duration(10) * time.Minute)
		if err != nil {
			fmt.Println("bad connection", cs)
			panic(err)
		}
	}
	return db
}
