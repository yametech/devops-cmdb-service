package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
)

var defaultDb *sql.DB
var dbMutex sync.Mutex

func getSelfDomainConn() *sql.DB {
	log.Println("GetSelfDomainConn")
	DB, err := sql.Open("mysql", "root:!0161fF9346b97f768536A80@tcp(10.202.16.77:3306)/domain_prod")
	if err != nil {
		log.Println("Open  failed,err:", err)
		return nil
	}
	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(20)                   //设置最大连接数
	DB.SetMaxIdleConns(5)                    //设置闲置连接数
	return DB
}

func GetDefaultDBConn() *sql.DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if defaultDb == nil {
		defaultDb = getSelfDomainConn()
	}
	return defaultDb
}
