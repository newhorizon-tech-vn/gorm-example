package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

func initDBClient() (*gorm.DB, error) {
	username := "root"
	password := "my-secret-pw"
	addr := "127.0.0.1:33061"
	dbname := "warehousedb"
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True",
		username, password, addr, dbname)

	client, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	client = client.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")
	return client, nil
}

var DBClient *gorm.DB

func main() {
	var err error
	DBClient, err = initDBClient()
	if err != nil {
		fmt.Println("ERROR", "init db client failed", err)
		return
	}

	// testHasMany()

	testManyToMany()
}
