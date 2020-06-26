package main

import (
	"fmt"
	"log"
	"os"
	md "testEFishery/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])
	fmt.Println(psqlInfo)
	db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	conf := make(map[string]string)
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	if port == "" {
		panic("DBPORT environment variable required but not set")
	}
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

func migrate() {
	db.AutoMigrate(&md.User{})
}
