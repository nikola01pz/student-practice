package http

import (
	"bettingAPI/internal/mysql"
	"encoding/json"
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
	var offers = mysql.NewDB().GetOfferByID()
	for _, item := range *offers {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode("No offer found with given id")
}
