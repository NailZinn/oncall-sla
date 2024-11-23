package sla

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func ConnectToDb() {
	mysqlCfg := viper.GetStringMapString("mysql")
	cfg := mysql.Config{
		User:                 mysqlCfg["user"],
		Passwd:               mysqlCfg["password"],
		Net:                  mysqlCfg["net"],
		Addr:                 mysqlCfg["address"],
		DBName:               mysqlCfg["dbname"],
		AllowNativePasswords: viper.GetBool("mysql.allownativepasswords"),
	}

	var openErr error
	db, openErr = sql.Open("mysql", cfg.FormatDSN())
	if openErr != nil {
		log.Fatalln(openErr)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalln(pingErr)
	}

	fmt.Println("[INFO]: connected to db")
}

func InitSliTable() {
	if db == nil {
		log.Fatalln("[FATAL]: connection to db is not established")
	}

	dbName := viper.GetString("mysql.dbname")
	tableName := viper.GetString("mysql.tablename")

	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))

	db.Exec(fmt.Sprintf("USE %s;", dbName))

	_, createErr := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			value FLOAT(4) NOT NULL,
			slo FLOAT(4) NOT NULL,
			datetime DATETIME NOT NULL DEFAULT NOW(),
			isBad BOOL NOT NULL DEFAULT false
		);
	`, tableName))

	if createErr != nil {
		log.Fatalln(createErr)
	}
}

func saveSli(name string, value, slo float32, datetime int32, isBad bool) {
	datetimeFormatted := time.Unix(int64(datetime), 0).Format(time.DateTime)
	tableName := viper.GetString("mysql.tableName")
	res, insertErr := db.Exec(fmt.Sprintf(`
		INSERT INTO %s (name, value, slo, datetime, isBad)
		VALUES (?, ?, ?, ?, ?);
	`, tableName), name, value, slo, datetimeFormatted, isBad)

	if insertErr != nil {
		fmt.Printf("[ERROR]: %s\n", insertErr.Error())
		return
	}

	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		fmt.Println("[ERROR]: no rows affected")
	}
}
