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
		log.Fatal()
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connectiong to DB: %s", err)
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
	query := "INSERT INTO `bettingdb`.`offers`(offer_id, game, time_played, tv_channel, has_statistics) SELECT ?,?,?,?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`offers` WHERE offer_id=?)"

	for i := range offers {
		_, err := d.conn.Exec(query, offers[i].ID, offers[i].Name, offers[i].Time, offers[i].TvChannel, offers[i].HasStatistics, offers[i].ID)
		if err != nil {
			log.Fatalf("impossible to insert offers: %s", err)
		}
	}
	fmt.Printf("Leagues inserted into DB\n")
}

func (d *db) InsertLeagues(leagues *LeaguesData) {
	query := "INSERT INTO `bettingdb`.`leagues`(title) SELECT ? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`leagues` WHERE `bettingdb`.`leagues`.`title`=?)"

	for i := range leagues.Leagues {
		res, err := d.conn.Exec(query, leagues.Leagues[i].Name, leagues.Leagues[i].Name)
		if err != nil {
			log.Fatalf("Impossible to insert leagues: %s", err)
		}
		q := "INSERT INTO `bettingdb`.`league_offers`(league_id, offer_id) SELECT ?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`league_offers` WHERE `bettingdb`.`league_offers`.`offer_id`=?)"
		league_id, err := res.LastInsertId()
		if err != nil {
			log.Fatalf("Impossible to get last insert id: %s", err)
		}
		for j := range leagues.Leagues[i].Elaborations {
			for k := range leagues.Leagues[i].Elaborations[j].ID {
				_, err := d.conn.Exec(q, league_id, leagues.Leagues[i].Elaborations[j].ID[k], leagues.Leagues[i].Elaborations[j].ID[k])
				if err != nil {
					log.Fatalf("Impossible to insert league offer: %s", err)
				}
			}
		}
	}
	fmt.Printf("Offers inserted into DB\n")
}

func (d *db) InsertTips(offers []Offer) {
	query := "INSERT INTO `bettingdb`.`offer_tips`(offer_id, tip, coefficient) SELECT ?,?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`offer_tips` WHERE `bettingdb`.`offer_tips`.`offer_id`=? AND `bettingdb`.`offer_tips`.`tip`=?)"

	for i := range offers {
		for j := range offers[i].Tips {
			_, err := d.conn.Exec(query, offers[i].ID, offers[i].Tips[j].Name, offers[i].Tips[j].Value, offers[i].ID, offers[i].Tips[j].Name)
			if err != nil {
				log.Fatalf("Impossible to insert tips: %s", err)
			}
		}
	}
}

func (d *db) GetLeagueOffers() []LeagueOffers {
	rowsLeagues, err := d.conn.Query("SELECT * FROM `bettingdb`.`leagues`")
	if err != nil {
		log.Fatalf("Impossible to scan from leagues table: %s", err)
	}
	defer rowsLeagues.Close()
	var leagueOffers []LeagueOffers
	for rowsLeagues.Next() {
		var league LeagueOffers
		if err := rowsLeagues.Scan(&league.ID, &league.Title); err != nil {
			log.Fatalf("Impossible to scan from leagues table: %s", err)
		}
		leagueOffers = append(leagueOffers, league)
	}
	for i := range leagueOffers {
		rowsOffers, err := d.conn.Query("SELECT offer_id FROM `bettingdb`.`league_offers` WHERE `bettingdb`.`league_offers`.`league_id`=?", leagueOffers[i].ID)
		if err != nil {
			log.Fatalf("Impossible to select from league_offers table: %s", err)
		}
		defer rowsOffers.Close()
		for rowsOffers.Next() {
			var Offer int
			if err := rowsOffers.Scan(&Offer); err != nil {
				log.Fatalf("Impossible to scan from league_offers table: %s", err)
			}
			leagueOffers[i].Offers = append(leagueOffers[i].Offers, Offer)
		}
	}
	return leagueOffers
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
