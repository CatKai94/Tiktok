package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func init() {
	// 配置MySql连接参数
	username := "root"   //账号
	password := "123456" //密码
	host := "127.0.0.1"  //数据库地址
	port := 3306         //数据库端口
	Dbname := "gorm"
	dsn := fmt.Sprintf("%s:%s@tcp($s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln("db connected error", err)
	}

	DB = db
}
