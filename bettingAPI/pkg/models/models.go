package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Leagues []League `json:"lige"`
}

type League struct {
	Name         string        `json:"naziv"`
	Elaborations []Elaboration `json:"razrade"`
}

type Elaboration struct {
	Tips []Tip `json:"tipovi,omitempty"`
	ID   []int `json:"ponude,omitempty"`
}

type Offer struct {
	Number        string `json:"broj"`
	ID            int    `json:"id"`
	Name          string `json:"naziv"`
	Time          string `json:"vrijeme"`
	Tips          []Tip  `json:"tecajevi"`
	TvChannel     string `json:"tv_kanal,omitempty"`
	HasStatistics bool   `json:"ima_statistiku,omitempty"`
}

type Tip struct {
	Name  string  `json:"naziv"`
	Value float64 `json:"tecaj,omitempty"`
}

type LeagueOffers struct {
	Name   string `json:"title"`
	Offers []int  `json:"ponude"`
}

type LeagueOffer struct {
	Name  string `json:"title"`
	Offer int    `json:"ponude"`
}

var DB *sql.DB

func GetJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func GetAllLeagues() *Data {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"
	var leagues Data
	err := GetJson(url, &leagues)
	if err != nil {
		fmt.Printf("Error getting leagues: %s\n", err.Error())
	} else {
		fmt.Printf("Leagues from web have been reached\n")
	}
	return &leagues
}

func GetAllOffers() []Offer {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json"
	var offers []Offer
	err := GetJson(url, &offers)
	if err != nil {
		fmt.Printf("Error getting offers: %s\n", err.Error())
	} else {
		fmt.Printf("Offers from web have been reached\n")
	}
	return offers
}

func InsertOffers(offers []Offer, db *sql.DB) {
	query := "INSERT INTO `bettingdb`.`offers`(offer_id, game, time_played, tv_channel, has_statistics) VALUES (?,?,?,?,?)"

	for i := range offers {
		_, err := db.Exec(query, offers[i].ID, offers[i].Name, offers[i].Time, offers[i].TvChannel, offers[i].HasStatistics)
		if err != nil {
			log.Fatalf("impossible to insert leagues: %s", err)
		}
	}
	fmt.Printf("Leagues inserted into DB\n")
}

func InsertLeagues(leagues *Data, db *sql.DB) {
	query := "INSERT INTO `bettingdb`.`leagues`(title) VALUES (?)"

	for i := range leagues.Leagues {
		_, err := db.Exec(query, leagues.Leagues[i].Name)
		if err != nil {
			log.Fatalf("impossible to insert offers: %s", err)
		}
	}
	fmt.Printf("Offers inserted into DB\n")
}

func InsertTips(offers []Offer, db *sql.DB) {
	query := "INSERT INTO `bettingdb`.`offer_tips` VALUES (?,?,?)"

	for i := range offers {
		for j := range offers[i].Tips {
			_, err := db.Exec(query, offers[i].ID, offers[i].Tips[j].Name, offers[i].Tips[j].Value)
			if err != nil {
				log.Fatalf("impossible to insert tips: %s", err)
			}
		}
	}
}

func InsertLeagueOffers(leagues *Data, offers []Offer, db *sql.DB) {
	query := "INSERT INTO `bettingdb`.`league_offers` VALUES (?,?)"
	for i := range leagues.Leagues {
		for j := range leagues.Leagues[i].Elaborations {
			for k := range leagues.Leagues[i].Elaborations[j].ID {
				tempID := leagues.Leagues[i].Elaborations[j].ID[k]
				for n := range offers {
					if tempID == offers[n].ID {
						_, err := db.Exec(query, i+1, tempID)
						if err != nil {
							log.Fatalf("impossible to insert league offers: %s", err)
						}
					}
				}
			}
		}
	}

}
