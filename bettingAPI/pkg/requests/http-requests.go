package requests

import (
	"bettingAPI/pkg/config"
	"bettingAPI/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var db = config.GetDB()

func GetLeagueOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rowsLeagueOffers, err := db.Query("SELECT title  FROM `bettingdb`.`leagues`")
	if err != nil {
		panic(err.Error())
	}
	defer rowsLeagueOffers.Close()

	rowsLeagueOffer, err := db.Query("SELECT title, offer_id FROM `bettingdb`.`leagues`, `bettingdb`.`league_offers` where leagues.league_id = league_offers.league_id")
	if err != nil {
		panic(err.Error())
	}
	defer rowsLeagueOffer.Close()

	var leagueOffer []models.LeagueOffer
	var leagueOffers []models.LeagueOffers

	for rowsLeagueOffers.Next() {
		var l models.LeagueOffers
		if err := rowsLeagueOffers.Scan(&l.Name); err != nil {
			panic(err.Error())
		}
		leagueOffers = append(leagueOffers, l)
	}

	for rowsLeagueOffer.Next() {
		var temp models.LeagueOffer
		if err := rowsLeagueOffer.Scan(&temp.Name, &temp.Offer); err != nil {
			panic(err.Error())
		}
		leagueOffer = append(leagueOffer, temp)
	}
	for i := range leagueOffer {
		for j := range leagueOffers {
			if leagueOffer[i].Name == leagueOffers[j].Name {
				leagueOffers[j].Offers = append(leagueOffers[j].Offers, leagueOffer[i].Offer)
			}
		}
	}
	json.NewEncoder(w).Encode(leagueOffers)

}

func GetOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	rowsOffers, err := db.Query("SELECT * from `bettingdb`.`offers`")
	if err != nil {
		panic(err.Error())
	}
	defer rowsOffers.Close()
	var offers []models.Offer

	for rowsOffers.Next() {
		var offer models.Offer
		if err := rowsOffers.Scan(&offer.ID, &offer.Name, &offer.Time, &offer.TvChannel, &offer.HasStatistics); err != nil {
			panic(err.Error())
		}
		offers = append(offers, offer)
	}

	for _, item := range offers {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode("No offer found with given id")
}
