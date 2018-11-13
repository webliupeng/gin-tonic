package db_test

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/webliupeng/gin-tonic/db"
)

func TestDB(t *testing.T) {
	db.DB()
	db.DB()
}

type TestItem struct {
	gorm.Model
	List db.JSONArray `gorm:"type:text"`
}

func TestJSONArray(t *testing.T) {
	ti := TestItem{}

	db.DB().AutoMigrate(TestItem{})
	ti.List = db.JSONArray{"1", "2"}
	db.DB().Save(&ti)

	t2 := TestItem{}

	db.DB().First(&t2)

	fmt.Printf("t2 %v", t2)
}
