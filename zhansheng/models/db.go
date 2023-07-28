package models

import (
	"danfwing.com/m/zhansheng/config/global"
	config2 "danfwing.com/m/zhansheng/models/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	GDB *gorm.DB
)

func InitDB() error {
	var err error
	if global.Env == "dev" {
		GDB, err = gorm.Open("sqlite3", "./ipsec.sqlite")
		if err != nil {
			logrus.Errorln("gorm openDb is fail err")
			return err
		}
	} else {
		key := "danfwing.com"
		// set DB name
		dbname := "./ipsec.sqlite"
		dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=x'%s'",
			key)
		GDB, err = gorm.Open("sqlite3", dbnameWithDSN)
		if err != nil {
			logrus.Errorln("gorm openDb is fail err")
			return err
		}
		//@@TODO 判断是否加密 否则退出
	}

	global.GDB = GDB
	GDB.LogMode(true)
	db := global.GDB.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(1024)
	db.SetMaxOpenConns(256)
	db.SetConnMaxLifetime(time.Hour)

	if !global.GDB.HasTable(&config2.SystemConfig{}) {
		global.GDB.AutoMigrate(&config2.SystemConfig{})
		err = global.GDB.Create(&config2.SystemConfig{
			Key:          "frpc",
			Name:         "frpc配置",
			Value:        "{\"frpc\":{\"common\":{\"server_addr\":\"127.0.0.1\",\"server_port\":\"17555\",\"admin_addr\":\"127.0.0.1\",\"admin_port\":\"30005\",\"token\":\"123456\"},\"ssh\":{\"type\":\"tcp\",\"local_ip\":\"127.0.0.1\",\"local_port\":\"17111\",\"remote_port\":\"5700\"}}}",
			ValueExplain: "frpc配置的json格式",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}).Error
		if err != nil {
			panic("create config is wrong")
		}
	}
	return nil
}
