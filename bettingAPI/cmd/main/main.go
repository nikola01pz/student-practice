package main

import (
	"bettingAPI/pkg/config"
	"bettingAPI/pkg/models"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var db = config.GetDB()
	var leagues = models.GetAllLeagues()
	var offers = models.GetAllOffers()
	models.InsertLeagues(leagues, db)
	models.InsertOffers(offers, db)

}
