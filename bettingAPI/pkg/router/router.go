package router

import (
	"bettingAPI/pkg/requests"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/league-offers", requests.GetLeagueOffers).Methods("GET")
	router.HandleFunc("/offer/{id}", requests.GetOffer).Methods("GET")

	return router
}
