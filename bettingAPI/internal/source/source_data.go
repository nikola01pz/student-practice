package source

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LeaguesData struct {
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

type Tip struct {
	Name  string  `json:"naziv"`
	Value float64 `json:"tecaj,omitempty"`
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

func GetAllLeaguesFromSource() *LeaguesData {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"
	var leagues *LeaguesData
	err := getJson(url, &leagues)
	if err != nil {
		log.Printf("Impossible to get leagues from source: %s", err)
	}
	return leagues
}

func GetAllOffersFromSource() []Offer {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json"
	var offers []Offer
	err := getJson(url, &offers)
	if err != nil {
		log.Printf("Impossible to get offers from source: %s", err)
	}
	return offers
}

func getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error getting json from source: %s", err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
