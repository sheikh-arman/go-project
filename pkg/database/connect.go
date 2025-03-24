package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

var (
	dataInserted int
	mutex        sync.Mutex
)

// User model
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
	Age  int
}

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
	DatabaseLoadTest(db)
}

func showDatabases(db *gorm.DB) {
	var databaseList []string
	db.Raw("show databases").Scan(&databaseList)
	fmt.Println("databaseList:-")
	for i, database := range databaseList {
		fmt.Println(i, database)
	}
}

func DatabaseLoadTest(db *gorm.DB) {

	//ReadTest(db)
	WriteTest(db)
}

func ReadTest(db *gorm.DB) {
	// Number of concurrent requests
	numWorkers := 100000
	wg := sync.WaitGroup{}

	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			testQuery(db, id)
		}(i)
	}

	wg.Wait()
	fmt.Println("Total time taken for read:", time.Since(start))
}

func testQuery(db *gorm.DB, id int) {
	var count int64
	result := db.Raw("SELECT COUNT(*) FROM user").Scan(&count)
	if result.Error != nil {
		log.Println("Query failed:", result.Error)
	} else {
		log.Printf("Worker %d: User count = %d\n", id, count)
	}
}

func WriteTest(db *gorm.DB) {

	err := db.Exec("create database if not exists testDB").Error
	if err != nil {
		log.Fatal(err)
		return
	}

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
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, "testDB")
	fmt.Println("connecting to database...", connectionString)
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Number of concurrent writes
	numWorkers := 100
	numRecords := 1000 // Total records to insert
	wg := sync.WaitGroup{}

	start := time.Now()
	dataInserted = 0
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			insertRecords(db, workerID, numRecords)
		}(i)
	}

	wg.Wait()
	fmt.Println("Total time taken for write:", time.Since(start))
	fmt.Println("Total dataInserted", dataInserted)
}

func insertRecords(db *gorm.DB, workerID, num int) {
	for i := 0; i < num; i++ {
		user := User{Name: fmt.Sprintf("User_%d_%d", workerID, i), Age: 20 + (i % 10)}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Worker %d: Insert failed: %v", workerID, err)
		}
		mutex.Lock()
		dataInserted = dataInserted + 1
		mutex.Unlock()

	}
}
