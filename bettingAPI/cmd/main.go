package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Offer struct {
	Number        string    `json:"broj"`
	ID            int       `json:"id"`
	Name          string    `json:"naziv"`
	Time          time.Time `json:"vrijeme"`
	Tips          []Tip     `json:"tecajevi"`
	TvChannel     string    `json:"tv_kanal,omitempty"`
	HasStatistics bool      `json:"ima_statistiku,omitempty"`
}

type Tip struct {
	Value float64 `json:"tecaj,omitempty"`
	Name  string  `json:"naziv"`
}

type Elaboration struct {
	Tips []Tip `json:"tipovi,omitempty"`
	ID   []int `json:"ponude,omitempty"`
}

type League struct {
	Name         string        `json:"naziv"`
	Elaborations []Elaboration `json:"razrade"`
}

type Data struct {
	Leagues []League `json:"lige"`
}

var client *http.Client

func main() {
	urlOffers := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/ponude.json"
	urlLeagues := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"

	client = &http.Client{Timeout: 10 * time.Second}

	resOffer, err := http.Get(urlOffers)
	CheckNilError(err)
	resData, err := http.Get(urlLeagues)
	CheckNilError(err)

	bodyOffer, err := io.ReadAll(resOffer.Body)
	CheckNilError(err)
	bodyLeague, err := io.ReadAll(resData.Body)
	CheckNilError(err)

	var offers []Offer
	var data Data

	json.Unmarshal(bodyOffer, &offers) //isto vraca error
	json.Unmarshal(bodyLeague, &data)
	// DisplayBets(offers, data)
	// GetAllLeagues()
	ConnectDB()
}

func CheckNilError(err error) {
	if err != nil {
		panic(err)
	}
}

func DisplayBets(offers []Offer, data Data) {
	for i := range data.Leagues {
		fmt.Printf("==== %s ===\n", data.Leagues[i].Name)
		for j := range data.Leagues[i].Elaborations {
			for k := range data.Leagues[i].Elaborations[j].ID {
				for n := range offers {
					if offers[n].ID == data.Leagues[i].Elaborations[j].ID[k] {
						fmt.Printf("\t%s", offers[n].Name)
						fmt.Printf(" | %s\n", offers[n].Time.Format("Monday 15:04"))
						for m := range data.Leagues[i].Elaborations[j].Tips {
							fmt.Printf("\t%s   ", data.Leagues[i].Elaborations[j].Tips[m].Name)
						}
						fmt.Printf("\n")
					}
				}
			}
		}
	}
}

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	CheckNilError(err)

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func GetAllLeagues() {
	url := "https://minus5-dev-test.s3.eu-central-1.amazonaws.com/lige.json"

	var allLeagues Data

	err := GetJson(url, &allLeagues)
	if err != nil {
		fmt.Printf("Error getting JSON: %s\n", err.Error())
	} else {
		for i := range allLeagues.Leagues {
			fmt.Println(allLeagues.Leagues[i].Name)
			for j := range allLeagues.Leagues[i].Elaborations {
				fmt.Println(allLeagues.Leagues[i].Elaborations[j].ID)
			}
		}
	}
}

func ConnectDB() {
	db, err := sql.Open("mysql", "root:Lozinka123#@tcp(localhost:3306)/bettingdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("connected")
}
