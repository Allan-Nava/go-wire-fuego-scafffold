package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Allan-Nava/go-wire-fuego-scafffold/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(config *env.Configuration) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connection on %s, err: %s", dsn, err.Error())
	}
	db, _ := conn.DB()
	idle, _ := strconv.Atoi(config.DbIdleConn)
	max, _ := strconv.Atoi(config.DbMaxConn)

	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(max)

	return conn
}
