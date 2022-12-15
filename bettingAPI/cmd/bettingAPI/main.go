package main

import (
	bethttp "bettingAPI/internal/http"
	"bettingAPI/internal/mysql"
	"bettingAPI/internal/source"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
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
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/league-offers", hdl.GetLeagueOffers).Methods("GET")
	router.HandleFunc("/offer/{id}", hdl.GetOffer).Methods("GET")
	router.HandleFunc("/register", hdl.HandleRegisterRequest).Methods("POST")
	router.HandleFunc("/login", hdl.HandleLoginRequest).Methods("POST")
	router.HandleFunc("/bet_slip", hdl.HandleBetSlipRequest).Methods("POST")
	log.Fatal(http.ListenAndServe("localhost:5000", handlers.CORS(headers, methods, origins)(router)))
}
