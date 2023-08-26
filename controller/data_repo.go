package controller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Init() {
	NgrokHost = "https://d24a-111-49-156-134.ngrok-free.app"
	dsn := "root:123456@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open("mysql", dsn)

	if err != nil {
		fmt.Println("failed to connect database err:", err)
		panic("failed to connect database")
	}
	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Message{})
}

func removeVideoID(slice []int64, id int64) []int64 {
	for i, v := range slice {
		if v == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (v *Video) AddCommentCount(n int64) {
	//user.FavoriteCount += n(n may be 1 or -1)
	v.CommentCount += n
	db.Save(&v) // 写入JSON字符串
}
