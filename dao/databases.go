package dao

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Server string
	Port string
	DbName     string
	DbUser     string
	DbPassword string
}

var db *gorm.DB

func InitDb() bool {
    var err error
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var dbConfiguration DatabaseConfiguration
	if err = viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file.", err)
		return false
	}
	err = viper.Unmarshal(&dbConfiguration)
	if err != nil {
		log.Println("Unable to decode config file.", err)
		return false
	}
	db, err = gorm.Open("mysql", dbConfiguration.DbUser + ":" + dbConfiguration.DbPassword + "@tcp(" + dbConfiguration.Server + ":" + dbConfiguration.Port + ")/" + dbConfiguration.DbName + "?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection Failed to Open.", err)
		return false
	} else {
		log.Println("Connection Established.")
		db.AutoMigrate(&Book{}, &Comment{})
		return true
	}
}