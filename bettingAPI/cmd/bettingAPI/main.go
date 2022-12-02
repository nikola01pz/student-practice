package main

import (
	bethttp "bettingAPI/internal/http"
	"bettingAPI/internal/mysql"
	"bettingAPI/internal/source"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db := mysql.NewDB()
	hdl := bethttp.NewHandler(db)

	offers := source.GetAllOffersFromSource()
	db.InsertOffers(offers)
	db.InsertTips(offers)
	db.InsertLeagues(source.GetAllLeaguesFromSource())

	router := mux.NewRouter()
	router.HandleFunc("/league-offers", hdl.GetLeagueOffers).Methods("GET")
	router.HandleFunc("/offer/{id}", hdl.GetOffer).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:5000", router))
}
