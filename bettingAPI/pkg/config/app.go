package config

import (
	"bettingAPI/pkg/models"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() {
	db, err := sql.Open("mysql", "root:Lozinka123#@tcp(localhost:3306)/bettingdb")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error connecting")
		panic(err.Error())
	} else {
		fmt.Println("Connected to DB")
	}

	DBClient = db
}

var DBClient *sql.DB

func GetDB() *sql.DB {
	return DBClient
}

func init() {
	ConnectDB()
	var db = GetDB()
	// defer DBClient.Close()
	var leagues = models.GetAllLeagues()
	var offers = models.GetAllOffers()
	models.InsertLeagues(leagues, db)
	models.InsertOffers(offers, db)
	models.InsertTips(offers, db)
	models.InsertLeagueOffers(leagues, offers, db)
}
