package main

import (
	bethttp "bettingAPI/internal/http"
	"bettingAPI/internal/mysql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	GetAllOffersFromHTTP()
	var leagues = GetAllLeaguesFromHTTP()
	mysql.NewDB().InsertLeagueOffers(leagues)

	router := mux.NewRouter()
	router.HandleFunc("/league-offers", bethttp.GetLeagueOffers).Methods("GET")
	router.HandleFunc("/offer/{id}", bethttp.GetOffer).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:5000", router))

}

func GetAllLeaguesFromHTTP() *mysql.LeaguesData {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"
	var leagues mysql.LeaguesData
	err := bethttp.GetJson(url, &leagues)
	if err != nil {
		fmt.Printf("Error getting leagues: %s\n", err.Error())
	} else {
		fmt.Printf("Leagues from web have been reached\n")
	}

	mysql.NewDB().InsertLeagues(&leagues)
	return &leagues
}

func GetAllOffersFromHTTP() {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json"
	var offers []mysql.Offer
	err := bethttp.GetJson(url, &offers)
	if err != nil {
		fmt.Printf("Error getting offers: %s\n", err.Error())
	} else {
		fmt.Printf("Offers from web have been reached\n")
	}
	mysql.NewDB().InsertOffers(offers)
	mysql.NewDB().InsertTips(offers)
}
