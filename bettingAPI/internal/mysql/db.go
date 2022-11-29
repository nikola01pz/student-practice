package mysql

import (
	"database/sql"
	"fmt"
	"log"
)

type db struct {
	conn *sql.DB
}

func NewDB() *db {
	return &db{
		conn: ConnectDB(),
	}
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Lozinka123#@tcp(localhost:3306)/bettingdb")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error connecting")
		panic(err.Error())
	} else {
		fmt.Println("Connected to DB")
	}
	return db
}

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

type LeagueOffers struct {
	ID     int    `db:"league_id" json:"id"`
	Title  string `db:"title" json:"title"`
	Offers []int  `json:"offers"`
}

type OfferByID struct {
	ID            int    `json:"id"`
	Name          string `json:"game"`
	Time          string `json:"time"`
	TvChannel     string `json:"tv_channel"`
	HasStatistics bool   `json:"statistics"`
}

func (d *db) InsertOffers(offers []Offer) {
	query := "INSERT INTO `bettingdb`.`offers`(offer_id, game, time_played, tv_channel, has_statistics) VALUES (?,?,?,?,?)"

	for i := range offers {
		_, err := d.conn.Exec(query, offers[i].ID, offers[i].Name, offers[i].Time, offers[i].TvChannel, offers[i].HasStatistics)
		if err != nil {
			log.Fatalf("impossible to insert leagues: %s", err)
		}
	}
	fmt.Printf("Leagues inserted into DB\n")
}

func (d *db) InsertLeagues(leagues *LeaguesData) {
	query := "INSERT INTO `bettingdb`.`leagues`(title) VALUES (?)"

	for i := range leagues.Leagues {
		_, err := d.conn.Exec(query, leagues.Leagues[i].Name)
		if err != nil {
			log.Fatalf("impossible to insert offers: %s", err)
		}
	}
	fmt.Printf("Offers inserted into DB\n")
}

func (d *db) InsertTips(offers []Offer) {
	query := "INSERT INTO `bettingdb`.`offer_tips` VALUES (?,?,?)"

	for i := range offers {
		for j := range offers[i].Tips {
			_, err := d.conn.Exec(query, offers[i].ID, offers[i].Tips[j].Name, offers[i].Tips[j].Value)
			if err != nil {
				log.Fatalf("impossible to insert tips: %s", err)
			}
		}
	}
}

func (d *db) InsertLeagueOffers(leagues *LeaguesData) {
	query := "INSERT INTO `bettingdb`.`league_offers` (league_id, offer_id) SELECT league_id, offer_id FROM `bettingdb`.`leagues`,`bettingdb`.`offers` WHERE `bettingdb`.`offers`.`offer_id`=? and `bettingdb`.`leagues`.`title`=?"
	for i := range leagues.Leagues {
		for j := range leagues.Leagues[i].Elaborations {
			for k := range leagues.Leagues[i].Elaborations[j].ID {
				_, err := d.conn.Exec(query, leagues.Leagues[i].Elaborations[j].ID[k], leagues.Leagues[i].Name)
				if err != nil {
					log.Fatalf("impossible to insert league_offers: %s", err)
				}
			}
		}
	}
}

func (d *db) GetLeagueOffers() *[]LeagueOffers {
	rowsLeagues, err := d.conn.Query("SELECT * FROM `bettingdb`.`leagues`")
	if err != nil {
		panic(err.Error())
	}
	defer rowsLeagues.Close()
	var leagueOffers []LeagueOffers
	for rowsLeagues.Next() {
		var league LeagueOffers
		if err := rowsLeagues.Scan(&league.ID, &league.Title); err != nil {
			panic(err.Error())
		}
		leagueOffers = append(leagueOffers, league)
	}
	for i := range leagueOffers {
		rowsOffers, err := d.conn.Query("SELECT offer_id FROM `bettingdb`.`league_offers` WHERE `bettingdb`.`league_offers`.`league_id`=?", leagueOffers[i].ID)
		if err != nil {
			panic(err.Error())
		}
		defer rowsOffers.Close()
		for rowsOffers.Next() {
			var Offer int
			if err := rowsOffers.Scan(&Offer); err != nil {
				panic(err.Error())
			}
			leagueOffers[i].Offers = append(leagueOffers[i].Offers, Offer)
		}
	}
	return &leagueOffers
}

func (d *db) GetOfferByID() *[]OfferByID {
	rowsOffers, err := d.conn.Query("SELECT * from `bettingdb`.`offers`")
	if err != nil {
		panic(err.Error())
	}
	defer rowsOffers.Close()
	var offers []OfferByID

	for rowsOffers.Next() {
		var offer OfferByID
		if err := rowsOffers.Scan(&offer.ID, &offer.Name, &offer.Time, &offer.TvChannel, &offer.HasStatistics); err != nil {
			panic(err.Error())
		}
		offers = append(offers, offer)
	}
	return &offers
}
