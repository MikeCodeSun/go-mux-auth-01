package model

import (
	"fmt"
	"log"
	"os"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type User struct {
	ID int `json:"id"`
	Name string `json:"name" validate:"required,min=3,max=20"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Created_at time.Time `json:"created_at"`
}

func DbInitilize() *sql.DB{
// env varible
	err := godotenv.Load()
	if err != nil {
		fmt.Println("go dot env")
		log.Fatal(err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
  // format connect mysql address
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName )
	
	// connect mysql
	db , err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		fmt.Println("db open")
		log.Fatal(err)
	}
	fmt.Println("db ok")
	return db
}
// get db
var Db *sql.DB = DbInitilize()
