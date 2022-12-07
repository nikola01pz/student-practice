package http

import (
	"bettingAPI/internal/mysql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct {
	db *mysql.DB
}

func NewHandler(d *mysql.DB) *handler {
	return &handler{
		db: d,
	}
}

func (h *handler) GetLeagueOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var leagueOffers = h.db.GetLeagueOffers()
	json.NewEncoder(w).Encode(leagueOffers)
}

func (h *handler) GetOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	offerID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Error converting id from string to int: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offer := h.db.GetOfferByID(offerID)
	json.NewEncoder(w).Encode(offer)
}
