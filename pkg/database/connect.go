package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func DbConnect() {
	user := os.Getenv("MYSQL_ROOT_USER")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	host := os.Getenv("MYSQL_ROOT_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MYSQL_ROOT_PORT")
	if port == "" {
		port = "3306"
	}
	fmt.Println("MYSQL_ROOT_USER:", user, "MYSQL_PASSWORD", password, "MYSQL_HOST", host, "MYSQL_PORT", port)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, "mysql")
	fmt.Println("connecting to database...", connectionString)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	//set pooling parameter
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(time.Hour)
	//query := fmt.Sprintf("create database if not exists arman2")
	//query := fmt.Sprintf("show databases")
	//var val []string
	//db.Exec(query).Scan(&val)
	//fmt.Println(val)
	//for _, val := range val {
	//	fmt.Println(val) =
	//}
	//err = sqlDB.Ping()
	//if err != nil {
	//	fmt.Println("error in ping-> ", err)
	//}
	showDatabases(db)
	fmt.Println("connected to database, hurrah!!!")
}

func showDatabases(db *gorm.DB) {
	var databaseList []string
	db.Raw("show databases").Scan(&databaseList)
	fmt.Println("databaseList:-")
	for i, database := range databaseList {
		fmt.Println(i, database)
	}
}
