package http

import (
	"bettingAPI/internal/mysql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetLeagueOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var leagueOffers = mysql.NewDB().GetLeagueOffers()
	json.NewEncoder(w).Encode(leagueOffers)

}

func GetOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	offerID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Error converting id from string to int: %s", err)
	}
	offer := mysql.NewDB().GetOfferByID(offerID)
	json.NewEncoder(w).Encode(offer)
}
